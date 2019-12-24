package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/arl/gitstatus"
)

// A Formater formats git status to JSON.
type Formater struct{}

// Format writes st as json into w.
func (Formater) Format(w io.Writer, st *gitstatus.Status) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	if err := enc.Encode(st); err != nil {
		return fmt.Errorf("can't format status to json: %v", err)
	}

	return nil
}
