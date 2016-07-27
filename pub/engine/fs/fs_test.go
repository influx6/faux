package fs_test

import (
	"fmt"
	"testing"

	"github.com/influx6/faux/pub/engine/fs"
)

func TestFileList(t *testing.T) {
	item := fs.New()

	listChan := item.ReadDir("../../..").
		SkipStat(fs.IsDir).
		UnwrapStats().
		DataPort(true)

	item.Read(true)

	list := <-listChan
	fmt.Printf("List: %s\n", list)
}

// func TestFileEnding(t *testing.T) {
// 	var b bytes.Buffer
//
// 	item := fs.New()
//
// 	item.Mkdir("fixtures", false).
// 		Mkdir("configs", true).
// 		CreateFile("boot.cfg", true).
// 		WriteBytes([]byte("Just got school,")).
// 		WriteBytes([]byte("from soccer practice.")).
// 		Close().
// 		OpenFile("fixtures/configs/boot.cfg", false).
// 		ReadReader().
// 		WriteWriter(&b).
// 		RemoveAll("fixtures").
// 		Signal(func(err error) {
// 			t.Errorf("Error occured: %s", err)
// 		}).
// 		Read("")
//
// 	expected := []byte("Just got school,from soccer practice.")
// 	if bytes.Compare(b.Bytes(), expected) == -1 {
// 		t.Logf("Expected: %s", expected)
// 		t.Logf("Recieved: %s", b.Bytes())
// 		t.Fatalf("Should have recieved matching values")
// 	}
// 	t.Logf("Should have recieved matching values")
// }
