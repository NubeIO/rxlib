package rxlib

type ObjectTypeTag string

const (
	DriversTag        ObjectTypeTag = "drivers"
	BACnetTag         ObjectTypeTag = "bacnet"
	ModbusTag         ObjectTypeTag = "modbus"
	MQTTTag           ObjectTypeTag = "mqtt"
	RestTag           ObjectTypeTag = "rest"
	RubixTag          ObjectTypeTag = "rubix"
	LogicTag          ObjectTypeTag = "logic"
	TimeTag           ObjectTypeTag = "time"
	DateTag           ObjectTypeTag = "date"
	ServicesTag       ObjectTypeTag = "service"
	UsersTag          ObjectTypeTag = "users"
	NetworkingTag     ObjectTypeTag = "networking"
	IpAddressTag      ObjectTypeTag = "ipAddress"
	ProtocolTag       ObjectTypeTag = "protocol"
	IpProtocolTag     ObjectTypeTag = "ipProtocol"
	SerialProtocolTag ObjectTypeTag = "serialProtocol"

	NetworkTag ObjectTypeTag = "network"
	DeviceTag  ObjectTypeTag = "device"
	PointTag   ObjectTypeTag = "point"

	HistTag    ObjectTypeTag = "hist"
	AlarmTag   ObjectTypeTag = "alarm"
	MappingTag ObjectTypeTag = "mapping"
)
