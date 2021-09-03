package four

import (
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/jncornett/vec"
	"github.com/stretchr/testify/assert"
)

func TestNode_Next(t *testing.T) {
	root := new(Node)
	_, ok := root.Next(vec.Scalar(0))
	assert.False(t, ok)
	root.Split()
	next, ok := root.Next(vec.Scalar(0))
	assert.True(t, ok)
	assert.NotNil(t, next)
}

func TestNode_Len(t *testing.T) {
	root := new(Node)
	assert.Equal(t, 0, root.Len())
	root.Add("a", vec.NewValue(1, 2))
	assert.Equal(t, 1, root.Len())
	root.Add("b", vec.NewValue(2, 2))
	assert.Equal(t, 2, root.Len())
	root.Add("c", vec.NewValue(3, 3))
	assert.Equal(t, 3, root.Len())
	_ = root.Balance(0, 2)
	assert.Equal(t, 3, root.Len())
}

func TestNode_Split(t *testing.T) {
	root := new(Node)
	root.Split()
	root.Split()

	root = new(Node)
	root.Add("a", vec.NewValue(1, 2))
	root.Add("b", vec.NewValue(1, 3))
	root.Add("c", vec.NewValue(-1, 3))

	var items []Item
	root.Each(func(i Item) bool {
		items = append(items, i)
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, []Item{
		{"a", vec.NewValue(1, 2)},
		{"b", vec.NewValue(1, 3)},
		{"c", vec.NewValue(-1, 3)},
	}, items)
	root.Split()

	items = nil
	root.Each(func(i Item) bool {
		items = append(items, i)
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, []Item{
		{"a", vec.NewValue(1, 2)},
		{"b", vec.NewValue(1, 3)},
		{"c", vec.NewValue(-1, 3)},
	}, items)

	next, ok := root.Next(vec.NewValue(1, 2))
	assert.True(t, ok)
	items = next.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Contains(t, items, Item{"a", vec.NewValue(1, 2)})
}

func TestNode_Merge(t *testing.T) {
	root := new(Node)
	root.Split()
	root.Merge()
	root.Merge()

	root = new(Node)
	root.Add("a", vec.NewValue(1, 2))
	root.Add("b", vec.NewValue(1, 3))
	root.Add("c", vec.NewValue(-1, 3))

	var items []Item
	root.Each(func(i Item) bool {
		items = append(items, i)
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, []Item{
		{"a", vec.NewValue(1, 2)},
		{"b", vec.NewValue(1, 3)},
		{"c", vec.NewValue(-1, 3)},
	}, items)

	root.Split()

	items = nil
	root.Each(func(i Item) bool {
		items = append(items, i)
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, []Item{
		{"a", vec.NewValue(1, 2)},
		{"b", vec.NewValue(1, 3)},
		{"c", vec.NewValue(-1, 3)},
	}, items)

	root.Merge()

	items = nil
	root.Each(func(i Item) bool {
		items = append(items, i)
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, []Item{
		{"a", vec.NewValue(1, 2)},
		{"b", vec.NewValue(1, 3)},
		{"c", vec.NewValue(-1, 3)},
	}, items)
}

func TestNode_Add(t *testing.T) {
	root := new(Node)
	root.Add("a", vec.NewValue(1, 1))
	root.Add("a", vec.NewValue(1, 1))
	root.Add("a", vec.NewValue(1, 2))

	items := root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 1, root.Len())
	assert.Equal(t, []Item{{"a", vec.NewValue(1, 2)}}, items)

	root.Split()
	root.Add("b", vec.NewValue(-1, -1))
	items = root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 2, root.Len())
	assert.Equal(t, []Item{{"a", vec.NewValue(1, 2)}, {"b", vec.NewValue(-1, -1)}}, items)
}

func TestNode_Del(t *testing.T) {
	root := new(Node)
	root.Add("a", vec.NewValue(1, 1))
	root.Add("b", vec.NewValue(1, 1))
	root.Add("c", vec.NewValue(1, 2))

	items := root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 3, root.Len())
	assert.Equal(t, []Item{{"a", vec.NewValue(1, 1)}, {"b", vec.NewValue(1, 1)}, {"c", vec.NewValue(1, 2)}}, items)

	root.Del("d", vec.NewValue(1, 1))
	root.Del("a", vec.NewValue(1, 2))

	items = root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 3, root.Len())
	assert.Equal(t, []Item{{"a", vec.NewValue(1, 1)}, {"b", vec.NewValue(1, 1)}, {"c", vec.NewValue(1, 2)}}, items)

	root.Del("a", vec.NewValue(1, 1))

	items = root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 2, root.Len())
	assert.Equal(t, []Item{{"b", vec.NewValue(1, 1)}, {"c", vec.NewValue(1, 2)}}, items)

	root.Split()

	root.Del("d", vec.NewValue(1, 1))
	root.Del("b", vec.NewValue(1, 2))

	items = root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 2, root.Len())
	assert.Equal(t, []Item{{"b", vec.NewValue(1, 1)}, {"c", vec.NewValue(1, 2)}}, items)

	root.Del("b", vec.NewValue(1, 1))

	items = root.Slice()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key.(string) < items[j].Key.(string)
	})
	assert.Equal(t, 1, root.Len())
	assert.Equal(t, []Item{{"c", vec.NewValue(1, 2)}}, items)
}

func TestNode_Query(t *testing.T) {
	root := new(Node)
	root.Split()

	root.Add(0, vec.NewValue(-1, -1))
	root.Add(1, vec.NewValue(-1, 1))
	root.Add(2, vec.NewValue(1, 1))
	root.Add(3, vec.NewValue(1, -1))

	tests := []struct {
		r    vec.Rect
		want []Item
	}{
		{
			r: vec.Rect{
				vec.Scalar(math.Inf(-1)),
				vec.Scalar(0),
			},
			want: []Item{{0, vec.Scalar(-1)}},
		},
		{
			r: vec.Rect{
				vec.Scalar(math.Inf(-1)),
				vec.Scalar(-0.9),
			},
			want: []Item{{0, vec.Scalar(-1)}},
		},
		{
			r: vec.Inf(),
			want: []Item{
				{0, vec.Scalar(-1)},
				{1, vec.NewValue(-1, 1)},
				{2, vec.Scalar(1)},
				{3, vec.NewValue(1, -1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.r), func(t *testing.T) {
			var got []Item
			root.Query(tt.r, func(i Item) bool {
				got = append(got, i)
				return true
			})
			sort.Slice(got, func(i, j int) bool { return got[i].Key.(int) < got[j].Key.(int) })
			assert.Equal(t, tt.want, got)
		})
	}

	// var query []Item
	// root.Query(vec.Rect{vec.Scalar(math.Inf(-1)), vec.Scalar(0)}, func(i Item) bool {
	// 	query = append(query, i)
	// 	return true
	// })
	// assert.Equal(t, []Item{{0, vec.NewValue(-1, -1)}}, query)

	// query = nil
	// root.Query(vec.Rect{vec.NewValue(math.Inf(-1), 0), vec.NewValue(0, math.Inf(1))}, func(i Item) bool {
	// 	query = append(query, i)
	// 	return true
	// })
	// assert.Equal(t, []Item{{1, vec.NewValue(-1, 1)}}, query)

	// query = nil
	// root.Query(vec.Inf(), func(i Item) bool {
	// 	query = append(query, i)
	// 	return false
	// })
	// assert.Equal(t, 1, len(query))
}
