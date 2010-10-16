package main

import (
	"bytes"
)

func AnalyzeBids(dat *ClientData, evt []string) string {
	buf := bytes.NewBuffer(make([]byte, 0, 4*1024))
	dat.bidbox(buf, evt)
	return buf.String();
}
