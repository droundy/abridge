package main

import (
	"fmt"
	"http"
	"json"
	"strings"
	"github.com/droundy/bridge"
)

type TransitoryData struct {
	Bids string
	Hands bridge.Table
	Dealer bridge.Seat
	Bidfor bridge.Seat
	AmBidding bool
	NScard, EWcard string
	Url string
}

func getTransitoryData(req *http.Request) (d *TransitoryData) {
	d = new(TransitoryData)
	req.ParseForm()
	if xx,ok := req.Form["transitorydata"]; ok {
		json.Unmarshal([]byte(strings.Replace(xx[0],"'","\"",-1)), &d) // I don't care about errors!
	}
	d.Url = req.URL.Path
	return
}

func (t *TransitoryData) Save(c http.ResponseWriter) {
	bytes,_ := json.Marshal(t)
	fmt.Fprintf(c, `<input type="hidden" name="transitorydata" value="%s" />`, strings.Replace(string(bytes),"\"","'",-1))
	fmt.Fprintln(c)
}
