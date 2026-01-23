package mycache
// lfu

import (
	"errors"
	"iter"

	
	// тут импорт linked list моего еще был, но с ссылкой на основной репо, его пока не оставляю
)

var ErrKeyNotFound = errors.New("key not found")

const DefaultCapacity = 5

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Put(key K, value V)
	All() iter.Seq2[K, V]
	Size() int
	Capacity() int
	GetKeyFrequency(key K) (int, error)
}
type nodeLfu[K comparable, V any] struct {
	key              K
	value            V
	freq             int
	freqClassElement *linkedlist.Node[*nodeLfu[K, V]]
}

type cacheImpl[K comparable, V any] struct {
	capacity     int
	minFreq      int
	keyToElement map[K]*nodeLfu[K, V]
	freqToClass  map[int]*linkedlist.List[*nodeLfu[K, V]]
	activeFreqs  []int
	freqsSet     map[int]struct{}
	needFreqSort bool
}

func New[K comparable, V any](capacity ...int) *cacheImpl[K, V] {
	nCapacity := DefaultCapacity
	if len(capacity) > 0 {
		if capacity[0] >= 0 {
			nCapacity = capacity[0]
		} else {
			panic("capacity must be >= 0")
		}
	}

	return &cacheImpl[K, V]{
		capacity:     nCapacity,
		minFreq:      0,
		keyToElement: make(map[K]*nodeLfu[K, V]),
		freqToClass:  make(map[int]*linkedlist.List[*nodeLfu[K, V]]),
		activeFreqs:  make([]int, 0),
		freqsSet:     make(map[int]struct{}),
		needFreqSort: false,
	}
}

func (l *cacheImpl[K, V]) Get(key K) (V, error) {
	node, exists := l.keyToElement[key]
	if !exists {
		var nothing V
		return nothing, ErrKeyNotFound
	}

	l.updNodeFreqClass(node)
	return node.value, nil
}

func (l *cacheImpl[K, V]) Put(key K, value V) {
	if node, exists := l.keyToElement[key]; exists {
		node.value = value
		l.updNodeFreqClass(node)
		return
	}

	if l.Size() >= l.capacity {
		l.reuseLFUnode(key, value)
	} else {
		node := &nodeLfu[K, V]{
			key:   key,
			value: value,
			freq:  1,
		}
		node.freqClassElement = &linkedlist.Node[*nodeLfu[K, V]]{Val: node}

		class := l.getFreqClass(1)
		class.PushFront(node.freqClassElement)
		l.keyToElement[key] = node
		l.minFreq = 1
	}
}

func (l *cacheImpl[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		l.sortFreqs()
		for _, freq := range l.activeFreqs {
			class := l.freqToClass[freq]
			if class == nil || class.Len() == 0 {
				continue
			}
			node := class.Front()
			for node != nil {
				lfuNode := node.Val
				if !yield(lfuNode.key, lfuNode.value) {
					return
				}
				node = node.Next()
			}
		}
	}
}

func (l *cacheImpl[K, V]) Size() int {
	return len(l.keyToElement)
}

func (l *cacheImpl[K, V]) Capacity() int {
	return l.capacity
}

func (l *cacheImpl[K, V]) GetKeyFrequency(key K) (int, error) {
	node, exists := l.keyToElement[key]
	if !exists {
		return 0, ErrKeyNotFound
	}
	return node.freq, nil
}

func (l *cacheImpl[K, V]) updNodeFreqClass(node *nodeLfu[K, V]) {
	if node.freqClassElement != nil {
		prevClass := l.getFreqClass(node.freq)
		prevClass.Remove(node.freqClassElement)
		if prevClass.Len() == 0 {
			l.forgetFreq(node.freq)
			delete(l.freqToClass, node.freq)
			if node.freq == l.minFreq {
				l.minFreq++
			}
		}
	}
	node.freq++
	newClass := l.getFreqClass(node.freq)
	if node.freqClassElement == nil {
		node.freqClassElement = &linkedlist.Node[*nodeLfu[K, V]]{Val: node}
	}
	newClass.InsertToFront(node.freqClassElement)
}

func (l *cacheImpl[K, V]) getFreqClass(freq int) *linkedlist.List[*nodeLfu[K, V]] {
	if class, exists := l.freqToClass[freq]; exists {
		return class
	}
	class := linkedlist.New[*nodeLfu[K, V]]()
	l.freqToClass[freq] = class
	l.registerFreq(freq)
	return class
}

func (l *cacheImpl[K, V]) reuseLFUnode(newKey K, newVal V) {
	if l.minFreq == 0 {
		return
	}
	nodeToReuse := l.extractLFUnode()
	if nodeToReuse == nil {
		return
	}
	l.reinitNode(nodeToReuse, newKey, newVal)
	newClass := l.getFreqClass(1)
	newClass.InsertToFront(nodeToReuse.freqClassElement)
	l.keyToElement[newKey] = nodeToReuse
	l.minFreq = 1
}

func (l *cacheImpl[K, V]) extractLFUnode() *nodeLfu[K, V] {
	class := l.getFreqClass(l.minFreq)
	extractingNode := class.Back()
	if extractingNode == nil {
		return nil
	}
	delete(l.keyToElement, extractingNode.Val.key)
	class.Remove(extractingNode)
	if class.Len() == 0 {
		l.forgetFreq(l.minFreq)
		l.updMinFreq()
	}
	return extractingNode.Val
}

func (l *cacheImpl[K, V]) updMinFreq() {
	if len(l.freqsSet) == 0 {
		l.minFreq = 0
		return
	}
	l.sortFreqs()
	for _, freq := range l.activeFreqs {
		if class, exists := l.freqToClass[freq]; exists && class.Len() != 0 && l.minFreq > freq {
			l.minFreq = freq
			return
		}
	}
	l.minFreq = 0
}
