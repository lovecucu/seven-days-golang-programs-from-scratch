package gencache

import (
	"fmt"
	"gencache/singleflight"
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
	peers     PeerPicker
	loader    *singleflight.Group
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
		loader:    &singleflight.Group{},
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

// 注册peers
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
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
	// loader.Do保证只触发一次
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			// key从其他节点获取
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err) // 远程获取失败，会从本地再获取一次
			}
		}
		return g.getLocally(key) // 本地获取
	})

	if err == nil {
		return viewi.(ByteView), nil
	}

	return
}

// 远程获取
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
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
