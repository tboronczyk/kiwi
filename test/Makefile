include ../Makefile.inc

INCLUDE+=-I../include -I../src

LDFLAGS_DBG=-licuio -licuuc

all: scanner

symtab:
	$(MAKE) -C ../src symtab-debug
	$(CC) $(CFLAGS_DBG) -o symtab ../src/symtab.o symtab.c $(LDFLAGS_DBG)

scanner: ytab
	$(MAKE) -C ../src scanner-debug
	$(CC) $(CFLAGS_DBG) -o scanner ../src/scanner.o scanner.c $(LDFLAGS_DBG)

ytab:
	$(MAKE) -C ../src ytab-debug

clean:

	rm -Rf symtab scanner

