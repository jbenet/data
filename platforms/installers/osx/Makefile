ifneq ($(shell uname),Darwin)
	$(error "Making osx installers is only supported in osx.")
endif

all: pkg

# for now, amd64. lipo both together later.
BIN=../../darwin_amd64/data
VERSION=$(shell $(BIN) version --number)
PKGNAME=data-v$(VERSION)-OSX-Installer

# setup the package contents
root: root/usr/bin/data

root/usr/bin/data: $(BIN)
	-mkdir -p $(dir $@)
	cp $< $@

# build the packages
data.pkg: root root/usr/bin/data
	pkgbuild \
		--root root \
		--identifier io.datadex.data \
		--version $(VERSION) \
		--ownership recommended \
		data.pkg

$(PKGNAME).pkg: data.pkg
	productbuild \
		--distribution distribution.xml \
		--resources resources \
		--package-path data.pkg \
		--version $(VERSION) \
		$(PKGNAME).pkg

pkg: $(PKGNAME).pkg

clean:
	rm -rf -- root/
	rm -f -- *.pkg
