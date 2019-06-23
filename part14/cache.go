package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

var PreAllocation = 10 //预分配的Data元素个数

//定义Cache接口的方法
type Cache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}, time.Duration) error
	Del(string) error
	Count() int
	Incr(key string, step int64) (int64, error)
	Expire(key string, duration time.Duration) error
}

type LocalCache struct {
	Data map[string]Item
	sync.Mutex
}

type Item struct {
	Object     interface{}
	Expiration int64
}

func (l *LocalCache) Get(key string) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()

	item, ok := l.Data[key]
	if !ok {
		return nil, false
	}

	//check item expired
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		delete(l.Data, key)
		return nil, false
	}

	return item.Object, true
}

func (l *LocalCache) Set(key string, val interface{}, d time.Duration) error {
	l.Lock()
	defer l.Unlock()

	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	l.Data[key] = Item{
		Object:     val,
		Expiration: e,
	}

	return nil
}

func (l *LocalCache) Del(key string) error {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.Data[key]; ok {
		delete(l.Data, key)
	}

	return nil
}

func (l *LocalCache) Count() int {
	l.Lock()
	defer l.Unlock()
	return len(l.Data)
}

func (l *LocalCache) Incr(key string, step int64) (int64, error) {
	if key == "" {
		return 0, errors.New("key is empty")
	}

	l.Lock()
	defer l.Unlock()

	if item, ok := l.Data[key]; ok {
		cnt := item.Object.(int64)
		cnt += step
		item.Object = cnt
		l.Data[key] = item

		return cnt, nil
	}

	l.Data[key] = Item{
		Object:     1,
		Expiration: 0,
	}

	return 1, nil
}

func NewCache(cacheType string) (Cache, error) {
	var c Cache

	switch cacheType {
	case "local":
		items := make(map[string]Item, PreAllocation)
		c = &LocalCache{
			Data: items,
		}
	default:
		return nil, errors.New("invalid cache type")
	}

	return c, nil
}

func (l *LocalCache) Expire(key string, d time.Duration) error {
	l.Lock()
	defer l.Unlock()

	item, ok := l.Data[key]

	if !ok {
		return errors.New("current key: " + key + " not exist!")
	}

	var e int64
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	item.Expiration = e
	l.Data[key] = item
	return nil
}

func main() {
	c, _ := NewCache("local")
	c.Set("myname", "daheige", 30*time.Second)

	time.Sleep(2 * time.Second)
	log.Println(c.Get("myname"))
	log.Println(c.Count())
}
