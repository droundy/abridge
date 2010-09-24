package main

import (
	"fmt"
	"http"
	"json"
)

type Settings struct {
	Style string
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	p = Settings{"default"}
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	// fmt.Println("Trying to unmarshall string", prefstr)
	json.Unmarshal([]byte(prefstr), &p) // I don't care about errors!
	// fmt.Println("Unmarshal gave: %v", e)
	return p
}

func (s *Settings) Set(c *http.Conn) {
	bytes,_ := json.Marshal(*s)
	c.SetHeader("Set-Cookie", string(bytes))
}

func checkRadio(c bool) string {
	if c {
		return ` checked="checked"`
	}
	return ""
}

func settings(c *http.Conn, req *http.Request) {
	p := getSettings(req)
	if s,ok := req.Form["style"]; ok {
		p.Style = s[0]
	}
	p.Set(c)
	defer header(c, req, "aBridge settings")()
	fmt.Fprintln(c, `<div class="textish">`)
	fmt.Fprint(c, `

<p> I am planning to add various configuration options here.  Ideally,
we'll even support a nice convention card input interface.  For now,
I'll probably just add a toggle to switch between two-color and
four-color suits (which will also test that I'm using CSS consistently).

</p>

`)
	fmt.Fprintf(c, `<form method="post" action="%s"><div>`, req.URL.Path)
	for _,s := range []string{"two color", "four color"} {
		fmt.Fprintf(c, `<input type="radio" name="style" value="%s"%s/>%s<br/>`,
			s, checkRadio(p.Style == s), s)
	}
	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `</div></form>`)
	fmt.Fprintln(c, `</div>`)
}
