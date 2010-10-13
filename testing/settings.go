package main

func checkRadio(c bool) string {
	if c {
		return ` checked="checked"`
	}
	return ""
}

func SettingsPage(dat *ClientData, evt []string) string {
	return `This is a stand-in for a real settings page`

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
