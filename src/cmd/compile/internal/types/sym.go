// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"cmd/internal/obj"
	"cmd/internal/src"
)

// Sym represents an object name. Most commonly, this is a Go identifier naming
// an object declared within a package, but Syms are also used to name internal
// synthesized objects.
//
// As an exception, field and method names that are exported use the Sym
// associated with localpkg instead of the package that declared them. This
// allows using Sym pointer equality to test for Go identifier uniqueness when
// handling selector expressions.
type Sym struct {
	Importdef *Pkg   // where imported definition was found
	Linkname  string // link name

	Pkg  *Pkg
	Name string // object name

	// saved and restored by dcopy
	Def        *Node    // definition: ONAME OTYPE OPACK or OLITERAL
	Block      int32    // blocknumber to catch redeclaration
	Lastlineno src.XPos // last declaration for diagnostic

	flags   bitset8
	Label   *Node // corresponding label (ephemeral)
	Origpkg *Pkg  // original package for . import
}

const (
	symOnExportList = 1 << iota // added to exportlist (no need to add again)
	symUniq
	symSiggen
	symAsm
	symAlgGen
)

func (sym *Sym) OnExportList() bool { return sym.flags&symOnExportList != 0 }
func (sym *Sym) Uniq() bool         { return sym.flags&symUniq != 0 }
func (sym *Sym) Siggen() bool       { return sym.flags&symSiggen != 0 }
func (sym *Sym) Asm() bool          { return sym.flags&symAsm != 0 }
func (sym *Sym) AlgGen() bool       { return sym.flags&symAlgGen != 0 }

func (sym *Sym) SetOnExportList(b bool) { sym.flags.set(symOnExportList, b) }
func (sym *Sym) SetUniq(b bool)         { sym.flags.set(symUniq, b) }
func (sym *Sym) SetSiggen(b bool)       { sym.flags.set(symSiggen, b) }
func (sym *Sym) SetAsm(b bool)          { sym.flags.set(symAsm, b) }
func (sym *Sym) SetAlgGen(b bool)       { sym.flags.set(symAlgGen, b) }

func (sym *Sym) IsBlank() bool {
	return sym != nil && sym.Name == "_"
}

func (sym *Sym) LinksymName() string {
	if sym.IsBlank() {
		return "_"
	}
	if sym.Linkname != "" {
		return sym.Linkname
	}
	return sym.Pkg.Prefix + "." + sym.Name
}

func (sym *Sym) Linksym() *obj.LSym {
	if sym == nil {
		return nil
	}
	return Ctxt.Lookup(sym.LinksymName())
}
