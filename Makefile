CXX=g++
CXXFLAGS=-g -Wall -std=c++11 -Iinclude

%.o: %.cc
	@$(CXX) $(CXXFLAGS) -o $@ $^

%.run: %.cc
	@rm -f $*.o
	@$(CXX) $(CXXFLAGS) -o $*.o $^
	@./$*.o

all: avl.o

avl.o: avl.cc

clean:
	rm -rf *.o