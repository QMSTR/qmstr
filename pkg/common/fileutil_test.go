package common_test

import (
	"log"
	"testing"

	"github.com/QMSTR/qmstr/pkg/common"
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
