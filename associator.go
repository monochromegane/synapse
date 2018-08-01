package synapse

type Associator interface {
	Initialize() error
	Associate(Context, Profile) (Association, error)
	Finalize() error
}

type Association map[string]string

type NoOpeAssociator struct {
}

func (m NoOpeAssociator) Associate(ctx Context, p Profile) (Association, error) {
	return Association{}, nil
}

func (m NoOpeAssociator) Initialize() error { return nil }
func (m NoOpeAssociator) Finalize() error   { return nil }
