package main

import (
	"fmt"
	"http"
)

func link(c *http.Conn, req *http.Request, url, label string) {
	if req.URL.Path == url {
		fmt.Fprintf(c, `<font color="#aaaaaa">%s</font>`, label)
	} else {
		fmt.Fprintf(c, `<a href="%s">%s</a>`, url, label)
	}
}

func header(c *http.Conn, req *http.Request, title string) (footer func()) {
	c.SetHeader("Content-Type", "text/html")
	fmt.Fprintf(c, `
<html>
<head>

<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>%s</title>

</head>

<body>`, title)
	link(c, req, "/", "Analyze bids")
	fmt.Fprint(c, ` | `)
	link(c, req, "/bidder", "Bid third hand")
	fmt.Fprint(c, ` | `)
	link(c, req, "/about", "About aBridge")
	fmt.Fprintln(c, `<br/>`)
	fmt.Fprintln(c, `<br/>`)

	return func() {
		// This is the footer... which is intended to be deferred.
		fmt.Fprintln(c, `</body></html>`)
	}
}

