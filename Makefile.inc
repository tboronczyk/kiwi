CC=gcc

# sub makefiles specify their own INCLUDE
INCLUDE=

WFLAGS=-Wall -Wextra -Winit-self -Wmissing-include-dirs -Wunused -Wfloat-equal \
 -Wshadow -Wunsafe-loop-optimizations -Wbad-function-cast -Wlogical-op \
 -Wmissing-prototypes -Wmissing-declarations -Wmissing-noreturn \
 -Wmissing-format-attribute

CFLAGS=-O2 -std=c99 $(INCLUDE)
CFLAGS_DBG=-O2 -std=c99 -pedantic $(WFLAGS) -DDEBUG -g $(INCLUDE)

