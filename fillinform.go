// Copyright 2015 The fillinform Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fillinform

import (
	_ "fmt"
	"log"
	"regexp"
	_ "time"

	"github.com/k0kubun/pp"
)

const (
	Form     = `[fF][oO][rR][mM]`
	Input    = `[iI][nN][pP][uU][tT]`
	Select   = `[sS][eE][lL][eE][cC][tT]`
	Option   = `[oO][pP][tT][iI][oO][nN]`
	Textarea = `[tT][eE][xX][tT][aA][rR][eE][aA]`
	Checked  = `[cC][hH][eE][cC][kK][eE][dD]`
	Selected = `[sS][eE][lL][eE][cC][tT][eE][dD]`
	Multiple = `[mM][uU][lL][tT][iI][pP][lL][eE]`
	Id       = `[iI][dD]`
	Type     = `[tT][yY][pP][eE]`
	Name     = `[nN][aA][mM][eE]`
	Value    = `[vV][aA][lL][uU][eE]`

	SPACE      = `\s`
	ATTR_NAME  = `[\w\-]+`
	ATTR_VALUE = `(?:"[^"]*"|'[^']*'|[^'"/>\s]+|[\w\-]+)`
	ATTR       = `(?:` + SPACE + `+(?:` + ATTR_NAME + `(?:=` + ATTR_VALUE + `)?))`

	FORM     = `(?:<` + Form + ATTR + `+` + SPACE + `*>)`     // <form>
	INPUT    = `(?:<` + Input + ATTR + `+` + SPACE + `*/?>)`  // <input>
	SELECT   = `(?:<` + Select + ATTR + `+` + SPACE + `*>)`   // <select>
	OPTION   = `(?:<` + Option + ATTR + `*` + SPACE + `*>)`   // <option>
	TEXTAREA = `(?:<` + Textarea + ATTR + `+` + SPACE + `*>)` // <textarea>

	EndFORM     = `(?:</` + Form + `>)`
	EndSELECT   = `(?:</` + Select + `>)`
	EndOPTION   = `(?:</` + Option + `>)`
	EndTEXTAREA = `(?:</` + Textarea + `>)`

	CHECKED  = `(?:` + Checked + ` (?:=(?:"` + Checked + `"|'` + Checked + `'|` + Checked + `))?)`
	SELECTED = `(?:` + Selected + ` (?:=(?:"` + Selected + `"|'` + Selected + `'|` + Selected + `))?)`
	MULTIPLE = `(?:` + Multiple + ` (?:=(?:"` + Multiple + `"|'` + Multiple + `'|` + Multiple + `))?)`
)

type FillinFormOptions struct {
	FillPassword bool
	IgnoreFields []string
	IgnoreTypes  map[string]bool
	Target       string
	Escape       bool
	DecodeEntity bool
}

type Filler struct {
	FillinFormOptions
}

func (f Filler) multilineRegexp(regstr string) *regexp.Regexp {
	return regexp.MustCompile(`(?ms:` + regstr + `)`)
}

func Fill(body *[]byte, data map[string]interface{}, options map[string]interface{}) ([]byte, error) {

	filler := &Filler{}

	return filler.fill(body, data)
}

func (f Filler) fill(body *[]byte, data map[string]interface{}) ([]byte, error) {
	log.Println("-- Start")

	// pp.Println(string(*body))
	pp.Println(data)

	reg := f.multilineRegexp(FORM + `(.*?)` + EndFORM)

	filled := reg.ReplaceAllFunc(*body, f.fillForm)

	log.Println("-- End")

	return filled, nil
}

func (f Filler) fillForm(formbody []byte) []byte {
	log.Println("-- -- Start")

	reg := f.multilineRegexp(INPUT)
	replaced := reg.ReplaceAllFunc(formbody, f.fillInput)

	reg = f.multilineRegexp(SELECT)
	replaced = reg.ReplaceAllFunc(replaced, f.fillSelect)

	reg = f.multilineRegexp(TEXTAREA)
	replaced = reg.ReplaceAllFunc(replaced, f.fillTextarea)

	log.Println("-- -- End")
	return replaced
}

func (f Filler) getInputType(tag []byte) []byte {
	reg := f.multilineRegexp(Type + `=(` + ATTR_VALUE + `)`)
	inputType := reg.Find(tag)
	log.Println(string(inputType))
	return inputType
}

func (f Filler) fillInput(tag []byte) []byte {
	log.Println("INPUT" + string(tag))

	inputType := f.getInputType(tag)
	if inputType == nil {
		inputType = []byte("text")
	}

	// ignore
	if _, ok := f.IgnoreTypes[string(inputType)]; ok {
		return tag
	}

	return tag
}

func (f Filler) fillSelect(tag []byte) []byte {
	log.Println("SELECT" + string(tag))

	return tag
}

func (f Filler) fillTextarea(tag []byte) []byte {
	log.Println("TEXTAREA" + string(tag))

	return tag
}
