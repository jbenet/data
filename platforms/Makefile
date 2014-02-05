
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

all: $(ARCHIVES)

%s: %s/data

linux_%/data:
	-rm $@
	cd $(dir $@) && \
	  vagrant up && \
	  vagrant ssh -c "source ~/.bashrc; cd data; make deps; make;" && \
	  vagrant suspend

darwin_%/data:
	cd ../ && make

windows_%/data:
	$(error not implemented)

archives/data-v$(VERSION)-%.tar.gz: %/data
	./archive.py $(<D) $(VERSION)

bins: $(BINS)

clean:
	rm -f -- archives/*
	rm -f -- $(BINS)
