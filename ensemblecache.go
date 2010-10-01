package bridge

var cache_ensemble = make(chan struct{ bid string; e *Ensemble; cc *[2]ConventionCard })
var seek_ensemble = make(chan struct{ bid string; cc *[2]ConventionCard; answer chan<- *Ensemble })

func init() {
	go func() {
		mycache := make(map[string]struct{e *Ensemble; cc *[2]ConventionCard})
		for {
			select {
			case wr := <- cache_ensemble:
				if wr.bid == "" && wr.e == nil {
					mycache = make(map[string]struct{e *Ensemble; cc *[2]ConventionCard})
				} else {
					mycache[wr.bid] = struct{e *Ensemble; cc *[2]ConventionCard}{wr.e, wr.cc}, wr.e != nil
				}
			case r := <- seek_ensemble:
				if e,ok := mycache[r.bid]; ok && e.cc[0].SameAs(&r.cc[0]) && e.cc[1].SameAs(&r.cc[1]) {
					r.answer <- e.e
				} else {
					r.answer <- nil
				}
			}
			if len(mycache) > 4096 {
				// Avoid leaking resources forever...
				toclear := 1024
				for k := range mycache {
					if toclear == 0 {
						break
					}
					if len(k) > 8 {
						// Only clear out longer bidding sequences, which will
						// normally be more rare.
						mycache[k] = struct{e *Ensemble; cc *[2]ConventionCard}{nil, nil}, false
						toclear--
					}
				}
			}
		}
	}()
}

func cacheEnsemble(bid string, cc [2]ConventionCard, e *Ensemble) {
	cache_ensemble <- struct{ bid string; e *Ensemble; cc *[2]ConventionCard }{bid, e, &cc}
}

func lookupEnsembleFromCache(bid string, cc [2]ConventionCard) (*Ensemble, bool) {
	c := make(chan *Ensemble)
	seek_ensemble <- struct{ bid string; cc *[2]ConventionCard; answer chan<- *Ensemble }{bid, &cc, c}
	e := <- c
	return e, e != nil
}

func ClearBid(bid string) {
	cacheEnsemble(bid, DefaultConventions(), nil)
}

func ClearCache() {
	cacheEnsemble("", DefaultConventions(), nil)
}
