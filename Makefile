GOCMD=go
GO_BUILD=$(GOCMD) build
GO_RUN=$(GOCMD) run
GO_BUILD_PLUGIN=$(GO_BUILD) -buildmode=plugin

.PHONY: example
example:
	$(GO_BUILD_PLUGIN) -o example/plugins/profiler.so.0.0.1   example/plugins/profiler.go
	$(GO_BUILD_PLUGIN) -o example/plugins/associator.so.0.0.1 example/plugins/associator.go
	$(GO_BUILD_PLUGIN) -o example/plugins/searcher.so.0.0.1   example/plugins/searcher.go
	$(GO_RUN) example/main.go
