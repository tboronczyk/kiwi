AM_CFLAGS = -std=c99 -I$(top_builddir)/include -I$(top_builddir)/src
if DEBUG
AM_CFLAGS += -Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused \
 -Wfloat-equal -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast \
 -Wlogical-op -Wmissing-prototypes -Wmissing-declarations -pedantic -DDEBUG -g
else
AM_CFLAGS += -O2
endif

TESTS = check_scanner

check_PROGRAMS = check_scanner
check_scanner_SOURCES = check_scanner.c
check_scanner_CFLAGSADD = @CHECK_CFLAGS@
check_scanner_LDADD = $(top_builddir)/src/scanner.o -licuio -licuuc \
 @CHECK_LIBS@ -lpthread

