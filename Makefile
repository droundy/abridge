# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=github.com/droundy/bridge

GOFILES=\
	suit.go\
	hand.go\
	table.go\
	ensemble.go\
	conventions.go\
	bridge.go\
	opening.go\
	overcall.go\
	notrumpopening.go\
	response.go\
	rebid.go\
	natural.go\
	ensemblecache.go\
	conventioncard.go\
	slam.go\

include $(GOROOT)/src/Make.pkg

demo: demo.go $(pkgdir)/$(TARG).a
	$(GC) -o demo.$(O) demo.go
	$(LD) -o demo demo.$(O)

$(pkgdir)/$(TARG)/speech.a: $(pkgdir)/$(TARG).a speech/*.go
	cd speech && make install

abridge/server: $(pkgdir)/$(TARG).a $(pkgdir)/$(TARG)/speech.a abridge/*.go
	cd abridge && make

testing/server: $(pkgdir)/$(TARG).a $(pkgdir)/$(TARG)/speech.a testing/*.go
	cd testing && make

benchmark: benchmark.go $(pkgdir)/$(TARG).a
	$(GC) -o benchmark.$(O) benchmark.go
	$(LD) -o benchmark benchmark.$(O)
