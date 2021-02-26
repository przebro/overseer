package pool

import (
	"context"
	"errors"
	"overseer/common/types/date"
	"overseer/datastore"
	"overseer/overseer/internal/unique"
	"sync"

	"github.com/przebro/databazaar/collection"
)

const base62Str = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const sequenceName = "orderid_sequence"
const seqMax uint16 = 3843 //base62 62^2 = 3844

type sequenceModel struct {
	Key   string `json:"_id" bson:"_id"`
	Value int    `json:"value" bson:"value"`
}

type sequenceGenerator struct {
	hi     uint16
	lo     uint16
	lock   sync.Mutex
	update chan uint16
	col    collection.DataCollection
}

//SequenceGenerator - generates an unique ordered sequence of a TaskOrderID
type SequenceGenerator interface {
	Next() unique.TaskOrderID
}

//Next - returns next TaskOrderID value
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

func (seq *sequenceGenerator) watch() {
	go func() {
		for {
			sval := <-seq.update
			model := sequenceModel{Key: sequenceName, Value: int(sval)}
			seq.col.Update(context.Background(), &model)
		}
	}()
}

//NewSequenceGenerator - creates a new sequence generator
func NewSequenceGenerator(colname string, provider *datastore.Provider) (SequenceGenerator, error) {

	var value uint16
	var model = sequenceModel{}
	var col collection.DataCollection
	var err error
	var create bool

	if col, err = provider.GetCollection(colname); err != nil {
		return nil, err
	}

	if err = col.Get(context.Background(), sequenceName, &model); err != nil {
		model.Key = sequenceName
		model.Value = 3843
		create = true
	}

	if model.Value < 0 || model.Value > 3843 {
		return nil, errors.New("sequence hi value is not in range")
	}

	model.Value++
	if model.Value > 3843 {
		model.Value = 0
	}
	value = uint16(model.Value)

	if create == true {
		col.Create(context.Background(), &model)
	} else {
		col.Update(context.Background(), &model)
	}

	seq := &sequenceGenerator{lock: sync.Mutex{}, hi: value, update: make(chan uint16), col: col}
	seq.watch()

	return seq, nil
}
