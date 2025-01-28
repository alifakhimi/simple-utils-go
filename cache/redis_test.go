package cache

import (
	"context"
	"testing"
	"time"

	simutils "github.com/alifakhimi/simple-utils-go"
)

type TestStruct struct {
	ID   int `sim:"primaryKey"`
	Name string
}

func TestRedisCache(t *testing.T) {
	_, err := Connect("redis://localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	ctx := context.Background()
	item := TestStruct{ID: 1, Name: "Test Item"}
	key := simutils.GetTKey(item)

	// Test Set
	if err := Set(ctx, item, 10*time.Minute); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test Get
	var result TestStruct
	if err := Get(ctx, key, &result); err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if result != item {
		t.Errorf("Expected %+v, got %+v", item, result)
	}

	// Test Del
	if err := Del(ctx, key); err != nil {
		t.Fatalf("Failed to delete cache: %v", err)
	}

	if err := Get(ctx, key, &result); err == nil {
		t.Fatalf("Expected error for missing key, got none")
	}
}
