package rxlib

// Extension represents an extension.
type Extension struct {
	AutoAddExtension bool              `json:"autoAddExtension"` // if auto add is true it will add the Extension automatically if the parent object is added
	ExtensionName    string            `json:"extensionName"`
	FromPlugin       string            `json:"fromPlugin"`
	ParentObjectUUID string            `json:"parentObjectUUID"`
	ExtensionPorts   []*ExtensionPorts `json:"extensionPorts"`
}

func HistoryExtension(parentObjectUUID, fromPortID string, autoAdd bool) []*Extension {
	extensions := NewExtensionBuilder().
		NewExtension().
		WithExtensionName("history").
		AutoAddExtension(autoAdd).
		WithFromPlugin("main").
		WithParentObjectUUID(parentObjectUUID).
		AddExtensionAutoConnect(fromPortID, "input").Build()
	return extensions
}

func AlarmExtension(parentObjectUUID, fromPortID string, autoAdd bool) []*Extension {
	extensions := NewExtensionBuilder().
		NewExtension().
		WithExtensionName("alarm").
		AutoAddExtension(autoAdd).
		WithFromPlugin("main").
		WithParentObjectUUID(parentObjectUUID).
		AddExtensionAutoConnect(fromPortID, "input").Build()
	return extensions
}

// NewExtensionBuilder creates a new ExtensionBuilder instance.
func NewExtensionBuilder() *ExtensionBuilder {
	return &ExtensionBuilder{
		extensions: []*Extension{},
	}
}

// ExtensionPorts represents an auto port connection.
type ExtensionPorts struct {
	FromPortID string
	ToPortID   string
}

// ExtensionBuilder is a chain builder for creating extensions and managing them in an array.
type ExtensionBuilder struct {
	extensions []*Extension
}

// NewExtension creates a new Extension instance and adds it to the array.
func (builder *ExtensionBuilder) NewExtension() *ExtensionBuilder {
	extension := &Extension{}
	builder.extensions = append(builder.extensions, extension)
	return builder
}

func (builder *ExtensionBuilder) AutoAddExtension(autoAdd bool) *ExtensionBuilder {
	if len(builder.extensions) == 0 {
		builder.NewExtension()
	}
	builder.extensions[len(builder.extensions)-1].AutoAddExtension = autoAdd
	return builder
}

// WithExtensionName sets the extension name for the current extension in the array.
func (builder *ExtensionBuilder) WithExtensionName(extensionName string) *ExtensionBuilder {
	if len(builder.extensions) == 0 {
		builder.NewExtension()
	}
	builder.extensions[len(builder.extensions)-1].ExtensionName = extensionName
	return builder
}

// WithFromPlugin sets the from plugin for the current extension in the array.
func (builder *ExtensionBuilder) WithFromPlugin(fromPlugin string) *ExtensionBuilder {
	if len(builder.extensions) == 0 {
		builder.NewExtension()
	}
	builder.extensions[len(builder.extensions)-1].FromPlugin = fromPlugin
	return builder
}

// WithParentObjectUUID sets the parent object UUID for the current extension in the array.
func (builder *ExtensionBuilder) WithParentObjectUUID(parentObjectUUID string) *ExtensionBuilder {
	if len(builder.extensions) == 0 {
		builder.NewExtension()
	}
	builder.extensions[len(builder.extensions)-1].ParentObjectUUID = parentObjectUUID
	return builder
}

// AddExtensionAutoConnect adds an extension auto-connect to the current extension in the array.
func (builder *ExtensionBuilder) AddExtensionAutoConnect(fromPortID, toPortID string) *ExtensionBuilder {
	if len(builder.extensions) == 0 {
		builder.NewExtension()
	}
	autoConnect := &ExtensionPorts{
		FromPortID: fromPortID,
		ToPortID:   toPortID,
	}
	builder.extensions[len(builder.extensions)-1].ExtensionPorts = append(
		builder.extensions[len(builder.extensions)-1].ExtensionPorts, autoConnect)
	return builder
}

// Build builds the final array of Extension instances.
func (builder *ExtensionBuilder) Build() []*Extension {
	return builder.extensions
}
