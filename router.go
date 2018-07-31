package synapse

type Router struct {
	routes map[string]Matcher
}

func (r Router) Match(path string, ctx Context) (Hits, error) {
	matcher, ok := r.routes[path]
	if !ok {
		return Hits{}, NoRouteMatchError{Path: path}
	}
	return matcher.Match(ctx)
}
