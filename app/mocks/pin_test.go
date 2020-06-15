package mocks

import (
	"app/models"
	"reflect"
	"testing"
	"time"
)

func TestPinMock(t *testing.T) {
	ID := 0
	UserID := 0
	boardID := 0
	pins := &PinMock{}
	pins.ExpectedPins = make([]*models.Pin, 0)
	pins.BoardPinMapper = make(map[int][]int)
	pin := &models.Pin{
		ID:          ID,
		UserID:      UserID,
		Title:       "test title",
		Description: "test description",
		URL:         "test url",
		ImageURL:    "test image url",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	pins.CreatePin(pin, boardID)
	got, err := pins.GetPin(ID)
	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}
	if !reflect.DeepEqual(*pin, *got) {
		t.Fatalf("Not equal pin")
	}
}

func TestPinMockRepository(t *testing.T) {
	pins := NewPinRepository()
	ID := 0
	UserID := 0
	boardID := 0
	pin := &models.Pin{
		ID:          ID,
		UserID:      UserID,
		Title:       "test title",
		Description: "test description",
		URL:         "test url",
		ImageURL:    "test image url",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	pins.CreatePin(pin, boardID)
	got, err := pins.GetPin(ID)
	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}
	if !reflect.DeepEqual(*pin, *got) {
		t.Fatalf("Not equal pin")
	}
}
