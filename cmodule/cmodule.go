// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cmodule.
*/
package cmodule

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jlubawy/go-ctext/cmacro"
)

const MacroFuncName = "CMODULE_DEFINE"

type Module struct {
	// Index is the index in the sorted module slice.
	Index int `json:"index"`

	// Name is the name of the module. Modules are sorted by name.
	Name string `json:"name"`

	// Path is the absolute path to the C source file.
	Path string `json:"path"`
}

type modulesByName []Module

func (x modulesByName) Less(i, j int) bool {
	return x[i].Name < x[j].Name
}

func (x modulesByName) Len() int {
	return len(x)
}

func (x modulesByName) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// WalkDirs walks multiples directories and finds all modules within the given
// directories.
func WalkDirs(roots ...string) (modules []Module, err error) {
	modules = make([]Module, 0)

	for _, root := range roots {
		var ms []Module
		ms, err = walkDir(root)
		if err != nil {
			return
		}
		modules = append(modules, ms...)
	}

	sort.Sort(modulesByName(modules))
	for i := 0; i < len(modules); i++ {
		modules[i].Index = i
	}

	return
}

// WalkDir walks a directory and finds all modules within that given directory.
func WalkDir(root string) (modules []Module, err error) {
	modules, err = walkDir(root)
	if err != nil {
		return
	}

	sort.Sort(modulesByName(modules))
	for i := 0; i < len(modules); i++ {
		modules[i].Index = i
	}

	return
}

func walkDir(root string) (modules []Module, err error) {
	modules = make([]Module, 0)

	walkFn := func(path string, info os.FileInfo, err1 error) (err error) {
		// Return any errors
		if err1 != nil {
			err = err1
			return
		}

		if info.IsDir() {
			return // skip directories
		}

		ext := filepath.Ext(path)
		if ext != ".c" {
			return // skip files that aren't C source
		}

		// Convert path to absolute path
		path, err = PathAbsToSlash(path)
		if err != nil {
			return
		}

		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()

		count := 0
		err = cmacro.ScanInvocations(f, func(inv cmacro.Invocation) {
			if len(inv.Args) != 1 {
				err = fmt.Errorf("expected a single argument in the module definition but got %d", len(inv.Args))
				return
			}

			modules = append(modules, Module{
				Name: inv.Args[0],
				Path: path,
			})
		}, MacroFuncName)
		if err != nil {
			return
		}

		if count > 1 {
			err = fmt.Errorf("more than one module definition found in %s", path)
			return
		}

		return
	}
	err = filepath.Walk(root, walkFn)
	if err != nil {
		return
	}

	return
}

// PathAbsToSlash returns an absolute path with / slash characters.
func PathAbsToSlash(p string) (cp string, err error) {
	cp, err = filepath.Abs(p)
	if err != nil {
		return
	}
	cp = filepath.ToSlash(cp)
	return
}
