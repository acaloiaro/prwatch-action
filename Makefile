all: compile
.PHONY : all

compile:
	$(MAKE) -C go all

mac:
	$(MAKE) -C go mac

clean:
	rm -rf bin/*

