package fillinform

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
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

	FORM     = `(?:<` + Form + ATTR + `+` + SPACE + `*>)`
	INPUT    = `(?:<` + Input + ATTR + `+` + SPACE + `*/?>)`
	SELECT   = `(?:<` + Select + ATTR + `+` + SPACE + `*>)`
	OPTION   = `(?:<` + Option + ATTR + `*` + SPACE + `*>)`
	TEXTAREA = `(?:<` + Textarea + ATTR + `+` + SPACE + `*>)`

	EndFORM     = `(?:</` + Form + `>)`
	EndSELECT   = `(?:</` + Select + `>)`
	EndOPTION   = `(?:</` + Option + `>)`
	EndTEXTAREA = `(?:</` + Textarea + `>)`

	CHECKED  = `(?:` + Checked + `(?:=(?:"` + Checked + `"|'` + Checked + `'|` + Checked + `))?)`
	SELECTED = `(?:` + Selected + `(?:=(?:"` + Selected + `"|'` + Selected + `'|` + Selected + `))?)`
	MULTIPLE = `(?:` + Multiple + `(?:=(?:"` + Multiple + `"|'` + Multiple + `'|` + Multiple + `))?)`
)

var BACheckbox = []byte{'c', 'h', 'e', 'c', 'k', 'b', 'o', 'x'}
var BARadio = []byte{'r', 'a', 'd', 'i', 'o'}
var BAAmp = []byte{'&', 'a', 'm', 'p', ';'}
var BALt = []byte{'&', 'l', 't', ';'}
var BAGt = []byte{'&', 'g', 't', ';'}
var BAQuot = []byte{'&', 'q', 'u', 'o', 't', ';'}

var CompiledRegexpMap map[string]*regexp.Regexp

func init() {
	createRegexpMap()
}

func compileMultiLine(regstr string) *regexp.Regexp {
	return regexp.MustCompile(`(?ms:` + regstr + `)`)
}

func createRegexpMap() {
	CompiledRegexpMap = make(map[string]*regexp.Regexp)
	CompiledRegexpMap["form"] = compileMultiLine(FORM + `.*?` + EndFORM)
	CompiledRegexpMap["input"] = compileMultiLine(INPUT)
	CompiledRegexpMap["select"] = compileMultiLine(SELECT + `.*?` + EndSELECT)
	CompiledRegexpMap["textarea"] = compileMultiLine(TEXTAREA + `.*?` + EndTEXTAREA)
	CompiledRegexpMap["type"] = compileMultiLine(Type + `=(` + ATTR_VALUE + `)`)
	CompiledRegexpMap["value"] = compileMultiLine(Value + `=(` + ATTR_VALUE + `)`)
	CompiledRegexpMap["name"] = compileMultiLine(Name + `=(` + ATTR_VALUE + `)`)
	CompiledRegexpMap["checked"] = compileMultiLine(CHECKED)
	CompiledRegexpMap["space+>"] = compileMultiLine(SPACE + `*(/?)>\z`)
	CompiledRegexpMap["space+checked"] = compileMultiLine(SPACE + CHECKED)
	CompiledRegexpMap["value(nocapture)"] = compileMultiLine(Value + `=` + ATTR_VALUE)
	CompiledRegexpMap["textarea(3capture)"] = compileMultiLine(`(` + TEXTAREA + `)(.*?)(` + EndTEXTAREA + `)`)
	CompiledRegexpMap["multiple"] = compileMultiLine(MULTIPLE)
	CompiledRegexpMap["option"] = compileMultiLine(OPTION + `(.*?)` + EndOPTION)
	CompiledRegexpMap["option(nocapture)"] = compileMultiLine(OPTION + `.*?` + EndOPTION)
	CompiledRegexpMap["selected"] = compileMultiLine(SELECTED)
	CompiledRegexpMap["start option"] = compileMultiLine(OPTION)
	CompiledRegexpMap["tagend"] = compileMultiLine(SPACE + `*>\z`)
	CompiledRegexpMap["space+selected"] = compileMultiLine(SPACE + SELECTED)

	CompiledRegexpMap["&"] = regexp.MustCompile(`&`)
	CompiledRegexpMap["<"] = regexp.MustCompile(`<`)
	CompiledRegexpMap[">"] = regexp.MustCompile(`>`)
	CompiledRegexpMap[`"`] = regexp.MustCompile(`"`)
}

func (f Filler) compiledRegexp(key string) *regexp.Regexp {
	if reg, ok := CompiledRegexpMap[key]; ok {
		return reg
	}
	panic(fmt.Sprintf("no such compiled exp for: %v\n", key))
}

type FillinFormOptions struct {
	FillPassword bool
	IgnoreFields map[string]bool
	IgnoreTypes  map[string]bool
	Target       string
	Params       map[string][]byte
	Data         map[string]interface{}
}

type Filler struct {
	FillinFormOptions
}

type Writer struct {
	filler *Filler
	wr     io.Writer
}

func newFiller(data map[string]interface{}) *Filler {
	return &Filler{FillinFormOptions{Data: data, Params: make(map[string][]byte)}}
}

func FillWriter(wr io.Writer, data map[string]interface{}, options map[string]interface{}) io.Writer {
	filler := newFiller(data)
	return Writer{filler: filler, wr: wr}
}
func (w Writer) Write(p []byte) (int, error) {
	filled := w.filler.fill(p)
	return w.wr.Write(filled)
}

func Fill(body []byte, data map[string]interface{}, options map[string]interface{}) []byte {
	filler := newFiller(data)
	return filler.fill(body)
}

func (f Filler) fill(body []byte) []byte {
	return f.compiledRegexp("form").ReplaceAllFunc(body, f.fillForm)
}

func (f Filler) fillForm(formbody []byte) []byte {
	replaced := f.compiledRegexp("input").ReplaceAllFunc(formbody, f.fillInput)
	replaced = f.compiledRegexp("select").ReplaceAllFunc(replaced, f.fillSelect)
	replaced = f.compiledRegexp("textarea").ReplaceAllFunc(replaced, f.fillTextarea)

	return replaced
}

func (f Filler) unquote(tag []byte) []byte {
	return bytes.Trim(tag, `'"`)
}

func (f Filler) getType(tag []byte) []byte {
	itype := f.compiledRegexp("type").FindSubmatch(tag)
	if cap(itype) == 2 {
		return f.unquote(itype[1])
	}
	return []byte{}
}

func (f Filler) getValue(tag []byte) []byte {
	value := f.compiledRegexp("value").FindSubmatch(tag)
	if cap(value) == 2 {
		return f.unquote(value[1])
	}
	return []byte{}
}

func (f Filler) getName(tag []byte) []byte {
	name := f.compiledRegexp("name").FindSubmatch(tag)
	if cap(name) == 2 {
		return f.unquote(name[1])
	}
	return []byte{}
}

func (f Filler) escapeHTML(tag []byte) []byte {
	tag = f.compiledRegexp("&").ReplaceAll(tag, BAAmp)
	tag = f.compiledRegexp("<").ReplaceAll(tag, BALt)
	tag = f.compiledRegexp(">").ReplaceAll(tag, BAGt)
	tag = f.compiledRegexp(`"`).ReplaceAll(tag, BAQuot)

	return tag
}

func (f Filler) getParam(name []byte) ([]byte, bool) {
	// ignore
	nameStr := string(name)
	if _, ok := f.IgnoreFields[nameStr]; ok {
		return []byte{}, false
	}
	if param, ok := f.Params[nameStr]; ok {
		return param, true
	}
	if param, ok := f.Data[nameStr]; ok {
		if casted, ok := param.(string); ok {
			f.Params[nameStr] = []byte(casted)
			return f.Params[nameStr], true
		} else {
			fmt.Printf("!!!cannot cast to []byte for %v\n", param)
		}
	}

	return []byte{}, false
}

func (f Filler) fillInput(tag []byte) []byte {
	inputType := f.getType(tag)
	if bytes.Equal(inputType, []byte{}) {
		inputType = []byte("text")
	}

	// ignore
	if _, ok := f.IgnoreTypes[string(inputType)]; ok {
		return tag
	}

	paramValue, exists := f.getParam(f.getName(tag))
	if !exists {
		return tag
	}

	if bytes.Equal(inputType, BACheckbox) || bytes.Equal(inputType, BARadio) {
		value := f.getValue(tag)

		if bytes.Equal(paramValue, value) {
			if !f.compiledRegexp("checked").Match(tag) {
				tag = f.compiledRegexp("space+>").ReplaceAll(tag, []byte(` checked="checked"$1>`))
			}
		} else {
			tag = f.compiledRegexp("space+checked").ReplaceAll(tag, []byte(``))
		}
	} else { // text
		escapedValue := f.escapeHTML(paramValue)
		reg := f.compiledRegexp("value(nocapture)")
		if reg.Match(tag) {
			tag = reg.ReplaceAll(tag, append([]byte(`value="`), append(escapedValue, []byte(`"`)...)...))
		} else {
			tag = f.compiledRegexp("space+>").ReplaceAll(tag, append([]byte(` value="`), append(escapedValue, []byte(`"$1>`)...)...))
		}
	}

	return tag
}

func (f Filler) fillTextarea(tag []byte) []byte {
	paramValue, exists := f.getParam(f.getName(tag))
	if !exists {
		return tag
	}
	escapedValue := f.escapeHTML(paramValue)
	replaced := append([]byte(`${1}`), append(escapedValue, []byte(`${3}`)...)...)
	tag = f.compiledRegexp("textarea(3capture)").ReplaceAll(tag, replaced)

	return tag
}

func (f Filler) fillSelect(tag []byte) []byte {
	paramValue, exists := f.getParam(f.getName(tag))
	if !exists {
		return tag
	}

	if f.compiledRegexp("multiple").Match(tag) {
		return tag
	}

	return f.compiledRegexp("option(nocapture)").ReplaceAllFunc(tag, func(tag []byte) []byte { return f.fillOption(tag, paramValue) })
}

func (f Filler) fillOption(tag, paramValue []byte) []byte {
	value := f.getValue(tag)
	if bytes.Equal(value, []byte{}) {
		value = f.compiledRegexp("option").ReplaceAll(tag, []byte(`$1`))
	}

	if bytes.Equal(paramValue, value) {
		if !f.compiledRegexp("selected").Match(tag) {
			tag = f.compiledRegexp("start option").ReplaceAllFunc(tag, func(tag []byte) []byte {
				return f.compiledRegexp("tagend").ReplaceAll(tag, []byte(` selected="selected">`))
			})
		}
	} else {
		tag = f.compiledRegexp("space+selected").ReplaceAll(tag, []byte(``))
	}

	return tag
}
