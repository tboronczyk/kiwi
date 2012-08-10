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
	byacc -d kiwi.y

ytab-debug:
	byacc -d -v -t kiwi.y

clean: 
	rm -Rf *.o y.tab.* y.output kiwi kiwi-debug
