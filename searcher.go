package synapse

type Searcher interface {
	Search(Context, Profile, Association) ([]string, []float64, error)
}
