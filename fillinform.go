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

	spaceRxp     = `\s`
	attrNameRxp  = `[\w\-]+`
	attrValueRxp = `(?:"[^"]*"|'[^']*'|[^'"/>\s]+|[\w\-]+)`
	attrRxp      = `(?:` + spaceRxp + `+(?:` + attrNameRxp + `(?:=` + attrValueRxp + `)?))`

	startFormRxp     = `(?:<` + _Form + attrRxp + `+` + spaceRxp + `*>)`
	inputRxp         = `(?:<` + _Input + attrRxp + `+` + spaceRxp + `*/?>)`
	startSelectRxp   = `(?:<` + _Select + attrRxp + `+` + spaceRxp + `*>)`
	startOptionRxp   = `(?:<` + _Option + attrRxp + `*` + spaceRxp + `*>)`
	startTextareaRxp = `(?:<` + _Textarea + attrRxp + `+` + spaceRxp + `*>)`

	endFormRxp     = `(?:</` + _Form + `>)`
	endSelectRxp   = `(?:</` + _Select + `>)`
	endOptionRxp   = `(?:</` + _Option + `>)`
	endTextareaRxp = `(?:</` + _Textarea + `>)`

	checkedRxp  = `(?:` + _Checked + `(?:=(?:"` + _Checked + `"|'` + _Checked + `'|` + _Checked + `))?)`
	selectedRxp = `(?:` + _Selected + `(?:=(?:"` + _Selected + `"|'` + _Selected + `'|` + _Selected + `))?)`
	multipleRxp = `(?:` + _Multiple + `(?:=(?:"` + _Multiple + `"|'` + _Multiple + `'|` + _Multiple + `))?)`
)

var (
	checkboxBytes   = []byte(`checkbox`)
	radioBytes      = []byte(`radio`)
	textBytes       = []byte(`text`)
	ampBytes        = []byte(`&amp;`)
	ltBytes         = []byte(`&lt;`)
	gtBytes         = []byte(`&gt;`)
	quotBytes       = []byte(`&quot;`)
	checkedBytes    = []byte(` checked="checked"$1>`)
	selectedBytes   = []byte(` selected="selected">`)
	blankBytes      = []byte(``)
	doubleQuotBytes = []byte(`"`)
)

var CompiledRegexpMap map[string]*regexp.Regexp

func init() {
	createRegexpMap()
}

func compileMultiLine(regstr string) *regexp.Regexp {
	return regexp.MustCompile(`(?msi:` + regstr + `)`)
}

func createRegexpMap() {
	CompiledRegexpMap = make(map[string]*regexp.Regexp)
	CompiledRegexpMap["form"] = compileMultiLine(startFormRxp + `.*?` + endFormRxp)
	CompiledRegexpMap["start form"] = compileMultiLine(`(` + startFormRxp + `)`)

	CompiledRegexpMap["input"] = compileMultiLine(inputRxp)
	CompiledRegexpMap["select"] = compileMultiLine(startSelectRxp + `.*?` + endSelectRxp)
	CompiledRegexpMap["textarea"] = compileMultiLine(startTextareaRxp + `.*?` + endTextareaRxp)

	CompiledRegexpMap["id"] = compileMultiLine(_Id + `=(` + attrValueRxp + `)`)
	CompiledRegexpMap["type"] = compileMultiLine(_Type + `=(` + attrValueRxp + `)`)
	CompiledRegexpMap["value"] = compileMultiLine(_Value + `=(` + attrValueRxp + `)`)
	CompiledRegexpMap["name"] = compileMultiLine(_Name + `=(` + attrValueRxp + `)`)

	CompiledRegexpMap["checked"] = compileMultiLine(checkedRxp)
	CompiledRegexpMap["space+>"] = compileMultiLine(spaceRxp + `*(/?)>\z`)
	CompiledRegexpMap["space+checked"] = compileMultiLine(spaceRxp + checkedRxp)
	CompiledRegexpMap["value(nocapture)"] = compileMultiLine(_Value + `=` + attrValueRxp)
	CompiledRegexpMap["textarea(3capture)"] = compileMultiLine(`(` + startTextareaRxp + `)(.*?)(` + endTextareaRxp + `)`)
	CompiledRegexpMap["multiple"] = compileMultiLine(multipleRxp)
	CompiledRegexpMap["option"] = compileMultiLine(startOptionRxp + `(.*?)` + endOptionRxp)
	CompiledRegexpMap["option(nocapture)"] = compileMultiLine(startOptionRxp + `.*?` + endOptionRxp)
	CompiledRegexpMap["selected"] = compileMultiLine(selectedRxp)
	CompiledRegexpMap["start option"] = compileMultiLine(startOptionRxp)
	CompiledRegexpMap["tag end"] = compileMultiLine(spaceRxp + `*>\z`)
	CompiledRegexpMap["space+selected"] = compileMultiLine(spaceRxp + selectedRxp)
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
	FillPassword bool
	Target       string
}

type Filler struct {
	FillInFormOptions
	Params map[string][][]byte
	Data   map[string][]string
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
	ffo.IgnoreTypes["submit"] = true
	ffo.IgnoreTypes["image"] = true
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

func newFiller(data map[string][]string, options map[string]interface{}) *Filler {
	ffo := setOptions(options)
	return &Filler{Data: data, Params: make(map[string][][]byte), FillInFormOptions: *ffo}
}

// return writer implement interface io.Writer.
func FillWriter(wr io.Writer, data map[string][]string, options map[string]interface{}) io.Writer {
	filler := newFiller(data, options)
	return Writer{filler: filler, wr: wr}
}
func (w Writer) Write(p []byte) (int, error) {
	filled := w.filler.fill(p)
	return w.wr.Write(filled)
}

// return filled formed html.
func Fill(body []byte, data map[string][]string, options map[string]interface{}) []byte {
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
		if len(formTag) == 2 {
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
	if len(id) == 2 {
		return f.unquote(id[1])
	}
	return []byte{}
}

func (f Filler) getType(tag []byte) []byte {
	itype := f.compiledRegexp("type").FindSubmatch(tag)
	if len(itype) == 2 {
		return f.unquote(itype[1])
	}
	return []byte{}
}

func (f Filler) getValue(tag []byte) []byte {
	value := f.compiledRegexp("value").FindSubmatch(tag)
	if len(value) == 2 {
		return f.unquote(value[1])
	}
	return []byte{}
}

func (f Filler) getName(tag []byte) []byte {
	name := f.compiledRegexp("name").FindSubmatch(tag)
	if len(name) == 2 {
		return f.unquote(name[1])
	}
	return []byte{}
}

func (f Filler) escapeHTML(tag []byte) []byte {
	return bytes.Replace(bytes.Replace(bytes.Replace(bytes.Replace(tag, []byte{'&'}, ampBytes, -1), []byte{'<'}, ltBytes, -1), []byte{'>'}, gtBytes, -1), []byte{'"'}, quotBytes, -1)
}

func (f Filler) getParam(name string) ([][]byte, bool) {
	// like cache
	if param, ok := f.Params[name]; ok {
		return param, true
	}
	if param, ok := f.Data[name]; ok {
		vals := make([][]byte, len(param))
		for i, val := range param {
			vals[i] = []byte(val)
		}
		f.Params[name] = vals
		return f.Params[name], true
	}

	return [][]byte{}, false
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

	name := string(f.getName(tag))
	if _, ok := f.IgnoreFields[name]; ok {
		return tag
	}
	paramValues, exists := f.getParam(name)

	if bytes.Equal(inputType, checkboxBytes) || bytes.Equal(inputType, radioBytes) {
		value := f.getValue(tag)

		tag = f.compiledRegexp("space+checked").ReplaceAll(tag, blankBytes)
		for _, paramValue := range paramValues {
			if bytes.Equal(paramValue, value) {
				if !f.compiledRegexp("checked").Match(tag) {
					tag = f.compiledRegexp("space+>").ReplaceAll(tag, checkedBytes)
				}
			}
		}
	} else { // if bytes.Equal(inputType, textBytes)
		var paramValue []byte
		if !exists {
			paramValue = []byte("")
		} else {
			paramValue = paramValues[0]
		}
		escapedValue := f.escapeHTML(paramValue)
		reg := f.compiledRegexp("value(nocapture)")
		if reg.Match(tag) {
			tag = reg.ReplaceAll(tag, append([]byte(`value="`), append(escapedValue, doubleQuotBytes...)...))
		} else {
			tag = f.compiledRegexp("space+>").ReplaceAll(tag, append([]byte(` value="`), append(escapedValue, []byte(`"$1>`)...)...))
		}
	}

	return tag
}

func (f Filler) fillTextarea(tag []byte) []byte {
	name := string(f.getName(tag))
	if _, ok := f.IgnoreFields[name]; ok {
		return tag
	}
	paramValues, exists := f.getParam(name)
	var paramValue []byte
	if !exists {
		paramValue = []byte("")
	} else {
		paramValue = paramValues[0]
	}
	tag = f.compiledRegexp("textarea(3capture)").ReplaceAll(tag, append([]byte(`${1}`), append(f.escapeHTML(paramValue), []byte(`${3}`)...)...))

	return tag
}

func (f Filler) fillSelect(tag []byte) []byte {
	name := string(f.getName(tag))
	if _, ok := f.IgnoreFields[name]; ok {
		return tag
	}
	paramValues, exists := f.getParam(name)

	if exists {
		if !f.compiledRegexp("multiple").Match(tag) {
			paramValues = paramValues[:1]
		}
	}

	return f.compiledRegexp("option(nocapture)").ReplaceAllFunc(tag,
		func(tag []byte) []byte {
			return f.fillOption(tag, paramValues)
		})
}

func (f Filler) fillOption(tag []byte, paramValues [][]byte) []byte {
	value := f.getValue(tag)
	if bytes.Equal(value, []byte{}) {
		value = f.compiledRegexp("option").ReplaceAll(tag, []byte(`$1`))
	}

	tag = f.compiledRegexp("space+selected").ReplaceAll(tag, blankBytes)
	for _, paramValue := range paramValues {
		if bytes.Equal(paramValue, value) {
			if !f.compiledRegexp("selected").Match(tag) {
				tag = f.compiledRegexp("start option").ReplaceAllFunc(tag,
					func(tag []byte) []byte {
						return f.compiledRegexp("tag end").ReplaceAll(tag, selectedBytes)
					})
			}
		}
	}
	return tag
}
