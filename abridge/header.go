package main

import (
	"fmt"
	"http"
)

func link(c *http.Conn, req *http.Request, url, label string) {
	if req.URL.Path == url {
		//fmt.Fprintf(c, `<font color="#aaaaaa">%s</font>`, label)
		fmt.Fprintf(c, `<a class="x">%s</a>`, label)
	} else {
		fmt.Fprintf(c, `<a href="%s">%s</a>`, url, label)
	}
}

func header(c *http.Conn, req *http.Request, title string) (footer func()) {
	c.SetHeader("Content-Type", "text/html")
	fmt.Fprintf(c, `
<html>
<head>

<link rel="stylesheet" href="style.css">
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>%s</title>

</head>

<body>`, title)
	fmt.Fprintln(c, `<div id="everything">
<div id="links">
<ul class="navbar"><li>`)
	link(c, req, "/", "Analyze bids")
	fmt.Fprintln(c, `</li><li>`)
	link(c, req, "/bidder", "Bid fourth hand")
	fmt.Fprintln(c, `</li><li>`)
	link(c, req, "/about", "About aBridge")
	fmt.Fprintln(c, `</li></ul>`)
	fmt.Fprintln(c, `</div>`)

	return func() {
		// This is the footer... which is intended to be deferred.
		fmt.Fprintln(c, `</div></body></html>`)
	}
}

