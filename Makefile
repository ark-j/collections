all: test-smap test-reqwest

test-reqwest:
	@go test -v -count 1 collections/reqwest

test-smap:
	@go test -v -count 1 collections/smap
