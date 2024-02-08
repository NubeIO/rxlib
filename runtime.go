package rxlib

type Runtime interface {
	Get() map[string]Object
	GetByUUID(uuid string) Object
	GetByName(name string) Object
}

func NewRuntime(objs map[string]Object) Runtime {
	r := &runtime{}
	r.objects = objs
	if r.objects == nil {
		panic("NewRuntime object map can not be empty")
	}
	return r
}

type runtime struct {
	objects map[string]Object
}

func (r *runtime) Get() map[string]Object {
	return r.objects
}

func (r *runtime) GetByUUID(uuid string) Object {
	obj, _ := r.objects[uuid]
	return obj
}

func (r *runtime) GetByName(name string) Object {
	for _, obj := range r.objects {
		if obj.GetName() == name {
			return obj
		}
	}
	return nil
}
