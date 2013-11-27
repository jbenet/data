


build:
	go build
	cd data && go build

install: build
	go install
	cd data && go install

pkg:
	go build
	go install

tool:
	cd data && go build
	cd data && go install
