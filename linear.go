package wavx

import "fmt"

type linearBufferElt struct {
	next  *linearBufferElt
	value float64
}

type LinearBuffer struct {
	first   *linearBufferElt
	last    *linearBufferElt
	maxSize int
	size    int
}

func NewLinearBuffer(maxSize int) *LinearBuffer {
	return &LinearBuffer{
		maxSize: maxSize,
		size:    0,
	}
}

func (b *LinearBuffer) IsEmpty() bool {
	return b.first == nil
}

func (b *LinearBuffer) FirstValue() float64 {
	if b.first == nil {
		return 0
	}
	return b.first.value
}

func (b *LinearBuffer) ChangeMaxSize(n int) {
	if n < 2 {
		panic(fmt.Sprintf("max size < 2: %d", n))
	}

	b.maxSize = n
	for b.size > b.maxSize {
		b.first = b.first.next
		b.size--
	}
}

func (b *LinearBuffer) Push(v float64) {
	if b.first == nil {
		b.first = &linearBufferElt{
			next:  nil,
			value: v,
		}
		b.last = b.first
	} else {
		elt := &linearBufferElt{
			next:  nil,
			value: v,
		}
		b.last.next = elt
		b.last = elt
	}
	b.size++
	for b.size > b.maxSize {
		b.first = b.first.next
		b.size--
	}
}
