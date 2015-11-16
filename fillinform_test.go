package fillinform

import (
	"bytes"
	"testing"
)

func TestUnquote(t *testing.T) {
	filler := &Filler{}

	hoge := filler.unquote([]byte(`"hoge"`))
	if string(hoge) != "hoge" {
		t.Errorf("double unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`'hoge'`))
	if string(hoge) != "hoge" {
		t.Errorf("single unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`'hoge"`))
	if string(hoge) != "hoge" {
		t.Errorf("single double unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`"hoge'`))
	if string(hoge) != "hoge" {
		t.Errorf("double single unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`hoge`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`"'hoge'"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`"'"hoge"'"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`'''hoge'''`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`''"hoge''"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error $v", hoge)
	}
	hoge = filler.unquote([]byte(`"''"`))
	if string(hoge) != "" {
		t.Errorf("no unquote error $v", hoge)
	}
}

func BenchmarkUnquote(b *testing.B) {
	str := `"hoge"`
	bstr := []byte(str)
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.unquote(bstr)
	}
}

var TesteeArray = map[string]string{
	`double_quote`: `<input type="hoge" name="hoge" value="hoge">`,
	`single_quote`: `<input type='hoge' name='hoge' value='hoge'>`,
	`no_quote`:     `<input type=hoge   name=hoge   value=hoge>`,
	`uppercase`:    `<input TYPE="hoge" NAME="hoge" VALUE="hoge">`,
	`capcase`:      `<input Type="hoge" Name="hoge" Value="hoge">`,
}

func TestGetType(t *testing.T) {
	filler := &Filler{}
	for key, val := range TesteeArray {
		hoge := filler.getType([]byte(val))
		if string(hoge) != "hoge" {
			t.Errorf("error in ", key)
		}
	}
}

func BenchmarkGetType(b *testing.B) {
	str := `<input type="hoge" value="hoge" name="hoge">`
	bstr := []byte(str)
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.getType(bstr)
	}
}

func TestGetValue(t *testing.T) {
	filler := &Filler{}
	for key, val := range TesteeArray {
		hoge := filler.getValue([]byte(val))
		if string(hoge) != "hoge" {
			t.Errorf("error in ", key)
		}
	}
}

func BenchmarkGetValue(b *testing.B) {
	str := `<input value="hoge" type="hoge" name="hoge">`
	bstr := []byte(str)
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.getValue(bstr)
	}
}

func TestGetName(t *testing.T) {
	filler := &Filler{}
	for key, val := range TesteeArray {
		hoge := filler.getName([]byte(val))
		if string(hoge) != "hoge" {
			t.Errorf("error in ", key)
		}
	}
}

func BenchmarkGetName(b *testing.B) {
	str := `<input name="hoge" type="hoge" value="hoge">`
	bstr := []byte(str)
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.getName(bstr)
	}
}

var TesteeArrayESCHTML = map[string]map[string]string{
	`and`:   {"arg": `Hi, & is ampasand`, "return": `Hi, &amp; is ampasand`},
	`lt`:    {"arg": `Hi, < is less than`, "return": `Hi, &lt; is less than`},
	`gt`:    {"arg": `Hi, > is greater than`, "return": `Hi, &gt; is greater than`},
	`quot`:  {"arg": `Hi, " is quotation`, "return": `Hi, &quot; is quotation`},
	`mixed`: {"arg": `html special char is <, >, &, and """.`, "return": `html special char is &lt;, &gt;, &amp;, and &quot;&quot;&quot;.`},
}

func TestEscapeHTML(t *testing.T) {
	filler := &Filler{}
	for key, mapval := range TesteeArrayESCHTML {
		val := []byte(mapval["arg"])
		res := []byte(mapval["return"])
		hoge := filler.escapeHTML(val)
		if !bytes.Equal(hoge, res) {
			t.Errorf("error in ", key, string(res))
		}
	}
}

func BenchmarkEscapeHTML(b *testing.B) {
	str := `<input & type="hoge" & value="hoge" & name="hoge">`
	ba := []byte(str)
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.escapeHTML(ba)
	}
}

func TestFillInput(t *testing.T) {
	formData := map[string]interface{}{
		"title": `hoge & Hoge <"Title">`,
		"chk":   "chkval",
		"rdo":   "rdoval2",
	}
	filler := newFiller(formData)

	htmlstr := filler.fillInput([]byte(`<input type="text" name="title"/>`))
	if string(htmlstr) != `<input type="text" name="title" value="hoge &amp; Hoge &lt;&quot;Title&quot;&gt;"/>` {
		t.Errorf("fillInput error: ", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="checkbox" name="chk" value="chkval" checked=checked/>`))
	if string(htmlstr) != `<input type="checkbox" name="chk" value="chkval" checked=checked/>` {
		t.Errorf("no affect error: ", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval1" checked=checked/>`))
	if string(htmlstr) != `<input type="radio" name="rdo" value="rdoval1"/>` {
		t.Errorf("fillout error: ", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval2" />`))
	if string(htmlstr) != `<input type="radio" name="rdo" value="rdoval2" checked="checked"/>` {
		t.Errorf("fillin error: ", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="submit" value="Send">`))
	if string(htmlstr) != `<input type="submit" value="Send">` {
		t.Errorf("no fill error: ", string(htmlstr))
	}
}

func BenchmarkFillInput(b *testing.B) {
	formData := map[string]interface{}{
		"title": `hoge & Hoge <"Title">`,
		"chk":   "chkval",
		"rdo":   "rdoval2",
	}
	filler := newFiller(formData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval2" />`))
	}
}

func TestFillTextarea(t *testing.T) {
	formData := map[string]interface{}{
		"body": "hoge & hoge <hoge@hogehoge>",
	}
	filler := newFiller(formData)

	htmlstr := filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge"></textarea>`))
	if string(htmlstr) != `<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hoge &amp; hoge &lt;hoge@hogehoge&gt;</textarea>` {
		t.Errorf("fillTextarea error: ", string(htmlstr))
	}
	htmlstr = filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>`))
	if string(htmlstr) != `<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hoge &amp; hoge &lt;hoge@hogehoge&gt;</textarea>` {
		t.Errorf("fillTextarea error: ", string(htmlstr))
	}
	htmlstr = filler.fillTextarea([]byte(`<textarea id="body" name="bodyX" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>`))
	if string(htmlstr) != `<textarea id="body" name="bodyX" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>` {
		t.Errorf("no affect error: ", string(htmlstr))
	}
}

func BenchmarkFillTextarea(b *testing.B) {
	formData := map[string]interface{}{
		"body": "hoge & hoge <hoge@hogehoge>",
	}
	filler := newFiller(formData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge"></textarea>`))
	}
}

func TestFillSelect(t *testing.T) {
	formData := map[string]interface{}{
		"select": "1",
	}
	filler := newFiller(formData)

	htmlstr := filler.fillSelect([]byte(`<select name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>`))

	if string(htmlstr) != `<select name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>` {
		t.Errorf("fillTextarea error: ", string(htmlstr))
	}

	htmlstr = filler.fillSelect([]byte(`<select name="select">
    <option value="1">1</option>
    <option value="2" selected="selected">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>`))
	if string(htmlstr) != `<select name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>` {
		t.Errorf("fillTextarea error: ", string(htmlstr))
	}

	htmlstr = filler.fillSelect([]byte(`<select name="selectX">
    <option value="1">1</option>
    <option value="2" selected="selected">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>`))
	if string(htmlstr) != `<select name="selectX">
    <option value="1">1</option>
    <option value="2" selected="selected">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>` {
		t.Errorf("no affect error: ", string(htmlstr))
	}
}

func BenchmarkFillSelect(b *testing.B) {
	formData := map[string]interface{}{
		"select": "1",
	}
	filler := newFiller(formData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillSelect([]byte(`<select name="select">
    <option value="1">1</option>
    <option value="2" selected="selected">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>`))
	}
}

func TestFillOption(t *testing.T) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "chkval",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}
	filler := newFiller(formData)

	htmlstr := filler.fillOption([]byte(`<option value="1">1</option>`), []byte(`1`))
	if string(htmlstr) != `<option value="1" selected="selected">1</option>` {
		t.Errorf("fillOption error: ", string(htmlstr))
	}

	htmlstr = filler.fillOption([]byte(`<option value="1">1</option>`), []byte(`2`))
	if string(htmlstr) != `<option value="1">1</option>` {
		t.Errorf("fillOption error: ", string(htmlstr))
	}
	htmlstr = filler.fillOption([]byte(`<option>1</option>`), []byte(`1`))
	if string(htmlstr) != `<option selected="selected">1</option>` {
		t.Errorf("fillOption error: ", string(htmlstr))
	}

}

func BenchmarkFillOption(b *testing.B) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "chkval",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}
	filler := newFiller(formData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillOption([]byte(`<option value="1">1</option>`), []byte(`1`))
	}
}

var HTML = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="text" name="title"/>
  <input type="checkbox" name="chk" value="chkval" checked=checked/>
  <input type="radio" name="rdo" value="rdoval1" checked=checked/>
  <input type="radio" name="rdo" value="rdoval2" />
  <select name="select">
    <option value="1">1</option>
    <option value="2" selected=selected>2</option>
  </select>
  <textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>
  <input type="submit" value="Send">
</form>
</body></html>
`

var HTMLSuccess = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="text" name="title" value="hogeTitle"/>
  <input type="checkbox" name="chk" value="chkval"/>
  <input type="radio" name="rdo" value="rdoval1"/>
  <input type="radio" name="rdo" value="rdoval2" checked="checked"/>
  <select name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
  </select>
  <textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hogehoge</textarea>
  <input type="submit" value="Send">
</form>
</body></html>
`

func TestFillinForm(t *testing.T) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "1",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}

	filler := newFiller(formData)

	htmlstr := filler.fill([]byte(HTML))

	if string(htmlstr) != HTMLSuccess {
		t.Errorf("fillinform error: ", string(htmlstr))
	}
}

func BenchmarkFillinForm(b *testing.B) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "1",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}

	filler := newFiller(formData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTML))
	}
}
