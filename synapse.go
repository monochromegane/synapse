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
	for _, p := range m.Profiles {
		sym, err := pluginSymbol(path, p, "Profiler")
		if err != nil {
			return nil, err
		}
		profiler, ok := sym.(Profiler)
		if !ok {
			return nil, UnexpectedTypeError{}
		}
		matcher.AddProfiler(profiler)
	}

	if p := m.Associator; p.Name != "" && p.Version != "" {
		sym, err := pluginSymbol(path, p, "Associator")
		if err != nil {
			return nil, err
		}
		associator, ok := sym.(Associator)
		if !ok {
			return nil, UnexpectedTypeError{}
		}
		matcher.SetAssociator(associator)
	}

	if p := m.Searcher; p.Name != "" && p.Version != "" {
		sym, err := pluginSymbol(path, p, "Searcher")
		if err != nil {
			return nil, err
		}
		searcher, ok := sym.(Searcher)
		if !ok {
			return nil, UnexpectedTypeError{}
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
