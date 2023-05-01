package pool

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/datastore"
	"github.com/rs/zerolog/log"
)

const base62Str = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const sequenceName = "orderid_sequence"
const seqMax uint16 = 3843 //base62 62^2 = 3844
const sequenFileName = "sequence.json"

type sequenceModel struct {
	Key   string `json:"_id" bson:"_id"`
	Value int    `json:"value" bson:"value"`
}

type sequenceGenerator struct {
	hi     uint16
	lo     uint16
	lock   sync.Mutex
	update chan uint16
}

// SequenceGenerator - generates an unique ordered sequence of a TaskOrderID
type SequenceGenerator interface {
	Next() unique.TaskOrderID
}

// Next - returns next TaskOrderID value
func (seq *sequenceGenerator) Next() unique.TaskOrderID {
	defer seq.lock.Unlock()
	seq.lock.Lock()

	odate := date.CurrentOdate()
	week := odate.Woyear()

	result := make([]byte, 4)
	seq.lo++
	if seq.lo > 3843 {
		seq.lo = 0
		seq.hi++
		if seq.hi > 3843 {
			seq.hi = 0
		}
		seq.update <- seq.hi
	}

	hi := seq.hi
	lo := seq.lo

	result[0] = byte((lo % 62))
	result[1] = byte((lo / 62) % 62)
	result[2] = byte((hi % 62))
	result[3] = byte((hi / 62) % 62)

	return unique.TaskOrderID(string([]byte{base62Str[week], base62Str[result[3]], base62Str[result[2]], base62Str[result[1]], base62Str[result[0]]}))

}

func (seq *sequenceGenerator) watch(path string) {
	go func() {
		var data []byte
		var err error

		for {
			sval := <-seq.update

			model := sequenceModel{Key: sequenceName, Value: int(sval)}
			if data, err = json.Marshal(&model); err != nil {
				log.Error().Err(err).Msg("marshal sequence failed")
			}

			if err = os.WriteFile(path, data, 0644); err != nil {
				log.Error().Err(err).Msg("marshal sequence failed")
			}
		}
	}()
}

// NewSequenceGenerator - creates a new sequence generator
func NewSequenceGenerator(provider *datastore.Provider) (SequenceGenerator, error) {

	var value uint16
	var model = sequenceModel{}
	var data []byte
	var err error

	path := filepath.Join(provider.Directory(), sequenFileName)

	if data, err = os.ReadFile(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err = json.Unmarshal(data, &model); err != nil {
		model.Key = sequenceName
		model.Value = 3843
	}

	if model.Value < 0 || model.Value > 3843 {
		return nil, errors.New("sequence hi value is not in range")
	}

	model.Value++
	if model.Value > 3843 {
		model.Value = 0
	}
	value = uint16(model.Value)

	if data, err = json.Marshal(&model); err != nil {
		return nil, err
	}

	if err = os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	seq := &sequenceGenerator{lock: sync.Mutex{}, hi: value, update: make(chan uint16)}
	seq.watch(path)

	return seq, nil
}
