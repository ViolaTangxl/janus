.PHONY: all
all: build

.PHONY: build
build:
	cd app; go build -o janusd

.PHONY: run
run: build
	cd app; ./janusd -conf ../config/config.yml
