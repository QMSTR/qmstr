package cli

import (
	"log"
	"testing"

	"github.com/spf13/pflag"
)

func TestStringFlag(t *testing.T) {
	flags := &pflag.FlagSet{}
	err := generateFlags(&Foo{}, flags)
	if err != nil {
		t.FailNow()
	}
	err = flags.Parse([]string{"--FooBar", "dafuq"})
	if err != nil {
		log.Printf("Error %v", err)
		t.FailNow()
	}
}

func TestInt64Flag(t *testing.T) {
	flags := &pflag.FlagSet{}
	err := generateFlags(&Foo{}, flags)
	if err != nil {
		t.FailNow()
	}
	err = flags.Parse([]string{"--FooBarInt", "1337"})
	if err != nil {
		log.Printf("Error %v", err)
		t.FailNow()
	}
}

func TestBoolFlag(t *testing.T) {
	flags := &pflag.FlagSet{}
	err := generateFlags(&Foo{}, flags)
	if err != nil {
		t.FailNow()
	}
	err = flags.Parse([]string{"--FooBarBool"})
	if err != nil {
		log.Printf("Error %v", err)
		t.FailNow()
	}
}
