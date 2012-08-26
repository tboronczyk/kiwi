all: kiwi

debug: kiwi-debug

kiwi:
	$(MAKE) -C src kiwi

kiwi-debug:
	$(MAKE) -C src kiwi-debug

test: test-scanner

test-scanner:
	$(MAKE) -C test scanner

clean:
	$(MAKE) -C src clean
	$(MAKE) -C test clean

