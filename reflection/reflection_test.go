package reflection_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/influx6/faux/reflection"
	"github.com/influx6/faux/tests"
)

type bull string

type speaker interface {
	Speak() string
}

// mosnter provides a basic struct test case type.
type monster struct {
	Name  string
	Items []bull
}

// Speak returns the sound the monster makes.
func (m *monster) Speak() string {
	return "Raaaaaaarggg!"
}

func get(t *testing.T, sm speaker) {
	name, embedded, err := reflection.ExternalTypeNames(sm)
	if err != nil {
		tests.FailedWithError(err, "Should be able to retrieve field names arguments lists")
	}
	tests.Info("Name: %s", name)
	tests.Info("Fields: %+q", embedded)
	tests.Passed("Should be able to retrieve function arguments lists")
}

type Addrs struct {
	Addr string
}

type addrFunc func(Addrs) error

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func TestIsSettableType(t *testing.T) {
	m2 := errors.New("invalid error")
	if !reflection.IsSettableType(errorType, reflect.TypeOf(m2)) {
		tests.Failed("Should have error matching type")
	}
}

func TestIsSettable(t *testing.T) {
	m2 := errors.New("invalid error")
	if !reflection.IsSettable(errorType, reflect.ValueOf(m2)) {
		tests.Failed("Should have error matching type")
	}

	m1 := mo("invalid error")
	if !reflection.IsSettable(errorType, reflect.ValueOf(m1)) {
		tests.Failed("Should have error matching type")
	}
}

type mo string

func (m mo) Error() string {
	return string(m)
}

func TestValidateFunc_Bad(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	err := reflection.ValidateFunc(testFunc, []reflection.TypeValidation{
		func(types []reflect.Type) bool {
			return len(types) == 1
		},
	}, []reflection.TypeValidation{
		func(types []reflect.Type) bool {
			return len(types) == 0
		},
	})

	if err == nil {
		tests.Failed("Should have function invalid to conditions")
	}
}

func TestValidateFunc(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	err := reflection.ValidateFunc(testFunc, []reflection.TypeValidation{
		func(types []reflect.Type) bool {
			return len(types) == 1
		},
	}, []reflection.TypeValidation{
		func(types []reflect.Type) bool {
			return len(types) == 1
		},
	})

	if err != nil {
		tests.FailedWithError(err, "Should have function valid to conditions")
	}
}

func TestFunctionApply_OneArgument(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	res, err := reflection.CallFunc(testFunc, "Alex")
	if err != nil {
		tests.FailedWithError(err, "Should have executed function")
	}
	tests.Passed("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{"Hello Alex"}) {
		tests.Info("Received: %q", res)
		tests.Failed("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgumentWithError(t *testing.T) {
	bad := errors.New("bad")
	var testFunc = func(v string, i int, d bool) ([]interface{}, error) {
		return []interface{}{v, i, d}, bad
	}

	res, err := reflection.CallFunc(testFunc, "Alex", 1, false)
	if err != nil {
		tests.FailedWithError(err, "Should have executed function")
	}
	tests.Passed("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{[]interface{}{"Alex", 1, false}, bad}) {
		tests.Info("Expected: %q", []interface{}{[]interface{}{"Alex", 1, false}, bad})
		tests.Info("Received: %q", res)
		tests.Failed("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgumentWithVariadic(t *testing.T) {
	var testFunc = func(v string, i int, d ...bool) []interface{} {
		return []interface{}{v, i, d}
	}

	res, err := reflection.CallFunc(testFunc, "Alex", 1, []bool{false})
	if err != nil {
		tests.FailedWithError(err, "Should have executed function")
	}
	tests.Passed("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{[]interface{}{"Alex", 1, []bool{false}}}) {
		tests.Info("Expected: %q", []interface{}{[]interface{}{"Alex", 1, []bool{false}}})
		tests.Info("Received: %q", res)
		tests.Failed("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgument(t *testing.T) {
	var testFunc = func(v string, i int, d bool) string {
		return "Hello " + v
	}

	res, err := reflection.CallFunc(testFunc, "Alex", 1, false)
	if err != nil {
		tests.FailedWithError(err, "Should have executed function")
	}
	tests.Passed("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{"Hello Alex"}) {
		tests.Info("Received: %q", res)
		tests.Failed("Expected value unmatched")
	}
}

func TestMatchFunction(t *testing.T) {
	var addr1 = func(_ Addrs) error { return nil }
	var addr2 = func(_ Addrs) error { return nil }

	if !reflection.MatchFunction(addr1, addr2) {
		tests.Failed("Should have matched argument types successfully")
	}
	tests.Passed("Should have matched argument types successfully")

	if !reflection.MatchFunction(&addr1, &addr2) {
		tests.Failed("Should have matched argument types successfully")
	}
	tests.Passed("Should have matched argument types successfully")

	if reflection.MatchFunction(&addr1, addr2) {
		tests.Failed("Should have failed matched argument types successfully")
	}
	tests.Passed("Should have failed matched argument types successfully")
}

func TestMatchElement(t *testing.T) {
	if !reflection.MatchElement(Addrs{}, Addrs{}, false) {
		tests.Failed("Should have matched argument types successfully")
	}
	tests.Passed("Should have matched argument types successfully")

	if !reflection.MatchElement(new(Addrs), new(Addrs), false) {
		tests.Failed("Should have matched argument types successfully")
	}
	tests.Passed("Should have matched argument types successfully")

	if reflection.MatchElement(new(Addrs), Addrs{}, false) {
		tests.Failed("Should have failed matched argument types successfully")
	}
	tests.Passed("Should have failed matched argument types successfully")
}

func TestStructMapperWithSlice(t *testing.T) {
	mapper := reflection.NewStructMapper()

	profile := struct {
		List []Addrs
	}{
		List: []Addrs{{Addr: "Tokura 20"}},
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully converted struct")
	}
	tests.Passed("Should have successfully converted struct")

	tests.Info("Map of Struct: %+q", mapped)

	profile2 := struct {
		List []Addrs
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		tests.FailedWithError(err, "Should have successfully mapped data back to struct")
	}
	tests.Passed("Should have successfully mapped data back to struct")

	if len(profile.List) != len(profile2.List) {
		tests.Failed("Mapped struct should have same length: %d - %d ", len(profile.List), len(profile2.List))
	}
	tests.Passed("Mapped struct should have same length: %d - %d ", len(profile.List), len(profile2.List))

	for ind, item := range profile.List {
		nxItem := profile2.List[ind]
		if item.Addr != nxItem.Addr {
			tests.Failed("Item at %d should have equal value %+q -> %+q", ind, item.Addr, nxItem.Addr)
		}
	}

	tests.Passed("All items should be exactly the same")
}

func TestStructMapperWthFieldStruct(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := reflection.NewStructMapper()
	mapper.AddAdapter(timeType, reflection.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, reflection.TimeInverseMapper(layout))

	profile := struct {
		Addr Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addr: Addrs{Addr: "Tokura 20"},
		Name: "Johnson",
		Date: time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully converted struct")
	}
	tests.Passed("Should have successfully converted struct")

	tests.Info("Map of Struct: %+q", mapped)

	profile2 := struct {
		Addr Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		tests.FailedWithError(err, "Should have successfully mapped data back to struct")
	}
	tests.Passed("Should have successfully mapped data back to struct")

	if profile2.Addr.Addr != profile.Addr.Addr {
		tests.Failed("Mapped struct should have same %q value", "Addr.Addr")
	}
	tests.Passed("Mapped struct should have same %q value", "Addr.Addr")
}

func TestGetFieldByTagAndValue(t *testing.T) {
	profile := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addrs: Addrs{Addr: "Tokura 20"},
		Name:  "Johnson",
		Date:  time.Now(),
	}

	_, err := reflection.GetFieldByTagAndValue(profile, "json", "name")
	if err != nil {
		tests.FailedWithError(err, "Should have successfully converted struct")
	}
}

func TestStructMapperWthEmbeddedStruct(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := reflection.NewStructMapper()
	mapper.AddAdapter(timeType, reflection.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, reflection.TimeInverseMapper(layout))

	profile := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addrs: Addrs{Addr: "Tokura 20"},
		Name:  "Johnson",
		Date:  time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully converted struct")
	}
	tests.Passed("Should have successfully converted struct")

	tests.Info("Map of Struct: %+q", mapped)

	profile2 := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		tests.FailedWithError(err, "Should have successfully mapped data back to struct")
	}
	tests.Passed("Should have successfully mapped data back to struct")

	if profile2.Addr != profile.Addr {
		tests.Failed("Mapped struct should have same %q value", "Addr.Addr")
	}
	tests.Passed("Mapped struct should have same %q value", "Addr.Addr")
}

func TestStructMapper(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := reflection.NewStructMapper()
	mapper.AddAdapter(timeType, reflection.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, reflection.TimeInverseMapper(layout))

	profile := struct {
		Addr        string
		CountryName string
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
	}{
		Addr:        "Tokura 20",
		Name:        "Johnson",
		CountryName: "Nigeria",
		Date:        time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully converted struct")
	}
	tests.Passed("Should have successfully converted struct")

	tests.Info("Map of Struct: %+q", mapped)

	if _, ok := mapped["name"]; !ok {
		tests.Failed("Map should have %q field", "name")
	}
	tests.Passed("Map should have %q field", "name")

	if _, ok := mapped["date"]; !ok {
		tests.Failed("Map should have %q field", "date")
	}
	tests.Passed("Map should have %q field", "date")

	if _, ok := mapped["addr"]; !ok {
		tests.Failed("Map should have %q field", "addr")
	}
	tests.Passed("Map should have %q field", "addr")

	if _, ok := mapped["data"].(string); ok {
		tests.Failed("Map should have field %q be a string", "date")
	}
	tests.Passed("Map should have field %q be a string", "date")

	profile2 := struct {
		Addr        string
		CountryName string
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		tests.FailedWithError(err, "Should have successfully mapped data back to struct")
	}
	tests.Passed("Should have successfully mapped data back to struct")

	tests.Info("Mapped Struct: %+q", profile2)

	if profile2.Name != profile.Name {
		tests.Failed("Mapped struct should have same %q value", "Name")
	}
	tests.Passed("Mapped struct should have same %q value", "Name")

	if profile2.Date.Format(layout) != profile.Date.Format(layout) {
		tests.Failed("Mapped struct should have same %q value", "Date")
	}
	tests.Passed("Mapped struct should have same %q value", "Date")

	if profile2.CountryName != profile.CountryName {
		tests.Failed("Mapped struct should have same %q value", "CountryName")
	}
	tests.Passed("Mapped struct should have same %q value", "CountryName")

	if profile2.Addr != profile.Addr {
		tests.Failed("Mapped struct should have same %q value", "Addr")
	}
	tests.Passed("Mapped struct should have same %q value", "Addr")
}

// TestGetArgumentsType validates reflection API GetArgumentsType functions
// results.
func TestGetArgumentsType(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	args, err := reflection.GetFuncArgumentsType(f)
	if err != nil {
		tests.FailedWithError(err, "Should be able to retrieve function arguments lists")
	}
	tests.Passed("Should be able to retrieve function arguments lists")

	name, embedded, err := reflection.ExternalTypeNames(monster{Name: "Bob"})
	if err != nil {
		tests.FailedWithError(err, "Should be able to retrieve field names arguments lists")
	}
	tests.Info("Name: %s", name)
	tests.Info("Fields: %+q", embedded)
	tests.Passed("Should be able to retrieve function arguments lists")

	get(t, &monster{Name: "Bob"})

	newVals := reflection.MakeArgumentsValues(args)
	if nlen, alen := len(newVals), len(args); nlen != alen {
		tests.Failed("Should have matching new values lists for arguments")
	}
	tests.Passed("Should have matching new values lists for arguments")

	mstring := reflect.TypeOf((*monster)(nil)).Elem()

	if mstring.Kind() != newVals[0].Kind() {
		tests.Failed("Should be able to match argument kind")
	}
	tests.Passed("Should be able to match argument kind")

}

func TestMatchFUncArgumentTypeWithValues(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf(monster{Name: "FireHouse"}))

	if index := reflection.MatchFuncArgumentTypeWithValues(f, vals); index != -1 {
		tests.Failed("Should have matching new values lists for arguments: %d", index)
	}
	tests.Passed("Should have matching new values lists for arguments")
}
