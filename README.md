# Synapse

A Distributed RESTful and Plugable Matching Engine.

## Usage

TODO

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
Plugins must be published with permitted variable names, "Profiler", "Associator" and "Searcher".

```go
package main

import "github.com/monochromegane/synapse"

type profiler struct {
}

func (p profiler) Profile(ctx synapse.Context) (synapse.Profile, error) {
	return synapse.Profile{}, nil
}

// Export variable
var Profiler profiler
```

And build as plugin.

```sh
$ go build -buildmode=plugin -o PLUGIN_NAME.so.x.x.x PLUGIN_NAME.go
```
