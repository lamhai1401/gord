package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	client := NewClient(&Config{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
		Pwd:   "",
	})

	assert.NotNil(t, client, "Redis client is nil")
}

func TestSet(t *testing.T) {
	client := NewClient(&Config{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
		Pwd:   "",
	})

	err := client.Set("test", "test", 3*time.Second)
	assert.Nil(t, err, fmt.Sprintf("Err is not nil %v", err))
}

func TestGet(t *testing.T) {
	client := NewClient(&Config{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
		Pwd:   "",
	})

	err := client.Set("test", "test", 3*time.Second)
	v := client.Get("test")

	assert.NotEmpty(t, v, fmt.Sprintf("%s is empty values", "test"))
	assert.Nil(t, err, fmt.Sprintf("Err is not nil %v", err))
}
