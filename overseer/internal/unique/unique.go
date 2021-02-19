package unique

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type unique struct {
	m sync.Mutex
}

//MsgID - Unique message ID.
type MsgID [12]byte

//TaskOrderID - unique order id of a task
type TaskOrderID string

var (
	unq            *unique = nil
	once           sync.Once
	errInvalidLen  = errors.New("TaskOrderID invalid length")
	errInvalidChar = errors.New("TaskOrderID contains invalid characters")
)

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

//Hex - Returns hex representation of a MsgID
func (msgid MsgID) Hex() string {
	bytes := make([]byte, 12)
	copy(bytes, msgid[:])
	return hex.EncodeToString(bytes)
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
