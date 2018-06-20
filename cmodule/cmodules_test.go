// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmodule

import (
	"reflect"
	"testing"

	"github.com/jlubawy/go-ctlog/internal"
)

func TestWalkDirs(t *testing.T) {
	var exp = []Module{
		{
			Name: "module_1",
			Path: "testdata/walkdirs/a/module_1.c",
		},
		{
			Name: "module_2",
			Path: "testdata/walkdirs/a/module_2.c",
		},
		{
			Name: "module_3",
			Path: "testdata/walkdirs/b/module_3.c",
		},
	}
	for i := 0; i < len(exp); i++ {
		exp[i].Index = i
		cp, err := internal.PathAbsToSlash(exp[i].Path)
		if err != nil {
			t.Fatal(err)
		}
		exp[i].Path = cp
	}

	modules, err := WalkDirs("testdata/walkdirs/a", "testdata/walkdirs/b")
	if err != nil {
		t.Fatal(err)
	}

	if len(exp) != len(modules) {
		t.Fatalf("length mismatch: exp=%d, act=%d", len(exp), len(modules))
	}

	for i := 0; i < len(exp); i++ {
		if !reflect.DeepEqual(modules[i], exp[i]) {
			t.Error("data mismatch")
			t.Log(modules[i])
			t.Log(exp[i])
		}
	}
}

func TestWalkDir(t *testing.T) {
	var exp = []Module{
		{
			Name: "module_1",
			Path: "testdata/walkdir/module_1.c",
		},
		{
			Name: "module_2",
			Path: "testdata/walkdir/module_2.c",
		},
		{
			Name: "module_3",
			Path: "testdata/walkdir/module_3.c",
		},
	}
	for i := 0; i < len(exp); i++ {
		exp[i].Index = i
		cp, err := internal.PathAbsToSlash(exp[i].Path)
		if err != nil {
			t.Fatal(err)
		}
		exp[i].Path = cp
	}

	modules, err := WalkDir("testdata/walkdir")
	if err != nil {
		t.Fatal(err)
	}

	if len(exp) != len(modules) {
		t.Logf("%q", exp)
		t.Logf("%q", modules)
		t.Fatalf("length mismatch: exp=%d, act=%d", len(exp), len(modules))
	}

	for i := 0; i < len(exp); i++ {
		if !reflect.DeepEqual(modules[i], exp[i]) {
			t.Error("data mismatch")
			t.Log(modules[i])
			t.Log(exp[i])
		}
	}
}
