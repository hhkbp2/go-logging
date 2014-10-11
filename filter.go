package logging

import (
    "github.com/deckarep/golang-set"
    "strings"
    "sync"
)

type Filter interface {
    Filter(record *LogRecord) bool
}

type NameFilter struct {
    name string
}

func NewNameFilter(name string) *NameFilter {
    return &NameFilter{
        name: name,
    }
}

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

type Filterer interface {
    AddFilter(filter Filter)
    RemoveFilter(filter Filter)
    Filter(record *LogRecord) int
}

type StandardFilterer struct {
    filters mapset.Set
    lock    sync.RWMutex
}

func NewStandardFilterer() *StandardFilterer {
    return &StandardFilterer{
        filters: mapset.NewThreadUnsafeSet(),
    }
}

func (self *StandardFilterer) AddFilter(filter Filter) {
    self.lock.Lock()
    defer self.lock.Unlock()
    if !self.filters.Contains(filter) {
        self.filters.Add(filter)
    }
}

func (self *StandardFilterer) RemoveFilter(filter Filter) {
    self.lock.Lock()
    defer self.lock.Unlock()
    if self.filters.Contains(filter) {
        self.filters.Remove(filter)
    }
}

func (self *StandardFilterer) Filter(record *LogRecord) int {
    self.lock.RLock()
    defer self.lock.RUnlock()
    recordVote := 1
    for i := range self.filters.Iter() {
        filter, _ := i.(Filter)
        if !filter.Filter(record) {
            recordVote = 0
            break
        }
    }
    return recordVote
}

func (self *StandardFilterer) GetFilters() []Filter {
    self.lock.Lock()
    defer self.lock.Unlock()
    result := make([]Filter, 0, self.filters.Cardinality())
    for i := range self.filters.Iter() {
        filter, _ := i.(Filter)
        result = append(result, filter)
    }
    return result
}
