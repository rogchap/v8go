package v8go

// RegisterCallback is exported for testing only.
func (i *Isolate) RegisterCallback(cb FunctionCallback) int {
	return i.registerCallback(cb)
}

// GetCallback is exported for testing only.
func (i *Isolate) GetCallback(ref int) FunctionCallback {
	return i.getCallback(ref)
}

// Register is exported for testing only.
func (c *Context) Register() {
	c.register()
}

// Deregister is exported for testing only.
func (c *Context) Deregister() {
	c.deregister()
}

// GetContext is exported for testing only.
var GetContext = getContext

// Ref is exported for testing only.
func (c *Context) Ref() int {
	return c.ref
}
