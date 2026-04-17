package filter

// Filter holds rules for ignoring specific ports or ranges.
type Filter struct {
	ignored map[int]struct{}
}

// New creates a Filter from a list of port numbers to ignore.
func New(ignoredPorts []int) *Filter {
	m := make(map[int]struct{}, len(ignoredPorts))
	for _, p := range ignoredPorts {
		m[p] = struct{}{}
	}
	return &Filter{ignored: m}
}

// Allow returns true if the port should be processed (not filtered out).
func (f *Filter) Allow(port int) bool {
	_, ignored := f.ignored[port]
	return !ignored
}

// FilterPorts returns only the ports that pass the filter.
func (f *Filter) FilterPorts(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if f.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}
