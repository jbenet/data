
VERSION=$(shell data version --number)

PLATFORMS= \
	darwin_amd64 \
	linux_386 \
	linux_amd64 \
	# darwin_386 \
	# windows_386 \
	# windows_amd64 \

BINS=$(addsuffix /data,$(PLATFORMS))

ARCHIVES=$(patsubst %,archives/data-v$(VERSION)-%.tar.gz,$(PLATFORMS))

OTHER= \
	installers/osx/data-v$(VERSION)-OSX-Installer.pkg

all: $(ARCHIVES) $(OTHER)

%s: %s/data

linux_%/data:
	-rm $@
	cd $(dir $@) && \
	  vagrant up && \
	  vagrant ssh -c "source ~/.bashrc; cd data; make deps; make;" && \
	  vagrant suspend

darwin_%/data:
	cd ../ && $(MAKE)

windows_%/data:
	$(error not implemented)

archives/data-v$(VERSION)-%.tar.gz: %/data
	./archive.py $(<D) $(VERSION)

installers/osx/data-v$(VERSION)-OSX-Installer.pkg: darwin_amd64/data
	cd installers/osx && $(MAKE)

bins: $(BINS)

clean:
	rm -f -- archives/*
	rm -f -- $(BINS)
