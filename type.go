package skiplist

import (
	"math/rand"
	"sync"
)

type elementNode struct {
	next []*Element
}

type Element struct {
	elementNode
	key   float64
	value interface{}
}

// 返回节点的key
func (elem *Element) Key() float64 {
	return elem.key
}

// 返回节点的值
func (elem *Element) Value() interface{} {
	return elem.value
}

// 返回下一个节点的下一个节点，当为末尾节点时，返回nil,注意：返回的是list最底层的下一个
func (elem *Element) Next() *Element {
	return elem.next[0]
}

type SkipList struct {
	elementNode
	maxLevel       int
	Length         int
	randSource     rand.Source
	probability    float64
	probTable      []float64
	mutex          sync.RWMutex
	prevNodesCache []*elementNode
}
