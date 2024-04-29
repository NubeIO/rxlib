package rxlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

func (inst *RuntimeImpl) AddObject(object Object) {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	inst.objects = append(inst.objects, object)
}

func (inst *RuntimeImpl) GetAllByID(objectID string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetID() == objectID {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetByUUID(uuid string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, object := range inst.objects {
		if object.GetUUID() == uuid {
			return object
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetAllByName(name string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetName() == name {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetChildObjects(parentUUID string) []Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	var out []Object
	for _, obj := range inst.objects {
		if obj.GetParentUUID() == parentUUID {
			out = append(out, obj)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetFirstByID(objectID string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, obj := range inst.objects {
		if obj.GetID() == objectID {
			return obj
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetFirstByName(name string) Object {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	for _, obj := range inst.objects {
		if obj.GetName() == name {
			return obj
		}
	}
	return nil
}

func (inst *RuntimeImpl) GetAllObjectValues() []*ObjectValue {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()
	allObjects := inst.Get()
	nodeValues := make([]*ObjectValue, len(allObjects))
	for _, node := range allObjects {
		nv := node.GetAllPorts()
		if nv == nil {
			continue
		}
		portValue := &ObjectValue{
			ObjectId:   node.GetID(),
			ObjectUUID: node.GetUUID(),
			Ports:      nv,
		}
		nodeValues = append(nodeValues, portValue)
	}
	return nodeValues
}

func (inst *RuntimeImpl) objectsFilteredIsNil() bool {
	if inst.objectsFiltered == nil {
		return true
	}
	return false
}

// GetChildObjectsByWorkingGroup
// for example get all the childs object for working group "rubix"
func (inst *RuntimeImpl) GetChildObjectsByWorkingGroup(objectUUID, workingGroup string) []Object {
	var out []Object
	for _, object := range inst.Get() {
		if object.GetUUID() == objectUUID {
			if object.GetWorkingGroup() == workingGroup {
				out = append(out, object)
			}
		}
	}
	return out
}

func (inst *RuntimeImpl) GetObjectValues(objectUUID string) []*runtime.PortValue {
	obj := inst.GetByUUID(objectUUID)
	if obj == nil {
		return nil
	}
	var out []*runtime.PortValue
	inputs := obj.GetInputs()
	for _, port := range inputs {
		out = append(out, obj.GetPortValue(port.GetID()))
	}
	outputs := obj.GetOutputs()
	for _, port := range outputs {
		out = append(out, obj.GetPortValue(port.GetID()))
	}
	return out
}

func (inst *RuntimeImpl) GetObjectsValues(parentUUID string) []*runtime.PortValue {
	var out []*runtime.PortValue
	if parentUUID == "" {
		for _, object := range inst.Get() {
			out = append(out, inst.GetObjectValues(object.GetUUID())...)
		}
	} else {
		for _, object := range inst.GetChildObjects(parentUUID) {
			out = append(out, inst.GetObjectValues(object.GetUUID())...)
		}
	}
	return out
}

func (inst *RuntimeImpl) GetObjectsConfig() []*runtime.ObjectConfig {
	return inst.SerializeObjects(false, inst.Get())
}

func (inst *RuntimeImpl) ToObjectsConfig(objects []Object) []*runtime.ObjectConfig {
	return inst.SerializeObjects(false, objects)
}

func (inst *RuntimeImpl) GetObjectsUUIDs(objects []Object) []string {
	var out []string
	for _, object := range objects {
		out = append(out, object.GetUUID())
	}
	return out
}

func (inst *RuntimeImpl) GetObjectConfig(uuid string) *runtime.ObjectConfig {
	return inst.serializeObject(false, inst.GetByUUID(uuid))
}

func (inst *RuntimeImpl) GetObjectConfigByID(objectID string) *runtime.ObjectConfig {
	object := inst.GetFirstByID(objectID)
	if object == nil {
		return nil
	}
	return inst.serializeObject(false, object)
}

func (inst *RuntimeImpl) GetTreeMapRoot() *runtime.ObjectsRootMap {
	inst.tree.addObjects(inst.objects)
	return inst.tree.GetTreeMapRoot()
}

func (inst *RuntimeImpl) GetAncestorTreeByUUID(objectUUID string) *runtime.AncestorObjectTree {
	return inst.tree.GetAncestorTreeByUUID(objectUUID)
}

func (inst *RuntimeImpl) GetTreeChilds(objectUUID string) *runtime.AncestorObjectTree {
	return inst.tree.GetChilds(objectUUID)
}
