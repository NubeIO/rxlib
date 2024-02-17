package rxlib

import "github.com/NubeIO/rxlib/priority"

func (d *DataPayload) PriorityValue() *priority.Value {
	return d.Data
}

func (d *DataPayload) PriorityData() *priority.PriorityData {
	return d.Data.PriorityData()
}

func (d *DataPayload) GetData() any {
	return d.Data.GetPriority()
}

func (d *DataPayload) GetFloatPointer() *float64 {
	return d.Data.GetFloatPointer()
}

func (d *DataPayload) GetType() priority.Type {
	return d.Data.GetType()
}
