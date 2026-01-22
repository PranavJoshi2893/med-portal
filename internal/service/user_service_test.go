package service

import (
	"testing"
)

func TestRegister(t *testing.T) {

	got := ""
	want := "test"

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

}
