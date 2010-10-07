package main

import (
	"http"
	"fmt"
	"github.com/droundy/bridge/speech"
)

func wavServer(c http.ResponseWriter, req *http.Request) {
	c.SetHeader("Content-Type", "audio/x-wav")
	dat := getTransitoryData(req)
	text := dat.Url[:len(dat.Url)-4]
	if len(text) > 8 {
		text = text[8:]
	}
	x,err := speech.Speak(text)
	if err == nil {
		fmt.Fprint(c, x)
	} else {
		fmt.Println("Error speaking:", err)
	}
}
