package main

import (
	"fmt"
	"github.com/droundy/bridge"
)

func checkRadio(c bool) string {
	if c {
		return ` checked="checked"`
	}
	return ""
}

func selected(c bool) string {
	if c {
		return ` selected="selected"`
	}
	return ""
}

func SettingsPage(dat *ClientData, evt []string) string {
	switch evt[1] {
	case "set style ":
		dat.Cookie.Style = evt[2]
		dat.WriteCookie()
	case "set card ":
		if _,ok := dat.Cookie.Cards[evt[2]]; !ok {
			c := bridge.DefaultConvention()
			dat.Cookie.Cards[evt[2]] = &c
		}
		dat.Cookie.WhichCard = evt[2]
		dat.WriteCookie()
	}
	out := `
<div class="textish">
<div>
<fieldset><legend>Suit color style</legend>
`
	for _,s := range []string{"two color", "four color"} {
		out += fmt.Sprintf(`<input type="radio" name="style" onchange="say('set style %s')" value="%s"%s/>%s`,
			s, s, checkRadio(dat.Cookie.Style == s), s)
	}
	out += `
</fieldset>

<select onchange="say('set card '+this.value)" name="whichcard">
`
	for k := range dat.Cookie.Cards {
		out += fmt.Sprintf(`<option value="%s"%s> %s </option>`, k, selected(k == dat.Cookie.WhichCard), k)
	}
	out += `<option value="New convention card">New convention card...</option>
</select>
</div>
</div>

`
	out += conventionCard(dat.Cookie)
	return out
	/*
	defer header(c, getTransitoryData(req), "aBridge settings")()
	fmt.Fprintln(c, `<div class="textish">`)
	fmt.Fprintf(c, `<div>`)
	fmt.Fprintln(c, `<fieldset><legend>Suit color style</legend>`)
	for _,s := range []string{"two color", "four color"} {
		fmt.Fprintf(c, `<input type="radio" name="style" onchange="submitform()" title="foobar" value="%s"%s/>%s`,
			s, checkRadio(p.Style == s), s)
	}
	fmt.Fprintln(c, `</fieldset>`)

	fmt.Fprintln(c, `<select onchange="submitform()" name="whichcard">`)
	if _,ok := p.Cards[p.WhichCard]; !ok {
		fmt.Fprint(c, `<option selected="selected" value="`, p.WhichCard, `">`, p.WhichCard, `</option>`)
	}
	for k := range p.Cards {
		fmt.Fprint(c, `<option value="`, k, `"`)
		if k == p.WhichCard {
			fmt.Fprint(c, ` selected="selected"`)
		}
		fmt.Fprintln(c, `>`, k, `</option>`)
	}
	fmt.Fprintln(c, `<option value="New convention card">New convention card...</option>`)
	fmt.Fprintln(c, `</select>`)

	conventionCard(c, p)

	fmt.Fprintln(c, `<input type="submit" value="Save settings" />`)
	fmt.Fprintln(c, `<input type="submit" name="revert" value="Revert to default" />`)
	fmt.Fprintln(c, `</div></div>`)
	 */
}
