// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ctlog

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestFindLines(t *testing.T) {
	var cases = []struct {
		Input string
		Lines []Line
	}{
		{
			Input: `
#include <stdbool.h>
#include <stdio.h>
#include "ctlog.h"

CMODULE_DEFINE( main );

int
main( void )
{
    ctlog_setEnabled( true );
    CTLOG_INFO( "Test" );
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 123 ) );              // CTLOG_TYPE_UINT
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 456 ) );              // CTLOG_TYPE_UINT
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_UINT( 789 ) );              // CTLOG_TYPE_UINT
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -123 ) );              // CTLOG_TYPE_INT
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -456 ) );              // CTLOG_TYPE_INT
    CTLOG_VAR_INFO( "%d", 1, CTLOG_TYPE_INT( -678 ) );              // CTLOG_TYPE_INT
    CTLOG_VAR_INFO( "%s", 1, CTLOG_TYPE_STRING( "Hello World" ) );  // CTLOG_TYPE_STRING
    CTLOG_VAR_INFO( "%t", 1, CTLOG_TYPE_BOOL( true ) );             // CTLOG_TYPE_BOOL
    CTLOG_VAR_INFO( "%c", 1, CTLOG_TYPE_CHAR( 'J' ) );              // CTLOG_TYPE_CHAR
    return 0;
}
`,
			Lines: []Line{
				{
					Number:       12,
					FormatString: "Test",
				},
				{
					Number:       13,
					FormatString: "%d",
				},
				{
					Number:       14,
					FormatString: "%d",
				},
				{
					Number:       15,
					FormatString: "%d",
				},
				{
					Number:       16,
					FormatString: "%d",
				},
				{
					Number:       17,
					FormatString: "%d",
				},
				{
					Number:       18,
					FormatString: "%d",
				},
				{
					Number:       19,
					FormatString: "%s",
				},
				{
					Number:       20,
					FormatString: "%t",
				},
				{
					Number:       21,
					FormatString: "%c",
				},
			},
		},
	}

	for i, tc := range cases {
		t.Logf("Test case %d", i)

		lines, err := FindLines(strings.NewReader(tc.Input))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(lines, tc.Lines) {
			t.Error("data mismatch")
			t.Error(tc.Lines)
			t.Error(lines)
		}
	}
}

func TestHasTlogLine(t *testing.T) {
	var cases = []struct {
		Input     string
		Ok        bool
		ExpectErr bool
	}{
		// No errors
		{
			Input:     "$TL00, ",
			Ok:        true,
			ExpectErr: false,
		},
		{
			Input:     "$TL00",
			Ok:        false,
			ExpectErr: false,
		},
		{
			Input:     "abcd $TL00,",
			Ok:        false,
			ExpectErr: false,
		},

		// Errors
		{
			Input:     "$TL01,",
			Ok:        false,
			ExpectErr: true,
		},
	}

	for i, tc := range cases {
		t.Logf("Test case %d", i)

		ok, err := HasTlogLine([]byte(tc.Input))
		if err != nil {
			if !tc.ExpectErr {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if tc.ExpectErr {
				t.Error("expected error")
			} else {
				if ok != tc.Ok {
					t.Error("expected match")
				}
			}
		}
	}
}

func TestScanner(t *testing.T) {
	var cases = []struct {
		Input     string
		Lines     []string
		ExpectErr bool
	}{
		{
			Input: "abcdef $TL00,0,I,1,14,1,6,^\x00Enter main$\x00,\n" +
				"$TL00,1,I,0,23,2,5,0,5,1,\n" +
				"$TL00,2,I,0,23,2,5,1,5,1,\n" +
				"$TL00,3,I,0,23,2,5,1,5,2,\n" +
				"$TL00,4,I,0,23,2,5,2,5,3,\n" +
				"$TL00,5,I,0,23,2,5,3,5,5,\n" +
				"$TL00,6,I,0,23,2,5,5,5,8,\n" +
				"$TL00,7,I,0,23,2,5,8,5,13,\n" +
				"$TL00,8,I,0,23,2,5,13,5,21,\n" +
				"$TL00,9,I,0,23,2,5,21,5,34,\n" +
				"$TL00,10,I,0,23,2,5,34,5,55,\n" +
				"$TL00,11,I,1,16,1,6,^\x00Exit\n" +
				"main$\x00,\n",
			Lines: []string{
				"abcdef ",
				"$TL00,0,I,1,14,1,6,^\x00Enter main$\x00,",
				"$TL00,1,I,0,23,2,5,0,5,1,",
				"$TL00,2,I,0,23,2,5,1,5,1,",
				"$TL00,3,I,0,23,2,5,1,5,2,",
				"$TL00,4,I,0,23,2,5,2,5,3,",
				"$TL00,5,I,0,23,2,5,3,5,5,",
				"$TL00,6,I,0,23,2,5,5,5,8,",
				"$TL00,7,I,0,23,2,5,8,5,13,",
				"$TL00,8,I,0,23,2,5,13,5,21,",
				"$TL00,9,I,0,23,2,5,21,5,34,",
				"$TL00,10,I,0,23,2,5,34,5,55,",
				"$TL00,11,I,1,16,1,6,^\x00Exit\nmain$\x00,",
			},
			ExpectErr: false,
		},
	}

	for i, tc := range cases {
		t.Logf("Test case %d", i)

		var (
			s      = NewScanner(strings.NewReader(tc.Input))
			nLines = 0
			ls     = make([]string, 0)
		)
		for s.Scan() {
			l := s.Text()
			ls = append(ls, l)
			nLines += 1
		}

		if err := s.Err(); err != nil {
			if !tc.ExpectErr {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if tc.ExpectErr {
				t.Error("expected error")
			} else {
				if nLines == len(tc.Lines) {
					if !reflect.DeepEqual(ls, tc.Lines) {
						t.Error("data mismatch")
						logLines(t, ls)
						logLines(t, tc.Lines)
					}
				} else {
					t.Errorf("expected %d lines but got %d", len(tc.Lines), nLines)
					logLines(t, ls)
					logLines(t, tc.Lines)
				}
			}
		}
	}
}

func logLines(t *testing.T, lines []string) {
	t.Log("Lines:")
	for i, line := range lines {
		t.Logf("line[%d]='%s'\n", i, line)
		t.Logf("         '% X'\n", line)
	}
}

func TestParseOutput(t *testing.T) {
	var cases = []struct {
		Input     string
		Ok        bool
		Output    *Output
		ExpectErr bool
	}{
		{
			Input: "$TL00,2,I,12,34,3,4,123,2,-1,1,74,\n",
			Ok:    true,
			Output: &Output{
				Sequence:    uint16(2),
				Level:       LevelInfo,
				ModuleIndex: uint32(12),
				LineNumber:  uint32(34),
				Args: []Arg{
					{
						Type:  TypeUint,
						Value: uint32(123),
					},
					{
						Type:  TypeInt,
						Value: int32(-1),
					},
					{
						Type:  TypeChar,
						Value: byte('J'),
					},
				},
			},
			ExpectErr: false,
		},
		{
			Input: "$TL00,2,I,12,34,1,3,^\x00Exit\nfibonacci_log$\x00,\n",
			Ok:    true,
			Output: &Output{
				Sequence:    uint16(2),
				Level:       LevelInfo,
				ModuleIndex: uint32(12),
				LineNumber:  uint32(34),
				Args: []Arg{
					{
						Type:  TypeString,
						Value: "Exit\nfibonacci_log",
					},
				},
			},
			ExpectErr: false,
		},
	}

	for i, tc := range cases {
		t.Logf("Test case %d", i)

		out, ok, err := ParseOutput([]byte(tc.Input))
		if err != nil {
			if !tc.ExpectErr {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if tc.ExpectErr {
				t.Error("expected error")
			} else {
				if ok != tc.Ok {
					t.Error("expected to be able to parse")
				} else {
					if !reflect.DeepEqual(out, tc.Output) {
						t.Error("unexpected output")
						t.Errorf("%+v\n", tc.Output)
						t.Errorf("%+v\n", out)
					}
				}
			}
		}
	}
}

func TestTranslator(t *testing.T) {
	var cases = []struct {
		Module  Module
		Outputs []Output

		Exp       []string
		ExpectErr bool
	}{
		{
			Module: Module{
				Index: 0,
				Name:  "module_0",
				Path:  "/path/to/module_0.c",
				Lines: []Line{
					{
						Number:       123,
						FormatString: "string=%s",
					},
					{
						Number:       345,
						FormatString: "uint32=%d",
					},
				},
			},
			Outputs: []Output{
				{
					Sequence:    0,
					Level:       LevelInfo,
					ModuleIndex: 0,
					LineNumber:  123,
					Args: []Arg{
						{
							Type:  TypeString,
							Value: string("test"),
						},
					},
				},
				{
					Sequence:    0,
					Level:       LevelInfo,
					ModuleIndex: 0,
					LineNumber:  345,
					Args: []Arg{
						{
							Type:  TypeUint,
							Value: uint32(123456),
						},
					},
				},
			},
			Exp: []string{
				"string=test",
				"uint32=123456",
			},
			ExpectErr: false,
		},
	}

	for i, tc := range cases {
		t.Logf("Test case %d", i)

		modules := []Module{tc.Module}
		tx := NewTranslator(modules)
		for j, out := range tc.Outputs {
			s, err := tx.Translate(&out)
			if err != nil {
				if !tc.ExpectErr {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if tc.ExpectErr {
					t.Error("expected error")
				} else {
					if !reflect.DeepEqual(tc.Exp[j], s) {
						t.Error("data mismatch")
						t.Error(tc.Exp[j])
						t.Error(s)
					}
				}
			}
		}
	}
}

func TestUnmarshalOutput(t *testing.T) {
	var out Output
	if err := json.Unmarshal([]byte(`{"ctlog":0,"seq":7,"lvl":"I","mi":0,"ml":20,"args":[{"t":3,"v":"Hello \\ \" \u0001 World"}]}`), &out); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", out)
}
