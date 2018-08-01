package synapse

type Searcher interface {
	Initialize() error
	Search(Context, Profile, Association) ([]string, []float64, error)
	Finalize() error
}
