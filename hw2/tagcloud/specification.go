package tagcloud

import (
	"sort"
)

const (
	tagOccurOnce = 1
)

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

type tagsArray []TagStat

func (t tagsArray) reOrder(tags map[string]*TagStat) {
	for i, tag := range t {
		t[i] = *tags[tag.Tag]
	}
	sort.Slice(
		t, func(i, j int) bool {
			return t[i].OccurrenceCount > t[j].OccurrenceCount
		},
	)
}

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tags    map[string]*TagStat
	ordered tagsArray
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{
		tags:    make(map[string]*TagStat, 0),
		ordered: make(tagsArray, 0),
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (c *TagCloud) AddTag(tag string) {
	tStat, ok := c.tags[tag]
	if !ok {
		c.tags[tag] = &TagStat{
			Tag:             tag,
			OccurrenceCount: tagOccurOnce,
		}
		c.ordered = append(c.ordered, *c.tags[tag])

		return
	}
	tStat.OccurrenceCount++
	c.ordered.reOrder(c.tags)
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (c *TagCloud) TopN(n int) []TagStat {
	if n > len(c.ordered) {
		return c.ordered
	}
	return c.ordered[:n]
}
