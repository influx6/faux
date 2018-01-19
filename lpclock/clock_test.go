package lpclock_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/lpclock"
	"github.com/influx6/faux/tests"
)

func TestResetableClock(t *testing.T) {
	uuids := lpclock.Unix("localhost")
	lastTime := uuids.Now()

	uuid2 := lpclock.Unix("localhust")
	if err := uuid2.Reset(lastTime); err != nil {
		tests.FailedWithError(err, "Should have successfully reset lpclock")
	}
	tests.Passed("Should have successfully reset lpclock")

	nextTime := uuid2.Now()
	if nextTime.LessThan(lastTime) {
		tests.Failed("Should have next clock be upper version of last")
	}
	tests.Passed("Should have next clock be upper version of last")

	if nextTime.Type != lastTime.Type {
		tests.Failed("Should have Type of last and next clock time equal")
	}
	tests.Passed("Should have Type of last and next clock time equal")

	if nextTime.ID != lastTime.ID {
		tests.Failed("Should have ID of last and next clock time equal")
	}
	tests.Passed("Should have ID of last and next clock time equal")

	if nextTime.Origin != lastTime.Origin {
		tests.Failed("Should have origin of last and next clock time equal")
	}
	tests.Passed("Should have origin of last and next clock time equal")
}

func TestMonotonicClock(t *testing.T) {
	uuids := lpclock.Unix("localhost")

	last := uuids.Now()
	for i := 1000; i > 0; i-- {
		now := uuids.Now()
		if last.ExactEqual(now) {
			tests.Failed("Should never match any other other timestamp")
		}

		if now.LessThan(last) {
			tests.Failed("Should always have increasing uuid values")
		}

		if last.GreaterThan(now) {
			tests.Failed("Should always have increasing uuid values")
		}
		last = now
	}
	tests.Passed("Should never match any other other timestamp")
}

func TestUUIDMarshaling(t *testing.T) {
	uuid := lpclock.Unix("localhost").Now()

	uuidBytes, err := uuid.MarshalText()
	if err != nil {
		tests.FailedWithError(err, "Should have successfully marshalled UUID")
	}
	tests.Passed("Should have successfully marshalled UUID")

	tests.Info("UUID: %q", uuid.String())
	var uid lpclock.UUID
	if err := uid.UnmarshalText(uuidBytes); err != nil {
		tests.FailedWithError(err, "Should have successfully unmarshalled uuid")
	}
	tests.Passed("Should have successfully unmarshalled uuid")

	tests.Info("UUID: %q", uid.String())
	if uid.String() != uuid.String() {
		tests.Failed("Should have uuid match successfully")
	}
	tests.Passed("Should have uuid match successfully")

	uidBytes, _ := uid.MarshalText()
	if !bytes.Equal(uidBytes, uuidBytes) {
		tests.Failed("Should have uuid match successfully")
	}
	tests.Passed("Should have uuid match successfully")
}
