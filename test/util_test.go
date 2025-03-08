package test

import (
	"mailbox-api/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCalculatePagination tests the pagination calculation utility
func TestCalculatePagination(t *testing.T) {
	// Test normal case
	page, pageSize, totalPages := util.CalculatePagination(2, 10, 25)
	assert.Equal(t, 2, page)
	assert.Equal(t, 10, pageSize)
	assert.Equal(t, 3, totalPages)

	// Test zero page case
	page, pageSize, totalPages = util.CalculatePagination(0, 10, 25)
	assert.Equal(t, 1, page) // Should default to 1
	assert.Equal(t, 10, pageSize)
	assert.Equal(t, 3, totalPages)

	// Test zero pageSize case
	page, pageSize, totalPages = util.CalculatePagination(2, 0, 25)
	assert.Equal(t, 2, page)
	assert.Equal(t, 10, pageSize) // Should default to 10
	assert.Equal(t, 3, totalPages)

	// Test empty result case
	page, pageSize, totalPages = util.CalculatePagination(1, 10, 0)
	assert.Equal(t, 1, page)
	assert.Equal(t, 10, pageSize)
	assert.Equal(t, 1, totalPages) // Should be at least 1
}

// TestParseInt tests the ParseInt utility
func TestParseInt(t *testing.T) {
	// Test valid number
	result := util.ParseInt("123", 0)
	assert.Equal(t, 123, result)

	// Test empty string
	result = util.ParseInt("", 42)
	assert.Equal(t, 42, result) // Should return default value

	// Test invalid number
	result = util.ParseInt("not-a-number", 42)
	assert.Equal(t, 42, result) // Should return default value
}
