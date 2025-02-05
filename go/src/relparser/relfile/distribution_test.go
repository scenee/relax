package relfile

import (
	"testing"
)

func TestCheck(t *testing.T) {
	path := "../../../../sample/Relfile"
	r, err := LoadRelfile(path)
	if err != nil {
		t.Fatal(err)
	}
	d := r.Distributions["adhoc"]
	//infopath := "../../../../sample/SampleApp/Info.plist"
	//d.Check(infopath)
	d.Check()

	if !d.ExportOptions.CompileBitcode && !d.ExportOptions.UploadBitcode {
	} else {
		t.Fatal("Invalid ExportOptions.UploadBitcode default")
	}
}
