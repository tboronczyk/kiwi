CC=gcc

INCLUDE=-I../include

WFLAGS=-Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused -Wfloat-equal \
 -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast -Wlogical-op \
 -Wmissing-prototypes -Wmissing-declarations -Wmissing-noreturn \
 -Wmissing-format-attribute

CFLAGS=-O2 -std=c99 $(INCLUDE)
CFLAGS_DBG=-O2 -std=c99 -pedantic $(WFLAGS) -DDEBUG -g $(INCLUDE)

LDFLAGS=-licuio -licuuc
LDFLAGS_DBG=-licuio -licuuc

all: parser
	$(CC) $(CFLAGS) -o kiwi  token.o scanner.o ast.o parser.o main.c $(LDFLAGS)

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
