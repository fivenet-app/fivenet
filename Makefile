.DEFAULT: run-server

.PHONY: run-server
run-server:
	swag init
	go run . server

.PHONY: build
build:
	swag init
	go build -o rphub .

.PHONY: watch-server
watch-server:
	while true; do { git ls-files; git ls-files . --exclude-standard --others; } | entr -d $(MAKE) run-server; sleep .75; done
