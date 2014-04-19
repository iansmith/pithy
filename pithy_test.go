package pithy_test

import (
	"github.com/iansmith/pithy"
	"testing"
)

func TestBasicNoPanic(t *testing.T) {

	//none of these should panice
	pithy.DIV("myid")
	pithy.DIV("myid", []string{"foo", "bar"})
	pithy.DIV([]string{"foo", "bar"})
	pithy.DIV("id", []string{"foo", "bar"}, map[string]string{
		"foo": "bar",
		"baz": "grik",
	})
	pithy.DIV([]string{"foo", "bar"}, map[string]string{
		"foo": "bar",
		"baz": "grik",
	})
	pithy.DIV(map[string]string{
		"foo": "bar",
		"baz": "grik",
	})
}

func expectPanic(fn func()) (result bool) {
	result = false
	defer func() {
		if r := recover(); r != nil {
			result = true
		}
	}()
	fn()
	return
}
func TestBasicPanic(t *testing.T) {

	if !expectPanic(func() {
		pithy.DIV([]string{"foo"}, "id")
	}) {
		t.Errorf("expected panic when classes are before id")
	}
	if !expectPanic(func() {
		pithy.DIV(map[string]string{"foo": "bar"}, "id")
	}) {
		t.Errorf("expected panic when attributes are before id")
	}
	if !expectPanic(func() {
		pithy.DIV(nil)
	}) {
		t.Errorf("expected panic with a nil value")
	}
}

func TestPithyStruct(t *testing.T) {
	div := pithy.DIV("myid")
	if div.Id != "myid" {
		t.Error("id not created properly: %s", div.Id)
	}

	div = pithy.DIV([]string{"foo", "bar"})
	if len(div.Classes) != 2 {
		t.Error("classes not created properly: %d", len(div.Classes))
	}
	if div.Classes[1] != "bar" {
		t.Error("classes not created properly: %v", div.Classes)
	}
}

func TestPithyNesetd(t *testing.T) {
	div := pithy.DIV("foo", []string{"bar", "baz"},
		pithy.DIV("frik", map[string]string{"grik": "grak"}),
		pithy.DIV("wumpus"),
		pithy.DIV([]string{"smelly"},
			pithy.DIV("itsdark"),
		),
	)

	if len(div.Children) != 3 {
		t.Errorf("children not built correctly: %d", len(div.Children))
	}

	if len(div.Children[2].Children) != 1 {
		t.Errorf("nested children not built correctly: %d", len(div.Children[2].Children))
	}
	if len(div.Children[1].Children) != 0 {
		t.Errorf("unexpected children not built correctly: %d", len(div.Children[1].Children))
	}
}
