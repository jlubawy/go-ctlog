// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package ctlog.
*/
package ctlog

import (
	"bufio"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jlubawy/go-ctlog/cmacros"
	"github.com/jlubawy/go-ctlog/ctoken"
)

type Level byte

const (
	LevelDebug Level = 'D'
	LevelError Level = 'E'
	LevelInfo  Level = 'I'
	LevelWarn  Level = 'W'
)

var _ encoding.TextUnmarshaler = (*Level)(nil)

func (lvl *Level) UnmarshalText(data []byte) (err error) {
	if len(data) == 1 {
		switch data[0] {
		case 'D':
			*lvl = LevelDebug
		case 'E':
			*lvl = LevelError
		case 'I':
			*lvl = LevelInfo
		case 'W':
			*lvl = LevelWarn
		default:
			err = fmt.Errorf("unsupported level '%s'", string(data))
		}
	} else {
		err = fmt.Errorf("unsupported level '%s'", string(data))
	}
	return
}

var MacroFuncNames = []string{
	"CTLOG_ERROR",
	"CTLOG_VAR_ERROR",
	"CTLOG_INFO",
	"CTLOG_VAR_INFO",
	"CTLOG_DEBUG",
	"CTLOG_VAR_DEBUG",
	"CTLOG_WARN",
	"CTLOG_VAR_WARN",
}

type Type int

const (
	TypeBool   Type = 0x00
	TypeChar   Type = 0x01
	TypeInt    Type = 0x02
	TypeString Type = 0x03
	TypeUint   Type = 0x04
)

type Module struct {
	// Index is the index in the sorted module slice.
	Index int `json:"index"`

	// Name is the name of the module. Modules are sorted by name.
	Name string `json:"name"`

	// Path is the absolute path to the C source file.
	Path string `json:"path"`

	Lines []Line `json:"lines"`
}

type Line struct {
	// Number is the line number of the tokenized logging output.
	Number uint32 `json:"number"`

	// FormatString is the format string that should be used for formatting the
	// tokenized logging variable output.
	FormatString string `json:"formatString"`
}

// FindLines finds all tokenized logging lines within the given io.Reader.
func FindLines(r io.Reader) (lines []Line, err error) {
	lines = make([]Line, 0)

	z := ctoken.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case ctoken.TokenTypeError:
			err = z.Err()
			if err == io.EOF {
				err = nil
			}
			return

		case ctoken.TokenTypeText:
			var (
				tok = z.Text()
				mfs []cmacros.MacroFunc
			)
			mfs, err = cmacros.FindMacroFuncs(&tok, MacroFuncNames...)
			if err != nil {
				return
			}

			if len(mfs) == 0 {
				continue // skip tokens with no tokenized logging invocation
			}

			for _, mf := range mfs {
				rs := mf.Args[0]
				if rs[0] != '"' {
					err = fmt.Errorf("format string missing opening quote")
					return
				}
				if rs[len(rs)-1] != '"' {
					err = fmt.Errorf("format string missing closing quote")
					return
				}

				lines = append(lines, Line{
					Number:       mf.LineEnd,
					FormatString: rs[1 : len(rs)-1],
				})
			}
		}
	}

	return
}

const MagicString = "$TL"

const (
	Version0            = uint8(0x00)
	MaxSupportedVersion = Version0
)

// HasTlogLine returns true if the provided byte slice might contain a tokenized
// logging line. If it does it would match the form '$TL00,'.
func HasTlogLine(data []byte) (ok bool, err error) {
	if len(data) >= 6 {
		if string(data[0:3]) == MagicString && data[5] == ',' {
			var n uint64
			n, err = strconv.ParseUint(string(data[3:5]), 16, 8)
			if err != nil {
				return
			}
			if uint8(n) > MaxSupportedVersion {
				err = fmt.Errorf("version 0x%02X exceeds max supported version 0x%02X", n, MaxSupportedVersion)
				return
			}
			ok = true
		}
	}
	return
}

// NewScanner returns a *bufio.Scanner with the SplitFunc set to ScanLines.
func NewScanner(r io.Reader) (s *bufio.Scanner) {
	s = bufio.NewScanner(r)
	s.Split(ScanLines)
	return
}

// ScanLines is like bufio.ScanLines except that it takes into account that there
// might be \r or \n characters inside tokenized log string literals.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var (
		inTlogLine bool
		inTlogStr  bool
	)

	for i := 0; i < len(data); i++ {
		b := data[i]

		switch b {
		case '\n':
			if !inTlogStr {
				// If not in a tlog string return the current data

				var end int
				if i >= 1 && data[i-1] == '\r' {
					end = i - 1
				} else {
					end = i
				}
				advance = i + 1
				token = data[0:end]
				return
			}

		case '$':
			if !inTlogLine && !inTlogStr {
				if i > 0 {
					// If we see a $ char that isn't at the beginning of a line then
					// return the existing data
					advance = i
					token = data[0:i]
					return
				}

				// Else this could be the start of a tlog line so find
				// the start and version if there is one
				inTlogLine, err = HasTlogLine(data)
				if err != nil {
					return
				}
			}

		case '\x00':
			// If NUL char then probably the start or end of a tlog string
			if inTlogLine && i >= 1 {
				lc := data[i-1]
				if lc == '^' {
					inTlogStr = true
				} else if lc == '$' {
					inTlogStr = false
				}
			}
		}
	}

	if atEOF {
		// If at the EOF return the current data, if any
		advance = len(data) + 1
		token = data[0:]
		err = io.EOF
	}

	return
}

type Output struct {
	// Sequence is the current log lines sequence number. It is useful for
	// determining if lines have been dropped.
	Sequence uint16 `json:"seq"`

	// Level is this output's logging level used for filtering.
	Level Level `json:"lvl"`

	// ModuleIndex is the module index that this output belongs to.
	ModuleIndex uint32 `json:"mi"`

	// LineNumber is the line number within the module that this output belongs to.
	LineNumber uint32 `json:"ml"`

	// Args is a slice of typed arguments to go with the format string.
	Args []Arg `json:"args"`
}

func (o *Output) Vals() []interface{} {
	vs := make([]interface{}, len(o.Args))
	for i := 0; i < len(vs); i++ {
		vs[i] = o.Args[i].Value
	}
	return vs
}

type Arg struct {
	Type  Type        `json:"t"`
	Value interface{} `json:"v"`
}

func (a *Arg) UnmarshalJSON(data []byte) (err error) {
	var v struct {
		Type  Type        `json:"t"`
		Value interface{} `json:"v"`
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return
	}

	switch v.Type {
	case TypeBool:
		x, ok := v.Value.(bool)
		if ok {
			v.Value = x
		}
	case TypeChar:
		x, ok := v.Value.(string)
		if ok {
			if len(x) == 0 {
				err = fmt.Errorf("empty character found")
				return
			}
			v.Value = x[0]
		}
	case TypeInt:
		x, ok := v.Value.(float64)
		if ok {
			v.Value = int32(x)
		}
	case TypeString:
		// string doesn't require casting
	case TypeUint:
		x, ok := v.Value.(float64)
		if ok {
			v.Value = uint32(x)
		}
	default:
		err = fmt.Errorf("unsupported type %d", v.Type)
	}
	a.Type = v.Type
	a.Value = v.Value
	return
}

type state int

const (
	psInit state = iota
	psSeq
	psLevel
	psModuleIdx
	psLine
	psNArgs
	psType
	psVal
)

func ParseOutput(data []byte) (output *Output, ok bool, err error) {
	output = new(Output)

	var (
		s     state
		nArgs int
		iArg  int
	)

	for {
		switch s {
		case psInit:
			ok, err = HasTlogLine(data)
			if err != nil {
				return
			}
			if !ok {
				return
			}
			if len(data) < 8 {
				err = fmt.Errorf("missing data after magic string")
				return
			}
			data = data[6:]
			s = psSeq

		case psSeq:
			ci := strings.IndexByte(string(data), ',')
			if ci == -1 {
				err = fmt.Errorf("expected sequence number comma but found none")
				return
			}
			var n uint64
			n, err = strconv.ParseUint(string(data[0:ci]), 10, 16)
			if err != nil {
				err = fmt.Errorf("error parsing sequence number: %v", err)
				return
			}
			output.Sequence = uint16(n)
			if len(data) < ci+1 {
				err = fmt.Errorf("missing data after sequence number")
				return
			}
			data = data[ci+1:]
			s = psLevel

		case psLevel:
			lvl := Level(data[0])
			switch lvl {
			case LevelDebug:
				output.Level = LevelDebug
			case LevelError:
				output.Level = LevelError
			case LevelInfo:
				output.Level = LevelInfo
			case LevelWarn:
				output.Level = LevelWarn
			default:
				err = fmt.Errorf("unsupported logging level '%c'", lvl)
				return
			}
			if len(data) < 3 {
				err = fmt.Errorf("missing data after logging level")
				return
			}
			data = data[2:]
			s = psModuleIdx

		case psModuleIdx:
			ci := strings.IndexByte(string(data), ',')
			if ci == -1 {
				err = fmt.Errorf("expected module index comma but found none")
				return
			}
			var n uint64
			n, err = strconv.ParseUint(string(data[0:ci]), 10, 32)
			if err != nil {
				err = fmt.Errorf("error parsing module index: %v", err)
				return
			}
			output.ModuleIndex = uint32(n)
			if len(data) < ci+1 {
				err = fmt.Errorf("missing data after module index")
				return
			}
			data = data[ci+1:]
			s = psLine

		case psLine:
			ci := strings.IndexByte(string(data), ',')
			if ci == -1 {
				err = fmt.Errorf("expected line number comma but found none")
				return
			}
			var n uint64
			n, err = strconv.ParseUint(string(data[0:ci]), 10, 32)
			if err != nil {
				err = fmt.Errorf("error parsing line number: %v", err)
				return
			}
			output.LineNumber = uint32(n)
			if len(data) < ci+1 {
				err = fmt.Errorf("missing data after line number")
				return
			}
			data = data[ci+1:]
			s = psNArgs

		case psNArgs:
			ci := strings.IndexByte(string(data), ',')
			if ci == -1 {
				err = fmt.Errorf("expected argument count comma but found none")
				return
			}
			var n uint64
			n, err = strconv.ParseUint(string(data[0:ci]), 10, 8)
			if err != nil {
				err = fmt.Errorf("error parsing argument count: %v", err)
				return
			}
			nArgs = int(n)

			if nArgs == 0 {
				// If no args then return
				return
			}
			if len(data) < ci+1 {
				err = fmt.Errorf("missing data after argument count")
				return
			}
			output.Args = make([]Arg, nArgs)
			data = data[ci+1:]
			s = psType

		case psType:
			ci := strings.IndexByte(string(data), ',')
			if ci == -1 {
				err = fmt.Errorf("expected argument %d type comma but found none", iArg)
				return
			}
			var n uint64
			n, err = strconv.ParseUint(string(data[0:ci]), 10, 8)
			if err != nil {
				err = fmt.Errorf("error parsing argument %d type: %v", iArg, err)
				return
			}
			output.Args[iArg].Type = Type(n)

			if len(data) < ci+1 {
				err = fmt.Errorf("missing data after argument %d type", iArg)
				return
			}
			data = data[ci+1:]
			s = psVal

		case psVal:
			t := output.Args[iArg].Type
			if t == TypeString {
				// If the argument type is string
				if string(data[0:2]) != "^\x00" {
					err = fmt.Errorf("missing start of argument %d string '% X'", iArg, data[0:2])
					return
				}
				end := strings.Index(string(data[0:]), "$\x00")
				if end == -1 {
					err = fmt.Errorf("missing end of argument %d string", iArg)
					return
				}
				output.Args[iArg].Value = string(data[2:end])

				if len(data) < end+3 {
					err = fmt.Errorf("missing data after argument %d string", iArg)
					return
				}
				data = data[end+3:]

			} else {
				// Else numerical arguments
				ci := strings.IndexByte(string(data), ',')
				if ci == -1 {
					err = fmt.Errorf("expected argument %d value comma but found none", iArg)
					return
				}

				switch t {
				case TypeChar, TypeUint:
					var (
						n    uint64
						size int
					)
					if t == TypeChar {
						size = 8
					} else { // if t == TypeUint {
						size = 32
					}
					n, err = strconv.ParseUint(string(data[0:ci]), 10, size)
					if err != nil {
						err = fmt.Errorf("error parsing argument %d value: %v", iArg, err)
						return
					}
					if t == TypeChar {
						output.Args[iArg].Value = byte(n)
					} else { // if t == TypeUint {
						output.Args[iArg].Value = uint32(n)
					}

				case TypeInt:
					var n int64
					n, err = strconv.ParseInt(string(data[0:ci]), 10, 32)
					if err != nil {
						err = fmt.Errorf("error parsing argument %d value: %v", iArg, err)
						return
					}
					output.Args[iArg].Value = int32(n)

				//case TypeString:

				case TypeBool:
					var b bool
					b, err = strconv.ParseBool(string(data[0:ci]))
					if err != nil {
						err = fmt.Errorf("error parsing argument %d value: %v", iArg, err)
						return
					}
					output.Args[iArg].Value = b

				default:
					err = fmt.Errorf("unknown argument type 0x%02X", t)
					return
				}

				if len(data) < ci+1 {
					err = fmt.Errorf("missing data after argument value")
					return
				}
				data = data[ci+1:]
			}

			s = psType
			iArg += 1

			if iArg >= len(output.Args) {
				// If done then return
				return
			}
		}
	}

	return
}

type Translator struct {
	modules []Module
}

func NewTranslator(modules []Module) *Translator {
	return &Translator{
		modules: modules,
	}
}

func (t *Translator) Translate(output *Output) (s string, err error) {
	// Find module first
	if output.ModuleIndex >= uint32(len(t.modules)) {
		err = fmt.Errorf("could not find module %d", output.ModuleIndex)
		return
	}
	module := t.modules[output.ModuleIndex]

	// Find the line within the module
	for _, line := range module.Lines {
		if line.Number == output.LineNumber {
			s = fmt.Sprintf(line.FormatString, output.Vals()...)
			return
		}
	}

	err = fmt.Errorf("could not find line %d in module %d", output.LineNumber, output.ModuleIndex)
	return
}
