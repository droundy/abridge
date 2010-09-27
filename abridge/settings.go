package main

import (
	"fmt"
	"http"
	"json"
	"strconv"
	"github.com/droundy/bridge"
)

type Settings struct {
	Style string
	Card bridge.ConventionCard
	Cards map[string]bridge.ConventionCard
}

func DefaultSettings() Settings {
	return Settings{
	Style: "two color",
	Card: bridge.DefaultConvention(),
	Cards: make(map[string]bridge.ConventionCard),
	}
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	p = DefaultSettings()
	p.Card = bridge.DefaultConvention() // so we have fresh maps!
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	// fmt.Println("Trying to unmarshall string", prefstr)
	json.Unmarshal([]byte(prefstr), &p) // I don't care about errors!
	for k,v := range bridge.DefaultConvention().Pts {
		if _,ok := p.Card.Pts[k]; !ok {
			p.Card.Pts[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention().Options {
		if _,ok := p.Card.Options[k]; !ok {
			p.Card.Options[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention().Radio {
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
		for k := range bridge.DefaultConvention().Pts {
			if x,ok := req.Form[k]; ok {
				pts,e := strconv.Atoi(x[0])
				if e == nil && bridge.Points(pts) != p.Card.Pts[k] {
					fmt.Println("Got pts", k, "of", pts)
					p.Card.Pts[k] = bridge.Points(pts)
				}
			}
		}
		for k := range bridge.DefaultConvention().Options {
			// There are two Jacobi checkboxes, so I treat it specially...
			vs,ok := req.Form[k]
			if ok && len(vs) == 1 && vs[0][0] == '2' {
				// There are two checkboxes for this one, so checking just
				// one of them means we're trying to change it!
				fmt.Println(k, "stuff is", vs)
				p.Card.Options[k] = !p.Card.Options[k]
			} else {
				p.Card.Options[k] = ok
			}
		}
		for k := range bridge.DefaultConvention().Radio {
			if x,ok := req.Form[k]; ok {
				p.Card.Radio[k] = x[0]
			}
		}
	}
	if _,ok := req.Form["revert"]; ok {
		fmt.Println("Reverting to defaults...")
		p = DefaultSettings()
	}

	p.Set(c)
	defer header(c, req, "aBridge settings")()
	fmt.Fprintln(c, `<div class="textish">`)
	fmt.Fprintf(c, `<form id="settings" method="post" action="%s"><div>`, req.URL.Path)
	fmt.Fprintln(c, `<fieldset><legend>Suit color style</legend>`)
	for _,s := range []string{"two color", "four color"} {
		fmt.Fprintf(c, `<input type="radio" name="style" onchange="submitform()" title="foobar" value="%s"%s/>%s`,
			s, checkRadio(p.Style == s), s)
	}
	fmt.Fprintln(c, `</fieldset>`)

	conventionCard(c, p)

	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `<input type="submit" name="revert" value="Revert to default" />`)
	fmt.Fprintln(c, `</div></form>`)
	fmt.Fprintln(c, `</div>`)
}
