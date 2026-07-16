package main

type Cleaner interface {
	Name() string
	Clean()
}
