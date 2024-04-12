package rxlib

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
	totalPages := (totalCount + pageSize - 1) / pageSize // Ceiling division
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
	filteredObjects := inst.GetChildObjects(parentUUID)
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
