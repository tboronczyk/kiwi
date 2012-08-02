# build kiwi and tests
all: kiwi test

# build kiwi for debugging
debug: kiwi-debug

# clean
clean: kiwi-clean test-clean

kiwi:
	$(MAKE) -C src all

kiwi-debug:
	$(MAKE) -C src debug

kiwi-clean: 
	$(MAKE) -C src clean

test: test-scanner

test-scanner:
	$(MAKE) -C test scanner

test-clean:
	$(MAKE) -C test clean
