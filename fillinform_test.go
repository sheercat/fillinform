package fillinform

import "testing"

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
	str := `<input type="hoge" value="hoge">`
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
	str := `<input type="hoge" value="hoge">`
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
	str := `<input type="hoge" value="hoge" name="hoge">`
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
	`mixed`: {"arg": `html special char is <, >, &, and ".`, "return": `html special char is &lt;, &gt;, &amp;, and &quot;.`},
}

func TestEscapeHTML(t *testing.T) {
	filler := &Filler{}
	for key, mapval := range TesteeArrayESCHTML {
		val := mapval["arg"]
		res := mapval["return"]
		hoge := filler.escapeHTML(val)
		if hoge != res {
			t.Errorf("error in ", key)
		}
	}
}

func BenchmarkEscapeHTML(b *testing.B) {
	str := `<input & type="hoge" & value="hoge" & name="hoge">`
	filler := &Filler{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.escapeHTML(str)
	}
}

func TestFillinForm(t *testing.T) {
	formData := map[string]interface{}{
		"title":  "hogeTitle",
		"chk":    "1",
		"rdo":    "rdoval2",
		"select": "1",
		"body":   "hogehoge",
	}
	html := `
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

	success := `
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

	filler := &Filler{FillinFormOptions{Data: formData}}

	htmlstr := filler.fill([]byte(html))

	if string(htmlstr) != success {
		t.Errorf("fillinform error: ")
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
	html := `
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
	filler := &Filler{FillinFormOptions{Data: formData}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filler.fill([]byte(html))
	}
}
