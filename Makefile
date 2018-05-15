#
# Makefile for CS 131 Assignment
# 
# As usual, you can build this software by typing
#
#     make
#

# ----- Make Macros -----

OPTFLAGS = -O3
WARNINGS = -Wall -Wextra
DEFINES  =
INCLUDES =
LDLIBS =

CPPFLAGS = $(INCLUDES)
LDFLAGS  = -g
CFLAGS =   -g -Wall $(OPTFLAGS) -pedantic $(WARNINGS) $(DEFINES) -std=c99
CXXFLAGS = -g -Wall $(OPTFLAGS) -pedantic $(WARNINGS) $(DEFINES) -std=c++1y

GCD_CC = $(CC)
GCD_LDLIBS = $(LDLIBS)
GCD_CFLAGS = $(CFLAGS)

OMP_CC = $(CC)
OMP_LDLIBS = $(LDLIBS)
OMP_CFLAGS = $(CFLAGS) -fopenmp

GO_CC = go

TARGETS =       nth-prime nth-prime-disp nth-prime-omp nth-prime-go

include Makefile.$(shell uname)
-include Makefile.local

.PHONY: all clean

all: $(TARGETS)

clean:
	rm -f $(TARGETS)

nth-prime: nth-prime.c
	$(CC) $(CFLAGS) -o nth-prime nth-prime.c $(LDLIBS)

nth-prime-disp: nth-prime-disp.c
	$(GCD_CC) $(GCD_CFLAGS) -o nth-prime-disp nth-prime-disp.c $(GCD_LDLIBS)

nth-prime-omp: nth-prime-omp.c
	$(OMP_CC) $(OMP_CFLAGS) -o nth-prime-omp nth-prime-omp.c $(OMP_LDLIBS)

nth-prime-go: nth-prime.go
	$(GO_CC) build -o nth-prime-go nth-prime.go

