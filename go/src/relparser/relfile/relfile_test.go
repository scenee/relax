package relfile

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenSource(t *testing.T) {
	in := "../../../../sample/Relfile"
	r, err := LoadRelfile(in)
	if err != nil {
		t.Fatal(err)
	}

	var (
		out  *os.File
		data []byte
	)

	out = r.GenSource("development")
	t.Logf("%s", out.Name())
	if data, err = ioutil.ReadFile(out.Name()); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", data)

	out = r.GenSource("adhoc")
	t.Logf("%s", out.Name())

	if data, err = ioutil.ReadFile(out.Name()); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", data)
}
