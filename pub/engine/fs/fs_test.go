package fs_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/influx6/faux/pub/engine/fs"
)

func TestFileEnding(t *testing.T) {
	var b bytes.Buffer

	item := fs.New()

	item.Mkdir("fixtures", false).
		Mkdir("configs", true).
		CreateFile("boot.cfg", true).
		WriteBytes([]byte("Just got school,")).
		WriteBytes([]byte("from soccer practice.")).
		Close().
		OpenFile("fixtures/configs/boot.cfg", false).
		ReadReader().
		WriteWriter(&b).
		Signal(func(err error) {
			fmt.Printf("Errors: %s\n", err)
		}).
		Read("")

	// RemoveAll("fixtures").Read("")

	fmt.Printf("Recieved: %+q\n", b.Bytes())
}
