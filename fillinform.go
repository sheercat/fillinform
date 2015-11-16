package fillinform

import (
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

type FillinFormOptions struct {
	FillPassword bool
	IgnoreFields map[string]bool
	IgnoreTypes  map[string]bool
	Target       string
	Data         map[string]interface{}
}

type Filler struct {
	FillinFormOptions
}

type Writer struct {
	filler *Filler
	wr     io.Writer
}

func (f Filler) compileMultiLine(regstr string) *regexp.Regexp {
	return regexp.MustCompile(`(?ms:` + regstr + `)`)
}

func FillWriter(wr io.Writer, data map[string]interface{}, options map[string]interface{}) io.Writer {
	filler := &Filler{FillinFormOptions{Data: data}}
	return Writer{filler: filler, wr: wr}
}
func (w Writer) Write(p []byte) (int, error) {
	filled := w.filler.fill(p)
	return w.wr.Write(filled)
}

func Fill(body []byte, data map[string]interface{}, options map[string]interface{}) []byte {
	filler := &Filler{FillinFormOptions{Data: data}}

	return filler.fill(body)
}

func (f Filler) fill(body []byte) []byte {
	return f.compileMultiLine(FORM+`(.*?)`+EndFORM).ReplaceAllFunc(body, f.fillForm)
}

func (f Filler) fillForm(formbody []byte) []byte {
	replaced := f.compileMultiLine(INPUT).ReplaceAllFunc(formbody, f.fillInput)

	replaced = f.compileMultiLine(SELECT+`(.*?)`+EndSELECT).ReplaceAllFunc(replaced, f.fillSelect)

	replaced = f.compileMultiLine(TEXTAREA+`(.*?)`+EndTEXTAREA).ReplaceAllFunc(replaced, f.fillTextarea)

	return replaced
}

func (f Filler) unquote(tag []byte) []byte {
	newTag := f.compileMultiLine(`['"](.*)['"]`).FindSubmatch(tag)
	if cap(newTag) == 2 {
		return newTag[1]
	}
	return tag
}

func (f Filler) getType(tag []byte) string {
	itype := f.compileMultiLine(Type + `=(` + ATTR_VALUE + `)`).FindSubmatch(tag)
	if cap(itype) == 2 {
		return string(f.unquote(itype[1]))
	}
	return string(tag)
}

func (f Filler) getValue(tag []byte) string {
	value := f.compileMultiLine(Value + `=(` + ATTR_VALUE + `)`).FindSubmatch(tag)
	if cap(value) == 2 {
		return string(f.unquote(value[1]))
	}
	return ""
}

func (f Filler) getName(tag []byte) string {
	name := f.compileMultiLine(Name + `=(` + ATTR_VALUE + `)`).FindSubmatch(tag)
	if cap(name) == 2 {
		return string(f.unquote(name[1]))
	}
	return string(tag)
}

func (f Filler) escapeHTML(tag string) string {
	tag = regexp.MustCompile(`&`).ReplaceAllString(tag, `&amp;`)
	tag = regexp.MustCompile(`<`).ReplaceAllString(tag, `&lt;`)
	tag = regexp.MustCompile(`>`).ReplaceAllString(tag, `&gt;`)
	tag = regexp.MustCompile(`"`).ReplaceAllString(tag, `&quot;`)
	return tag
}

func (f Filler) getParam(name string) (string, bool) {
	// ignore
	if _, ok := f.IgnoreFields[name]; ok {
		return "", false
	}
	if param, ok := f.Data[name]; ok {
		if casted, ok := param.(string); ok {
			return casted, true
		}
	}

	return "", false
}

func (f Filler) fillInput(tag []byte) []byte {
	inputType := f.getType(tag)
	if inputType == "" {
		inputType = "text"
	}

	// ignore
	if _, ok := f.IgnoreTypes[inputType]; ok {
		return tag
	}

	paramValue, exists := f.getParam(f.getName(tag))
	if !exists {
		return tag
	}

	if inputType == "checkbox" || inputType == "radio" {
		value := f.getValue(tag)

		if paramValue == value {
			if !f.compileMultiLine(CHECKED).Match(tag) {
				tag = f.compileMultiLine(SPACE+`*(/?)>\z`).ReplaceAll(tag, []byte(` checked="checked"$1>`))
			}
		} else {
			tag = f.compileMultiLine(SPACE+CHECKED).ReplaceAll(tag, []byte(``))
		}
	} else { // text
		escapedValue := f.escapeHTML(paramValue)
		reg := f.compileMultiLine(Value + `=` + ATTR_VALUE)
		if reg.Match(tag) {
			tag = reg.ReplaceAll(tag, []byte(`value="`+escapedValue+`"`))
		} else {
			tag = f.compileMultiLine(SPACE+`*(/?)>\z`).ReplaceAll(tag, []byte(` value="`+escapedValue+`"$1>`))
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
	replaced := `${1}` + escapedValue + `${3}`
	tag = f.compileMultiLine(`(`+TEXTAREA+`)(.*?)(`+EndTEXTAREA+`)`).ReplaceAll(tag, []byte(replaced))
	// matched := f.compileMultiLine(`(` + TEXTAREA + `).*?(` + EndTEXTAREA + `)`).FindSubmatch(tag)
	// if cap(matched) == 3 {
	// 	tag = []byte(string(matched[1]) + escapedValue + string(matched[2]))
	// }

	return tag
}

func (f Filler) fillSelect(tag []byte) []byte {
	paramValue, exists := f.getParam(f.getName(tag))
	if !exists {
		return tag
	}

	if f.compileMultiLine(MULTIPLE).Match(tag) {
		return tag
	}

	return f.compileMultiLine(OPTION+`.*?`+EndOPTION).ReplaceAllFunc(tag, func(tag []byte) []byte { return f.fillOption(tag, paramValue) })
}

func (f Filler) fillOption(tag []byte, paramValue string) []byte {
	value := f.getValue(tag)
	if value == "" {
		value = string(f.compileMultiLine(OPTION+`(.*?)`+EndOPTION).ReplaceAll(tag, []byte(`$1`)))
	}

	if paramValue == value {
		if !f.compileMultiLine(SELECTED).Match(tag) {
			tag = f.compileMultiLine(OPTION).ReplaceAllFunc(tag, func(tag []byte) []byte {
				return f.compileMultiLine(SPACE+`*>\z`).ReplaceAll(tag, []byte(` selected="selected">`))
			})
		}
	} else {
		tag = f.compileMultiLine(SPACE+SELECTED).ReplaceAll(tag, []byte(``))
	}

	return tag
}
