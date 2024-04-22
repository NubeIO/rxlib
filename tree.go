package rxlib

import "github.com/NubeIO/rxlib/protos/runtimebase/runtime"

type tree struct {
	objects []Object
}

//type ObjectExtractedDetails struct {
//	ID         string                    `json:"objectID,omitempty"`
//	Name       string                    `json:"name,omitempty"`
//	UUID       string                    `json:"uuid,omitempty"`
//	ParentUUID string                    `json:"parentUUID"`
//	Category   string                    `json:"category,omitempty"`
//	ObjectType string                    `json:"objectType,omitempty"`
//	IsParent   bool                      `json:"isParent,omitempty"`
//	Children   []*ObjectExtractedDetails `json:"children,omitempty"`
//}
//
//type ObjectsRootMap struct {
//	RubixNetworkName string                    `json:"rubixNetworkName"`
//	RubixNetworkDesc string                    `json:"RubixNetworkDesc"`
//	RubixNetwork     []*ObjectExtractedDetails `json:"rubixNetwork"`
//	DriversName      string                    `json:"driversName"`
//	DriversDesc      string                    `json:"driversDesc"`
//	Drivers          []*ObjectExtractedDetails `json:"drivers"`
//	ServicesName     string                    `json:"servicesName"`
//	ServicesDesc     string                    `json:"servicesDesc"`
//	Services         []*ObjectExtractedDetails `json:"services"`
//	LogicName        string                    `json:"logicName"`
//	LogicDesc        string                    `json:"logicDesc"`
//	Logic            []*ObjectExtractedDetails `json:"logic"`
//}

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

type AncestorTreeNode struct {
	UUID       string              `json:"uuid"`
	Name       string              `json:"name,omitempty"`
	ID         string              `json:"id,omitempty"`
	ParentUUID string              `json:"parentUUID,omitempty"`
	Category   string              `json:"category,omitempty"`
	Children   []*AncestorTreeNode `json:"children,omitempty"`
}

func (t *tree) GetAncestorTreeByUUID(uuid string) *AncestorTreeNode {
	return t.buildAncestorTree(uuid)
}

func (t *tree) GetChilds(uuid string) *AncestorTreeNode {
	return t.buildChildTree(uuid)
}

func (t *tree) buildChildTree(parentUUID string) *AncestorTreeNode {
	for _, obj := range t.objects {
		if obj.GetParentUUID() == parentUUID {
			node := &AncestorTreeNode{
				UUID:       obj.GetUUID(),
				Name:       obj.GetName(),
				ID:         obj.GetID(),
				ParentUUID: obj.GetUUID(),
				Category:   obj.GetCategory(),
				Children:   []*AncestorTreeNode{},
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

func (t *tree) buildAncestorTree(childUUID string) *AncestorTreeNode {
	for _, obj := range t.objects {
		if obj.GetUUID() == childUUID {
			node := &AncestorTreeNode{
				UUID:       obj.GetUUID(),
				Name:       obj.GetName(),
				ID:         obj.GetID(),
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
