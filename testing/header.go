package main

func (dat *ClientData) Header() string {
	link := func(x string) string {
		if x == dat.MyPage {
			return " " + x + "\n"
		}
		return `<a href="javascript:say('go ` + x + `')">` + x + `</a>` + "\n"
	}
	return `
<link href="style.css" rel="stylesheet" type="text/css"/>

<div id="links">
  <ul class="navbar">
    <li>` + link("Home") + ` </li>
    <li>` + link("Analyze bids") + ` </li>
    <li>` + link("Bid fourth hand") + ` </li>
    <li>` + link("Bid for me") + ` </li>
    <li>` + link("Settings") + ` </li>
  </ul>
</div>
`
}
