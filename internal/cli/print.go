package cli

import (
	"fmt"
	"io"
)

// printer writes to an io.Writer and retains the first write error.
type printer struct {
	w   io.Writer
	err error
}

func newPrinter(w io.Writer) *printer {
	return &printer{w: w}
}

func (p *printer) println(a ...any) {
	if p.err != nil {
		return
	}
	_, p.err = fmt.Fprintln(p.w, a...)
}

func (p *printer) printf(format string, a ...any) {
	if p.err != nil {
		return
	}
	_, p.err = fmt.Fprintf(p.w, format, a...)
}
