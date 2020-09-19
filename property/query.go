package property

// EmptyQuery contain an empty Query. Can be used to call for no filtering.
var EmptyQuery = Query{}

// Query contains parameters that can be used to filter or search certain results.
type Query struct {
	ID  string
	Set string
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
