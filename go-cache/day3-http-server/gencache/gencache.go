package gencache

import (
	"fmt"
	"log"
	"sync"
)

// 用于load某个key的数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// 实现了Getter接口的函数，这种称为接口型函数
type GetterFunc func(key string) ([]byte, error)

// Get实现了Getter接口的方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group是cache的命名空间，
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 获取group中key的值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GenCache] hit")
		return v, nil
	}

	return g.load(key)
}

// 加载缓存（可能本地获取，也可能从其它节点获取）
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 根据getter获取缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 写入缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
