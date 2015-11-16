package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/sheercat/fillinform"
)

func main() {
	cpuprofile := "mycpu.prof"
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting...", sig)
			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

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
	for {
		filled := fillinform.Fill([]byte(html), formData, nil)
		fmt.Println(filled)
	}
}
