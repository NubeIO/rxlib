package rxlib

type Settings struct {
	Value interface{} `json:"value"` // comes from UI from the JSON schema
}

// Connection defines a structure for input subscriptions.
type Connection struct {
	SourceUUID    string        `json:"source"`
	SourcePort    string        `json:"sourceHandle"`
	TargetUUID    string        `json:"target"`
	TargetPort    string        `json:"targetHandle"`
	FlowDirection FlowDirection `json:"flowDirection"` // subscriber is if it's in an input and publisher if It's for an output

}
