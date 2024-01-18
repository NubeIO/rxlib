package plugins

import (
	"errors"
)

type Object struct {
	ID       string    `json:"id"`
	Export   string    `json:"export"`
	Children []*Object `json:"children,omitempty"`
}

type Category struct {
	Name    string    `json:"name"`
	Objects []*Object `json:"objects,omitempty"`
}

type Export struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Path        string      `json:"path"` // its file name
	Description string      `json:"description"`
	Categories  []*Category `json:"categories,omitempty"`
}

func NewPlugin(name, version, description string) *Export {
	return &Export{
		Name:        name,
		Version:     version,
		Description: description,
		Categories:  make([]*Category, 0),
	}
}

func (p *Export) AddCategory(categoryName string) {
	category := &Category{Name: categoryName, Objects: make([]*Object, 0)}
	p.Categories = append(p.Categories, category)
}

func (p *Export) GetCategory(categoryName string) (*Category, error) {
	for _, category := range p.Categories {
		if category.Name == categoryName {
			return category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (p *Export) AddObject(categoryName, objectID, export string) error {
	category, err := p.GetCategory(categoryName)
	if err != nil {
		return err
	}

	object := &Object{ID: objectID, Export: export}
	category.Objects = append(category.Objects, object)
	return nil
}

func (p *Export) AddChildObject(categoryName, parentID, childID, childExport string) error {
	category, err := p.GetCategory(categoryName)
	if err != nil {
		return err
	}

	parentObject := p.findObjectByID(category.Objects, parentID)
	if parentObject == nil {
		return errors.New("parent object not found")
	}

	childObject := &Object{ID: childID, Export: childExport}
	parentObject.Children = append(parentObject.Children, childObject)
	return nil
}

func (p *Export) findObjectByID(objects []*Object, id string) *Object {
	for _, object := range objects {
		if object.ID == id {
			return object
		}
		if foundObject := p.findObjectByID(object.Children, id); foundObject != nil {
			return foundObject
		}
	}
	return nil
}

func (p *Export) GetAllObjects() []*Object {
	var allObjects []*Object

	for _, category := range p.Categories {
		allObjects = append(allObjects, p.getAllObjectsInCategory(category.Objects)...)
	}

	return allObjects
}

func (p *Export) getAllObjectsInCategory(objects []*Object) []*Object {
	var allObjects []*Object

	for _, object := range objects {
		allObjects = append(allObjects, object)
		allObjects = append(allObjects, p.getAllObjectsInCategory(object.Children)...)
	}

	return allObjects
}

func (p *Export) GetObjects(export *Export) []*Object {
	var allObjects []*Object

	for _, category := range export.Categories {
		allObjects = append(allObjects, p.getAllObjectsInCategory(category.Objects)...)
	}

	return allObjects
}
