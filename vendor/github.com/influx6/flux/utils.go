package flux

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var elapso = regexp.MustCompile(`(\d+)(\w+)`)

// Capitalize capitalizes the first character in a string
func Capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToCamelCase takes a string set and returns the camelcase version
func ToCamelCase(set []string) string {
	var camel []string
	for ind, word := range set {
		cameld := word
		if ind > 0 {
			cameld = (strings.ToUpper(word[:1]) + word[1:])
		}

		camel = append(camel, cameld)
	}
	return strings.Join(camel, "")
}

//MakeDuration allows you to make create a duration from a string
func MakeDuration(target string, def int) time.Duration {
	if !elapso.MatchString(target) {
		return time.Duration(def)
	}

	matchs := elapso.FindAllStringSubmatch(target, -1)

	if len(matchs) <= 0 {
		return time.Duration(def)
	}

	match := matchs[0]

	if len(match) < 3 {
		return time.Duration(def)
	}

	dur := time.Duration(ConvertToInt(match[1], def))

	mtype := match[2]

	switch mtype {
	case "s":
		return dur * time.Second
	case "mcs":
		return dur * time.Microsecond
	case "ns":
		return dur * time.Nanosecond
	case "ms":
		return dur * time.Millisecond
	case "m":
		return dur * time.Minute
	case "h":
		return dur * time.Hour
	default:
		return time.Duration(dur) * time.Second
	}
}

//ConvertToInt wraps the internal int coverter
func ConvertToInt(target string, def int) int {
	fo, err := strconv.Atoi(target)
	if err != nil {
		return def
	}
	return fo
}

func uInt16ToByteArray(value uint16, bufferSize int) []byte {
	toWriteLen := make([]byte, bufferSize)
	binary.LittleEndian.PutUint16(toWriteLen, value)
	return toWriteLen
}

// Formula for taking size in bytes and calculating # of bits to express that size
// http://www.exploringbinary.com/number-of-bits-in-a-decimal-integer/
func messageSizeToBitLength(messageSize int) int {
	bytes := float64(messageSize)
	header := math.Ceil(math.Floor(math.Log2(bytes)+1) / 8.0)
	return int(header)
}

// RandString generates a set of random numbers of a set length
func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

// RandAlpha generates a set of random numbers of a set length
func RandAlpha(n int) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
