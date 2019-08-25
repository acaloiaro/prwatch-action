all: test compile
.PHONY : all

compile:
	$(MAKE) -C go all

mac:
	$(MAKE) -C go mac

test:
	$(MAKE) -C go test

clean:
	rm -rf bin/*
