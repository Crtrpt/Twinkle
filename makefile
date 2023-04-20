all:default_config build build_test_server

build:
	go build -race -ldflags "-s -w"  -gcflags "-m -l"  -o bin/twinkle  cmd/twinkle/main.go
	chmod +x bin/twinkle

build_test_server:
	go build -ldflags "-s -w" -o bin/twinkle_test_server  cmd/twinkle_test_server/main.go
	chmod +x bin/twinkle_test_server

default_config:
	mkdir -p bin/key 
	mkdir -p bin/conf 
	cp conf/app.toml bin/conf/app.toml

run:
	go run -race -ldflags "-s -w"  -gcflags "-m -l"   cmd/twinkle/main.go