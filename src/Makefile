include ../Makefile.inc

INCLUDE+=-I../include

LDFLAGS=-licuio -licuuc
LDFLAGS_DBG=-licuio -licuuc

all: parser
	$(CC) $(CFLAGS) -o kiwi scanner.o y.tab.o main.c $(LDFLAGS)

debug: parser-debug
	$(CC) $(CFLAGS_DBG) -o kiwi-debug scanner.o y.tab.o main.c $(LDFLAGS_DBG)

parser: scanner
	$(CC) $(CFLAGS) -c y.tab.c

parser-debug: scanner-debug
	$(CC) $(CFLAGS_DBG) -c y.tab.c

scanner: ytab
	$(CC) $(CFLAGS) -c scanner.c

scanner-debug: 
	$(CC) $(CFLAGS_DBG) -c scanner.c

ytab:
	byacc -P -d -p scanner_ kiwi.y
	mv y.tab.h ../include/y.tab.h

ytab-debug:
	byacc -P -d -p scanner_ -v -t kiwi.y
	mv y.tab.h ../include/y.tab.h

clean: 
	rm -Rf *.o y.tab.* y.output ../include/y.tab.h kiwi kiwi-debug
