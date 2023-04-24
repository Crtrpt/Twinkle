all:default_config build build_test_server build_forward

build_forward:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -gcflags "-m -l"  -o bin/twinkle_udp_forward  cmd/twinkle_udp_forward/main.go
	chmod +x bin/twinkle_udp_forward


build:
	go build -race -ldflags "-s -w"  -gcflags "-m -l"  -o bin/twinkle  cmd/twinkle/main.go
	chmod +x bin/twinkle

build_doc:
	mdbook build doc
	mkdir -p bin/doc
	cp -rf doc/book/* bin/doc

build_test_server:
	go build -ldflags "-s -w" -o bin/twinkle_http_server  cmd/twinkle_http_server/main.go
	chmod +x bin/twinkle_http_server

default_config:
	mkdir -p bin/key 
	mkdir -p bin/conf 
	cp conf/app.toml bin/conf/app.toml

run:
	go run -race -ldflags "-s -w"  -gcflags "-m -l"   cmd/twinkle/main.go