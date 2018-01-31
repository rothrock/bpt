# A B+ Tree in Go

This is an effort to create an example B+Tree library in go.

[Here is a good B Tree visualization](https://www.cs.usfca.edu/~galles/visualization/BTree.html)

# Why do this?
* Interesting intellectual exercise
* Opportunity to think about design choices
* Good way to learn about how btrees work
** B+Trees and their associated methods are fundamental to databases and filesystems.
** SQL Databases need proper indexes (btrees). Go [here](http://use-the-index-luke.com/)
* Very good way to practice writing Go code
* Opportunity to learn how to do the tricky things the Go way
** Granular locking
** Smart disk IO
** Concurrency
** Message passing
* Practice rewriting big chunks of code. Good code is always the product of iteration.
* Maybe create something that someone finds useful

# Architecture and approach to design
It's fairly primitive. I read and studied algorithms and then made a stab at implementation.
* Layout like a very plain Go library
* Implement some testing early on
* Choose public-facing and private methods
* Use recursion (it's a tree, after all)
* Get something working.
** It doesn't have to be pretty to start with.
** It is just code. You can change it.

# To Do List
1. Get Find working (done)
1. Insert working (done)
1. Implement locking and concurrent access
1. Get delete working
1. Get update working
1. Implement some range queries
1. Some better encapsulation (get rid of Init())
1. Move to disk-based storage
