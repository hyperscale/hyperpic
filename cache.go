package main

type CacheProvider interface {
	Get(file string, options *ImageOptions) (*Resource, bool)

	Set(resource *Resource) error
}
