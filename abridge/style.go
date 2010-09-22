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
    font-family: arial,helvetica,"sans serif";
    font-size: 12pt;
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
font-size:12pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 12pt;
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

.bridgetable {
  padding: 5px;
  width: 500px;
}

.bridgehand {
  float: right;
  color: #666666;
  font-family: serif;
  font-size: 10pt;
}

.bridgecards {
  color: #000000;
  font-family: monospace;
  font-size: 12pt;
}

#conventions {
clear: right;
float: right;
width: 300px;
padding: 5px;
/* margin: 0 0 0 220px; */
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

#logo {
width:100%;
height:100%;
position:absolute;
top:0; left:0;
z-index:-1;"
}

`)
}
