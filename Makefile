CXX=g++
CXXFLAGS=-g -Wall -std=c++11 -Iinclude
CFLAGS += -Wall -Iinclude

%.o: %.cc
	@$(CXX) $(CXXFLAGS) -o $@ $^

%.run: %.cc
	@rm -f $*.o
	@$(CXX) $(CXXFLAGS) -o $*.o $^
	@./$*.o

%.run: %.c
	@rm -f $*.o
	@$(CC) $(CFLAGS) -o $*.o $^
	@./$*.o

all: avl.o

avl.o: avl.cc

clean:
	rm -rf *.o
