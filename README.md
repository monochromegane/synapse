# Synapse

A Distributed RESTful and Plugable Matching Engine.

## Usage

WIP.

```go
func main() {
	config := synapse.Config{
		PluginDir: "plugins",
		Matchers: []synapse.ConfigMatcher{
			synapse.ConfigMatcher{
				Name: "users_products",
				Profiles: []synapse.ConfigPlugin{
					synapse.ConfigPlugin{
						Name:    "sample_profiler",
						Version: "0.0.1",
					},
				},
				Associator: synapse.ConfigPlugin{
					Name:    "sample_associator",
					Version: "0.0.1",
				},
				Searcher: synapse.ConfigPlugin{
					Name:    "sample_searcher",
					Version: "0.0.1",
				},
			},
			synapse.ConfigMatcher{
				Name: "m2",
				Host: "http://127.0.0.1:8080",
			},
		},
	}

	syn, err := synapse.NewSynapse(config)
	if err != nil {
		panic(err)
	}

	ctx := synapse.Context{"account_id": "xxxx"}
	hits, err := syn.Match("users_products", ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", hits.IDs) // [10 15 30]
}
```

## Architecture

### Distributed matching engine

Synapse selects the matching algorithm by Router.
Matching algorithms are provided as plugins and distributed processing is also possible.

![synapse](https://user-images.githubusercontent.com/1845486/43449245-e83d55fc-94ea-11e8-8e64-227884db0adb.jpeg)

### Plugable matching engine

Matcher is composed of profilers, associator and searcher. This is provided as a [plugin](https://golang.org/pkg/plugin/).

![synapse](https://user-images.githubusercontent.com/1845486/43451492-05ea97e0-94f0-11e8-9214-4307ac189b9d.jpeg)

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
