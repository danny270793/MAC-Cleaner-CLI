package main

type Cleaner interface {
	Name() string
	Size() (int64, bool)
	Clean() (int64, bool)
	// Commands lists the commands Clean would run, for --verbose.
	Commands() []string
}
