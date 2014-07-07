package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnknownGithubAccount(t *testing.T) {
	assert.True(t, true, "pending...")
}

func TestBadConfigFile(t *testing.T) {
	// verify an error is returned when bad syntax is used in the config file
	assert.True(t, true, "pending...")
}

func TestNoConfigFile(t *testing.T) {
	// ensure default are used and things remain operational
	assert.True(t, true, "pending...")
}

func TestAuthorization(t *testing.T) {
	// verify in two separate authorization cycles that the right keys are returned
	assert.True(t, true, "pending...")
}

func TestPanicMode(t *testing.T) {
	// make the API offline and ensure authorization continues
	assert.True(t, true, "pending...")
}
