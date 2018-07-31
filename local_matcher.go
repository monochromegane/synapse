package synapse

type LocalMatcher struct {
	config     ConfigMatcher
	profilers  []Profiler
	associator Associator
	searcher   Searcher
}

func newLocalMatcher(config ConfigMatcher) *LocalMatcher {
	return &LocalMatcher{
		config:     config,
		profilers:  []Profiler{},
		associator: NoOpeAssociator{},
		searcher:   nil,
	}
}

func (m LocalMatcher) Name() string {
	return m.config.Name
}

func (m *LocalMatcher) AddProfiler(p Profiler) {
	m.profilers = append(m.profilers, p)
}

func (m *LocalMatcher) SetAssociator(a Associator) {
	m.associator = a
}

func (m *LocalMatcher) SetSeacher(s Searcher) {
	m.searcher = s
}

func (m LocalMatcher) Match(ctx Context) (Hits, error) {
	hits := m.newHits()
	profile := Profile{}
	for _, profiler := range m.profilers {
		p, err := profiler.Profile(ctx)
		if err != nil {
			return hits, err
		}
		profile.Merge(p)
	}

	association, err := m.associator.Associate(ctx, profile)
	if err != nil {
		return hits, err
	}

	ids, scores, err := m.searcher.Search(ctx, profile, association)
	if err != nil {
		return hits, err
	}
	hits.IDs = ids
	hits.Scores = scores

	return hits, nil
}

func (m LocalMatcher) newHits() Hits {
	return Hits{
		Matcher: m.config,
	}
}
