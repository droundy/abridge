package bridge

var cache_ensemble = make(chan struct{ bid string; e *Ensemble })
var seek_ensemble = make(chan struct{ bid string; answer chan<- *Ensemble })

func init() {
	go func() {
		mycache := make(map[string]*Ensemble)
		for {
			select {
			case wr := <- cache_ensemble:
				if wr.bid == "" && wr.e == nil {
					mycache = make(map[string]*Ensemble)
				} else {
					mycache[wr.bid] = wr.e, wr.e != nil
				}
			case r := <- seek_ensemble:
				if e,ok := mycache[r.bid]; ok {
					r.answer <- e
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
						mycache[k] = nil, false
						toclear--
					}
				}
			}
		}
	}()
}

func cacheEnsemble(bid string, e *Ensemble) {
	cache_ensemble <- struct{ bid string; e *Ensemble }{bid, e}
}

func lookupEnsembleFromCache(bid string) (*Ensemble, bool) {
	c := make(chan *Ensemble)
	seek_ensemble <- struct{ bid string; answer chan<- *Ensemble }{bid, c}
	e := <- c
	return e, e != nil
}

func ClearBid(bid string) {
	cacheEnsemble(bid, nil)
}

func ClearCache() {
	cacheEnsemble("", nil)
}
