package eventdispatcher

import (
	"fmt"
	"testing"
)

func func1(e *Event) {
	fmt.Println("---func1")
}
func func2(e *Event) {
	fmt.Println("---func2")
}

func Test1(t *testing.T) {
	dispatcher := NewEventDispatcher()

	var test EventHandler = func1
	var test2 EventHandler = func2

	dispatcher.AddListener("abc", &test)
	dispatcher.AddListener("abc", &test2)

	dispatcher.Dispatch("abc", nil)

	dispatcher.RemoveListener("abc", &test)

	dispatcher.Dispatch("abc", nil)
}
