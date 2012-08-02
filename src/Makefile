include ../Makefile.inc

INCLUDE+=-I../include

LDFLAGS=-licuio -licuuc
LDFLAGS_DBG=-licuio -licuuc

all: parser
	$(CC) $(CFLAGS) -o kiwi token.o scanner.o ast.o parser.o main.c $(LDFLAGS)

debug: parser-debug
	$(CC) $(CFLAGS_DBG) -o kiwi-debug token.o scanner.o ast.o \
 parser.o main.c $(LDFLAGS_DBG)

parser: scanner
	$(CC) $(CFLAGS) -c ast.c parser.c

parser-debug: scanner-debug
	$(CC) $(CFLAGS_DBG) -c ast.c parser.c

scanner: 
	$(CC) $(CFLAGS) -c token.c scanner.c

scanner-debug: 
	$(CC) $(CFLAGS_DBG) -c token.c scanner.c

clean: 
	rm -Rf *.o kiwi kiwi-debug
