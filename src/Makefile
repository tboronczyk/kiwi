include ../Makefile.inc

INCLUDE+=-I../include

LDFLAGS=-licuio -licuuc
LDFLAGS_DBG=-licuio -licuuc

all: kiwi vm

kiwi: parser vm
	$(CC) $(CFLAGS) -o kiwi ast.o scanner.o symtable.o y.tab.o main.c $(LDFLAGS)

kiwi-debug: parser-debug vm-debug
	$(CC) $(CFLAGS_DBG) -o kiwi-debug ast.o scanner.o symtable.o y.tab.o main.c $(LDFLAGS_DBG)

vm:
	$(CC) $(CFLAGS) -o vm vm.c $(LDFLAGS)

vm-debug:
	$(CC) $(CFLAGS_DBG) -o vm-debug vm.c $(LDFLAGS_DBG)

parser: symtable scanner ast
	$(CC) $(CFLAGS) -c y.tab.c

parser-debug: symtable-debug scanner-debug ast-debug
	$(CC) $(CFLAGS_DBG) -c y.tab.c

symtable:
	$(CC) $(CFLAGS) -c symtable.c

symtable-debug:
	$(CC) $(CFLAGS_DBG) -c symtable.c

scanner: ytab
	$(CC) $(CFLAGS) -c scanner.c

scanner-debug: ytab-debug
	$(CC) $(CFLAGS_DBG) -c scanner.c

ast:
	$(CC) $(CFLAGS) -c ast.c

ast-debug:
	$(CC) $(CFLAGS_DBG) -c ast.c

ytab:
	byacc -P -d -p yy kiwi.y
	mv y.tab.h ../include/y.tab.h

ytab-debug:
	byacc -P -d -p yy -v -t kiwi.y
	mv y.tab.h ../include/y.tab.h

clean: 
	rm -Rf  ../include/y.tab.h y.tab.c y.output *.o kiwi kiwi-debug \
	vm vm-debug
