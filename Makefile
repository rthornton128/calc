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

ifeq ($(OS),Window_NT)
	RM=del
	RMFLAGS=/f
	SRC=$(SRC:/=\\)
endif

OBJ=$(SRC:.c=.o)

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
