package proxy

type ProxyNetwork struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	IsLocal     bool          `json:"isLocal"`
	Type        string        `json:"type"` // edgeToCloud, local, internal, rest, mqtt
	Connections []*Connection `json:"connections,omitempty"`
}

type Connection struct {
	UUID             string `json:"uuid"`
	ProxyNetworkUUID string `json:"proxyNetworkUUID,omitempty"`
	IP               string `json:"ip"`
	Port             int    `json:"port"`
	UserName         string `json:"userName"`
	Password         string `json:"password"`
	Token            string `json:"token"`
}
