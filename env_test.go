package envx

import (
	"reflect"
	"testing"
	"time"
)

func TestEnvSet(t *testing.T) {
	var d time.Duration

	fset := NewEnvSet("ENVX")
	fset.Duration(&d, "TIMEOUT", 5*time.Second, "just a timeout")
	// fset.IntSlice(&ids, "ids", "", wantIDs, "just a timeout")
	// fset.Float64Set(&offsets, "offsets", "", wantOffsets, "just a timeout")

	err := fset.Parse([]string{"ENVX_TIMEOUT=10m", "TIMEOUT=12345"})
	failIfErr(t, err)
	mustEqual(t, fset.IsParsed(), true)
	mustEqual(t, int64(d), int64(10*time.Minute))
}

func failIfErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustEqual(t testing.TB, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
