package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct{
	hash Hash
	replicas int
	keys []int
	hashMap map[int]string
}

func New(relicas int,fn Hash) *Map{
	m:=&Map{
		replicas: relicas,
		hash:fn,
		hashMap: make(map[int]string),
	}
	if m.hash==nil{
		m.hash=crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string){
	for _,key:=range keys{
		for i:=0;i<m.replicas;i++{
			hash:=int( m.hash([]byte(strconv.Itoa(i)+key)))
			m.keys = append(m.keys,hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)

}
func (m *Map) Get(key string) string {
	if len(m.keys) ==0 {
		return ""
	}
	hash:=int(m.hash([]byte(key)))
	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	idx:=sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i]>=hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

