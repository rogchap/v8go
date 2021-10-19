package v8go

type CPUProfileNode struct {
	// p *C.CPUProfileNode

	// The resource name for script from where the function originates.
	scriptResourceName string

	// The function name (empty string for anonymous functions.)
	functionName string

	// The number of the line where the function originates.
	lineNumber int

	// The number of the column where the function originates.
	columnNumber int

	// The children node of this node.
	children []*CPUProfileNode

	// The parent node of this node.
	parent *CPUProfileNode
}

// Returns function name (empty string for anonymous functions.)
func (c *CPUProfileNode) GetFunctionName() string {
	return c.functionName
}

// Returns resource name for script from where the function originates.
func (c *CPUProfileNode) GetScriptResourceName() string {
	return c.scriptResourceName
}

// Returns number of the line where the function originates.
func (c *CPUProfileNode) GetLineNumber() int {
	return c.lineNumber
}

// Returns number of the column where the function originates.
func (c *CPUProfileNode) GetColumnNumber() int {
	return c.columnNumber
}

// Retrieves the ancestor node, or nil if the root.
func (c *CPUProfileNode) GetParent() *CPUProfileNode {
	return c.parent
}

func (c *CPUProfileNode) GetChildrenCount() int {
	return len(c.children)
}

// Retrieves a child node by index.
func (c *CPUProfileNode) GetChild(index int) *CPUProfileNode {
	return c.children[index]
}
