package connection

import (
	"github.com/NubeIO/rxlib/helpers"
	"github.com/NubeIO/rxlib/libs/history"
)

type Supervisor interface {
	AddRemoteManager(serverName string, manager Manager)
	GetRemoteManager(serverName string) Manager
	All() map[string]Manager
}

type SupervisorImpl struct {
	supervisorMap map[string]Manager
}

func NewSupervisor() Supervisor {
	return &SupervisorImpl{
		supervisorMap: make(map[string]Manager),
	}
}

func (a *SupervisorImpl) AddRemoteManager(serverName string, manager Manager) {
	a.supervisorMap[serverName] = manager
}

func (a *SupervisorImpl) GetRemoteManager(serverName string) Manager {
	return a.supervisorMap[serverName]
}

func (a *SupervisorImpl) All() map[string]Manager {
	return a.supervisorMap
}

type Manager interface {
	Get(uuid string) string
}

type ManagerImpl struct {
	UUID     string             `json:"uuid"`
	Name     string             `json:"name"`
	IsLocal  bool               `json:"isLocal"`
	Type     string             `json:"type"` // edgeToCloud, local, internal, rest, mqtt
	Networks map[string]Network `json:"networks,omitempty"`
}

type AddManager struct {
	Name string `json:"name"`
}

func NewManager(add *AddManager) Network {
	checkManagerImpl(add)
	return &ManagerImpl{UUID: helpers.UUID()}
}

func checkManagerImpl(add *AddManager) {
	if add == nil {
		panic("add AddManager is empty")
	}
	if add.Name == "" {
		panic("add AddManager.Name is empty")
	}
}

func (m *ManagerImpl) Get(uuid string) string {
	panic("implement me")
}

type Network interface {
	Get(uuid string) string
}

type AddNetwork struct {
	Name string `json:"name"`
}

type NetworkImp struct {
	UUID             string          `json:"uuid"`
	ProxyNetworkUUID string          `json:"proxyNetworkUUID,omitempty"`
	IP               string          `json:"ip"`
	Port             int             `json:"port"`
	UserName         string          `json:"userName"`
	Password         string          `json:"password"`
	Token            string          `json:"token"`
	HistoryManager   history.Manager `json:"historySupervisor"`
}

func (c *NetworkImp) Get(uuid string) string {
	panic("implement me")
}

func NewNetwork(add *AddNetwork) Network {
	checksNetworkImp(add)
	return &NetworkImp{UUID: helpers.UUID()}
}

func checksNetworkImp(add *AddNetwork) {
	if add == nil {
		panic("add AddNetwork is empty")
	}
	if add.Name == "" {
		panic("add AddNetwork.Name is empty")
	}
}
