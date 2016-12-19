package stringSet

import (
	"encoding/gob"
)

type Set interface{
	Add(i string) bool
	Get(i string) bool
	Remove(i string)
	Count() int
}

func init() {
	gob.Register(set{})
}

type set struct {
	Set map[string]bool
}

func New() Set {
	return &set{make(map[string]bool)}
}

func (set *set) Add(i string) bool {
	_, found := set.Set[i]
	set.Set[i] = true
	return !found //False if it existed already
}

func (set *set) Get(i string) bool {
	_, found := set.Set[i]
	return found //true if it existed already
}

func (set *set) Remove(i string) {
	delete(set.Set, i)
}

func (set *set) Count() int {
	return len(set.Set)
}