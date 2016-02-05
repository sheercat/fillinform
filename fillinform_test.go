package fillinform

import (
	"bytes"
	"regexp"
	"testing"
)

func TestOnepass(t *testing.T) {
	re := regexp.MustCompile("x?")
	re.MatchString("y")
}

func TestUnquote(t *testing.T) {
	filler := &Filler{}

	hoge := filler.unquote([]byte(`"hoge"`))
	if string(hoge) != "hoge" {
		t.Errorf("double unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`'hoge'`))
	if string(hoge) != "hoge" {
		t.Errorf("single unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`'hoge"`))
	if string(hoge) != "hoge" {
		t.Errorf("single double unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`"hoge'`))
	if string(hoge) != "hoge" {
		t.Errorf("double single unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`hoge`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`"'hoge'"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`"'"hoge"'"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`'''hoge'''`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`''"hoge''"`))
	if string(hoge) != "hoge" {
		t.Errorf("no unquote error %v", hoge)
	}
	hoge = filler.unquote([]byte(`"''"`))
	if string(hoge) != "" {
		t.Errorf("no unquote error %v", hoge)
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
			t.Errorf("error in %v", key)
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
			t.Errorf("error in %v", key)
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
			t.Errorf("error in %v", key)
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
			t.Errorf("error in %v %v", key, string(res))
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
	formData := map[string][]string{
		"title": []string{`hoge & Hoge <"Title">`},
		"chk":   []string{"chkval"},
		"rdo":   []string{"rdoval2"},
	}
	filler := newFiller(formData, nil)

	htmlstr := filler.fillInput([]byte(`<input type="text" name="title"/>`))
	if string(htmlstr) != `<input type="text" name="title" value="hoge &amp; Hoge &lt;&quot;Title&quot;&gt;"/>` {
		t.Errorf("fillInput error: %v", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="checkbox" name="chk" value="chkval" checked=checked/>`))
	if string(htmlstr) != `<input type="checkbox" name="chk" value="chkval" checked="checked"/>` {
		t.Errorf("no affect error: %v", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval1" checked=checked/>`))
	if string(htmlstr) != `<input type="radio" name="rdo" value="rdoval1"/>` {
		t.Errorf("fillout error: %v", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval2" />`))
	if string(htmlstr) != `<input type="radio" name="rdo" value="rdoval2" checked="checked"/>` {
		t.Errorf("fillin error: %v", string(htmlstr))
	}

	htmlstr = filler.fillInput([]byte(`<input type="submit" value="Send">`))
	if string(htmlstr) != `<input type="submit" value="Send">` {
		t.Errorf("no fill error: %v", string(htmlstr))
	}
}

func BenchmarkFillInput(b *testing.B) {
	formData := map[string][]string{
		"title": []string{`hoge & Hoge <"Title">`},
		"chk":   []string{"chkval"},
		"rdo":   []string{"rdoval2"},
	}
	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval2" />`))
	}
}

func TestFillTextarea(t *testing.T) {
	formData := map[string][]string{
		"body": []string{"hoge & hoge <hoge@hogehoge>"},
	}
	filler := newFiller(formData, nil)

	htmlstr := filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge"></textarea>`))
	if string(htmlstr) != `<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hoge &amp; hoge &lt;hoge@hogehoge&gt;</textarea>` {
		t.Errorf("fillTextarea error: %v", string(htmlstr))
	}
	htmlstr = filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>`))
	if string(htmlstr) != `<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hoge &amp; hoge &lt;hoge@hogehoge&gt;</textarea>` {
		t.Errorf("fillTextarea error: %v", string(htmlstr))
	}
	htmlstr = filler.fillTextarea([]byte(`<textarea id="body" name="bodyX" cols="80" rows="20" placeholder="hoge">gakuburu</textarea>`))
	if string(htmlstr) != `<textarea id="body" name="bodyX" cols="80" rows="20" placeholder="hoge"></textarea>` {
		t.Errorf("no affect error: %v", string(htmlstr))
	}
}

func BenchmarkFillTextarea(b *testing.B) {
	formData := map[string][]string{
		"body": []string{"hoge & hoge <hoge@hogehoge>"},
	}
	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge"></textarea>`))
	}
}

func TestFillSelect(t *testing.T) {
	formData := map[string][]string{
		"select": []string{"1"},
	}
	filler := newFiller(formData, nil)

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
		t.Errorf("fillTextarea error: %v", string(htmlstr))
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
		t.Errorf("fillTextarea error: %v", string(htmlstr))
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
    <option value="2">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>` {
		t.Errorf("no affect error: %v", string(htmlstr))
	}

	formData2 := map[string][]string{
		"select": []string{"1", "3"},
	}
	filler2 := newFiller(formData2, nil)
	htmlstr = filler2.fillSelect([]byte(`<select multiple="multiple" name="select">
    <option value="1">1</option>
    <option value="2" selected="selected">2</option>
    <option value="3">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>`))
	if string(htmlstr) != `<select multiple="multiple" name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
    <option value="3" selected="selected">3</option>
    <option value="4">4</option>
    <option value="5">5</option>
    <option value="6">6</option>
  </select>` {
		t.Errorf("multiple error: %v", string(htmlstr))
	}

	htmlstr = filler2.fillSelect([]byte(`<select name="select">
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
		t.Errorf("no multiple error: %v", string(htmlstr))
	}
}

func BenchmarkFillSelect(b *testing.B) {
	formData := map[string][]string{
		"select": []string{"1"},
	}
	filler := newFiller(formData, nil)

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
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"chkval"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
	}
	filler := newFiller(formData, nil)

	htmlstr := filler.fillOption([]byte(`<option value="1">1</option>`), [][]byte{[]byte(`1`)})
	if string(htmlstr) != `<option value="1" selected="selected">1</option>` {
		t.Errorf("fillOption error: %v", string(htmlstr))
	}

	htmlstr = filler.fillOption([]byte(`<option value="1">1</option>`), [][]byte{[]byte(`2`)})
	if string(htmlstr) != `<option value="1">1</option>` {
		t.Errorf("fillOption error: %v", string(htmlstr))
	}
	htmlstr = filler.fillOption([]byte(`<option>1</option>`), [][]byte{[]byte(`1`)})
	if string(htmlstr) != `<option selected="selected">1</option>` {
		t.Errorf("fillOption error: %v", string(htmlstr))
	}

}

func BenchmarkFillOption(b *testing.B) {
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"chkval"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
	}
	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillOption([]byte(`<option value="1">1</option>`), [][]byte{[]byte(`1`)})
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
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"1"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
	}

	filler := newFiller(formData, nil)

	htmlstr := filler.fill([]byte(HTML))

	if string(htmlstr) != HTMLSuccess {
		t.Errorf("fillinform error: %v", string(htmlstr))
	}
}

func BenchmarkFillinForm(b *testing.B) {
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"1"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
	}

	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTML))
	}
}

func TestFillinForm2(t *testing.T) {
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"1"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
	}

	filler := newFiller(formData, nil)

	htmlstr := filler.fill([]byte(HTML))

	if string(htmlstr) != HTMLSuccess {
		t.Errorf("fillinform error: %v", string(htmlstr))
	}
}

var HTMLMulti = `
<html><head><title>title of test</title></head><body>
<form id="myform" action="./" method="POST">
  <input type="hidden" name="hidden"/>
  <input type="password" name="pass"/>
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

<form id="myform2" action="./" method="POST">
  <input type="hidden" name="hidden"/>
  <input type="password" name="pass"/>
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

var HTMLMultiSuccess = `
<html><head><title>title of test</title></head><body>
<form id="myform" action="./" method="POST">
  <input type="hidden" name="hidden"/>
  <input type="password" name="pass"/>
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

<form id="myform2" action="./" method="POST">
  <input type="hidden" name="hidden" value=""/>
  <input type="password" name="pass"/>
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

var HTMLPassword = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="password" name="pass"/>
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

var HTMLPasswordSuccess = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="password" name="pass" value="hogepass"/>
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

var HTMLFields = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="password" name="pass"/>
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

var HTMLFieldsSuccess = `
<html><head><title>title of test</title></head><body>
<form name="myform" action="./" method="POST">
  <input type="password" name="pass"/>
  <input type="text" name="title"/>
  <input type="checkbox" name="chk" value="chkval"/>
  <input type="radio" name="rdo" value="rdoval1" checked=checked/>
  <input type="radio" name="rdo" value="rdoval2" />
  <select name="select">
    <option value="1" selected="selected">1</option>
    <option value="2">2</option>
  </select>
  <textarea id="body" name="body" cols="80" rows="20" placeholder="hoge">hogehoge</textarea>
  <input type="submit" value="Send">
</form>
</body></html>
`

func TestFillinFormOptions(t *testing.T) {
	formData := map[string][]string{
		"title":  []string{"hogeTitle"},
		"chk":    []string{"1"},
		"rdo":    []string{"rdoval2"},
		"select": []string{"1"},
		"body":   []string{"hogehoge"},
		"pass":   []string{"hogepass"},
	}

	filler := newFiller(formData, map[string]interface{}{"Target": "myform2"})

	htmlstr := filler.fill([]byte(HTMLMulti))

	if string(htmlstr) != HTMLMultiSuccess {
		t.Errorf("fillinform error: %v", string(htmlstr))
	}

	filler = newFiller(formData, map[string]interface{}{"FillPassword": true})
	htmlstr = filler.fill([]byte(HTMLPassword))

	if string(htmlstr) != HTMLPasswordSuccess {
		t.Errorf("fillinform error: %v", string(htmlstr))
	}

	filler = newFiller(formData, map[string]interface{}{"IgnoreFields": []string{"title", "rdo"}})
	htmlstr = filler.fill([]byte(HTMLFields))

	if string(htmlstr) != HTMLFieldsSuccess {
		t.Errorf("fillinform error: %v", string(htmlstr))
	}
}

var HTMLBig = `
<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8" />
<title>会員登録(無料)</title>
<script type="text/javascript" charset="utf-8" src="/js/jquery.min.js"></script>
</head>
<body id="regist">
<form action="/register/simple/input" method="post" name="regist" id="input">
<input type="hidden" name=".site_token" value="t0.tSEbWYjOBNbOEfuBNV_Zqa6JVg4">
<dl><dt>性別</dt>
<dd><input type="radio" name="sex" value="0" id="member_sex_female"><label for="member_sex_female">女性</label>
<input type="radio" name="sex" value="1" id="member_sex_male"><label for="member_sex_male">男性</label></dd>
<dt>ニックネーム（8文字以内）</dt>
<dd><input type="text" name="user_name" size="25" istyle="1"></dd>
<dt>ログインID（半角英数字5〜20文字以内）</dt>
<dd><input id="loginid" type="text" name="loginid" size="25" autocapitalize="off" istyle="4">
</dd><dt>パスワード（半角英数字6～16文字以内）</dt>
<dd><input autocapitalize="off" type="password" name="user_pass" size="25" maxlength="16" istyle="4"></dd>
<dt>パスワード確認</dt>
<dd><input autocapitalize="off" type="password" name="user_pass_confirm" size="25" maxlength="16" istyle="4"></dd>
<dt>生年月日</dt>
<dd><select name="user_birth">
<option value="1910" >1910</option>
<option value="1911" >1911</option>
<option value="1912" >1912</option>
<option value="1913" >1913</option>
<option value="1914" >1914</option>
<option value="1915" >1915</option>
<option value="1916" >1916</option>
<option value="1917" >1917</option>
<option value="1918" >1918</option>
<option value="1919" >1919</option>
<option value="1920" >1920</option>
<option value="1921" >1921</option>
<option value="1922" >1922</option>
<option value="1923" >1923</option>
<option value="1924" >1924</option>
<option value="1925" >1925</option>
<option value="1926" >1926</option>
<option value="1927" >1927</option>
<option value="1928" >1928</option>
<option value="1929" >1929</option>
<option value="1930" >1930</option>
<option value="1931" >1931</option>
<option value="1932" >1932</option>
<option value="1933" >1933</option>
<option value="1934" >1934</option>
<option value="1935" >1935</option>
<option value="1936" >1936</option>
<option value="1937" >1937</option>
<option value="1938" >1938</option>
<option value="1939" >1939</option>
<option value="1940" >1940</option>
<option value="1941" >1941</option>
<option value="1942" >1942</option>
<option value="1943" >1943</option>
<option value="1944" >1944</option>
<option value="1945" >1945</option>
<option value="1946" >1946</option>
<option value="1947" >1947</option>
<option value="1948" >1948</option>
<option value="1949" >1949</option>
<option value="1950" >1950</option>
<option value="1951" >1951</option>
<option value="1952" >1952</option>
<option value="1953" >1953</option>
<option value="1954" >1954</option>
<option value="1955" >1955</option>
<option value="1956" >1956</option>
<option value="1957" >1957</option>
<option value="1958" >1958</option>
<option value="1959" >1959</option>
<option value="1960" >1960</option>
<option value="1961" >1961</option>
<option value="1962" >1962</option>
<option value="1963" >1963</option>
<option value="1964" >1964</option>
<option value="1965" >1965</option>
<option value="1966" >1966</option>
<option value="1967" >1967</option>
<option value="1968" >1968</option>
<option value="1969" >1969</option>
<option value="1970" >1970</option>
<option value="1971" >1971</option>
<option value="1972" >1972</option>
<option value="1973" >1973</option>
<option value="1974" >1974</option>
<option value="1975" >1975</option>
<option value="1976" >1976</option>
<option value="1977" >1977</option>
<option value="1978" >1978</option>
<option value="1979" >1979</option>
<option value="1980" >1980</option>
<option value="1981" >1981</option>
<option value="1982" >1982</option>
<option value="1983" >1983</option>
<option value="1984" >1984</option>
<option value="1985" selected="selected">1985</option>
<option value="1986" >1986</option>
<option value="1987" >1987</option>
<option value="1988" >1988</option>
<option value="1989" >1989</option>
<option value="1990" >1990</option>
<option value="1991" >1991</option>
<option value="1992" >1992</option>
<option value="1993" >1993</option>
<option value="1994" >1994</option>
<option value="1995" >1995</option>
<option value="1996" >1996</option>
<option value="1997" >1997</option>
</select>年
<select name="user_birth_month">
<option value="01" >1</option>
<option value="02" >2</option>
<option value="03" >3</option>
<option value="04" >4</option>
<option value="05" >5</option>
<option value="06" >6</option>
<option value="07" >7</option>
<option value="08" >8</option>
<option value="09" >9</option>
<option value="10" >10</option>
<option value="11" >11</option>
<option value="12" >12</option>
</select>月
<select name="user_birth_day">
<option value="01" >1</option>
<option value="02" >2</option>
<option value="03" >3</option>
<option value="04" >4</option>
<option value="05" >5</option>
<option value="06" >6</option>
<option value="07" >7</option>
<option value="08" >8</option>
<option value="09" >9</option>
<option value="10" >10</option>
<option value="11" >11</option>
<option value="12" >12</option>
<option value="13" >13</option>
<option value="14" >14</option>
<option value="15" >15</option>
<option value="16" >16</option>
<option value="17" >17</option>
<option value="18" >18</option>
<option value="19" >19</option>
<option value="20" >20</option>
<option value="21" >21</option>
<option value="22" >22</option>
<option value="23" >23</option>
<option value="24" >24</option>
<option value="25" >25</option>
<option value="26" >26</option>
<option value="27" >27</option>
<option value="28" >28</option>
<option value="29" >29</option>
<option value="30" >30</option>
<option value="31" >31</option>
</select>日</dd>
<dt>登録地域</dt>
<dd><select name="user_tdfk">
<option value="P01" >北海道</option>
<option value="P02" >青森県</option>
<option value="P03" >岩手県</option>
<option value="P04" >宮城県</option>
<option value="P05" >秋田県</option>
<option value="P06" >山形県</option>
<option value="P07" >福島県</option>
<option value="P08" >茨城県</option>
<option value="P09" >栃木県</option>
<option value="P10" >群馬県</option>
<option value="P11" >埼玉県</option>
<option value="P12" >千葉県</option>
<option value="P13" selected="selected">東京都</option>
<option value="P14" >神奈川県</option>
<option value="P15" >新潟県</option>
<option value="P16" >富山県</option>
<option value="P17" >石川県</option>
<option value="P18" >福井県</option>
<option value="P19" >山梨県</option>
<option value="P20" >長野県</option>
<option value="P21" >岐阜県</option>
<option value="P22" >静岡県</option>
<option value="P23" >愛知県</option>
<option value="P24" >三重県</option>
<option value="P25" >滋賀県</option>
<option value="P26" >京都府</option>
<option value="P27" >大阪府</option>
<option value="P28" >兵庫県</option>
<option value="P29" >奈良県</option>
<option value="P30" >和歌山県</option>
<option value="P31" >鳥取県</option>
<option value="P32" >島根県</option>
<option value="P33" >岡山県</option>
<option value="P34" >広島県</option>
<option value="P35" >山口県</option>
<option value="P36" >徳島県</option>
<option value="P37" >香川県</option>
<option value="P38" >愛媛県</option>
<option value="P39" >高知県</option>
<option value="P40" >福岡県</option>
<option value="P41" >佐賀県</option>
<option value="P42" >長崎県</option>
<option value="P43" >熊本県</option>
<option value="P44" >大分県</option>
<option value="P45" >宮崎県</option>
<option value="P46" >鹿児島県</option>
<option value="P47" >沖縄県</option>
<option value="P48" >海外</option>
</select></dd>
<dt>職業</dt>
<dd><select name="job_code">
<option value="1" >公務員</option>
<option value="2" >会社経営/自営業</option>
<option value="3" >役員/管理職</option>
<option value="4" >事務職/OL</option>
<option value="5" >受付/秘書</option>
<option value="6" >金融/不動産</option>
<option value="7" selected="selected">営業</option>
<option value="8" >企画/マーケティング</option>
<option value="9" >広報/広告宣伝</option>
<option value="10" >販売/飲食</option>
<option value="11" >旅行/宿泊/交通</option>
<option value="12" >技術者/コンピュータ関係</option>
<option value="13" >クリエイティブ/メディア</option>
<option value="14" >フリーランス</option>
<option value="15" >法律関係/弁護士</option>
<option value="16" >医療関係/医師</option>
<option value="17" >専門職</option>
<option value="18" >学生</option>
<option value="19" >パート/アルバイト</option>
<option value="20" >専業主婦/専業主夫</option>
<option value="21" >家事手伝い</option>
<option value="22" >無職</option>
<option value="-1" >その他</option>
</select></dd>
<dt>利用目的</dt>
<dd>
<select name="ctg_no" >
<option value="C01">友達募集</option>
<option value="C02">恋人募集</option>
<option value="C11">メル友募集</option>
</select>
</form>
</div>
</section>
</div>
</body>
</html>
`

var HTMLBigSuccess = `
<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8" />
<title>会員登録(無料)</title>
<script type="text/javascript" charset="utf-8" src="/js/jquery.min.js"></script>
</head>
<body id="regist">
<form action="/register/simple/input" method="post" name="regist" id="input">
<input type="hidden" name=".site_token" value="t0.tSEbWYjOBNbOEfuBNV_Zqa6JVg4">
<dl><dt>性別</dt>
<dd><input type="radio" name="sex" value="0" id="member_sex_female"><label for="member_sex_female">女性</label>
<input type="radio" name="sex" value="1" id="member_sex_male" checked="checked"><label for="member_sex_male">男性</label></dd>
<dt>ニックネーム（8文字以内）</dt>
<dd><input type="text" name="user_name" size="25" istyle="1" value="かわたん"></dd>
<dt>ログインID（半角英数字5〜20文字以内）</dt>
<dd><input id="loginid" type="text" name="loginid" size="25" autocapitalize="off" istyle="4" value="kawatan">
</dd><dt>パスワード（半角英数字6～16文字以内）</dt>
<dd><input autocapitalize="off" type="password" name="user_pass" size="25" maxlength="16" istyle="4"></dd>
<dt>パスワード確認</dt>
<dd><input autocapitalize="off" type="password" name="user_pass_confirm" size="25" maxlength="16" istyle="4"></dd>
<dt>生年月日</dt>
<dd><select name="user_birth">
<option value="1910" >1910</option>
<option value="1911" >1911</option>
<option value="1912" >1912</option>
<option value="1913" >1913</option>
<option value="1914" >1914</option>
<option value="1915" >1915</option>
<option value="1916" >1916</option>
<option value="1917" >1917</option>
<option value="1918" >1918</option>
<option value="1919" >1919</option>
<option value="1920" >1920</option>
<option value="1921" >1921</option>
<option value="1922" >1922</option>
<option value="1923" >1923</option>
<option value="1924" >1924</option>
<option value="1925" >1925</option>
<option value="1926" >1926</option>
<option value="1927" >1927</option>
<option value="1928" >1928</option>
<option value="1929" >1929</option>
<option value="1930" >1930</option>
<option value="1931" >1931</option>
<option value="1932" >1932</option>
<option value="1933" >1933</option>
<option value="1934" >1934</option>
<option value="1935" >1935</option>
<option value="1936" >1936</option>
<option value="1937" >1937</option>
<option value="1938" >1938</option>
<option value="1939" >1939</option>
<option value="1940" >1940</option>
<option value="1941" >1941</option>
<option value="1942" >1942</option>
<option value="1943" >1943</option>
<option value="1944" >1944</option>
<option value="1945" >1945</option>
<option value="1946" >1946</option>
<option value="1947" >1947</option>
<option value="1948" >1948</option>
<option value="1949" >1949</option>
<option value="1950" >1950</option>
<option value="1951" >1951</option>
<option value="1952" >1952</option>
<option value="1953" >1953</option>
<option value="1954" >1954</option>
<option value="1955" >1955</option>
<option value="1956" >1956</option>
<option value="1957" >1957</option>
<option value="1958" >1958</option>
<option value="1959" >1959</option>
<option value="1960" >1960</option>
<option value="1961" >1961</option>
<option value="1962" >1962</option>
<option value="1963" >1963</option>
<option value="1964" >1964</option>
<option value="1965" >1965</option>
<option value="1966" >1966</option>
<option value="1967" >1967</option>
<option value="1968" >1968</option>
<option value="1969" >1969</option>
<option value="1970" >1970</option>
<option value="1971" >1971</option>
<option value="1972" >1972</option>
<option value="1973" selected="selected">1973</option>
<option value="1974" >1974</option>
<option value="1975" >1975</option>
<option value="1976" >1976</option>
<option value="1977" >1977</option>
<option value="1978" >1978</option>
<option value="1979" >1979</option>
<option value="1980" >1980</option>
<option value="1981" >1981</option>
<option value="1982" >1982</option>
<option value="1983" >1983</option>
<option value="1984" >1984</option>
<option value="1985">1985</option>
<option value="1986" >1986</option>
<option value="1987" >1987</option>
<option value="1988" >1988</option>
<option value="1989" >1989</option>
<option value="1990" >1990</option>
<option value="1991" >1991</option>
<option value="1992" >1992</option>
<option value="1993" >1993</option>
<option value="1994" >1994</option>
<option value="1995" >1995</option>
<option value="1996" >1996</option>
<option value="1997" >1997</option>
</select>年
<select name="user_birth_month">
<option value="01" >1</option>
<option value="02" selected="selected">2</option>
<option value="03" >3</option>
<option value="04" >4</option>
<option value="05" >5</option>
<option value="06" >6</option>
<option value="07" >7</option>
<option value="08" >8</option>
<option value="09" >9</option>
<option value="10" >10</option>
<option value="11" >11</option>
<option value="12" >12</option>
</select>月
<select name="user_birth_day">
<option value="01" >1</option>
<option value="02" >2</option>
<option value="03" >3</option>
<option value="04" >4</option>
<option value="05" >5</option>
<option value="06" >6</option>
<option value="07" >7</option>
<option value="08" >8</option>
<option value="09" >9</option>
<option value="10" >10</option>
<option value="11" >11</option>
<option value="12" >12</option>
<option value="13" >13</option>
<option value="14" >14</option>
<option value="15" >15</option>
<option value="16" >16</option>
<option value="17" selected="selected">17</option>
<option value="18" >18</option>
<option value="19" >19</option>
<option value="20" >20</option>
<option value="21" >21</option>
<option value="22" >22</option>
<option value="23" >23</option>
<option value="24" >24</option>
<option value="25" >25</option>
<option value="26" >26</option>
<option value="27" >27</option>
<option value="28" >28</option>
<option value="29" >29</option>
<option value="30" >30</option>
<option value="31" >31</option>
</select>日</dd>
<dt>登録地域</dt>
<dd><select name="user_tdfk">
<option value="P01" >北海道</option>
<option value="P02" >青森県</option>
<option value="P03" >岩手県</option>
<option value="P04" >宮城県</option>
<option value="P05" >秋田県</option>
<option value="P06" >山形県</option>
<option value="P07" >福島県</option>
<option value="P08" >茨城県</option>
<option value="P09" >栃木県</option>
<option value="P10" >群馬県</option>
<option value="P11" >埼玉県</option>
<option value="P12" >千葉県</option>
<option value="P13">東京都</option>
<option value="P14" selected="selected">神奈川県</option>
<option value="P15" >新潟県</option>
<option value="P16" >富山県</option>
<option value="P17" >石川県</option>
<option value="P18" >福井県</option>
<option value="P19" >山梨県</option>
<option value="P20" >長野県</option>
<option value="P21" >岐阜県</option>
<option value="P22" >静岡県</option>
<option value="P23" >愛知県</option>
<option value="P24" >三重県</option>
<option value="P25" >滋賀県</option>
<option value="P26" >京都府</option>
<option value="P27" >大阪府</option>
<option value="P28" >兵庫県</option>
<option value="P29" >奈良県</option>
<option value="P30" >和歌山県</option>
<option value="P31" >鳥取県</option>
<option value="P32" >島根県</option>
<option value="P33" >岡山県</option>
<option value="P34" >広島県</option>
<option value="P35" >山口県</option>
<option value="P36" >徳島県</option>
<option value="P37" >香川県</option>
<option value="P38" >愛媛県</option>
<option value="P39" >高知県</option>
<option value="P40" >福岡県</option>
<option value="P41" >佐賀県</option>
<option value="P42" >長崎県</option>
<option value="P43" >熊本県</option>
<option value="P44" >大分県</option>
<option value="P45" >宮崎県</option>
<option value="P46" >鹿児島県</option>
<option value="P47" >沖縄県</option>
<option value="P48" >海外</option>
</select></dd>
<dt>職業</dt>
<dd><select name="job_code">
<option value="1" >公務員</option>
<option value="2" >会社経営/自営業</option>
<option value="3" >役員/管理職</option>
<option value="4" >事務職/OL</option>
<option value="5" >受付/秘書</option>
<option value="6" >金融/不動産</option>
<option value="7">営業</option>
<option value="8" >企画/マーケティング</option>
<option value="9" >広報/広告宣伝</option>
<option value="10" >販売/飲食</option>
<option value="11" >旅行/宿泊/交通</option>
<option value="12" selected="selected">技術者/コンピュータ関係</option>
<option value="13" >クリエイティブ/メディア</option>
<option value="14" >フリーランス</option>
<option value="15" >法律関係/弁護士</option>
<option value="16" >医療関係/医師</option>
<option value="17" >専門職</option>
<option value="18" >学生</option>
<option value="19" >パート/アルバイト</option>
<option value="20" >専業主婦/専業主夫</option>
<option value="21" >家事手伝い</option>
<option value="22" >無職</option>
<option value="-1" >その他</option>
</select></dd>
<dt>利用目的</dt>
<dd>
<select name="ctg_no" >
<option value="C01">友達募集</option>
<option value="C02">恋人募集</option>
<option value="C11" selected="selected">メル友募集</option>
</select>
</form>
</div>
</section>
</div>
</body>
</html>
`

func TestBigHTML(t *testing.T) {
	formData := map[string][]string{
		"sex":              []string{"1"},
		"user_name":        []string{"かわたん"},
		"loginid":          []string{"kawatan"},
		"user_birth":       []string{"1973"},
		"user_birth_month": []string{"02"},
		"user_birth_day":   []string{"17"},
		"user_tdfk":        []string{"P14"},
		"job_code":         []string{"12"},
		"ctg_no":           []string{"C11"},
		".site_token":      []string{"t0.tSEbWYjOBNbOEfuBNV_Zqa6JVg4"},
	}

	filler := newFiller(formData, nil)

	htmlstr := filler.fill([]byte(HTMLBig))

	if string(htmlstr) != HTMLBigSuccess {
		t.Errorf("fillinform error: %s", string(htmlstr))
	}
}
func BenchmarkBigHTML(b *testing.B) {
	formData := map[string][]string{
		"sex":              []string{"1"},
		"user_name":        []string{"かわたん"},
		"loginid":          []string{"kawatan"},
		"user_birth":       []string{"1973"},
		"user_birth_month": []string{"02"},
		"user_birth_day":   []string{"17"},
		"user_tdfk":        []string{"P14"},
		"job_code":         []string{"12"},
		"ctg_no":           []string{"C11"},
		".site_token":      []string{"t0.tSEbWYjOBNbOEfuBNV_Zqa6JVg4"},
	}

	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTMLBig))
	}
}
