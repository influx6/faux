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
		OpenFile("fixtures/configs/boot.cfg", false).
		WriteBytes([]byte("Just got school,")).
		WriteBytes([]byte("from soccer practice.")).
		ReadIncomingReader().
		WriteWriter(&b).
		Signal(func(err error) {
			fmt.Println("Errors: %s\n", err)
		}).
		Read("")

	// RemoveAll("fixtures").Read("")

	fmt.Printf("Recieved: %+q\n", b.Bytes())
}
