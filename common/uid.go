package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58" // Importing base58 encoding from btcutil
	"strconv"
	"strings"
)


// UID represents a unique identifier composed of three parts:
// - localID: 32-bit 
// - 10-bit objectType
// - shardID: 18-bit 
type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(localID uint32, objType int, shardID uint32) UID {
	return UID{
		localID:    localID,
		objectType: objType,
		shardID:    shardID,
	}
}


// Shard: 1, Object: 1, ID: 1 => 0001 0001 0001
// 1 << 8  = 0001 0000 0000
// 1 << 4  =      1 0000
// 1 << 0  =          1
// => 0001 0001 0001

func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)<<0
	return base58.Encode([]byte(fmt.Sprintf("%v", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func DecomposeUID(s string) (UID, error) {
	uid, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
			return UID{}, err
	}

	if (1 << 18) > uid {
			return UID{}, errors.New("wrong uid")
	}

// x = 1110 1110 0101 => x >> 4 = 1110 1110 & 0000 1111 = 1110
/* 
* Shift 4 bits to get 1110 1110
* Whatever bit you want, AND with 1, the rest you don't need, AND with 0
* Example: 111, extract first bit => 111 & 001 = 001
*/


	u := UID{
			localID:     uint32(uid >> 28),
			//get the 10 bits between local and shard
			objectType:  int(uid >> 18 & 0x3FF),
			shardID:     uint32(uid >> 0 & 0x3FFFF),
	}

	return u, nil 
}


func FromBase58(s string) (UID, error) {
    return DecomposeUID(string(base58.Decode(s)))
}

func (uid UID) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmarshalJSON(data []byte) error {
    decodeUID, err := FromBase58(strings.Replace(string(data), "\"", "", -1))

    if err != nil {
				return err
		}

    uid.localID = decodeUID.localID
    uid.shardID = decodeUID.shardID
    uid.objectType = decodeUID.objectType

    return nil
}

func (uid *UID) Value() (driver.Value, error) {
    if uid == nil {
				return nil, nil
		}
    return int64(uid.localID), nil
}


func (uid *UID) Scan(value interface{}) error {
	if value == nil {
			return nil
	}

	var i uint32

	switch t := value.(type) {
	case int:
			i = uint32(t)
	case int8:
			i = uint32(t) // standardizes across systems
	case int16:
			i = uint32(t) // standardizes across systems
	case int32:
			i = uint32(t) // standardizes across systems
	case int64:
			i = uint32(t) // standardizes across systems
	case uint8:
			i = uint32(t) // standardizes across systems
	case uint16:
			i = uint32(t) // standardizes across systems
	case uint32:
			i = t
	case uint64:
			i = uint32(t)
	case []byte:
			a, err := strconv.Atoi(string(t))
			if err != nil {
					return err
			}
			i = uint32(a)
	default:
			return errors.New("invalid Scan Source")
	}

	*uid = NewUID(i, 0, 1)

	return nil
}
