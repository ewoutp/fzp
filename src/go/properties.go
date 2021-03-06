package fzp

import (
	"errors"
)

// Properties array
type Properties []Property

// NewProperties return a new Properties object
func NewProperties(famname string) Properties {
	p := Properties{}
	p = make([]Property, 0)
	// the 'family' property is required by the fritzing app
	p = append(p, NewProperty("family", famname))
	return p
}

// Total return the total number of properties
func (p *Properties) Total() int {
	return len(*p)
}

// AddValue a new property
func (p *Properties) AddValue(name, val string) error {
	_, err := p.GetValue(name)
	if err != nil {
		newProp := NewProperty(name, val)
		*p = append(*p, newProp)
		return nil
	}
	return errors.New("exist")
}

// GetValue return a property
func (p *Properties) GetValue(name string) (string, error) {
	for _, v := range *p {
		if v.Name == name {
			return v.Value, nil
		}
	}
	return "", errors.New("property '" + name + "' does not exist")
}

// Exist return true is a property name exist
func (p *Properties) Exist(name string) error {
	// TODO:...
	return nil
}

// Check the properties
func (p *Properties) Check() error {
	// check if each property name only exist once a time
	var tmp map[string]bool
	tmp = make(map[string]bool, 0)
	for _, prop := range *p {
		//TODO: check if a property with name="family" exists and is not empty
		if !tmp[prop.Name] {
			tmp[prop.Name] = true
		} else {
			return errors.New("property name '" + prop.Name + "' already exist")
		}
	}
	return nil
}
