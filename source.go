package main

type SourceProvider interface {
	Get(file string) (*Resource, error)
}
