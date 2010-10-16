package main

func (dat *ClientData) Header() string {
	link := func(x string) string {
		if x == dat.MyPage {
			return " " + x + "\n"
		}
		return `<a href="javascript:say('go ` + x + `')">` + x + `</a>` + "\n"
	}
	stylesheet := "style.css"
	switch dat.Cookie.Style {
	case "four color": stylesheet = "style-fourcolor.css"
	}
	return `
<link href="`+stylesheet+`" rel="stylesheet" type="text/css"/>

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
