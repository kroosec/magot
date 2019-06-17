package object

import "testing"

func TestStringHasKey(t *testing.T) {
	hello1 := &String{Value: "Hello world"}
	hello2 := &String{Value: "Hello world"}
	diff1 := &String{Value: "Iz diff"}
	diff2 := &String{Value: "Iz diff"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() == hello1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
