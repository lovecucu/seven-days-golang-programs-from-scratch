package gencache

import (
	"lru/lru"
	"sync"
)

type cache struct {
	mu  sync.Mutex
	lru *lru.Cache
}
