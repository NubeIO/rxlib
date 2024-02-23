package rxlib

type tree struct {
	objects []Object
}

type ExtractedDetails struct {
	ID         string              `json:"objectID,omitempty"`
	Name       string              `json:"name,omitempty"`
	UUID       string              `json:"uuid,omitempty"`
	ParentUUID string              `json:"parentUUID"`
	Category   string              `json:"category,omitempty"`
	ObjectType string              `json:"objectType,omitempty"`
	IsParent   bool                `json:"isParent,omitempty"`
	Children   []*ExtractedDetails `json:"children,omitempty"`
}

type ObjectsRootMap struct {
	RubixNetworkName string              `json:"rubixNetworkName"`
	RubixNetworkDesc string              `json:"RubixNetworkDesc"`
	RubixNetwork     []*ExtractedDetails `json:"rubixNetwork"`
	DriversName      string              `json:"driversName"`
	DriversDesc      string              `json:"driversDesc"`
	Drivers          []*ExtractedDetails `json:"drivers"`
	ServicesName     string              `json:"servicesName"`
	ServicesDesc     string              `json:"servicesDesc"`
	Services         []*ExtractedDetails `json:"services"`
	LogicName        string              `json:"logicName"`
	LogicDesc        string              `json:"logicDesc"`
	Logic            []*ExtractedDetails `json:"logic"`
}

func (t *tree) addObjects(objects []Object) {
	t.objects = objects
}

func (t *tree) GetTreeMapRoot() *ObjectsRootMap {
	rootTreeMap := &ObjectsRootMap{
		RubixNetworkName: "Rubix Networks",
		RubixNetworkDesc: "A place to add rubix-networks",
		RubixNetwork:     []*ExtractedDetails{},
		DriversName:      "Protocol Drivers",
		DriversDesc:      "Network protocols",
		Drivers:          []*ExtractedDetails{},
		ServicesName:     "System Services",
		ServicesDesc:     "Services for manging things like a user or network settings",
		Services:         []*ExtractedDetails{},
		LogicName:        "Logic Programs",
		LogicDesc:        "Logic Wiresheet Programs",
		Logic:            []*ExtractedDetails{},
	}

	// Create a helper function to add children to the appropriate category
	addToCategory := func(category string, obj *ExtractedDetails) {
		switch category {
		case "driver":
			rootTreeMap.Drivers = append(rootTreeMap.Drivers, obj)
		case "service":
			rootTreeMap.Services = append(rootTreeMap.Services, obj)
		case "logic":
			rootTreeMap.Logic = append(rootTreeMap.Logic, obj)
		case "rubix-network":
			rootTreeMap.RubixNetwork = append(rootTreeMap.RubixNetwork, obj)
		}
	}

	// Build the tree for each root reqUUID
	for _, obj := range t.objects {
		if obj.GetParentUUID() == "" {
			// Create the root object
			details := &ExtractedDetails{
				ID:         obj.GetID(),
				Name:       obj.GetName(),
				UUID:       obj.GetUUID(),
				Category:   obj.GetCategory(),
				ObjectType: string(obj.GetObjectType()),
				IsParent:   true,
				Children:   []*ExtractedDetails{},
			}

			// Add root object to the appropriate category
			addToCategory(string(obj.GetObjectType()), details)

			// Recursively build the tree
			t.buildTreeForTreeMap(details, obj.GetUUID(), addToCategory)
		}
	}

	return rootTreeMap
}

// Updated buildTreeForTreeMap function with addToCategory callback
func (t *tree) buildTreeForTreeMap(details *ExtractedDetails, uuid string, addToCategory func(string, *ExtractedDetails)) {
	// Continue building the tree with children objects
	for _, obj := range t.objects {
		if obj.GetParentUUID() == uuid {
			childDetails := &ExtractedDetails{
				ID:         obj.GetID(),
				Name:       obj.GetName(),
				UUID:       obj.GetUUID(),
				Category:   obj.GetCategory(),
				ObjectType: string(obj.GetObjectType()),
				Children:   []*ExtractedDetails{},
			}
			details.Children = append(details.Children, childDetails)

			// Recursively build the tree for children
			t.buildTreeForTreeMap(childDetails, obj.GetUUID(), addToCategory)
		}
	}
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
