package navigation

// Router provides a stack-based view history for back navigation.
type Router struct {
	stack  []int
	params map[int]map[string]interface{} // per-view params
}

func NewRouter() *Router {
	return &Router{
		stack:  make([]int, 0, 16),
		params: make(map[int]map[string]interface{}),
	}
}

// Push saves the current view index for later back navigation.
func (r *Router) Push(view int) {
	r.stack = append(r.stack, view)
}

// Pop returns the previous view, or false if the stack is empty.
func (r *Router) Pop() (int, bool) {
	if len(r.stack) == 0 {
		return 0, false
	}
	prev := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return prev, true
}

// Depth returns the current stack depth.
func (r *Router) Depth() int {
	return len(r.stack)
}

// SetParam sets a parameter for a view.
func (r *Router) SetParam(view int, key string, value interface{}) {
	if r.params[view] == nil {
		r.params[view] = make(map[string]interface{})
	}
	r.params[view][key] = value
}

// GetParam gets a parameter for a view.
func (r *Router) GetParam(view int, key string) interface{} {
	if p := r.params[view]; p != nil {
		return p[key]
	}
	return nil
}

// ClearParams clears all params for a view.
func (r *Router) ClearParams(view int) {
	delete(r.params, view)
}
