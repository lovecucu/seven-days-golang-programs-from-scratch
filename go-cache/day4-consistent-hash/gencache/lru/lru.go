package lru

import "container/list"

type Cache struct {
	// 缓存实体占用的最大内存
	maxBytes int64

	// 缓存实体已占用的内存
	nbytes int64

	// 删除过期key时的回调（可选）
	OnEvicted func(key string, value Value)

	// 双向链表，用于存储缓存实体
	ll *list.List

	// map用于存储缓存key和对应的缓存值
	cache map[interface{}]*list.Element
}

// 定义缓存实体的结构（仅包内使用）
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		OnEvicted: onEvicted,
		cache:     make(map[interface{}]*list.Element),
	}
}

// 添加缓存，已存在则更新值
func (c *Cache) Add(key string, value Value) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		kv := ee.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) // 计算value的内存占用差值
		kv.value = value
		return
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len()) // 计算新key,value的内存占用
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes { // 淘汰最少访问的节点
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele) // 移动队首
		return ele.Value.(*entry).value, true
	}
	return
}

// 删除某个key
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// 删除最旧的缓存
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

// 删除list的指定Element（仅包内使用）
func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// 获取缓存个数
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// 清空缓存
func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
	c.maxBytes = 0
	c.nbytes = 0
}
