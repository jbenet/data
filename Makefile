all: build

deps:
	go get ./...

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

watch:
	-make install
	@echo "[watching *.go;*.html for recompilation]"
	# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="*.go;" --recursive \
		--command='make install' .
