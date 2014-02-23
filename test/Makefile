include ../Makefile.inc

INCLUDE+=-I../include -I../src

LDFLAGS_DBG=-licuio -licuuc

all: scanner

symtable:
	$(MAKE) -C ../src symtable-debug
	$(CC) $(CFLAGS_DBG) -o symtable ../src/symtable.o symtable.c $(LDFLAGS_DBG)

scanner: ytab
	$(MAKE) -C ../src scanner-debug
	$(CC) $(CFLAGS_DBG) -o scanner ../src/scanner.o scanner.c $(LDFLAGS_DBG)

ytab:
	$(MAKE) -C ../src ytab-debug

clean:

	rm -Rf symtable scanner

