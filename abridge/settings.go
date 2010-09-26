package main

import (
	"fmt"
	"io"
	"http"
	"template"
	"json"
	"strconv"
	"github.com/droundy/bridge"
)

type Settings struct {
	Style string
	Card bridge.ConventionCard
}

var DefaultSettings = Settings{
Style: "two color",
Card: bridge.DefaultConvention(),
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	p = DefaultSettings
	p.Card = bridge.DefaultConvention() // so we have fresh maps!
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	// fmt.Println("Trying to unmarshall string", prefstr)
	json.Unmarshal([]byte(prefstr), &p) // I don't care about errors!
	for k,v := range DefaultSettings.Card.Pts {
		if _,ok := p.Card.Pts[k]; !ok {
			p.Card.Pts[k] = v
		}
	}
	for k,v := range DefaultSettings.Card.Options {
		if _,ok := p.Card.Options[k]; !ok {
			p.Card.Options[k] = v
		}
	}
	for k,v := range DefaultSettings.Card.Radio {
		if _,ok := p.Card.Radio[k]; !ok {
			p.Card.Radio[k] = v
		}
	}
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
	// Since we may have cached things based on old card (or someone
	// else's card!), we should clear out the cache.
	bridge.ClearCache()
	if _,ok := req.Form["amsubmitting"]; ok {
		if s,ok := req.Form["style"]; ok {
			p.Style = s[0]
		}
		//for k,v := range req.Form {
		//	fmt.Println("Got key", k, "and value", v)
		//}
		if x,ok := req.Form["GeneralApproach"]; ok {
			p.Card.GeneralApproach = x[0]
		}
		if x,ok := req.Form["Name"]; ok {
			p.Card.Name = x[0]
		}
		for k := range DefaultSettings.Card.Pts {
			if x,ok := req.Form[k]; ok {
				pts,e := strconv.Atoi(x[0])
				if e == nil && bridge.Points(pts) != p.Card.Pts[k] {
					fmt.Println("Got pts", k, "of", pts)
					p.Card.Pts[k] = bridge.Points(pts)
				}
			}
		}
		if x,ok := req.Form["Jacobi"]; ok {
			fmt.Println("Jacobi stuff is", x, "with len", len(x))
			switch len(x) {
			case 2: p.Card.Options["Jacobi"] = true;
			case 1: p.Card.Options["Jacobi"] = !p.Card.Options["Jacobi"]
			case 0: p.Card.Options["Jacobi"] = false;
			}
		} else {
			p.Card.Options["Jacobi"] = false
		}
		for k := range DefaultSettings.Card.Options {
			// There are two Jacobi checkboxes, so I treat it specially...
			if k != "Jacobi" {
				_,ok := req.Form[k]
				p.Card.Options[k] = ok
			}
		}
		for k := range DefaultSettings.Card.Radio {
			if x,ok := req.Form[k]; ok {
				p.Card.Radio[k] = x[0]
			}
		}
	}
	if _,ok := req.Form["revert"]; ok {
		fmt.Println("Reverting to defaults...")
		p = DefaultSettings
		p.Card = bridge.DefaultConvention()
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
	fmt.Fprintln(c, `<fieldset><legend>Suit color style</legend>`)
	for _,s := range []string{"two color", "four color"} {
		fmt.Fprintf(c, `<input type="radio" name="style" title="foobar" value="%s"%s/>%s`,
			s, checkRadio(p.Style == s), s)
	}
	fmt.Fprintln(c, `</fieldset>`)

	conventionCard(c, p)

	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `<input type="submit" name="revert" value="Revert to default" />`)
	fmt.Fprintln(c, `</div></form>`)
	fmt.Fprintln(c, `</div>`)
}

func conventionCard(c *http.Conn, p Settings) {
	e := template.MustParse(cctemplate, myformatter).Execute(p.Card, c)
	if e != nil {
		fmt.Println("Template error:", e)
	}
}

var myformatter = template.FormatterMap(map[string]func(io.Writer, interface{}, string){
	"checked": func(c io.Writer, v interface{}, format string) {
		if b,ok := v.(bool); ok && b {
			fmt.Fprint(c, `checked="checked"`)
		}
	},
  "html": template.HTMLFormatter,
  "str":  template.StringFormatter,
	"": func(c io.Writer, v interface{}, format string) {
		if b,ok := v.(bool); ok {
			if b {
				fmt.Fprint(c, `checked="checked"`)
			}
			return
		}
		template.StringFormatter(c, v, format)
	},
	"Sound": compareStringThing,
	"Light": compareStringThing,
	"VeryLight": compareStringThing,
	"Natural": compareStringThing,
	"StrongTO": compareStringThing,
	"Michaels": compareStringThing,
	"NotForce": compareStringThing,
	"Force": compareStringThing,
	"Invitational": compareStringThing,
	"Weak": compareStringThing,
	"Intermediate": compareStringThing,
	"Strong": compareStringThing,
})

func compareStringThing(c io.Writer, v interface{}, format string) {
	if b,ok := v.(string); ok && b == format {
		fmt.Fprint(c, `checked="checked"`)
	}
}
