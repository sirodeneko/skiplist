package skiplist

import (
	"math"
	"math/rand"
	"time"
)

const (
	// Suitable for math.Floor(math.Pow(math.E, 18)) == 65659969 elements in list
	DefaultMaxLevel    int     = 16
	DefaultProbability float64 = 1 / math.E
)

// Front returns the head node of the list.
// 返回list的头节点
func (list *SkipList) Front() *Element {
	return list.next[0]
}

// 放入一个值，如果存在，则更新。同时，返回这个节点
func (list *SkipList) Set(key float64, value interface{}) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var element *Element
	prevs := list.getPrevElementNodes(key)

	// 通过getPrevElementNodes可查询到每一层的前驱节点，0即代表最底层，即链表上当前节点的前驱
	// 因为getPrevElementNodes查询出来的元素的.next的key一定是大于或等于key，这个判断是处理等于的情况，故可以直接覆盖
	// 同时，如果插入的值是位于第一个元素，element为nil,将不执行本条件，
	// (应该！！！，不可能出现element.key小于key)
	if element = prevs[0].next[0]; element != nil && element.key <= key {
		element.value = value
		return element
	}

	// 处理需要进行插值
	element = &Element{
		elementNode: elementNode{
			next: make([]*Element, list.randLevel()),
		},
		key:   key,
		value: value,
	}

	for i := range element.next {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}

	list.Length++
	return element
}

// 通过key查找元素。如果找到，返回元素指针，如果找不到，则返回nil
func (list *SkipList) Get(key float64) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	var prev = &list.elementNode
	var next *Element

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && key > next.key {
			prev = &next.elementNode
			next = next.next[i]
		}
	}

	if next != nil && next.key <= key {
		return next
	}

	return nil
}

// 从list中删除一个元素。
// 返回删除的元素指针（如果找到），否则返回nil。
func (list *SkipList) Remove(key float64) *Element {
	list.mutex.Lock()
	defer list.mutex.Unlock()
	prevs := list.getPrevElementNodes(key)

	// found the element, remove it
	if element := prevs[0].next[0]; element != nil && element.key <= key {
		for k, v := range element.next {
			prevs[k].next[k] = v
		}

		list.Length--
		return element
	}

	return nil
}

// getPrevElementNodes是其他函数使用的私有搜索机制。
// 查找相对于当前Element的每个级别上的先前节点，并将其缓存
func (list *SkipList) getPrevElementNodes(key float64) []*elementNode {
	var prev = &list.elementNode
	var next *Element

	prevs := list.prevNodesCache

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]

		// 只有key大于next.key才会向右，故返回的前驱节点一定小当前key(不可能等于，因为早在上一个循环就退出了)，
		// 后驱(.next)一定是大于等于当前key或者为nil
		for next != nil && key > next.key {
			prev = &next.elementNode
			next = next.next[i]
		}

		prevs[i] = prev
	}

	return prevs
}

// 改变probability
func (list *SkipList) SetProbability(newProbability float64) {
	list.probability = newProbability
	list.probTable = probabilityTable(list.probability, list.maxLevel)
}

// 产生随机层
func (list *SkipList) randLevel() (level int) {
	// 生成一个0~1的随机数
	r := float64(list.randSource.Int63()) / (1 << 63)

	level = 1
	for level < list.maxLevel && r < list.probTable[level] {
		level++
	}
	return
}

// 预先计算每一层的概率
// probability[0,1]  MaxLevel[1,64]
func probabilityTable(probability float64, MaxLevel int) (table []float64) {
	for i := 1; i <= MaxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

// 创建一个新的跳过列表，并将MaxLevel设置为提供的数字。
// maxLevel推荐为int(math.Ceil(math.Log(N))) (其中N是跳过列表中元素数量的上限)。
func NewWithMaxLevel(maxLevel int) *SkipList {
	if maxLevel < 1 || maxLevel > 64 {
		panic("maxLevel for a SkipList must be a positive integer <= 64")
	}

	return &SkipList{
		elementNode:    elementNode{next: make([]*Element, maxLevel)},
		prevNodesCache: make([]*elementNode, maxLevel),
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    DefaultProbability,
		probTable:      probabilityTable(DefaultProbability, maxLevel),
	}
}

// 采用默认maxLevel创建一个跳跃链表
func New() *SkipList {
	return NewWithMaxLevel(DefaultMaxLevel)
}
