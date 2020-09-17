package model

// PropertySet is the central model struct of the property set feature.
// A set contains a collection of property names and can be used for searching
// and processing only a set of properties.
type PropertySet struct {
	Name   string
	Values []string
}
