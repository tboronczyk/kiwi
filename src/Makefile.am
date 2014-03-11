AM_CFLAGS = -std=c99 -I$(top_builddir)/include 
if DEBUG
AM_CFLAGS += -Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused -Wfloat-equal \
 -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast -Wlogical-op -Wmissing-prototypes \
 -Wmissing-declarations -pedantic -DDEBUG -g
else
AM_CFLAGS += -O2
endif

AM_YFLAGS = -d

BUILT_SOURCES = parser.h

bin_PROGRAMS = kiwi

kiwi_LDADD = -licuio -licuuc
kiwi_SOURCES = ast.c  parser.y  main.c scanner.c

