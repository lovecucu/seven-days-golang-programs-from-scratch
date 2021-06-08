package gencache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db)) // 用于记录初始化的次数，从而明确缓存是否生效
	gen := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gen.Get(k); err != nil || view.String() != v { // 初始化失败或初始的值不及预期，则报错
			t.Fatalf("failed to get value of %s", k)
		}
		if _, err := gen.Get(k); err != nil || loadCounts[k] > 1 { // 获取缓存报错或重复初始化缓存，则报错
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := gen.Get("unknown"); err == nil { // 无法初始化的缓存，直接报错
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestGetGroup(t *testing.T) {
	groupName := "scores"
	NewGroup(groupName, 2<<10, GetterFunc(
		func(key string) (bytes []byte, err error) { return }))
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s not exist", groupName)
	}

	if group := GetGroup(groupName + "111"); group != nil {
		t.Fatalf("expect nil, but %s got", group.name)
	}
}

func TestMapDel(t *testing.T) {
	delete(db, "Tom")
	for k, v := range db {
		log.Printf("key %s, value %s", k, v)
	}
}
