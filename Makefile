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

include $(GOROOT)/src/Make.pkg

demo: demo.go $(pkgdir)/$(TARG).a
	$(GC) -o demo.$(O) demo.go
	$(LD) -o demo demo.$(O)

abridge/server: $(pkgdir)/$(TARG).a abridge/*.go
	cd abridge && make clean && make

benchmark: benchmark.go $(pkgdir)/$(TARG).a
	$(GC) -o benchmark.$(O) benchmark.go
	$(LD) -o benchmark benchmark.$(O)
