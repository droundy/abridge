package main

import (
	"http"
	"json"
	"github.com/droundy/bridge"
)

type Settings struct {
	Style string
	WhichCard string
	Cards map[string]*bridge.ConventionCard
}
func (s Settings) Card() *bridge.ConventionCard {
	if c,ok := s.Cards[s.WhichCard]; ok {
		return c
	}
	c := bridge.DefaultConvention()
	s.Cards[s.WhichCard] = &c
	return &c
}

func DefaultSettings() Settings {
	return Settings{
	Style: "two color",
	WhichCard: "New convention card",
	Cards: make(map[string]*bridge.ConventionCard),
	}
}

func getSettings(req *http.Request) (p Settings) {
	req.ParseForm()
	prefstr, _ := req.Header["Cookie"] // I don't care about errors!
	return readCookie(prefstr)
}

func readCookie(cookie string) (p Settings) {
	p = DefaultSettings()
	p.Cards = make(map[string]*bridge.ConventionCard)
	// fmt.Println("Trying to unmarshall string", cookie)
	json.Unmarshal([]byte(cookie), &p) // I don't care about errors!
	for k,v := range bridge.DefaultConvention().Pts {
		if _,ok := p.Card().Pts[k]; !ok {
			p.Card().Pts[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention().Options {
		if _,ok := p.Card().Options[k]; !ok {
			p.Card().Options[k] = v
		}
	}
	for k,v := range bridge.DefaultConvention().Radio {
		if _,ok := p.Card().Radio[k]; !ok {
			p.Card().Radio[k] = v
		}
	}
	// fmt.Println("Unmarshal gave: %v", e)
	return p
}

func (s *Settings) Write() string {
	bytes,_ := json.Marshal(*s)
	return string(bytes)
}

func (s *Settings) SetScript() string {
	return `
<script type="text/javascript">
  alert('I set the cookie');
  createCookie('WebSocketCookie', '` + s.Write() + `', 365);
  say('I just set the cookie');
</script>
`
}

func (s *Settings) Set(c http.ResponseWriter) {
	c.SetHeader("Set-Cookie", s.Write())
}
