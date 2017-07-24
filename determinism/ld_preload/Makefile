CFLAGS ?= -O0 -g -Wall -Wextra -Wconversion

all: main lib
main: main.c
lib: lib.c
	gcc $(CFLAGS) -fPIC -c lib.c -o lib.o
	gcc $(CFLAGS) -shared -o libxd.so lib.o
clean:
	rm -f main *.o *.so
.PHONY: clean all
