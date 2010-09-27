package bridge

var cache_ensemble = make(chan struct{ bid string; e *Ensemble; cc *ConventionCard })
var seek_ensemble = make(chan struct{ bid string; cc *ConventionCard; answer chan<- *Ensemble })

func init() {
	go func() {
		mycache := make(map[string]struct{e *Ensemble; cc *ConventionCard})
		for {
			select {
			case wr := <- cache_ensemble:
				if wr.bid == "" && wr.e == nil {
					mycache = make(map[string]struct{e *Ensemble; cc *ConventionCard})
				} else {
					mycache[wr.bid] = struct{e *Ensemble; cc *ConventionCard}{wr.e, wr.cc}, wr.e != nil
				}
			case r := <- seek_ensemble:
				if e,ok := mycache[r.bid]; ok && e.cc.SameAs(r.cc) {
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
						mycache[k] = struct{e *Ensemble; cc *ConventionCard}{nil, nil}, false
						toclear--
					}
				}
			}
		}
	}()
}

func cacheEnsemble(bid string, cc ConventionCard, e *Ensemble) {
	cache_ensemble <- struct{ bid string; e *Ensemble; cc *ConventionCard }{bid, e, &cc}
}

func lookupEnsembleFromCache(bid string, cc ConventionCard) (*Ensemble, bool) {
	c := make(chan *Ensemble)
	seek_ensemble <- struct{ bid string; cc *ConventionCard; answer chan<- *Ensemble }{bid, &cc, c}
	e := <- c
	return e, e != nil
}

func ClearBid(bid string) {
	cacheEnsemble(bid, DefaultConvention(), nil)
}

func ClearCache() {
	cacheEnsemble("", DefaultConvention(), nil)
}
