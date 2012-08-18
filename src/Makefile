include ../Makefile.inc

INCLUDE+=-I../include

LDFLAGS=-licuio -licuuc
LDFLAGS_DBG=-licuio -licuuc

all: kiwi vm vmclient vmserver

kiwi: parser vm
	$(CC) $(CFLAGS) -o kiwi scanner.o y.tab.o main.c $(LDFLAGS)

kiwi-debug: parser-debug vm-debug
	$(CC) $(CFLAGS_DBG) -o kiwi-debug scanner.o y.tab.o main.c $(LDFLAGS_DBG)

vm: vmprogbuf vmmach
	$(CC) $(CFLAGS) -o vm vminstr.o vmprogbuf.o vmop.o vmmach.o vm.c $(LDFLAGS)

vm-debug: vmprogbuf-debug vmmach-debug
	$(CC) $(CFLAGS_DBG) -o vm-debug vminstr.o vmprogbuf.o vmop.o vmmach.o vm.c $(LDFLAGS_DBG)

vmprogbuf: vminstr
	$(CC) $(CFLAGS) -c vmprogbuf.c

vmprogbuf-debug: vminstr-debug
	$(CC) $(CFLAGS_DBG) -c vmprogbuf.c

vminstr:
	$(CC) $(CFLAGS) -c vminstr.c

vminstr-debug:
	$(CC) $(CFLAGS_DBG) -c vminstr.c

vmmach: vmop
	$(CC) $(CFLAGS) -c vmmach.c

vmmach-debug: vmop-debug
	$(CC) $(CFLAGS_DBG) -c vmmach.c

vmop:
	$(CC) $(CFLAGS) -c vmop.c

vmop-debug:
	$(CC) $(CFLAGS_DBG) -c vmop.c

vmclient: mysocket
	$(CC) $(CFLAGS) -o vmclient mysocket.o vmclient.c $(LDFLAGS)

vmclient-debug: mysocket-debug
	$(CC) $(CFLAGS_DBG) -o vmclient mysocket.o vmclient.c $(LDFLAGS_DBG)

vmserver: mysocket
	$(CC) $(CFLAGS) -o vmserver mysocket.o vmserver.c $(LDFLAGS)

vmserver-debug: mysocket-debug
	$(CC) $(CFLAGS_DBG) -o vmserver mysocket.o vmserver.c $(LDFLAGS_DBG)

parser: scanner
	$(CC) $(CFLAGS) -c y.tab.c

parser-debug: scanner-debug
	$(CC) $(CFLAGS_DBG) -c y.tab.c

scanner: ytab
	$(CC) $(CFLAGS) -c scanner.c

scanner-debug: ytab-debug
	$(CC) $(CFLAGS_DBG) -c scanner.c

ytab:
	byacc -P -d -p scanner_ kiwi.y
	mv y.tab.h ../include/y.tab.h

ytab-debug:
	byacc -P -d -p scanner_ -v -t kiwi.y
	mv y.tab.h ../include/y.tab.h

mysocket:
	$(CC) $(CFLAGS) -c mysocket.c 

mysocket-debug:
	$(CC) $(CFLAGS_DBG) -c mysocket.c

clean: 
	rm -Rf  ../include/y.tab.h y.tab.c y.output *.o kiwi kiwi-debug \
	 vm vm-debug vmserver vmserver-debug vmclient vmclient-debug

