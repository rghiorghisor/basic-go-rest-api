package property

import (
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/util"
)

// EmptyQuery contain an empty Query. Can be used to call for no filtering.
var EmptyQuery = Query{}

// Query contains parameters that can be used to filter or search certain results.
type Query struct {
	ID     string
	Set    string
	Fields Fields
}

// Fields contains any field names that must be returned.
type Fields struct {
	enabled bool
	values  map[string]struct{}
}

// NewFields retrieves a new Fields struct populated with the given names.
func NewFields(fields []string) Fields {
	newFields := Fields{enabled: false}

	if len(fields) == 0 {
		return newFields
	}

	newFields.enabled = true
	newFields.values = util.ArrayToSetString(fields...)

	return newFields
}

// Disable specifies that no Fields filtering must be done.
func (f *Fields) Disable() {
	f.enabled = false
}

// IsEnabled returns true if Fields filtering must be used and false otherwise.
func (f *Fields) IsEnabled() bool {
	return f.enabled && len(f.values) > 0
}

// Contains verifies if the given string value is contained by this Fields.
func (f *Fields) Contains(value string) bool {
	value = strings.ToLower(value)
	_, has := f.values[value]

	return has
}

// HasSet retrieves true if the Query has a Set name define, false otherwise.
func (q Query) HasSet() bool {
	return q.Set != ""
}

// GetSet retrieves the defined set name. The call to this func should be preceded
// by a call to the Query.HasSet method to make sure that a Set name is really defined.
func (q Query) GetSet() string {
	return q.Set
}
