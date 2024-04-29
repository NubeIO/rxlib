package rxlib

import (
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
)

type ObjectValuesPagination struct {
	PortValues []*runtime.PortValue
	Count      int `json:"count"`      // how many per page
	PageNumber int `json:"pageNumber"` // which page number
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"` // all the objects count
}

func (inst *RuntimeImpl) GetObjectsValuesPaginate(parentUUID string, pageNumber, pageSize int) *ObjectValuesPagination {
	var out []*runtime.PortValue
	objects := inst.PaginateGetChildObjects(parentUUID, pageNumber, pageSize)
	for _, object := range objects.Objects {
		out = append(out, inst.GetObjectValues(object.GetUUID())...)
	}
	return &ObjectValuesPagination{
		PortValues: out,
		Count:      objects.Count,
		PageNumber: objects.PageNumber,
		PageSize:   objects.PageSize,
		TotalPages: objects.TotalPages,
		TotalCount: objects.TotalCount,
	}

}

type ObjectPagination struct {
	Objects    []Object `json:"-"`
	Count      int      `json:"count"`      // how many per page
	PageNumber int      `json:"pageNumber"` // which page number
	PageSize   int      `json:"pageSize"`
	TotalPages int      `json:"totalPages"`
	TotalCount int      `json:"totalCount"` // all the objects count
}

func (inst *RuntimeImpl) ObjectsPagination(pageNumber, pageSize int) *ObjectPagination {
	inst.mutex.Lock()
	defer inst.mutex.Unlock()

	totalCount := len(inst.objects)
	totalPages := (totalCount + pageSize - 1) / pageSize // Ceiling division

	start := (pageNumber - 1) * pageSize
	end := start + pageSize
	if end > len(inst.objects) {
		end = len(inst.objects)
	}
	pagedObjects := inst.objects[start:end]

	return &ObjectPagination{
		Objects:    pagedObjects,
		Count:      len(pagedObjects),
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}

func (inst *RuntimeImpl) PaginateObjects(objects []Object, pageNumber, pageSize int) *ObjectPagination {
	totalCount := len(objects)
	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalCount + pageSize - 1) / pageSize // Ceiling division
	}
	start := (pageNumber - 1) * pageSize
	end := start + pageSize
	if end > len(objects) {
		end = len(objects)
	}
	pagedObjects := objects[start:end]

	return &ObjectPagination{
		Objects:    pagedObjects,
		Count:      len(pagedObjects),
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}

func (inst *RuntimeImpl) PaginateGetAllByID(objectID string, pageNumber, pageSize int) *ObjectPagination {
	filteredObjects := inst.GetAllByID(objectID)
	return inst.PaginateObjects(filteredObjects, pageNumber, pageSize)
}

func (inst *RuntimeImpl) PaginateGetAllByName(name string, pageNumber, pageSize int) *ObjectPagination {
	filteredObjects := inst.GetAllByName(name)
	return inst.PaginateObjects(filteredObjects, pageNumber, pageSize)
}

func (inst *RuntimeImpl) PaginateGetChildObjects(parentUUID string, pageNumber, pageSize int) *ObjectPagination {
	var filteredObjects []Object
	if parentUUID == "" {
		filteredObjects = inst.Get()
	} else {
		filteredObjects = inst.GetChildObjects(parentUUID)
	}
	return inst.PaginateObjects(filteredObjects, pageNumber, pageSize)
}

func (inst *RuntimeImpl) PaginateGetChildObjectsByWorkingGroup(objectUUID, workingGroup string, pageNumber, pageSize int) *ObjectPagination {
	filteredObjects := inst.GetChildObjectsByWorkingGroup(objectUUID, workingGroup)
	return inst.PaginateObjects(filteredObjects, pageNumber, pageSize)
}

func (inst *RuntimeImpl) handlePagination(parsedArgs *ParsedCommand) (*ObjectPagination, error) {
	childs := parsedArgs.GetChilds()
	objectUUID := parsedArgs.GetUUID()
	objectID := parsedArgs.GetID()
	pageSize := parsedArgs.GetPaginationPageSize()
	pageNumber := parsedArgs.GetPaginationPageNumber()
	if childs {
		resp := inst.PaginateGetChildObjects(objectUUID, pageNumber, pageSize)
		return resp, nil
	}
	if objectID != "" {
		resp := inst.PaginateGetAllByID(objectID, pageNumber, pageSize)
		return resp, nil
	}
	resp := inst.ObjectsPagination(pageNumber, pageSize)
	return resp, nil
}
