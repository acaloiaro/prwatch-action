all: compile
.PHONY : all

compile:
	CGO=0 GOOS=linux go build -o ../bin/prwatch prwatch.go

mac:
	CGO=0 GOOS=darwin go build -o ../bin/prwatch prwatch.go

clean:
	rm -rf ../bin/*
