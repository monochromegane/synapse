package synapse

import (
	"fmt"
	"path/filepath"
	"plugin"
)

type Synapse struct {
	config Config
	router Router
}

func NewSynapse(config Config) (*Synapse, error) {
	router := Router{routes: map[string]Matcher{}}
	for _, m := range config.Matchers {
		matcher, err := matcherFromConfig(config.PluginDir, m)
		if err != nil {
			return nil, err
		}
		router.routes[m.Name] = matcher
	}
	return &Synapse{
		config: config,
		router: router,
	}, nil
}

func (s Synapse) Match(algorithm string, ctx Context) (Hits, error) {
	return s.router.Match(algorithm, ctx)
}

func matcherFromConfig(path string, m ConfigMatcher) (Matcher, error) {
	if m.Host == "" {
		return localMatcherFromConfig(path, m)
	} else {
		return newRemoteMatcher(m)
	}
}

func localMatcherFromConfig(path string, m ConfigMatcher) (Matcher, error) {
	matcher := newLocalMatcher(m)
	for _, p := range m.Profilers {
		sym, err := pluginSymbol(path, p, "NewProfiler")
		if err != nil {
			return nil, err
		}
		newer, ok := sym.(func() Profiler)
		if !ok {
			return nil, UnexpectedTypeError{}
		}
		profiler := newer()
		err = profiler.Initialize()
		if err != nil {
			return nil, err
		}
		matcher.AddProfiler(profiler)
	}

	if p := m.Associator; p.Name != "" && p.Version != "" {
		sym, err := pluginSymbol(path, p, "NewAssociator")
		if err != nil {
			return nil, err
		}
		newer, ok := sym.(func() Associator)
		if !ok {
			return nil, UnexpectedTypeError{}
		}
		associator := newer()
		err = associator.Initialize()
		if err != nil {
			return nil, err
		}
		matcher.SetAssociator(associator)
	}

	if p := m.Searcher; p.Name != "" && p.Version != "" {
		sym, err := pluginSymbol(path, p, "NewSearcher")
		if err != nil {
			return nil, err
		}
		newer, ok := sym.(func() Searcher)
		if !ok {
			return nil, UnexpectedTypeError{}
		}
		searcher := newer()
		err = searcher.Initialize()
		if err != nil {
			return nil, err
		}
		matcher.SetSeacher(searcher)
	}
	return matcher, nil
}

func pluginSymbol(path string, p ConfigPlugin, name string) (plugin.Symbol, error) {
	mod := filepath.Join(path, fmt.Sprintf("%s.so.%s", p.Name, p.Version))
	plug, err := plugin.Open(mod)
	if err != nil {
		return nil, err
	}
	return plug.Lookup(name)
}
