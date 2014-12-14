io.EOF != error{"EOF"}
======================

This repository demonstrates that "EOF" errors are not propagated as io.EOF
errors through gocircuit workers (current master: a6d1f33804). When the WriteCloser of a channel is closed,
the Reader receives an "EOF" error that is not equal to io.EOF. This makes
ioutil.ReadAll return an "EOF" error even though it's [documentation](http://golang.org/pkg/io/ioutil/#ReadAll) says, it should not do so.

The source for the behaviour is here: https://github.com/gocircuit/circuit/blob/7c97d0da1167d27da64ea8457a4e3d4766d084f2/kit/x/io/io.go#L50
