package watchdb

//import "fmt"

type IntSet struct {
	set map[int]bool
}

func NewIntSet() *IntSet {
	return &IntSet{make(map[int]bool)}
}

func (set *IntSet) Add(i int) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found //False if it existed already
}

func (set *IntSet) Contains(i int) bool {
	_, found := set.set[i]
	return found //true if it existed already
}

func (set *IntSet) Remove(i int) {
	delete(set.set, i)
}

func (set *IntSet) Size() int {
	return len(set.set)
}
