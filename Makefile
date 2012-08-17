all: kiwi vmclient

debug: kiwi-debug vmclient-debug

kiwi:
	$(MAKE) -C src kiwi
	$(MAKE) -C src vmserver

kiwi-debug:
	$(MAKE) -C src kiwi-debug
	$(MAKE) -C src vmserver-debug

vmclient:
	$(MAKE) -C src vmclient

vmclient-debug:
	$(MAKE) -C src vmclient-debug

test: test-scanner

test-scanner:
	$(MAKE) -C test scanner

clean:
	$(MAKE) -C src clean
	$(MAKE) -C test clean

