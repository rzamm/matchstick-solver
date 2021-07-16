// Package print extends the functionality of message printer
package io

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Printer struct {
	*message.Printer
}

func NewPrinter(t language.Tag, opts ...message.Option) *Printer {
	return &Printer{
		message.NewPrinter(t, opts...),
	}
}

// MustPrintf is like message.Printf, but does not return an error.
func (p *Printer) MustPrintf(key message.Reference, a ...interface{}) int {
	n, err := p.Printf(key, a...)
	if err != nil {
		panic(err)
	}
	return n
}
