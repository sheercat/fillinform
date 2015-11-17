package fillinform

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

const (
	_Form     = `form`
	_Input    = `input`
	_Select   = `select`
	_Option   = `option`
	_Textarea = `textarea`
	_Checked  = `checked`
	_Selected = `selected`
	_Multiple = `multiple`
	_Id       = `id`
	_Type     = `type`
	_Name     = `name`
	_Value    = `value`

	Space     = `\s`
	AttrName  = `[\w\-]+`
	AttrValue = `(?:"[^"]*"|'[^']*'|[^'"/>\s]+|[\w\-]+)`
	Attr      = `(?:` + Space + `+(?:` + AttrName + `(?:=` + AttrValue + `)?))`

	StartForm     = `(?:<` + _Form + Attr + `+` + Space + `*>)`
	Input         = `(?:<` + _Input + Attr + `+` + Space + `*/?>)`
	StartSelect   = `(?:<` + _Select + Attr + `+` + Space + `*>)`
	StartOption   = `(?:<` + _Option + Attr + `*` + Space + `*>)`
	StartTextarea = `(?:<` + _Textarea + Attr + `+` + Space + `*>)`

	EndForm     = `(?:</` + _Form + `>)`
	EndSelect   = `(?:</` + _Select + `>)`
	EndOption   = `(?:</` + _Option + `>)`
	EndTextarea = `(?:</` + _Textarea + `>)`

	Checked  = `(?:` + _Checked + `(?:=(?:"` + _Checked + `"|'` + _Checked + `'|` + _Checked + `))?)`
	Selected = `(?:` + _Selected + `(?:=(?:"` + _Selected + `"|'` + _Selected + `'|` + _Selected + `))?)`
	Multiple = `(?:` + _Multiple + `(?:=(?:"` + _Multiple + `"|'` + _Multiple + `'|` + _Multiple + `))?)`
)

var BACheckbox = []byte(`checkbox`)
var BARadio = []byte(`radio`)
var BAAmp = []byte(`&amp;`)
var BALt = []byte(`&lt;`)
var BAGt = []byte(`&gt;`)
var BAQuot = []byte(`&quot;`)
var BAChecked = []byte(` checked="checked"$1>`)
var BASelected = []byte(` selected="selected">`)
var BABlank = []byte(``)
var BADQuot = []byte(`"`)

var CompiledRegexpMap map[string]*regexp.Regexp

func init() {
	createRegexpMap()
}

func compileMultiLine(regstr string) *regexp.Regexp {
	return regexp.MustCompile(`(?msi:` + regstr + `)`)
}

func createRegexpMap() {
	CompiledRegexpMap = make(map[string]*regexp.Regexp)
	CompiledRegexpMap["form"] = compileMultiLine(StartForm + `.*?` + EndForm)
	CompiledRegexpMap["start form"] = compileMultiLine(`(` + StartForm + `)`)

	CompiledRegexpMap["input"] = compileMultiLine(Input)
	CompiledRegexpMap["select"] = compileMultiLine(StartSelect + `.*?` + EndSelect)
	CompiledRegexpMap["textarea"] = compileMultiLine(StartTextarea + `.*?` + EndTextarea)

	CompiledRegexpMap["id"] = compileMultiLine(_Id + `=(` + AttrValue + `)`)
	CompiledRegexpMap["type"] = compileMultiLine(_Type + `=(` + AttrValue + `)`)
	CompiledRegexpMap["value"] = compileMultiLine(_Value + `=(` + AttrValue + `)`)
	CompiledRegexpMap["name"] = compileMultiLine(_Name + `=(` + AttrValue + `)`)

	CompiledRegexpMap["checked"] = compileMultiLine(Checked)
	CompiledRegexpMap["space+>"] = compileMultiLine(Space + `*(/?)>\z`)
	CompiledRegexpMap["space+checked"] = compileMultiLine(Space + Checked)
	CompiledRegexpMap["value(nocapture)"] = compileMultiLine(_Value + `=` + AttrValue)
	CompiledRegexpMap["textarea(3capture)"] = compileMultiLine(`(` + StartTextarea + `)(.*?)(` + EndTextarea + `)`)
	CompiledRegexpMap["multiple"] = compileMultiLine(Multiple)
	CompiledRegexpMap["option"] = compileMultiLine(StartOption + `(.*?)` + EndOption)
	CompiledRegexpMap["option(nocapture)"] = compileMultiLine(StartOption + `.*?` + EndOption)
	CompiledRegexpMap["selected"] = compileMultiLine(Selected)
	CompiledRegexpMap["start option"] = compileMultiLine(StartOption)
	CompiledRegexpMap["tag end"] = compileMultiLine(Space + `*>\z`)
	CompiledRegexpMap["space+selected"] = compileMultiLine(Space + Selected)
}

func (f Filler) compiledRegexp(key string) *regexp.Regexp {
	if reg, ok := CompiledRegexpMap[key]; ok {
		return reg
	}
	panic(fmt.Sprintf("no such compiled exp for: %v\n", key))
}

// Options for fillin
// Set { "FillPassword": true } if fillin value to field type="password".
// Target is id for form tag.
type FillInFormOptions struct {
	IgnoreFields map[string]bool
	IgnoreTypes  map[string]bool
	Target       string
}

type Filler struct {
	FillInFormOptions
	Params map[string][]byte
	Data   map[string]interface{}
}

type Writer struct {
	filler *Filler
	wr     io.Writer
}

func setOptions(options map[string]interface{}) *FillInFormOptions {
	var ffo FillInFormOptions
	// default set
	ffo.IgnoreFields = make(map[string]bool)
	ffo.IgnoreTypes = make(map[string]bool)
	ffo.IgnoreTypes["password"] = true
	ffo.Target = ""

	for key, val := range options {
		switch key {
		case "IgnoreFields":
			if valArray, ok := val.([]string); ok {
				for _, val := range valArray {
					ffo.IgnoreFields[val] = true
				}
			}
		case "IgnoreTypes":
			if valArray, ok := val.([]string); ok {
				for _, val := range valArray {
					ffo.IgnoreTypes[val] = true
				}
			}
		case "FillPassword":
			if valBool, ok := val.(bool); ok {
				ffo.IgnoreTypes["password"] = !valBool
			}
		case "Target":
			if valStr, ok := val.(string); ok {
				ffo.Target = valStr
			}
		}
	}

	return &ffo
}

func newFiller(data map[string]interface{}, options map[string]interface{}) *Filler {
	ffo := setOptions(options)
	return &Filler{Data: data, Params: make(map[string][]byte), FillInFormOptions: *ffo}
}

// return writer implement interface io.Writer.
func FillWriter(wr io.Writer, data map[string]interface{}, options map[string]interface{}) io.Writer {
	filler := newFiller(data, options)
	return Writer{filler: filler, wr: wr}
}
func (w Writer) Write(p []byte) (int, error) {
	filled := w.filler.fill(p)
	return w.wr.Write(filled)
}

// return filled formed html.
func Fill(body []byte, data map[string]interface{}, options map[string]interface{}) []byte {
	filler := newFiller(data, options)
	return filler.fill(body)
}

func (f Filler) fill(body []byte) []byte {
	return f.compiledRegexp("form").ReplaceAllFunc(body, f.fillForm)
}

func (f Filler) fillForm(formbody []byte) []byte {
	// process only form with target id
	if f.FillInFormOptions.Target != "" {
		formTag := f.compiledRegexp("start form").FindSubmatch(formbody)
		if cap(formTag) == 2 {
			if id := f.getId(formTag[1]); !bytes.Equal(id, []byte{}) {
				if string(id) != f.FillInFormOptions.Target {
					return formbody
				}
			}
		}
	}

	replaced := f.compiledRegexp("input").ReplaceAllFunc(formbody, f.fillInput)
	replaced = f.compiledRegexp("select").ReplaceAllFunc(replaced, f.fillSelect)
	replaced = f.compiledRegexp("textarea").ReplaceAllFunc(replaced, f.fillTextarea)

	return replaced
}

func (f Filler) unquote(tag []byte) []byte {
	return bytes.Trim(tag, `'"`)
}

func (f Filler) getId(tag []byte) []byte {
	id := f.compiledRegexp("id").FindSubmatch(tag)
	if cap(id) == 2 {
		return f.unquote(id[1])
	}
	return []byte{}
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
	return bytes.Replace(bytes.Replace(bytes.Replace(bytes.Replace(tag, []byte{'&'}, BAAmp, -1), []byte{'<'}, BALt, -1), []byte{'>'}, BAGt, -1), []byte{'"'}, BAQuot, -1)
}

func (f Filler) getParam(name []byte) ([]byte, bool) {
	// ignore
	nameStr := string(name)
	if _, ok := f.IgnoreFields[nameStr]; ok {
		return []byte{}, false
	}
	// like cache
	if param, ok := f.Params[nameStr]; ok {
		return param, true
	}
	if param, ok := f.Data[nameStr]; ok {
		if casted, ok := param.(string); ok {
			f.Params[nameStr] = []byte(casted)
			return f.Params[nameStr], true
		}
		if casted, ok := param.(int); ok {
			f.Params[nameStr] = []byte(fmt.Sprintf("%d", casted))
			return f.Params[nameStr], true
		}
	}

	return []byte{}, false
}

func (f Filler) fillInput(tag []byte) []byte {
	inputType := f.getType(tag)
	if bytes.Equal(inputType, []byte{}) {
		inputType = []byte("text")
	}

	// ignore types (password is default true (not fillin))
	if flg, ok := f.IgnoreTypes[string(inputType)]; ok && flg {
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
				tag = f.compiledRegexp("space+>").ReplaceAll(tag, BAChecked)
			}
		} else {
			tag = f.compiledRegexp("space+checked").ReplaceAll(tag, BABlank)
		}
	} else { // text
		escapedValue := f.escapeHTML(paramValue)
		reg := f.compiledRegexp("value(nocapture)")
		if reg.Match(tag) {
			tag = reg.ReplaceAll(tag, append([]byte(`value="`), append(escapedValue, BADQuot...)...))
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
	tag = f.compiledRegexp("textarea(3capture)").ReplaceAll(tag, append([]byte(`${1}`), append(f.escapeHTML(paramValue), []byte(`${3}`)...)...))

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

	return f.compiledRegexp("option(nocapture)").ReplaceAllFunc(tag,
		func(tag []byte) []byte {
			return f.fillOption(tag, paramValue)
		})
}

func (f Filler) fillOption(tag, paramValue []byte) []byte {
	value := f.getValue(tag)
	if bytes.Equal(value, []byte{}) {
		value = f.compiledRegexp("option").ReplaceAll(tag, []byte(`$1`))
	}

	if bytes.Equal(paramValue, value) {
		if !f.compiledRegexp("selected").Match(tag) {
			tag = f.compiledRegexp("start option").ReplaceAllFunc(tag,
				func(tag []byte) []byte {
					return f.compiledRegexp("tag end").ReplaceAll(tag, BASelected)
				})
		}
	} else {
		tag = f.compiledRegexp("space+selected").ReplaceAll(tag, BABlank)
	}

	return tag
}
