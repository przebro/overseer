package unique

import (
	"encoding/binary"
	"encoding/hex"
	"goscheduler/overseer/internal/date"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type unique struct {
	m sync.Mutex
}

//MsgID - Unique message ID.
type MsgID [12]byte

//TaskOrderID - unique order id of a task
type TaskOrderID string

var oidSeq int32 = 0

const maxOidSeq = 238328
const base62Str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//Hex - Returns hex representation of a MsgID
func (msgid MsgID) Hex() string {
	bytes := make([]byte, 12)
	copy(bytes, msgid[:])
	return hex.EncodeToString(bytes)
}

var unq *unique = nil
var once sync.Once

func initialize() {
	unq = &unique{m: sync.Mutex{}}
	rand.Seed(time.Now().UnixNano())
}

//None - Creates an empty MsgID.
func None() MsgID {
	msgid := MsgID{}
	return msgid
}

//NewID - Creates a new MsgID.
func NewID() MsgID {
	once.Do(initialize)
	defer unq.m.Unlock()
	unq.m.Lock()

	bytes := getUniqueBytes()
	return bytes
}

func getUniqueBytes() MsgID {

	tpart := uint64(time.Now().UnixNano())
	rpart := rand.Uint32()

	lbytes := make([]byte, 6)
	nbytes := make([]byte, 10)

	lbytes[0] = byte(tpart >> 56)
	lbytes[1] = byte(tpart >> 48)
	lbytes[2] = byte(tpart >> 40)
	lbytes[3] = byte(tpart >> 32)
	lbytes[4] = byte(tpart >> 24)
	lbytes[5] = byte(tpart >> 16)

	rbytes := make([]byte, 4)
	binary.BigEndian.PutUint32(rbytes, rpart)

	rand.Shuffle(8, func(i, j int) {})

	bytes := append(lbytes, rbytes...)
	copy(nbytes, bytes)

	rand.Shuffle(8, func(i, j int) { nbytes[i], nbytes[j] = nbytes[j], nbytes[i] })
	bytes = append(bytes, []byte{nbytes[1], nbytes[3]}...)

	var mid MsgID = MsgID{}
	copy(mid[:], bytes)

	return mid

}

//NewOrderID - Generates a new unique OrderID
func NewOrderID() TaskOrderID {

	seq := atomic.AddInt32(&oidSeq, 1)

	pos := 0
	result := make([]byte, 3)

	odate := date.CurrentOdate()
	mth := odate.Omonth()

	for seq != 0 {
		result[pos] = byte((seq % 62))
		seq = seq / 62
		pos++

	}

	atomic.CompareAndSwapInt32(&oidSeq, maxOidSeq, 0)

	return TaskOrderID(string([]byte{mth[0], mth[1], base62Str[result[2]], base62Str[result[1]], base62Str[result[0]]}))
}
