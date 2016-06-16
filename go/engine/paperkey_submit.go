// Copyright 2016 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"github.com/keybase/client/go/libkb"
)

// PaperKeySubmit is an engine.
type PaperKeySubmit struct {
	libkb.Contextified
	paperPhrase string
}

// NewPaperKeySubmit creates a PaperKeySubmit engine.
func NewPaperKeySubmit(g *libkb.GlobalContext, paperPhrase string) *PaperKeySubmit {
	return &PaperKeySubmit{
		Contextified: libkb.NewContextified(g),
		paperPhrase:  paperPhrase,
	}
}

// Name is the unique engine name.
func (e *PaperKeySubmit) Name() string {
	return "PaperKeySubmit"
}

// Prereqs returns the engine prereqs.
func (e *PaperKeySubmit) Prereqs() Prereqs {
	return Prereqs{Device: true}
}

// RequiredUIs returns the required UIs.
func (e *PaperKeySubmit) RequiredUIs() []libkb.UIKind {
	return []libkb.UIKind{}
}

// SubConsumers returns the other UI consumers for this engine.
func (e *PaperKeySubmit) SubConsumers() []libkb.UIConsumer {
	return []libkb.UIConsumer{
		&PaperKeyGen{},
	}
}

// Run starts the engine.
func (e *PaperKeySubmit) Run(ctx *Context) error {
	me, err := libkb.LoadMe(libkb.NewLoadUserArg(e.G()))
	if err != nil {
		return err
	}

	kp, err := matchPaperKey(ctx, e.G(), me, e.paperPhrase)
	if err != nil {
		return err
	}

	aerr := e.G().LoginState().Account(func(a *libkb.Account) {
		err = a.SetUnlockedPaperKey(kp.sigKey, kp.encKey)
	}, "PaperKeySubmit - cache paper key")
	if aerr != nil {
		return aerr
	}
	if err != nil {
		return err
	}

	return nil
}
