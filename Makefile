all: kiwi kiwi-vmclient

debug: kiwi-debug kiwi-vmclient-debug

kiwi:
	$(MAKE) -C src kiwi
	$(MAKE) -C src kiwi-vmserver

kiwi-debug:
	$(MAKE) -C src kiwi-debug
	$(MAKE) -C src kiwi-vmserver-debug

kiwi-vmclient:
	$(MAKE) -C src kiwi-vmclient

kiwi-vmclient-debug:
	$(MAKE) -C src kiwi-vmclient-debug

test: test-scanner

test-scanner:
	$(MAKE) -C test scanner

clean:
	$(MAKE) -C src clean
	$(MAKE) -C test clean

