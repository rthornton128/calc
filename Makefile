CC=gcc
CFLAGS=-I ./runtime
LD=gcc
LDFLAGS=
AR=ar
ARFLAGS=crv
RM=rm
RMFLAGS=-vf

LIB=runtime/runtime.a
SRC=runtime/registers.c\
    runtime/stack.c
OBJ=$(SRC:.c=.o)
TEST_SRC=runtime/test.c
TEST_OBJ=runtime/test.o
TEST_BIN=runtime/test

ifeq ($(OS),Windows_NT)
RM=cmd /c del
RMFLAGS=/q
FixPath=$(subst /,\,$1)
X=.exe
else
FixPath=$1
endif

all: $(LIB)

.PHONY: install test test-all clean distclean

install: $(LIB)
	go install ./calcc

test-all: $(TEST)
	go test ./...

test: $(TEST_BIN)$(X)
	exec ./$(TEST_BIN)

$(TEST_BIN)$(X): $(TEST_OBJ) $(LIB)
	$(CC) $(LDFLAGS) -o $@ $^

$(LIB): $(OBJ)
	$(AR) $(ARFLAGS) $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) $(RMFLAGS) $(call FixPath,$(OBJ)) $(call FixPath,$(TEST_OBJ)) $(call FixPath,$(TEST_BIN)$(X))

distclean: clean
	$(RM) $(RMFLAGS) $(call FixPath,$(LIB))
