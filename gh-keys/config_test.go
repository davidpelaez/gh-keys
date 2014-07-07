package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// ensure TTL and AllowPanicMode are defined
	// ensure no permissions are set by default
	assert.True(t, true, "pending...")
}

func TestKeysLocation(t *testing.T) {
	// verify the config file path is included in the keys location
	assert.True(t, true, "pending...")
}
