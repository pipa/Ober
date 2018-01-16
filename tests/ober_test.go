package tests

import (
	"testing"
	"time"

	ober "github.com/pipa/Ober"
)

func TestOber(t *testing.T) {
	o := ober.New()

	if o.Router() == nil {
		t.Fatal("router not found")
	}
}

func TestStart(t *testing.T) {
	o := ober.New()
	go func() {
		t.Fatal(o.Start())
	}()
	time.Sleep(200 * time.Millisecond)
}
