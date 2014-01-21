
OS=$(shell uname -s)
ifeq ($(OS),Darwin)
OS=darwin
else ifeq ($(OS),Linux)
OS=linux
#else ifeq (windows check)
#	OS=windows
else
$(error os seemingly not supported yet)
endif

ARCH=$(shell uname -m)
ifeq ($(ARCH),x86_64)
ARCH=amd64
else ifeq ($(ARCH),i686)
ARCH=386
else ifeq ($(ARCH),i386)
ARCH=386
else ifeq ($(ARCH),amd64)
else ifeq ($(ARCH),386)
else
$(error arch seemingly not supported yet.)
endif


all: build

deps:
	go get ./...

build:
	go build
	cd data && go build && cd ..
	cp data/data platforms/$(OS)_$(ARCH)/data

install: build
	go install
	cd data && go install

watch:
	-make install
	@echo "[watching *.go;*.html for recompilation]"
	# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="*.go;" --recursive \
		--command='make install' .
