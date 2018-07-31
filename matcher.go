package synapse

type Matcher interface {
	Name() string
	Match(Context) (Hits, error)
}
