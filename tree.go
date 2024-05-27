package rxlib

import "github.com/NubeIO/rxlib/protos/runtimebase/runtime"

type tree struct {
	objects []Object
}

func (t *tree) addObjects(objects []Object) {
	t.objects = objects
}

func (t *tree) GetTreeMapRoot() *runtime.ObjectsRootMap {
	rootTreeMap := &runtime.ObjectsRootMap{
		RubixNetworkName: "Rubix Networks",
		RubixNetworkDesc: "A place to add rubix-networks",
		RubixNetwork:     []*runtime.ObjectExtractedDetails{},
		DriversName:      "Protocol Drivers",
		DriversDesc:      "Network protocols",
		Drivers:          []*runtime.ObjectExtractedDetails{},
		ServicesName:     "System Services",
		ServicesDesc:     "Services for managing things like a user or network settings",
		Services:         []*runtime.ObjectExtractedDetails{},
		LogicName:        "Logic Programs",
		LogicDesc:        "Logic Wiresheet Programs",
		Logic:            []*runtime.ObjectExtractedDetails{},
	}

	// Create a map to hold all objects for quick access by UUID
	objectMap := make(map[string]*runtime.ObjectExtractedDetails)

	// First pass: Create all objects and add them to the map
	for _, obj := range t.objects {
		details := &runtime.ObjectExtractedDetails{
			Id:         obj.GetID(),
			Name:       obj.GetName(),
			Uuid:       obj.GetUUID(),
			ParentUUID: obj.GetParentUUID(),
			Category:   obj.GetCategory(),
			ObjectType: string(obj.GetObjectType()),
			Children:   []*runtime.ObjectExtractedDetails{},
		}
		objectMap[obj.GetUUID()] = details
	}

	// Second pass: Build the tree by assigning children to their parents
	for _, details := range objectMap {
		if details.ParentUUID != "" {
			if parent, ok := objectMap[details.ParentUUID]; ok {
				parent.Children = append(parent.Children, details)
			}
		} else {
			// Root object, add it to the appropriate category
			switch details.ObjectType {
			case "driver":
				rootTreeMap.Drivers = append(rootTreeMap.Drivers, details)
			case "service":
				rootTreeMap.Services = append(rootTreeMap.Services, details)
			case "logic":
				rootTreeMap.Logic = append(rootTreeMap.Logic, details)
			case "rubix-network":
				rootTreeMap.RubixNetwork = append(rootTreeMap.RubixNetwork, details)
			}
		}
	}

	return rootTreeMap
}

// -------------------Ancestor-----------------------

//type AncestorObjectTree struct {
//	UUID       string              `json:"uuid"`
//	Name       string              `json:"name,omitempty"`
//	ID         string              `json:"id,omitempty"`
//	ParentUUID string              `json:"parentUUID,omitempty"`
//	Category   string              `json:"category,omitempty"`
//	Children   []*AncestorObjectTree `json:"children,omitempty"`
//}

func (t *tree) GetAncestorTreeByUUID(uuid string) *runtime.AncestorObjectTree {
	return t.buildAncestorTree(uuid)
}

func (t *tree) GetChilds(uuid string) *runtime.AncestorObjectTree {
	return t.buildChildTree(uuid)
}

func (t *tree) buildChildTree(parentUUID string) *runtime.AncestorObjectTree {
	for _, obj := range t.objects {
		if obj.GetParentUUID() == parentUUID {
			node := &runtime.AncestorObjectTree{
				Uuid:       obj.GetUUID(),
				Name:       obj.GetName(),
				Id:         obj.GetID(),
				ParentUUID: obj.GetUUID(),
				Category:   obj.GetCategory(),
				Children:   []*runtime.AncestorObjectTree{},
			}
			childNode := t.buildChildTree(obj.GetUUID())
			if childNode != nil {
				node.Children = append(node.Children, childNode)
			}
			return node
		}
	}
	return nil
}

func (t *tree) buildAncestorTree(childUUID string) *runtime.AncestorObjectTree {
	for _, obj := range t.objects {
		if obj.GetUUID() == childUUID {
			node := &runtime.AncestorObjectTree{
				Uuid:       obj.GetUUID(),
				Name:       obj.GetName(),
				Id:         obj.GetID(),
				ParentUUID: obj.GetUUID(),
				Category:   obj.GetCategory(),
			}
			if obj.GetParentUUID() != "" {
				parentNode := t.buildAncestorTree(obj.GetParentUUID())
				node.Children = append(node.Children, parentNode)
			}
			return node
		}
	}
	return nil
}
