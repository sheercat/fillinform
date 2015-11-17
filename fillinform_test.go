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
	filler := newFiller(formData, nil)

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
	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillInput([]byte(`<input type="radio" name="rdo" value="rdoval2" />`))
	}
}

func TestFillTextarea(t *testing.T) {
	formData := map[string]interface{}{
		"body": "hoge & hoge <hoge@hogehoge>",
	}
	filler := newFiller(formData, nil)

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
	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fillTextarea([]byte(`<textarea id="body" name="body" cols="80" rows="20" placeholder="hoge"></textarea>`))
	}
}

func TestFillSelect(t *testing.T) {
	formData := map[string]interface{}{
		"select": "1",
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
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "chkval",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}
	filler := newFiller(formData, nil)

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
	filler := newFiller(formData, nil)

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

	filler := newFiller(formData, nil)

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

	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTML))
	}
}

func TestFillinForm2(t *testing.T) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    1,
		"rdo":    "rdoval2",
		"select": 1,
		"body":   "hogehoge",
	}

	filler := newFiller(formData, nil)

	htmlstr := filler.fill([]byte(HTML))

	if string(htmlstr) != HTMLSuccess {
		t.Errorf("fillinform error: ", string(htmlstr))
	}
}

func BenchmarkFillinForm2(b *testing.B) {
	var i int64
	i = 1
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    1,
		"rdo":    "rdoval2",
		"select": i,
		"body":   "hogehoge",
	}

	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTML))
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
  <input type="hidden" name="hidden"/>
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
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "1",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
		"pass":   "hogepass",
	}

	filler := newFiller(formData, map[string]interface{}{"Target": "myform2"})

	htmlstr := filler.fill([]byte(HTMLMulti))

	if string(htmlstr) != HTMLMultiSuccess {
		t.Errorf("fillinform error: ", string(htmlstr))
	}

	filler = newFiller(formData, map[string]interface{}{"FillPassword": true})
	htmlstr = filler.fill([]byte(HTMLPassword))

	if string(htmlstr) != HTMLPasswordSuccess {
		t.Errorf("fillinform error: ", string(htmlstr))
	}

	filler = newFiller(formData, map[string]interface{}{"IgnoreFields": []string{"title", "rdo"}})
	htmlstr = filler.fill([]byte(HTMLFields))

	if string(htmlstr) != HTMLFieldsSuccess {
		t.Errorf("fillinform error: ", string(htmlstr))
	}
}

var HTMLBig = `
<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="utf-8" />
<title>会員登録(無料)</title>
<meta name="keywords" content="恋愛,結婚,恋人,メル友,掲示板,プロフィール,プロフィール検索,写真" />
<meta name="description" content="日本最大級のサイトはおかげさまで15周年､累計会員数1000万人以上の方にご利用いただいております｡あなたをまじめに応援します｡" />
<meta name="viewport" content="width=320, user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0">
<meta name="format-detection" content="telephone=no" /><link rel="apple-touch-icon-precomposed" href="/apple-touch-icon.png" />
<meta name="google-site-verification" content="-VqPPGh8eAbXl8zjNmXOQoiNj_GNYC3PZXkUXOFUas8" /><link rel="stylesheet" href="/css2/import.css?t=1340786666" />
<link rel="stylesheet" href="/css2/guest.css?t=1340786666" />
<link rel="stylesheet" href="/css2/lp.css" />
<link rel="stylesheet" href="/css2/docs.css" />
<link rel="stylesheet" href="/css2/affiliate.css" />
<script type="text/javascript" charset="utf-8" src="/js/jquery.min.js"></script>
<script type="text/javascript" charset="utf-8" src="/js/common.js?t=1418317158"></script>
</head>
<body id="regist">
<header id="guest-header" class="global-header">
<h1><a href="http://.co.jp/">応援サイト</a></h1>
</header>
<div id="container">
<div class="headline2"><h2>無料会員登録</h2></div>
<div id="regist-step">
<img src="/img2/regist/unidentified/step1.png" width="300"></div>
<section id="regist-input">
<div id="regist-input-body">
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
<p>※選択内容に応じてプロフィールの目的に反映されます。(男女共通)<br>
また、女性のみ掲示板に自動で投稿されます。</p>
</dd>




</dl>

<div class="btn"><input type="submit" class="block-btn l-spacing-l g-btn" value="登録内容を確認する"></div>

<div style="background-color:#eee;padding:5px;">
<div id="toggleTester" style="height:5px;">
</div>
</div>
<script>
<!--
    $(function() {
        $('.suggest-loginid').click(
            function(){
                $('#loginid').val($(this).text());
            }
        );
    });
-->
</script>

</form>
</div>
</section>

</div><!--/container-->
<footer id="global-footer"></footer>
<script type="text/javascript">
$(function() {
  var pageTop = $('.returnTop');
  pageTop.hide();
  $(window).scroll(function () {
    if ($(this).scrollTop() > 200) {
      pageTop.fadeIn();
    } else {
      pageTop.fadeOut();
    }
  });
    pageTop.click(function () {
    $('body, html').animate({scrollTop:0}, 500, 'swing');
    return false;
    });
});
</script>
<footer>
<div id="footer">
<nav>
<div><a href="http://.co.jp/hc/ja" target="_blank">ヘルプセンター</a></div>
<p>&copy; .Inc</p>
</nav>
</div>
</footer>
<div id="setting-list" class="slideOverlay">
<div id="floatbox-title"><h3>各種設定</h3><a href="javascript:void(0)" id="setting-list-close" class="block-btn-s">閉じる</a></div>
<section id="floating-config">
<table>
<tr>
<th><img src="/img2/icon/50-b-01.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/account/"><dl>
<dt>基本情報設定</dt>
<dd>メールアドレス、パスワードの設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-05.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/profile/"><dl>
<dt>プロフィール設定</dt>
<dd>プロフィールの設定・更新ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-27.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/design/"><dl>
<dt>デザイン設定</dt>
<dd>マイページ、プロフィール、今、なにしてるのページのデザインを変更できます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-10.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/config/auto_login"><dl>
<dt>オートログイン設定</dt>
<dd>ユーザーID・パスワードの入力なしで、ログイン可能な設定ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-02.png" width="25" height="25"></th>
<td><a href="https://.co.jp/my/config/mail_reception"><dl>
<dt>メール受信設定</dt>
<dd>受信数や受信時間、登録メールアドレス以外へのメール通知設定などメール受信に関する設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-30.png" width="25" height="25"></th>
<td><a href="http://.co.jp/bottlemail/"><dl>
<dt>ボトルメール設定<span>NEW!</span></dt>
<dd>ボトルメールの配信文章と送受信の設定ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-26.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/foot_print/config/"><dl>
<dt>足あとお知らせ設定</dt>
<dd>足あとがついた場合メールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/interest/icon_type_home.png" width="25" height="25"></th>
<td><a href="http://.co.jp/interested/config/"><dl>
<dt>タイプお知らせ設定</dt>
<dd>タイプされた時と両思いになった時の通知設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/interest/icon_type_home.png" width="25" height="25"></th>
<td><a href="http://.co.jp/private_photo/notifier_config/"><dl>
<dt>プライベート写真お知らせ設定</dt>
<dd>プライベート写真が許可された時の通知設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-06.png" width="25" height="25"></th>
<td><a href="http://.co.jp/bbs/config/"><dl>
<dt>最新投稿お知らせ設定</dt>
<dd>掲示板の最新投稿をメールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-07.png" width="25" height="25"></th>
<td><a href="http://.co.jp/chat/notice/"><dl>
<dt>チャットお知らせ設定</dt>
<dd>チャットの開設通知をメールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-23.png" width="25" height="25"></th>
<td><a href="http://.co.jp/diary/config/"><dl>
<dt>日記設定</dt>
<dd>日記のコメント通知やお気に入り通知設定が可能です。</dd>
</dl></a></td>
</tr>

<tr>
<th><img src="/img2/icon/50-b-04.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/mail_box/template/"><dl>
<dt>保存メッセージ設定</dt>
<dd>お気に入りの文章を設定し、サイト内からメールを送信する際に使用できる機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-29.png" width="25" height="25"></th>
<td><a href="http://.co.jp/fortune/config/"><dl>
<dt>占いお知らせ設定<span>NEW!</span></dt>
<dd>毎日の運勢をメールでお知らせする機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-28.png" width="25" height="25"></th>
<td><a href="http://.co.jp/logout/"><dl>
<dt>ログアウト</dt>
<dd>ログアウトします。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-08.png" width="25" height="25"></th>
<td><a href="https://ssl..co.jp/resign/"><dl>
<dt>退会</dt>
<dd>退会する方はこちら。</dd>
</dl></a></td>
</tr>
</table>
</section>
<div id="floatbox-footer">&copy; .Inc</div>
</div>

<div id="service-list" class="slideOverlay">
<div id="floatbox-title"><h3>サービス一覧</h3><a href="javascript:void(0)" id="service-list-close" class="block-btn-s">閉じる</a></div>

<section id="yb_bnr">
<a href="http://.jp/r/fsp/" target="_blank"><img src="/img2/bnr/you.png" border="0" width="100%" alt="" /></a>
</section>

<section id="community">
<dl>
<dt class="icon-label prof-search"><a href="http://.co.jp/profile/search/result?.xf=lite_nav&.xt=profile_search">プロフィール検索<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/profile/search/newcomer?.xf=lite_nav&.xt=profile_search_newcomer">&nbsp;新人検索</a></li>

<li><a href="http://.co.jp/profile/search/photo?.xf=lite_nav&.xt=profile_search_photo">&nbsp;写真検索</a></li>
<li><a href="http://.co.jp/profile/search/lite?.xf=lite_nav&.xt=profile_search_entry_tdfk" >&nbsp;設定地域で探す</a></li>
</ul></dd>
<dt class="icon-label bbs"><a href="http://.co.jp/bbs/pure/?.xf=lite_nav&.xt=bbs_pure">掲示板<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/bbs/pure/newcomer?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;新人掲示板</a></li>
<li><a href="http://.co.jp/bbs/pure/chance?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;チャンス掲示板</a></li>
<li><a href="http://.co.jp/bbs/pure/?.xf=lite_nav&.xt=bbs_pure_search_entry_tdfk">&nbsp;設定地域で見る</a></li>
<li><a href="http://.co.jp/bbs/pure/post/?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;投稿する</a></li>
</ul></dd>

<dt class="icon-label imanani"><a href="http://.co.jp/short_message/?.xf=lite_nav&.xt=short_message">今なにしてる？<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/short_message/?.xf=top&.xt=short_message">&nbsp;みんなのつぶやき一覧</a></li>
<li><a href="http://.co.jp/short_message/my?.xf=top&.xt=short_message">&nbsp;あなたのつぶやき一覧</a></li>
</ul></dd>

<dt class="icon-label gourmet"><a href="http://.co.jp/gourmet/?.xf=lite_nav&.xt=gourmet">グルメデート<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/gourmet/?.xf=top&.xt=gourmet_search">&nbsp;グルメデートを探す</a></li>
<li><a href="http://.co.jp/gourmet/entry/create?.xf=top&.xt=%2Fgourmet_entry_create">&nbsp;グルメデートを募集する</a></li>
</ul></dd>
<dt class="icon-label photo-chat"><a href="http://.co.jp/chat/?.xf=lite_nav&.xt=chat">フォトチャット<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/chat/twoshot/?.xf=lite_nav&.xt=chat_twoshot">&nbsp;2ショットチャット2部屋待ち</a></li>
<li><a href="http://.co.jp/chat/twoshot/create/?.xf=lite_nav&.xt=chat_create_twoshot">&nbsp;2ショットチャットを開設する</a></li>
<li><a href="http://.co.jp/chat/group/?.xf=lite_nav&.xt=chat_group">&nbsp;グループチャット14部屋会話中</a></li>
<li><a href="http://.co.jp/chat/group/create/?.xf=lite_nav&.xt=chat_create_group">&nbsp;グループチャットを開設する</a></li>
</ul></dd>
<dt class="icon-label diary"><a href="http://.co.jp/diary/?.xf=lite_nav&.xt=diary">みんなの日記<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/diary/search/?.xf=lite_nav&.xt=diary_search">&nbsp;日記を探す</a></li>
<li><a href="http://.co.jp/diary/article/post?.xf=lite_nav&.xt=diary_article_post">&nbsp;日記を書く</a></li>
</ul></dd>
<dt class="icon-label research"><a href="http://.co.jp/research/?.xf=lite_nav&.xt=research">まるばつ！<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/research/review/ranking?.xf=lite_nav&.xt=research_ranking">&nbsp;回答ランキング</a></li>
<li><a href="http://.co.jp/research/edit?.xf=lite_nav&.xt=research_edit">&nbsp;出題する</a></li>
</ul></dd>
<dt class="icon-label feedback"><a href="http://.co.jp/feedback/tl?.xf=lite_nav&.xt=feedback">みんなのイイネ！<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/feedback/actioned?.xf=lite_nav&.xt=feedback_actioned">&nbsp;自分のイイネ！</a></li>
<li><a href="http://.co.jp/feedback/receive?.xf=lite_nav&.xt=feedback_receive">&nbsp;自分へのイイネ！</a></li>
</ul></dd>
</dl>
</section>
<div id="floatbox-footer">&copy; .Inc</div>
</div>

<div id="case-iphone-lightbox" class="slideOverlay fromSlideOverlay">
<div id="floatbox-title"><h3>空メール送信</h3><a href="javascript:void(0)" id="case-iphone-lightbox-close">閉じる</a></div>
<dl class="mailConfigBox">
<dt><div class="confirmBtn"><a href="/register/simple/@?subject=CCC" class="photoBtn">通常のメールの方</a></div></dt>
<dd>i.softbank／gmail等のメールアドレス</dd>
<dt><div class="confirmBtn"><a href="/register/simple/@" class="photoBtn">SMS／MMSの方</a></div></dt>
<dd>「SMS／MMS」は本文が空のままではメール送信できません。<br />1文字でも何か入力して送信してください。</dd>
</dl>
<div id="floatbox-footer">&copy;&nbsp;.Inc</div>
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
<meta name="keywords" content="恋愛,結婚,恋人,メル友,掲示板,プロフィール,プロフィール検索,写真" />
<meta name="description" content="日本最大級のサイトはおかげさまで15周年､累計会員数1000万人以上の方にご利用いただいております｡あなたをまじめに応援します｡" />
<meta name="viewport" content="width=320, user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0">
<meta name="format-detection" content="telephone=no" /><link rel="apple-touch-icon-precomposed" href="/apple-touch-icon.png" />
<meta name="google-site-verification" content="-VqPPGh8eAbXl8zjNmXOQoiNj_GNYC3PZXkUXOFUas8" /><link rel="stylesheet" href="/css2/import.css?t=1340786666" />
<link rel="stylesheet" href="/css2/guest.css?t=1340786666" />
<link rel="stylesheet" href="/css2/lp.css" />
<link rel="stylesheet" href="/css2/docs.css" />
<link rel="stylesheet" href="/css2/affiliate.css" />
<script type="text/javascript" charset="utf-8" src="/js/jquery.min.js"></script>
<script type="text/javascript" charset="utf-8" src="/js/common.js?t=1418317158"></script>
</head>
<body id="regist">
<header id="guest-header" class="global-header">
<h1><a href="http://.co.jp/">応援サイト</a></h1>
</header>
<div id="container">
<div class="headline2"><h2>無料会員登録</h2></div>
<div id="regist-step">
<img src="/img2/regist/unidentified/step1.png" width="300"></div>
<section id="regist-input">
<div id="regist-input-body">
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
<p>※選択内容に応じてプロフィールの目的に反映されます。(男女共通)<br>
また、女性のみ掲示板に自動で投稿されます。</p>
</dd>




</dl>

<div class="btn"><input type="submit" class="block-btn l-spacing-l g-btn" value="登録内容を確認する"></div>

<div style="background-color:#eee;padding:5px;">
<div id="toggleTester" style="height:5px;">
</div>
</div>
<script>
<!--
    $(function() {
        $('.suggest-loginid').click(
            function(){
                $('#loginid').val($(this).text());
            }
        );
    });
-->
</script>

</form>
</div>
</section>

</div><!--/container-->
<footer id="global-footer"></footer>
<script type="text/javascript">
$(function() {
  var pageTop = $('.returnTop');
  pageTop.hide();
  $(window).scroll(function () {
    if ($(this).scrollTop() > 200) {
      pageTop.fadeIn();
    } else {
      pageTop.fadeOut();
    }
  });
    pageTop.click(function () {
    $('body, html').animate({scrollTop:0}, 500, 'swing');
    return false;
    });
});
</script>
<footer>
<div id="footer">
<nav>
<div><a href="http://.co.jp/hc/ja" target="_blank">ヘルプセンター</a></div>
<p>&copy; .Inc</p>
</nav>
</div>
</footer>
<div id="setting-list" class="slideOverlay">
<div id="floatbox-title"><h3>各種設定</h3><a href="javascript:void(0)" id="setting-list-close" class="block-btn-s">閉じる</a></div>
<section id="floating-config">
<table>
<tr>
<th><img src="/img2/icon/50-b-01.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/account/"><dl>
<dt>基本情報設定</dt>
<dd>メールアドレス、パスワードの設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-05.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/profile/"><dl>
<dt>プロフィール設定</dt>
<dd>プロフィールの設定・更新ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-27.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/design/"><dl>
<dt>デザイン設定</dt>
<dd>マイページ、プロフィール、今、なにしてるのページのデザインを変更できます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-10.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/config/auto_login"><dl>
<dt>オートログイン設定</dt>
<dd>ユーザーID・パスワードの入力なしで、ログイン可能な設定ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-02.png" width="25" height="25"></th>
<td><a href="https://.co.jp/my/config/mail_reception"><dl>
<dt>メール受信設定</dt>
<dd>受信数や受信時間、登録メールアドレス以外へのメール通知設定などメール受信に関する設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-30.png" width="25" height="25"></th>
<td><a href="http://.co.jp/bottlemail/"><dl>
<dt>ボトルメール設定<span>NEW!</span></dt>
<dd>ボトルメールの配信文章と送受信の設定ができます。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-26.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/foot_print/config/"><dl>
<dt>足あとお知らせ設定</dt>
<dd>足あとがついた場合メールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/interest/icon_type_home.png" width="25" height="25"></th>
<td><a href="http://.co.jp/interested/config/"><dl>
<dt>タイプお知らせ設定</dt>
<dd>タイプされた時と両思いになった時の通知設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/interest/icon_type_home.png" width="25" height="25"></th>
<td><a href="http://.co.jp/private_photo/notifier_config/"><dl>
<dt>プライベート写真お知らせ設定</dt>
<dd>プライベート写真が許可された時の通知設定が可能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-06.png" width="25" height="25"></th>
<td><a href="http://.co.jp/bbs/config/"><dl>
<dt>最新投稿お知らせ設定</dt>
<dd>掲示板の最新投稿をメールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-07.png" width="25" height="25"></th>
<td><a href="http://.co.jp/chat/notice/"><dl>
<dt>チャットお知らせ設定</dt>
<dd>チャットの開設通知をメールで通知する機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-23.png" width="25" height="25"></th>
<td><a href="http://.co.jp/diary/config/"><dl>
<dt>日記設定</dt>
<dd>日記のコメント通知やお気に入り通知設定が可能です。</dd>
</dl></a></td>
</tr>

<tr>
<th><img src="/img2/icon/50-b-04.png" width="25" height="25"></th>
<td><a href="http://.co.jp/my/mail_box/template/"><dl>
<dt>保存メッセージ設定</dt>
<dd>お気に入りの文章を設定し、サイト内からメールを送信する際に使用できる機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-29.png" width="25" height="25"></th>
<td><a href="http://.co.jp/fortune/config/"><dl>
<dt>占いお知らせ設定<span>NEW!</span></dt>
<dd>毎日の運勢をメールでお知らせする機能です。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-28.png" width="25" height="25"></th>
<td><a href="http://.co.jp/logout/"><dl>
<dt>ログアウト</dt>
<dd>ログアウトします。</dd>
</dl></a></td>
</tr>
<tr>
<th><img src="/img2/icon/50-b-08.png" width="25" height="25"></th>
<td><a href="https://ssl..co.jp/resign/"><dl>
<dt>退会</dt>
<dd>退会する方はこちら。</dd>
</dl></a></td>
</tr>
</table>
</section>
<div id="floatbox-footer">&copy; .Inc</div>
</div>

<div id="service-list" class="slideOverlay">
<div id="floatbox-title"><h3>サービス一覧</h3><a href="javascript:void(0)" id="service-list-close" class="block-btn-s">閉じる</a></div>

<section id="yb_bnr">
<a href="http://.jp/r/fsp/" target="_blank"><img src="/img2/bnr/you.png" border="0" width="100%" alt="" /></a>
</section>

<section id="community">
<dl>
<dt class="icon-label prof-search"><a href="http://.co.jp/profile/search/result?.xf=lite_nav&.xt=profile_search">プロフィール検索<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/profile/search/newcomer?.xf=lite_nav&.xt=profile_search_newcomer">&nbsp;新人検索</a></li>

<li><a href="http://.co.jp/profile/search/photo?.xf=lite_nav&.xt=profile_search_photo">&nbsp;写真検索</a></li>
<li><a href="http://.co.jp/profile/search/lite?.xf=lite_nav&.xt=profile_search_entry_tdfk" >&nbsp;設定地域で探す</a></li>
</ul></dd>
<dt class="icon-label bbs"><a href="http://.co.jp/bbs/pure/?.xf=lite_nav&.xt=bbs_pure">掲示板<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/bbs/pure/newcomer?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;新人掲示板</a></li>
<li><a href="http://.co.jp/bbs/pure/chance?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;チャンス掲示板</a></li>
<li><a href="http://.co.jp/bbs/pure/?.xf=lite_nav&.xt=bbs_pure_search_entry_tdfk">&nbsp;設定地域で見る</a></li>
<li><a href="http://.co.jp/bbs/pure/post/?.xf=lite_nav&.xt=bbs_pure_post">&nbsp;投稿する</a></li>
</ul></dd>

<dt class="icon-label imanani"><a href="http://.co.jp/short_message/?.xf=lite_nav&.xt=short_message">今なにしてる？<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/short_message/?.xf=top&.xt=short_message">&nbsp;みんなのつぶやき一覧</a></li>
<li><a href="http://.co.jp/short_message/my?.xf=top&.xt=short_message">&nbsp;あなたのつぶやき一覧</a></li>
</ul></dd>

<dt class="icon-label gourmet"><a href="http://.co.jp/gourmet/?.xf=lite_nav&.xt=gourmet">グルメデート<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/gourmet/?.xf=top&.xt=gourmet_search">&nbsp;グルメデートを探す</a></li>
<li><a href="http://.co.jp/gourmet/entry/create?.xf=top&.xt=%2Fgourmet_entry_create">&nbsp;グルメデートを募集する</a></li>
</ul></dd>
<dt class="icon-label photo-chat"><a href="http://.co.jp/chat/?.xf=lite_nav&.xt=chat">フォトチャット<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/chat/twoshot/?.xf=lite_nav&.xt=chat_twoshot">&nbsp;2ショットチャット2部屋待ち</a></li>
<li><a href="http://.co.jp/chat/twoshot/create/?.xf=lite_nav&.xt=chat_create_twoshot">&nbsp;2ショットチャットを開設する</a></li>
<li><a href="http://.co.jp/chat/group/?.xf=lite_nav&.xt=chat_group">&nbsp;グループチャット14部屋会話中</a></li>
<li><a href="http://.co.jp/chat/group/create/?.xf=lite_nav&.xt=chat_create_group">&nbsp;グループチャットを開設する</a></li>
</ul></dd>
<dt class="icon-label diary"><a href="http://.co.jp/diary/?.xf=lite_nav&.xt=diary">みんなの日記<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/diary/search/?.xf=lite_nav&.xt=diary_search">&nbsp;日記を探す</a></li>
<li><a href="http://.co.jp/diary/article/post?.xf=lite_nav&.xt=diary_article_post">&nbsp;日記を書く</a></li>
</ul></dd>
<dt class="icon-label research"><a href="http://.co.jp/research/?.xf=lite_nav&.xt=research">まるばつ！<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/research/review/ranking?.xf=lite_nav&.xt=research_ranking">&nbsp;回答ランキング</a></li>
<li><a href="http://.co.jp/research/edit?.xf=lite_nav&.xt=research_edit">&nbsp;出題する</a></li>
</ul></dd>
<dt class="icon-label feedback"><a href="http://.co.jp/feedback/tl?.xf=lite_nav&.xt=feedback">みんなのイイネ！<span class="block-btn-s">もっと見る</span></a></dt>
<dd><ul class="block-link">
<li><a href="http://.co.jp/feedback/actioned?.xf=lite_nav&.xt=feedback_actioned">&nbsp;自分のイイネ！</a></li>
<li><a href="http://.co.jp/feedback/receive?.xf=lite_nav&.xt=feedback_receive">&nbsp;自分へのイイネ！</a></li>
</ul></dd>
</dl>
</section>
<div id="floatbox-footer">&copy; .Inc</div>
</div>

<div id="case-iphone-lightbox" class="slideOverlay fromSlideOverlay">
<div id="floatbox-title"><h3>空メール送信</h3><a href="javascript:void(0)" id="case-iphone-lightbox-close">閉じる</a></div>
<dl class="mailConfigBox">
<dt><div class="confirmBtn"><a href="/register/simple/@?subject=CCC" class="photoBtn">通常のメールの方</a></div></dt>
<dd>i.softbank／gmail等のメールアドレス</dd>
<dt><div class="confirmBtn"><a href="/register/simple/@" class="photoBtn">SMS／MMSの方</a></div></dt>
<dd>「SMS／MMS」は本文が空のままではメール送信できません。<br />1文字でも何か入力して送信してください。</dd>
</dl>
<div id="floatbox-footer">&copy;&nbsp;.Inc</div>
</div>

</body>
</html>
`

func TestBigHTML(t *testing.T) {
	formData := map[string]interface{}{
		"sex":              "1",
		"user_name":        "かわたん",
		"loginid":          "kawatan",
		"user_birth":       "1973",
		"user_birth_month": "02",
		"user_birth_day":   "17",
		"user_tdfk":        "P14",
		"job_code":         "12",
		"ctg_no":           "C11",
	}

	filler := newFiller(formData, nil)

	htmlstr := filler.fill([]byte(HTMLBig))

	if string(htmlstr) != HTMLBigSuccess {
		t.Errorf("fillinform error: snip")
	}
}
func BenchmarkBigHTML(b *testing.B) {
	formData := map[string]interface{}{
		"sex":              "1",
		"user_name":        "かわたん",
		"loginid":          "kawatan",
		"user_birth":       "1973",
		"user_birth_month": "02",
		"user_birth_day":   "17",
		"user_tdfk":        "P14",
		"job_code":         "12",
		"ctg_no":           "C11",
	}

	filler := newFiller(formData, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(HTMLBig))
	}
}
