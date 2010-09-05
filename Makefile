# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=github.com/droundy/bridge

GOFILES=\
        suit.go\
        hand.go\
        table.go\
        conventions.go\
        bridge.go\
	opening.go\

include $(GOROOT)/src/Make.pkg

# ifneq ($(strip $(shell which gotgo)),)
# pkg/slice.go: $(srcpkgdir)/gotgo/slice.got
# 	gotgo --package-name goopt -o "$@" "$<" string
# endif

demo: demo.go install
	$(GC) -o demo.$(O) demo.go
	$(LD) -o demo demo.$(O)
