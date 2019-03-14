package logging

import (
	"strings"
)

// Filter interface is to perform arbitrary filtering of LogRecords.
// Loggers and handlers can optionally use filter instances to filter records
// as desired.
type Filter interface {
	Filter(record *LogRecord) bool
}

// The base filter allows events which are below a certain point in the logger
// hierarchy. For Example, a filter initialized with "A.B" will allow events
// logged by loggers "A.B", "A.B.C", "A.B.C.D", "A.B.D" etc. but not "A.BB",
// "B.A.B" etc. If initialized with the empty string, all events are passed.
type NameFilter struct {
	name string
}

// Initialize a name filter.
// The name of the logger/handler is specified, all the events of
// logger's children are allowed through the filter. If no name is specified,
// every event is allowed.
func NewNameFilter(name string) *NameFilter {
	return &NameFilter{
		name: name,
	}
}

// Determine if the specified record is to be logged.
// Is the specified record to be logged? Returns false for no, true for yes.
// If deemed appropriate, the record may be modified in-place.
func (self *NameFilter) Filter(record *LogRecord) bool {
	length := len(self.name)
	if length == 0 {
		return true
	} else if self.name == record.Name {
		return true
	} else if !strings.HasPrefix(record.Name, self.name) {
		return false
	}
	return (record.Name[length] == '.')
}

// An interface for managing filters.
type Filterer interface {
	AddFilter(filter Filter)
	RemoveFilter(filter Filter)
	Filter(record *LogRecord) int
}

// An base class for loggers and handlers which allows them to share common code
// of managing the filters.
type StandardFilterer struct {
	filters *ListSet
}

// Initialize the standard filterer, with no filter.
func NewStandardFilterer() *StandardFilterer {
	return &StandardFilterer{
		filters: NewListSet(),
	}
}

// AddFilter adds the specified filter.
func (self *StandardFilterer) AddFilter(filter Filter) {
	if !self.filters.SetContains(filter) {
		self.filters.SetAdd(filter)
	}
}

// RemoveFilter removes the specified filter.
func (self *StandardFilterer) RemoveFilter(filter Filter) {
	if self.filters.SetContains(filter) {
		self.filters.SetRemove(filter)
	}
}

// Determine if a record is loggable by consulting all the filters.
// The default is to allow the record to be logged: any filter can veto
// this and the record is then dropped. Returns a zero value if a record
// is to be dropped, else non-zero.
func (self *StandardFilterer) Filter(record *LogRecord) int {
	recordVote := 1
	for e := self.filters.Front(); e != nil; e = e.Next() {
		filter, _ := e.Value.(Filter)
		if !filter.Filter(record) {
			recordVote = 0
			break
		}
	}
	return recordVote
}

// GetFilters returns all the filter in this filterer.
func (self *StandardFilterer) GetFilters() *ListSet {
	return self.filters
}
