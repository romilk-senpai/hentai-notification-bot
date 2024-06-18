package tgcontroller

import "slices"

type TagMap struct {
	Tags []string                  `json:"tagGroup"`
	Data map[string]map[string]int `json:"parserInfo"`
}

func NewTagMap() *TagMap {
	return &TagMap{
		Tags: make([]string, 0),
		Data: make(map[string]map[string]int),
	}
}

func (m *TagMap) Get(key string) (map[string]int, bool) {
	value, exists := m.Data[key]
	return value, exists
}

func (m *TagMap) Set(key string, value map[string]int) {
	if _, exists := m.Data[key]; !exists {
		m.Tags = append(m.Tags, key)
	}
	m.Data[key] = value
}

func (m *TagMap) Delete(key string) {
	if _, exists := m.Data[key]; exists {
		delete(m.Data, key)
		for i, k := range m.Tags {
			if k == key {
				m.Tags = append(m.Tags[:i], m.Tags[i+1:]...)
				break
			}
		}
	}
}

func (m *TagMap) ForEach(f func(key string, value map[string]int) error) {
	for _, key := range m.Tags {
		if _, exists := m.Data[key]; exists {
			err := f(key, m.Data[key])
			if err != nil {
				return
			}
		}
	}
}

func (m *TagMap) SubscribedToTag(tagGroup string) bool {
	return slices.Contains(m.Tags, tagGroup)
}
