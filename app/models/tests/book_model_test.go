package models_test

import (
	"encoding/json"
	"testing"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/stretchr/testify/assert"
)

func TestBookAttrs_Value(t *testing.T) {
	attrs := models.BookAttrs{
		Picture:     "cover.jpg",
		Description: "A great book",
		Rating:      8,
	}

	val, err := attrs.Value()
	assert.NoError(t, err)

	bytes, ok := val.([]byte)
	assert.True(t, ok)

	var unmarshaled models.BookAttrs
	err = json.Unmarshal(bytes, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, attrs, unmarshaled)
}

func TestBookAttrs_Scan_Success(t *testing.T) {
	original := models.BookAttrs{
		Picture:     "image.png",
		Description: "description text",
		Rating:      5,
	}

	data, _ := json.Marshal(original)

	var scanned models.BookAttrs
	err := scanned.Scan(data)
	assert.NoError(t, err)
	assert.Equal(t, original, scanned)
}

func TestBookAttrs_Scan_InvalidType(t *testing.T) {
	var scanned models.BookAttrs
	err := scanned.Scan("not-a-byte-slice")
	assert.Error(t, err)
	assert.Equal(t, "type assertion to []byte failed", err.Error())
}

func TestBookAttrs_Scan_InvalidJSON(t *testing.T) {
	var scanned models.BookAttrs
	invalidJSON := []byte(`{"picture": "x", "rating": "not-a-number"}`)
	err := scanned.Scan(invalidJSON)
	assert.Error(t, err)
}
