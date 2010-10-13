package main

func (dat *ClientData) Header() string {
	link := func(x string) string {
		if x == dat.MyPage {
			return " " + x + "\n"
		}
		return `<a href="javascript:say('go ` + x + `')">` + x + `</a>` + "\n"
	}
	return `<div id="header">
` + link("Home") + link("Analyze bids") + link("Bid fourth hand") + `
</div>
`
}
