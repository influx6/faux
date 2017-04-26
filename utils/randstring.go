package utils

import (
	"crypto/rand"
	"math"
	mrand "math/rand"
	"strings"
	"time"
)

// RandomInt returns a random number between the provided min-max range.
func RandomInt(min, max int64) int64 {
	return int64(math.Floor(float64(mrand.Int63()*((max-min)+1)))) + min
}

// RandomFloat returns a random number between the provided min-max range.
func RandomFloat(min, max float64) float64 {
	return math.Floor(mrand.Float64()*((max-min)+1)) + min
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

// RandText generates random string based on type (Alphanum, Alpha, Number).
func RandText(strSize int, randType string) string {
	var dictionary string

	randType = strings.ToLower(randType)

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}

// TickToHighResTimer provides a method to transform requestAnimationFrame
// clock elapsed time into a appropriate time.Duration
func TickToHighResTimer(ms float64) time.Duration {
	return time.Duration(ms * float64(time.Millisecond))
}

// Times using the provided count, runs the function (n-1) number of times, since
// it starts from zero.
func Times(n int, fn func(int)) {
	for i := 0; i < n; i++ {
		fn(i + 1)
	}
}

//==================================================================================

// contains different constants of different accepted sclaes.
const (
	AugmentedFourth = 1.414
	MinorSecond     = 1.067
	MajorSecond     = 1.125
	MinorThird      = 1.200
	MajorThird      = 1.250
	PerfectFourth   = 1.333
	PerfectFifth    = 1.500
	GoldenRation    = 1.618
)

// GenerateValueScale returns a value scale which is produced from generating
// a slice of n values representing the given scale value and are multipled
// by the provided base values.
func GenerateValueScale(scale float64, base float64) []float64 {

	// Generate scale based on 1.0 scale using the provided scale.
	scales := GenerateScale(scale, 5, 10)

	// Multiply all scale value by the provided base.
	Times(len(scales), func(index int) {
		scales[index-1] = scales[index-1] * base
	})

	return scales
}

// GenerateEMScale returns a slice of values which are the a combination of
// a reducing + increasing scaled values of the provided scale generated from
// using the base initial 1.0 value against an ever incremental 1.0*(scale * n)
// or 1.0 / (scale *n) value, where n is the ever increasing index.
func GenerateScale(scale float64, minorCount int, majorCount int) []float64 {
	var scales []float64

	minorScales := make([]float64, minorCount)

	Times(minorCount, func(index int) {
		scaled := 1.0

		Times(index, func(_ int) {
			scaled *= scale
		})

		minorLen := len(minorScales)
		minorScales[minorLen-index] = 1.0 / scaled
	})

	scales = append(scales, minorScales...)
	scales = append(scales, 1.0)

	Times(majorCount, func(index int) {
		scaled := 1.0

		Times(index, func(_ int) {
			scaled *= scale
		})

		scales = append(scales, scaled)
	})

	return scales
}
