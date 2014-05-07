CC=gcc
CFLAGS=
LD=gcc
LDFLAGS=
AR=ar
ARFLAGS=crv
RM=rm
RMFLAGS=-vf

LIB=runtime.a
SRC=runtime/cmp.c\
    runtime/instructions.c\
    runtime/registers.c\
    runtime/stack.c
OBJ=$(SRC:.c=.o)

ifeq ($(OS),Windows_NT)
RM=cmd /c del
RMFLAGS=
OBJ=$(subst /,\,$(SRC:.c=.o))
endif

all: $(LIB)

.PHONY: install

install:
	go install ./calcc

$(LIB): $(OBJ)
	$(AR) $(ARFLAGS) $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) $(RMFLAGS) $(OBJ)

distclean: clean
	$(RM) $(RMFLAGS) $(LIB)
