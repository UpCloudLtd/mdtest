package id_test

import (
	"testing"

	"github.com/UpCloudLtd/mdtest/id"
	"github.com/stretchr/testify/assert"
)

func TestNewTestID(t *testing.T) {
	t.Parallel()
	assert.NotEqual(t, id.NewTestID(), id.NewTestID()) //nolint:testifylint // Validate both calls return unique values
}

func TestNewRunID(t *testing.T) {
	t.Parallel()
	assert.NotEqual(t, id.NewRunID(), id.NewRunID()) //nolint:testifylint // Validate both calls return unique values
}
