package calcserver

import (
	"reflect"
	"testing"
)

func TestNewCache(t *testing.T) {
	got := reflect.TypeOf(NewCache()).String()
	want := "*calcserver.CalcCache"
	if got != want {
		t.Errorf("NewCache() type got: %v, type want: %v", got, want)
	}
}
