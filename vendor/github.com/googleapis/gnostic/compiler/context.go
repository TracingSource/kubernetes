package compiler

// Context contains state of the compiler as it traverses a document.
type Context struct {
	Parent            *Context
	Name              string
	ExtensionHandlers *[]ExtensionHandler
}

// NewContextWithExtensions returns a new object representing the compiler state
func NewContextWithExtensions(name string, parent *Context, extensionHandlers *[]ExtensionHandler) *Context {
	return &Context{Name: name, Parent: parent, ExtensionHandlers: extensionHandlers}
}

// NewContext returns a new object representing the compiler state
func NewContext(name string, parent *Context) *Context {
	if parent != nil {
		return &Context{Name: name, Parent: parent, ExtensionHandlers: parent.ExtensionHandlers}
	}
	return &Context{Name: name, Parent: parent, ExtensionHandlers: nil}
}

// Description returns a text description of the compiler state
func (context *Context) Description() string {
	if context.Parent != nil {
		return context.Parent.Description() + "." + context.Name
	}
	return context.Name
}
