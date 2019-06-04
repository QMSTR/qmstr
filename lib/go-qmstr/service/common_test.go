package service 

import "testing"
import "log"

type Dummy struct {
	a string
	b int64
	c []Dummy
}

func TestCheckEmpty(t *testing.T) {
	d := Dummy{}
	err := checkEmpty(&d)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	d.a = "Test"
	err = checkEmpty(&d)
	if err == nil {
		t.Fail()
	}
}
