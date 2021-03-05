package api

import "fmt"

type Error struct {
	Message      string
	Status, Code int
}

func (e *Error) Error() string {
	var i int

	if e.Status != 0 && e.Code == 0 {
		i = e.Status
	} else {
		i = e.Code
	}

	return fmt.Sprintf("%d: %s", i, e.Message)
}
