package rubix

import (
	"errors"
	"github.com/NubeIO/rxlib/helpers"
	"sync"
	"time"
)

type Manager interface {
	All() map[string]Network
	AllByParent(objectUUID string) map[string]Network
	GetNetwork(uuid string) (Network, error)
	GetNetworkByObjectUUID(uuid string) (Network, error)
	GetNetworkByParentUUID(uuid string) (Network, error)
	AddNetwork(network Network)
	DeleteNetwork(uuid string) error
}

type ManagerImpl struct {
	Networks map[string]Network
	mu       sync.Mutex
}

func NewManager() Manager {
	return &ManagerImpl{
		Networks: make(map[string]Network),
	}
}

func (m *ManagerImpl) All() map[string]Network {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Networks
}

func (m *ManagerImpl) AllByParent(objectUUID string) map[string]Network {
	networks := make(map[string]Network)
	for _, network := range m.Networks {
		if network.GetParentObjectUUID() == objectUUID {
			networks[network.GetUUID()] = network
		}
	}
	return networks
}

func (m *ManagerImpl) GetNetwork(uuid string) (Network, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	network, exists := m.Networks[uuid]
	if !exists {
		return nil, errors.New("network not found")
	}
	return network, nil
}

func (m *ManagerImpl) GetNetworkByObjectUUID(uuid string) (Network, error) {
	for _, network := range m.Networks {
		if network.GetObjectUUID() == uuid {
			return network, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *ManagerImpl) GetNetworkByParentUUID(uuid string) (Network, error) {
	for _, network := range m.Networks {
		if network.GetParentObjectUUID() == uuid {
			return network, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *ManagerImpl) AddNetwork(network Network) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Networks[network.GetUUID()] = network
}

func (m *ManagerImpl) DeleteNetwork(uuid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.Networks[uuid]
	if !exists {
		return errors.New("network not found")
	}
	delete(m.Networks, uuid)
	return nil
}

type Network interface {
	GetUUID() string
	GetObjectUUID() string
	GetParentObjectUUID() string
	NewMapping(body *Mapping)
	GetMapping(uuid string) (*Mapping, error)
	DeleteMapping(uuid string) error
	AllMapping() []*Mapping
	UpdateLastOk()
	UpdateLastFail()
	AddConnection(mappingUUID string, connection *Connection) error
	EditConnection(mappingUUID string, connection *Connection) error
	DeleteConnection(mappingUUID, connectionUUID string) error
}

type NetworkImp struct {
	UUID             string     `json:"uuid"`
	ObjectUUID       string     `json:"objectUUID"`
	ParentObjectUUID string     `json:"parentObjectUUID"`
	Name             string     `json:"name"`
	IP               string     `json:"IP"`
	Port             int        `json:"port"`
	Created          time.Time  `json:"created"`
	IsError          bool       `json:"isError"`
	LastOk           *time.Time `json:"LastOk"`
	LastFail         *time.Time `json:"LastFail"`
	Mapping          []*Mapping `json:"mapping"`
}

type Mapping struct {
	UUID              string               `json:"uuid"`
	RemoteMappingUUID string               `json:"remoteMappingUUID"`
	Direction         MappingTypeDirection `json:"direction"`
	Type              MappingTypeSetting   `json:"type"`
	IntervalSetting   time.Duration        `json:"intervalSetting"`
	Enable            bool                 `json:"enable"`
	Created           time.Time            `json:"created"`
	Connections       []*Connection        `json:"connections"`
}

type Connection struct {
	UUID           string     `json:"uuid"`
	ConnectionUUID string     `json:"connectionUUID"`
	ObjectUUID     string     `json:"objectUUID"`
	PortID         string     `json:"portID"`
	Enable         bool       `json:"enable"`
	IsError        bool       `json:"isError"`
	Created        time.Time  `json:"created"`
	LastOk         *time.Time `json:"LastOk"`
	LastFail       *time.Time `json:"LastFail"`
	FailCount      int        `json:"failCount"`
}

type MappingTypeSetting string

const (
	MappingTypeCOV         = "cov"
	MappingTypeInterval    = "interval"
	MappingTypeCOVInterval = "cov-interval"
)

type MappingTypeDirection string

const (
	UpStream   = "upstream"
	DownStream = "down-stream"
	DualStream = "dual-stream"
)

type AddNetwork struct {
	Name             string
	ObjectUUID       string
	ParentObjectUUID string
}

func NewNetwork(body *AddNetwork) Network {
	return &NetworkImp{
		UUID:             helpers.UUID(),
		Name:             body.Name,
		ObjectUUID:       body.ObjectUUID,
		ParentObjectUUID: body.ParentObjectUUID,
	}
}

func (n *NetworkImp) GetUUID() string {
	return n.UUID
}

func (n *NetworkImp) GetObjectUUID() string {
	return n.ObjectUUID
}

func (n *NetworkImp) GetParentObjectUUID() string {
	return n.ParentObjectUUID
}

type AddMapping struct {
	UUID                 string `json:"uuid"`
	SourceUUID           string `json:"sourceUUID"`
	SourceConnectionUUID string `json:"sourceConnectionUUID"` // rxlib.Connection
	TargetUUID           string `json:"targetUUID"`
	TargetConnectionUUID string `json:"targetConnectionUUID"`
}

func NewMapping(body *AddMapping) *Mapping {
	return &Mapping{
		UUID:              helpers.UUID(),
		RemoteMappingUUID: "",
		Direction:         "",
		Type:              "",
		IntervalSetting:   0,
		Connections:       nil,
	}
}

func (n *NetworkImp) NewMapping(body *Mapping) {
	if body != nil {
		body.Created = time.Now()
	}
	n.Mapping = append(n.Mapping, body)
}

func (n *NetworkImp) GetMapping(uuid string) (*Mapping, error) {
	for _, mapping := range n.Mapping {
		if mapping.UUID == uuid {
			return mapping, nil
		}
	}
	return nil, errors.New("mapping not found")
}

func (n *NetworkImp) DeleteMapping(uuid string) error {
	for i, mapping := range n.Mapping {
		if mapping.UUID == uuid {
			n.Mapping = append(n.Mapping[:i], n.Mapping[i+1:]...)
			return nil
		}
	}
	return errors.New("mapping not found")
}

func (n *NetworkImp) AllMapping() []*Mapping {
	return n.Mapping
}

func (n *NetworkImp) UpdateLastOk() {
	n.LastOk = toTime(time.Now())
	n.IsError = false
}

func (n *NetworkImp) UpdateLastFail() {
	n.LastFail = toTime(time.Now())
	n.IsError = true
}

func toTime(t time.Time) *time.Time {
	return &t
}
func (n *NetworkImp) AddConnection(mappingUUID string, connection *Connection) error {
	for _, mapping := range n.Mapping {
		if mapping.UUID == mappingUUID {
			if mapping.Connections == nil {
				mapping.Connections = []*Connection{}
			}
			mapping.Connections = append(mapping.Connections, connection)
			return nil
		}
	}
	return errors.New("mapping not found")
}

func (n *NetworkImp) EditConnection(mappingUUID string, connection *Connection) error {
	for _, mapping := range n.Mapping {
		if mapping.UUID == mappingUUID {
			if mapping.Connections == nil {
				return errors.New("connection not found")
			}
			for i, conn := range mapping.Connections {
				if conn.UUID == connection.UUID {
					mapping.Connections[i] = connection
					return nil
				}
			}
			return errors.New("connection not found")
		}
	}
	return errors.New("mapping not found")
}

func (n *NetworkImp) DeleteConnection(mappingUUID string, connectionUUID string) error {
	for _, mapping := range n.Mapping {
		if mapping.UUID == mappingUUID {
			if mapping.Connections == nil {
				return errors.New("connection not found")
			}
			for i, conn := range mapping.Connections {
				if conn.UUID == connectionUUID {
					mapping.Connections = append(mapping.Connections[:i], mapping.Connections[i+1:]...)
					return nil
				}
			}
			return errors.New("connection not found")
		}
	}
	return errors.New("mapping not found")
}
