package synapse

type Hits struct {
	Matcher ConfigMatcher

	IDs    []string
	Scores []float64
}
