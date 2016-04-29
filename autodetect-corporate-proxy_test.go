package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAutodetectProxy(t *testing.T) {
    proxy, err := autodetectProxy("google.com")
    assert.Equal(t, "", proxy)
    assert.NotNil(t, err)
}