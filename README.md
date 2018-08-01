# Synapse

A Distributed RESTful and Plugable Matching Engine.

## Usage

WIP.

```go
func main() {
	config := synapse.Config{}
	err := config.LoadYAML("synapse.yaml")
	if err != nil {
		panic(err)
	}

	syn, err := synapse.NewSynapse(config)
	if err != nil {
		panic(err)
	}

	ctx := synapse.Context{"user_id": "1"}
	hits, err := syn.Match("users_products", ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", hits.IDs) // [10 15 30]
}
```

```yaml
plugin_dir: plugins
matchers:
  -
    name: users_products
    profilers:
      -
        name: profiler
        version: 0.0.1
    associator:
      name: associator
      version: 0.0.1
    searcher:
      name: searcher
      version: 0.0.1
```

## Architecture

### Distributed matching engine

Synapse selects the matching algorithm by Router.
Matching algorithms are provided as plugins and distributed processing is also possible.

![synapse](https://user-images.githubusercontent.com/1845486/43449245-e83d55fc-94ea-11e8-8e64-227884db0adb.jpeg)

### Plugable matching engine

Matcher is composed of profilers, associator and searcher. This is provided as a [plugin](https://golang.org/pkg/plugin/).

![synapse](https://user-images.githubusercontent.com/1845486/43510200-d53269a0-95af-11e8-8797-de73f1e9babe.jpeg)

Profilers are filter group that converts subject and context into profile.
Associator is a filter that associates profile and context with cluster.
Searcher is a filter that searches for cluster using context, profile and association.
This also sorts the search results.

## Plugin

You can prepare custom plugins.
Plugins must be implemented as Profiler or Associator or Searcher interface, and be published with permitted function names, "NewProfiler", "NewAssociator" and "NewSearcher".

```go
package main

import "github.com/monochromegane/synapse"

type profiler struct {
}

// Implement interface
func (p *profiler) Profile(ctx synapse.Context) (synapse.Profile, error) {
	return synapse.Profile{}, nil
}

func (p *profiler) Initialize() error {
	return nil
}

func (p *profiler) Finalize() error {
	return nil
}

// Export new func
func NewProfiler() synapse.Profiler {
	return &profiler{}
}
```

And build as plugin.

```sh
$ go build -buildmode=plugin -o PLUGIN_NAME.so.x.x.x PLUGIN_NAME.go
```
