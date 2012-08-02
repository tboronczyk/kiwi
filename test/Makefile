include ../Makefile.inc

INCLUDE+=-I../include

LDFLAGS_DBG=-licuio -licuuc

all: scanner

scanner:
	$(MAKE) -C ../src scanner-debug
	$(CC) $(CFLAGS_DBG) -o scanner  ../src/token.o ../src/scanner.o scanner.c $(LDFLAGS_DBG)

clean:
	rm -Rf *.o scanner
