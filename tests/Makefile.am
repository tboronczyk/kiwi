AM_CFLAGS = -std=c99 -I$(top_builddir)/include -I$(top_builddir)/src
if DEBUG
AM_CFLAGS += -Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused \
 -Wfloat-equal -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast \
 -Wlogical-op -Wmissing-prototypes -Wmissing-declarations -pedantic -DDEBUG -g
else
AM_CFLAGS += -O2
endif

TESTS = check_scanner \
        check_parser

check_PROGRAMS = check_scanner \
                 check_parser

check_scanner_SOURCES = $(top_builddir)/src/scanner.c check_scanner.c
check_scanner_CFLAGSADD = @CHECK_CFLAGS@
check_scanner_LDADD = -licuio -licuuc @CHECK_LIBS@ -lpthread

check_parser_SOURCES = $(top_builddir)/src/ast.c \
                       $(top_builddir)/src/parser.c \
                       $(top_builddir)/src/scanner.c \
                       astdump.c \
                       check_parser.c
check_parser_CFLAGSADD = @CHECK_CFLAGS@
check_parser_LDADD = -licuio -licuuc @CHECK_LIBS@ -lpthread

