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
	c.SetHeader("Content-Type", "application/xhtml+xml")
	fmt.Fprintf(c, `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN"
  "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">

<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<head>

<link href="style.css" rel="stylesheet" type="text/css"/>
<meta http-equiv="content-type" content="text/html; charset=utf-8"/>
<title>%s</title>

</head>

<body id="body">`, title)
	fmt.Fprintln(c, `

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
		fmt.Fprintln(c, `

<svg xmlns="http://www.w3.org/2000/svg" version="1.1"  
viewBox="0 0 100 100" preserveAspectRatio="xMidYMid slice"  
style="width:100%; height:100%; position:absolute; top:0; left:0; z-index:-1;">  
<linearGradient id="gradient">  
<stop stop-color="#ffeeff" offset="0%"/>  
<stop stop-color="#ffffee" offset="100%"/>  
</linearGradient>  
<rect x="0" y="0" width="100" height="100" style="fill:url(#gradient)" />  
<circle cx="50" cy="50" r="30" style="fill:url(#gradient)" />  
</svg>

</body></html>`)
	}
}

/*

*/
