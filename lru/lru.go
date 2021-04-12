package lru

import "container/list"

type Cache struct{
	//0为无限制
	maxBytes int64
	nBytes int64
	ll *list.List
	//在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
	cache map[string]*list.Element

	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}
type entry struct{
	key string
	value Value
}
type Value interface {
	Len() int
}

func New(maxBytes int64,onEvicted func(string,Value)) *Cache{
	return &Cache{
		maxBytes: maxBytes,
		ll:list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (v Value,ok bool){
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		return kv.value,true
	}
	return
}
func (c *Cache) RemoveOldest() {
	ele:=c.ll.Back()
	if ele!=nil{
		c.ll.Remove(ele)
		kv:=ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nBytes = c.nBytes - int64(int64(len(kv.key))+int64(kv.value.Len()))
		if c.OnEvicted!=nil{
			c.OnEvicted(kv.key,kv.value)
		}

	}
}

func (c *Cache) Add(key string, v Value) {
	if ele,ok:=c.cache[key];ok{
		c.ll.MoveToFront(ele)
		kv:=ele.Value.(*entry)
		c.nBytes += int64(v.Len()) - int64(kv.value.Len())
		kv.value =v
	}else{
		ele:=c.ll.PushFront(&entry{key,v})
		c.cache[key] = ele
		c.nBytes += int64(len(key))+int64(v.Len())
	}
	for c.maxBytes!=0 && c.maxBytes<c.nBytes{
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int{
	return c.ll.Len()
}

