CC=gcc
CFLAGS=-I ./runtime
LD=gcc
LDFLAGS=
AR=ar
ARFLAGS=crv
RM=rm
RMFLAGS=-vf

LIB=runtime/runtime.a
SRC=runtime/cmp.c\
    runtime/instructions.c\
    runtime/registers.c\
    runtime/stack.c
OBJ=$(SRC:.c=.o)
TEST_SRC=runtime/test.c
TEST_OBJ=runtime/test.o

ifeq ($(OS),Windows_NT)
RM=cmd /c del
RMFLAGS=
OBJ=$(subst /,\,$(SRC:.c=.o))
TEST_OBJ=$(subst /,\,$(TEST_SRC:.c=.o))
TEST_BIN=test.exe
else
TEST_BIN=test
endif

all: $(LIB)

.PHONY: install test-all clean distclean

install: $(LIB)
	go install ./calcc

test-all: $(TEST)
	exec ./$(TEST_BIN)
	go test ./...

$(TEST_BIN): $(TEST_OBJ) $(LIB)
	$(CC) $(LDFLAGS) -o $@ $^

$(LIB): $(OBJ)
	$(AR) $(ARFLAGS) $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) $(RMFLAGS) $(OBJ) $(TEST_OBJ) $(TEST_BIN)

distclean: clean
	$(RM) $(RMFLAGS) $(LIB)
