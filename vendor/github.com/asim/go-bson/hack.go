// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hack gives you some efficient functionality at the cost of
// breaking some Go rules.
package bson

import (
	"reflect"
	"unsafe"
)

// StringArena lets you consolidate allocations for a group of strings
// that have similar life length
type StringArena struct {
	buf []byte
	str string
}

// NewStringArena creates an arena of the specified size.
func NewStringArena(size int) *StringArena {
	sa := &StringArena{buf: make([]byte, 0, size)}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&sa.buf))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&sa.str))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Cap
	return sa
}

// NewString copies a byte slice into the arena and returns it as a string.
// If the arena is full, it returns a traditional go string.
func (sa *StringArena) NewString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if len(sa.buf)+len(b) > cap(sa.buf) {
		return string(b)
	}
	start := len(sa.buf)
	sa.buf = append(sa.buf, b...)
	return sa.str[start : start+len(b)]
}

// SpaceLeft returns the amount of space left in the arena.
func (sa *StringArena) SpaceLeft() int {
	return cap(sa.buf) - len(sa.buf)
}

// String force casts a []byte to a string.
// USE AT YOUR OWN RISK
func StringFn(b []byte) (s string) {
	if len(b) == 0 {
		return ""
	}
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// StringPointer returns &s[0], which is not allowed in go
func StringPointer(s string) unsafe.Pointer {
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(pstring.Data)
}
