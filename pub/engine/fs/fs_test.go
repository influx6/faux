package fs_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/influx6/faux/pub"
	"github.com/influx6/faux/pub/engine/fs"
)

func TestFS(t *testing.T) {
	var b bytes.Buffer

	item := fs.New()

	item.Signal(func(ctx pub.Ctx, d interface{}) {
		fmt.Println("Got: %+v", d)
		ctx.RW().Write(ctx, d)
	})

	item.Mkdir("fixtures", false).
		Mkdir("configs", true).
		WriteFile("boot.cfg").
		WriteWriterBytes([]byte("Just got school")).
		ReadFile("fixtures/configs/boot.cfg").
		WriteWriter(&b).Read("")
	// RemoveAll("fixtures").Read("")

	fmt.Printf("Recieved: %+q\n", b.Bytes())
}
