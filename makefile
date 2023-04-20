all:default_config build run

build:
	go build -ldflags "-s -w" -o bin/twinkle  cmd/twinkle/main.go
	chmod +x bin/twinkle

default_config:
	mkdir -p bin/key 
	mkdir -p bin/conf 
	cp conf/app.toml bin/conf/app.toml

run:
	cd bin && ./twinkle