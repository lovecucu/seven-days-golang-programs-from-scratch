package singleflight

import "sync"

// 代表请求体，进行中或已完成
type call struct {
	wg  sync.WaitGroup // 用于控制并发请求
	val interface{}    // fn获取的值
	err error          // fn返回的error
}

type Group struct {
	mu sync.Mutex       // 防止并发读写m
	m  map[string]*call // 用于存储进行中的请求
}

// 并发执行
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 相同key已经有进行中的，直接等待结果
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 给key初始化一个call
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 执行fn，等待返回
	c.val, c.err = fn()
	c.wg.Done()

	// 完成后解锁，并删除map中的数据
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
