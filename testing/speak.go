package main

import (
	"http"
	"fmt"
	"github.com/droundy/abridge/speech"
)

func wavServer(c http.ResponseWriter, req *http.Request) {
	c.SetHeader("Content-Type", "audio/x-wav")
	text := req.URL.Path[:len(req.URL.Path)-4]
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
