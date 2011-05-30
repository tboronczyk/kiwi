CC=gcc

INCLUDE=-I../include

WFLAGS=-Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused -Wfloat-equal \
 -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast -Wlogical-op \
 -Wmissing-prototypes -Wmissing-declarations -Wmissing-noreturn \
 -Wmissing-format-attribute

CFLAGS=-O2 -std=c99 $(INCLUDE)
CFLAGS_DBG=-O2 -std=c99 -pedantic $(WFLAGS) -DDEBUG -g $(INCLUDE)

LDFLAGS=-licuio
LDFLAGS_DBG=-licuio

all: scanner
	$(CC) $(LDFLAGS) $(CFLAGS) -o kiwi token.o scanner.o main.c

debug: scanner-debug
	$(CC) $(LDFLAGS_DBG) $(CFLAGS_DBG) -o kiwi token.o scanner.o main.c

scanner: 
	$(CC) $(CFLAGS) -c token.c scanner.c

scanner-debug: 
	$(CC) $(CFLAGS_DBG) -c token.c scanner.c

clean: 
	rm -Rf *.o kiwi

