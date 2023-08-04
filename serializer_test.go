package lqtserializer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"slug"`
	Thumbnail string    `json:"thumbnail"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
}

type ArticleSerializer struct {
	Model Article
}

func (ar ArticleSerializer) Fields() []string {
	return []string{"title", "content", "slug"}
}

func Contains(_slice []map[string]any, element map[string]any) bool {
	for i := 0; i < len(_slice); i++ {
		if reflect.DeepEqual(_slice[i], element) {
			return true
		}
	}

	return false
}

func isExpected(data []map[string]any, expected []map[string]any) bool {
	var equal bool = true

	if len(data) != len(expected) {
		equal = false
	} else {
		for i := 0; i < len(data); i++ {
			ele := data[i]

			if !Contains(expected, ele) {
				equal = false
				break
			}
		}
	}

	return equal
}

func TestSerialier(t *testing.T) {
	id, _ := uuid.NewRandom()
	articles := []Article{
		{
			ID:        id,
			Title:     "Testing",
			Content:   "Testing",
			Slug:      "Testing",
			Thumbnail: "abcdef.png",
			Type:      "normal",
			Status:    "draft",
		},
	}

	s := New(&[]ArticleSerializer{}, articles)
	data := s.GetData()

	fmt.Println("---data", data)

	expected := []map[string]any{
		{
			"id":        id,
			"title":     "Testing",
			"content":   "Testing",
			"slug":      "Testing",
			"thumbnail": "abcdef.png",
			"type":      "normal",
			"status":    "draft",
		},
	}

	var expected2 []map[string]any

	for _, ele := range expected {
		value := SliceMaps(ele, []string{"title", "content", "slug"})

		expected2 = append(expected2, value)
	}

	if isExpected(*data, expected) {
		t.Errorf("Output not equal")
	}

	if !isExpected(*data, expected2) {
		t.Errorf("Output not equal")
	}
}
