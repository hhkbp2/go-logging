package logging

import (
    "github.com/deckarep/golang-set"
)

type Filter interface {
    Filter(record *LogRecord) bool
}

type DefaultFilter struct {
    name string
}

func NewDefaultFilter(name string) *DefaultFilter {
    return &DefaultFilter{
        name: name,
    }
}

func (self *DefaultFilter) Filter(record *LogRecord) bool {
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

type Filterer struct {
    filters mapset.Set
}

func NewFilterer() *Filterer {
    return &Filterer{
        filters: mapset.NewThreadUnsafeSet(),
    }
}

func (self *Filterer) AddFilter(filter Filter) {
    if !self.filters.Contains(filter) {
        self.filters.Add(filter)
    }
}

func (self *Filterer) RemoveFilter(filter Filter) {
    if self.filters.Contains(filter) {
        self.filters.Remove(filter)
    }
}

func (self *Filterer) Filter(record *LogRecord) int {
    recordVote = 1
    for filter := range self.filters {
        if !filter.Filter(record) {
            recordVote = 0
            break
        }
    }
    return recordVote
}
