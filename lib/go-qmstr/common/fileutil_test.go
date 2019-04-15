package common_test

import (
	"log"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/spf13/afero"
)

func TestHash(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("test")
	if err != nil || f == nil {
		t.Fail()
	}
	i, err := f.WriteString("test")
	if err != nil {
		t.Fail()
	}
	log.Printf("%d", i)
	f.Sync()
	err = f.Close()
	if err != nil {
		t.Fail()
	}

	f, err = fs.Open("test")
	if err != nil {
		t.Fail()
	}
	chksum, err := common.Hash(f)
	if err != nil {
		t.Fail()
	}
	if chksum != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		log.Printf("%s != a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", chksum)
		t.Fail()
	}

}

func TestPosixFullyPortableFilename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Whitespace", args{filename: "There are whitespaces\there"}, "There_are_whitespaces_here"},
		{"Newline", args{filename: "There is a newline\nhere"}, "There_is_a_newline_here"},
		{"non-ascii", args{filename: "There is a non-ascii char ä here"}, "There_is_a_non-ascii_char___here"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := common.GetPosixFullyPortableFilename(tt.args.filename); got != tt.want {
				t.Errorf("posixFullyPortableFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
