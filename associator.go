package synapse

type Associator interface {
	Associate(Context, Profile) (Association, error)
}

type Association map[string]string

type NoOpeAssociator struct {
}

func (m NoOpeAssociator) Associate(ctx Context, p Profile) (Association, error) {
	return Association{}, nil
}
