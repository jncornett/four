package four

import (
	"fmt"
	"math"
	"sync"

	"github.com/jncornett/vec"
)

type Node struct {
	entries map[interface{}]vec.Value
	c       vec.Value
	μ       vec.Mass
	next    *[4]*Node
	once    sync.Once
}

type Item struct {
	Key interface{}
	Pos vec.Value
}

func (n *Node) Next(pos vec.Value) (next *Node, ok bool) {
	if n.next == nil {
		return nil, false
	}
	return n.next[dir(n.c, pos)], true
}

func (n *Node) Len() int {
	if n.next == nil {
		return len(n.entries)
	}
	len := 0
	for _, next := range n.next {
		len += next.Len()
	}
	return len
}

func (n *Node) Split() {
	if n.next != nil {
		return
	}
	n.c = n.μ.Center()
	if n.c.IsNaN() {
		n.c = vec.Scalar(0)
	}
	n.μ = vec.Mass{}
	n.next = new([4]*Node)
	for i := range n.next {
		n.next[i] = new(Node)
	}
	for k, pos := range n.entries {
		delete(n.entries, k)
		next := n.next[dir(n.c, pos)]
		next.init()
		next.μ = vec.Combine(next.μ, vec.Mass{Sum: pos, N: 1})
		next.entries[k] = pos
	}
	n.μ = vec.Mass{}
}

func (n *Node) Merge() {
	if n.next == nil {
		return
	}
	nexts := n.next
	n.next = nil
	n.μ = vec.Mass{}
	n.init()
	for _, next := range nexts {
		next.Each(func(i Item) bool {
			n.μ = vec.Combine(n.μ, vec.Mass{Sum: i.Pos, N: 1})
			n.entries[i.Key] = i.Pos
			return true
		})
	}
}

func (n *Node) Add(k interface{}, pos vec.Value) {
	if next, ok := n.Next(pos); ok {
		next.Add(k, pos)
		return
	}
	old, ok := n.entries[k]
	if ok {
		if pos == old {
			return
		}
		n.μ = vec.Subtract(n.μ, vec.Mass{Sum: old, N: 1})
	}
	n.init()
	n.μ = vec.Combine(n.μ, vec.Mass{Sum: pos, N: 1})
	n.entries[k] = pos
}

func (n *Node) Del(k interface{}, pos vec.Value) {
	if next, ok := n.Next(pos); ok {
		next.Del(k, pos)
		return
	}
	if old, ok := n.entries[k]; !ok || pos != old {
		return
	}
	delete(n.entries, k)
	n.μ = vec.Subtract(n.μ, vec.Mass{Sum: pos, N: 1})
}

func (n *Node) Mov(k interface{}, from, to vec.Value) {
	if from == to {
		n.Add(k, to) // intentionally not a no-op
		return
	}
	if next, ok := n.Next(from); ok {
		next2, ok := n.Next(to)
		if !ok {
			panic("four: invariant: node is either a leaf or nonleaf")
		}
		if next == next2 {
			next.Mov(k, from, to)
			return
		}
	}
	n.Del(k, from)
	n.Add(k, to)
}

func (n *Node) Query(box vec.Rect, fn func(Item) bool) bool {
	if n.next != nil {
		for i, quad := range quadrants(n.c) {
			if vec.Intersects(quad, box) {
				if !n.next[i].Query(box, fn) {
					return false
				}
			}
		}
		return true
	}
	for k, pos := range n.entries {
		if box.Contains(pos) {
			if !fn(Item{Key: k, Pos: pos}) {
				return false
			}
		}
	}
	return true
}

func (n *Node) Balance(min, max int) int {
	if n.next != nil {
		var items int
		for _, next := range n.next {
			items += next.Balance(min, max)
		}
		if items < min {
			n.Merge()
		}
		return items
	}
	items := len(n.entries)
	if items > max {
		n.Split()
	}
	return items
}

func (n *Node) Each(fn func(Item) bool) bool {
	if n.next == nil {
		for k, pos := range n.entries {
			if !fn(Item{Key: k, Pos: pos}) {
				return false
			}
		}
		return true
	}
	for _, next := range n.next {
		if !next.Each(fn) {
			return false
		}
	}
	return true
}

func (n *Node) Slice() []Item {
	var out []Item
	n.Each(func(i Item) bool {
		out = append(out, i)
		return true
	})
	return out
}

func (n *Node) init() {
	n.once.Do(func() {
		n.entries = make(map[interface{}]vec.Value)
	})
}

func quadrants(center vec.Value) [4]vec.Rect {
	return [4]vec.Rect{
		{vec.Scalar(math.Inf(-1)), center},
		{vec.NewValue(math.Inf(-1), center.Y()), vec.NewValue(center.X(), math.Inf(1))},
		{center, vec.Scalar(math.Inf(1))},
		{vec.NewValue(center.X(), math.Inf(-1)), vec.NewValue(math.Inf(1), center.Y())},
	}
}

type Dir int

const (
	NW Dir = iota
	SW
	SE
	NE
)

func dir(v, w vec.Value) Dir {
	w -= v // translate to the right basis
	switch {
	default:
		panic(fmt.Errorf("vec: invariant: all regions not covered for scaled value %v", w))
	case w.X() < 0 && w.Y() < 0:
		return NW
	case w.X() < 0 && w.Y() >= 0:
		return SW
	case w.X() >= 0 && w.Y() >= 0:
		return SE
	case w.X() >= 0 && w.Y() < 0:
		return NE
	}
}
