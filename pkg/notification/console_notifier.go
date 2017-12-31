package notification

import (
	"bytes"
	"fmt"
	"io"
)

type console struct {
	w io.Writer
}

func NewConsole(w io.Writer) Notifier {
	return &console{w}
}

func (c *console) Notify(content []byte, to string) error {
	buf := bytes.Buffer{}
	buf.WriteString("To: " + to)
	buf.WriteString("Content: " + string(content))
	fmt.Fprintf(c.w, buf.String())
	return nil
}
