package synapse

type Profiler interface {
	Initialize() error
	Profile(Context) (Profile, error)
	Finalize() error
}

type Profile map[string]string

func (p Profile) Merge(p2 Profile) {
	for k, v := range p2 {
		p[k] = v
	}
}
