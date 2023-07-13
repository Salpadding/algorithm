CXX=g++
CXXFLAGS=-g -Wall -std=c++11 -Iinclude

%.o: %.cc
	@$(CXX) $(CXXFLAGS) -o $@ $^

%.run: %.o
	@./$^
