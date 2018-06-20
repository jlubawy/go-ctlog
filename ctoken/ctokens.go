// Copyright 2018 Josh Lubawy. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package ctoken implements a naive C tokenizer.
*/
package ctoken

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/jlubawy/go-ctlog/internal"
)

// A TokenType is the type of token.
type TokenType int

const (
	// TokenTypeError is an error token type. The error can be retrieved by
	// calling the tokenizer Err() method.
	TokenTypeError TokenType = iota
	// TokenTypeComment is a comment token type. The comment token can be
	// retrieved by calling the tokenizer Comment() method.
	TokenTypeComment
	// TokenTypeText is a text token. The text data can be retrieved by calling
	// the tokenizer Text() method.
	TokenTypeText
)

func (tt TokenType) String() (s string) {
	switch tt {
	case TokenTypeError:
		s = "error"
	case TokenTypeComment:
		s = "comment"
	case TokenTypeText:
		s = "text"
	default:
		s = fmt.Sprintf("TokenType(%d)", tt)
	}
	return
}

// A Token is either of type Comment or Text.
type Token struct {
	// Type is the token type.
	Type string `json:"type"`

	// Data is the token content.
	Data string `json:"data"`

	// LineStart and LineEnd are the line numbers at which the token starts and
	// ends, respectively. The first line is always one.
	LineStart uint32 `json:"lineStart"`

	// ColumnStart and ColumnEnd are the line positions at which the token starts and
	// ends, respectively.
	ColumnStart uint32 `json:"columnStart"`
}

// A Tokenizer is used to split a C source file into comment and text tokens for
// further processing.
type Tokenizer struct {
	br     *bufio.Reader
	maxBuf int
	buf    []byte

	err error

	// Track lines and positions
	lineCurr, lineStart     uint32
	columnCurr, columnStart uint32

	// Should be reset every invocation of Next
	currentTT       TokenType
	inStringLiteral bool
	mlCommentCount  int
	inSLComment     bool
}

// NewTokenizer returns a pointer to a new tokenizer.
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		br:  bufio.NewReader(r),
		buf: make([]byte, 0, 4096),

		// Line number and column start at one.
		lineCurr:   1,
		columnCurr: 1,
	}
}

// Comment returns a comment token if the last token type was TokenTypeComment.
func (z *Tokenizer) Comment() Token {
	if z.currentTT != TokenTypeComment {
		panic("token type was not comment")
	}
	return z.token()
}

// Err returns the error associated with the most recent ErrorToken token.
// This is typically io.EOF, meaning the end of tokenization.
func (z *Tokenizer) Err() error {
	if z.currentTT == TokenTypeError && z.err == nil {
		panic("token type was error but there was no error")
	}
	return z.err
}

// Next returns the next token type to be processed.
func (z *Tokenizer) Next() TokenType {
	// Return error right away if one already exists
	if z.err != nil {
		return TokenTypeError
	}

	// Reset the buffer length to 0
	z.buf = z.buf[:0]

	// Reset tokenizer fields
	z.currentTT = TokenTypeError
	z.inStringLiteral = false
	z.mlCommentCount = 0
	z.inSLComment = false
	z.lineStart = 0
	z.columnStart = 0

	for done := false; !done; {
		// Peek one character first so we can skip any chars we don't want
		var bs []byte
		bs, z.err = z.br.Peek(1)
		if z.err != nil {
			if z.err == io.EOF {
				if len(z.buf) > 0 {
					// If EOF but there is data in the buffer then process it first,
					// the EOF will be returned on the next call to this function.
					if z.mlCommentCount > 0 {
						z.err = errors.New("unexpected end of multi-line comment")
						z.currentTT = TokenTypeError
					} else if z.inSLComment {
						z.currentTT = TokenTypeText
					} else {
						z.currentTT = TokenTypeText
					}
					return z.currentTT
				}
			}

			return TokenTypeError
		}

		if z.lineStart == 0 {
			z.lineStart = z.lineCurr
			if z.columnCurr == 0 {
				z.columnStart = 1
			} else {
				z.columnStart = z.columnCurr
			}
		}

		b := bs[0]
		switch b {
		case '/':
			if !z.inSLComment && z.mlCommentCount == 0 {
				// If not in a comment

				if !z.inStringLiteral {
					// If not in a string literal check if this is the start
					// of a single-line comment.
					lc, ok := internal.LastChar(z.buf)
					if ok && lc == '/' {
						// Check if this is the start of a comment
						z.inSLComment = true
						z.lineStart, z.columnStart = z.lineCurr, z.columnCurr-1
					} else if len(z.buf) > 0 {
						// If the buffer is not empty then process the text first
						z.currentTT = TokenTypeText
						return z.currentTT
					}
				}

			} else if z.mlCommentCount > 0 {
				// Else if in a multi-line comment
				lc, ok := internal.LastChar(z.buf)
				if ok && lc == '*' {
					z.mlCommentCount -= 1

					if z.mlCommentCount == 0 {
						z.currentTT = TokenTypeComment
						done = true
					}
				}
			} else {
				// Else if in a single-line comment do nothing
			}

		case '*':
			// Possible start or end of multi-line comment
			lc, ok := internal.LastChar(z.buf)
			if ok && lc == '/' {
				z.mlCommentCount += 1
				if z.mlCommentCount == 1 {
					z.lineStart, z.columnStart = z.lineCurr, z.columnCurr-1
				}
			}

		case '\r':
			// Discard and wait for the \n
			_, z.err = z.br.Discard(1)
			if z.err != nil {
				return TokenTypeError
			}

			z.columnCurr += 1

			continue

		case '\n':
			// Increment the line and reset the current column
			z.lineCurr += 1
			z.columnCurr = 0

			if z.mlCommentCount > 0 {
				// If in a multi-line comment then continue processing
			} else if z.inSLComment {
				z.inSLComment = false
				z.currentTT = TokenTypeComment
				done = true
			}

		case '"':
			if !z.inSLComment && z.mlCommentCount == 0 {
				lc, ok := internal.LastChar(z.buf)
				if ok && lc != '\\' {
					z.inStringLiteral = !z.inStringLiteral
				}
			}
		}

		b, z.err = z.br.ReadByte()
		if z.err != nil {
			// EOF is not expected since we already peeked successfully above
			return TokenTypeError
		}

		z.columnCurr += 1

		z.buf, z.err = internal.AddChar(&z.buf, z.maxBuf, b)
		if z.err != nil {
			return TokenTypeError
		}
	}

	return z.currentTT
}

// SetMaxBuf sets the maximum buffer allowed by the tokenizer. Zero is the default
// and it means an unlimited buffer size.
func (z *Tokenizer) SetMaxBuf(maxBuf uint) {
	z.maxBuf = int(maxBuf)
}

// Text returns a text token if the last token type was TokenTypeText.
func (z *Tokenizer) Text() Token {
	if z.currentTT != TokenTypeText {
		panic("token type was not text")
	}
	return z.token()
}

func (z *Tokenizer) token() Token {
	return Token{
		Type:        z.currentTT.String(),
		Data:        string(z.buf[:]),
		LineStart:   z.lineStart,
		ColumnStart: z.columnStart,
	}
}
