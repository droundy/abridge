package main

import (
	"fmt"
  "http"
)

func styleServer(c *http.Conn, req *http.Request) {
	c.SetHeader("Content-Type", "text/css")
	fmt.Fprint(c, `
html {
    margin: 0;
    padding: 0;
}

body {
    margin: 0;
    padding: 0;
    background: #ffffff;
}
h1 {
font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 16pt;
}
h2 { font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 14pt;
}
p {
font-family: arial,helvetica,"sans serif";
font-size:10pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 10pt;
}
a {
color: #555599;
}

#bidbox {
float: left;
}

.analysis {
float: left;
padding: 5px;
}

#bidtable {
clear:right;
float: right;
padding: 5px;
}

#conventions {
clear: right;
float: right;
width: 300px;
padding: 5px;
/* margin: 0 0 0 220px; */
}

.bridgetable {
clear: left;
float: left;
padding: 5px;
}

.navbar {
list-style: none;
}
.navbar li {
padding: 2px;
display: inline;
}

#links {
background: #eeeeff;
float: right;
/* width: 100px; */
margin: 0;
padding: 4px;
}

#links a.x {
/* background: #000000; */
color: #000000;
font-weight: bold;
}

`)
}
