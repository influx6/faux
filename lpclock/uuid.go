package lpclock

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// Errors ...
var (
	ErrNoHashPrefix          = errors.New("must have '#' prefix")
	ErrInvalidTickType       = errors.New("invalid tick type received")
	ErrInvalidUUIDFormat     = errors.New("data has invalid UUID format")
	ErrInvalidIdentityFormat = errors.New("identity data has invalid format")
	ErrInvalidIDTimeFormat   = errors.New("ID_Time data has invalid format")
)

var (
	dot        = []byte(".")
	hash       = []byte("#")
	underscore = []byte("_")
)

// TickType used to indicate type of tick value.
type TickType int64

// consts of tick types.
const (
	LAMPORTTICK TickType = 1
	UNIXTICK    TickType = 2
)

// UUID generates a uuid which runs with a giving length of encoded values.
type UUID struct {
	ID     string
	Origin string
	Tick   int64
	Type   TickType
}

// Validate returns an error if uuid is in a invalid state.
func (u UUID) Validate() error {
	if u.ID == "" {
		return errors.New("UUID.ID is required")
	}

	if u.Origin == "" {
		return errors.New("UUID.Origin is required")
	}

	if u.Tick == 0 {
		return errors.New("UUID.Tick can not be zero")
	}

	switch u.Type {
	case LAMPORTTICK, UNIXTICK:
		return nil
	}

	return errors.New("UUID.Type can not be 0 or invalid state")
}

// GreaterThan validates that the uuid is less than value of
// provided uuid.
func (u UUID) GreaterThan(n UUID) bool {
	if u.Equal(n) {
		return u.Tick > n.Tick
	}

	return false
}

// LessThan validates that the uuid is less than value of
// provided uuid.
func (u UUID) LessThan(n UUID) bool {
	if u.Equal(n) {
		return u.Tick < n.Tick
	}

	return false
}

// ExactEqual returns true/false if giving UUIDs a exact match.
func (u UUID) ExactEqual(n UUID) bool {
	if u.Equal(n) {
		return u.Tick == n.Tick
	}
	return false
}

// Equal returns true/false if giving UUIDs are a match in tick type, origin and id.
// It does not compare tick value.
func (u UUID) Equal(n UUID) bool {
	if u.Origin != n.Origin {
		return false
	}

	if u.ID != n.ID {
		return false
	}

	if u.Type != n.Type {
		return false
	}

	return true
}

// MarshalJSON marshals provided UUID into text version
// returning byte slice.
func (u UUID) MarshalJSON() ([]byte, error) {
	return u.MarshalText()
}

// UnmarshalJSON relies on UnmarshalText and expects
// to receive json string of uuid.
// It returns an error if the provided byte slice is not
// a quote string.
func (u *UUID) UnmarshalJSON(d []byte) error {
	content, err := strconv.Unquote(string(d))
	if err != nil {
		return err
	}

	return u.UnmarshalJSON([]byte(content))
}

// MarshalText returns byte slice of giving uuid.
func (u UUID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText unmarshals giving uuid into appropriate UUID struct.
func (u *UUID) UnmarshalText(dest []byte) error {
	if !bytes.HasPrefix(dest, hash) {
		return ErrNoHashPrefix
	}

	dest = bytes.TrimPrefix(dest, hash)
	areas := bytes.Split(dest, hash)
	if len(areas) != 2 {
		return ErrInvalidUUIDFormat
	}

	tickValue, err := strconv.ParseInt(string(areas[0]), 10, 64)
	if err != nil {
		return err
	}

	switch TickType(tickValue) {
	case LAMPORTTICK, UNIXTICK:
		u.Type = TickType(tickValue)
	default:
		return ErrInvalidTickType
	}

	identity := bytes.Split(areas[1], dot)
	if len(identity) != 2 {
		return ErrInvalidIdentityFormat
	}

	u.Origin = string(identity[0])

	idTime := bytes.Split(identity[1], underscore)
	if len(idTime) != 2 {
		return ErrInvalidIDTimeFormat
	}

	u.ID = string(idTime[0])

	timeTick, err := strconv.ParseInt(string(idTime[1]), 10, 64)
	if err != nil {
		return err
	}

	u.Tick = timeTick
	return nil
}

// String returns string version of uuid.
// Format: #TICK_TYPE#ID_LENGTH#OriginID_TIMETICK
func (u UUID) String() string {
	return fmt.Sprintf("#%d#%s.%s_%d", u.Type, u.Origin, u.ID, u.Tick)
}
