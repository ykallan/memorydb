package memorydb

import (
	"errors"
	"sync"
	"time"
)

type BaseInfo struct {
	Index int
	Age   time.Time
	Value interface{}
}

type MemoryDataBase struct {
	count         int
	useLock       bool
	mdLock        sync.Mutex
	baseInfoSlice []*BaseInfo
}

func (md *MemoryDataBase) Set(value interface{}, expire int) int {
	index := md.generateIndex()
	baseInfo := &BaseInfo{
		Index: index,
		Age:   time.Now().Add(time.Duration(expire) * time.Second),
		Value: value,
	}
	md.lock()
	md.baseInfoSlice = append(md.baseInfoSlice, baseInfo)
	md.unlock()
	return index
}

func (md *MemoryDataBase) SetBatch(valueSlice []interface{}, expire int) []int {
	var indexSlice []int
	md.lock()
	for _, value := range valueSlice {
		index := md.generateIndex()
		baseInfo := &BaseInfo{
			Index: index,
			Age:   time.Now().Add(time.Duration(expire) * time.Second),
			Value: value,
		}
		md.baseInfoSlice = append(md.baseInfoSlice, baseInfo)
	}
	md.unlock()
	return indexSlice
}

func (md *MemoryDataBase) Get(index int) interface{} {
	if len(md.baseInfoSlice) == 0 {
		return nil
	}
	for _, baseInfo := range md.baseInfoSlice {
		if baseInfo.Index == index {
			return baseInfo.Value
		}
	}
	return nil
}

func (md *MemoryDataBase) GetAll() []interface{} {
	var result []interface{}

	for _, baseInfo := range md.baseInfoSlice {
		result = append(result, baseInfo.Value)
	}

	return result
}

func (md *MemoryDataBase) Remove(index int) bool {
	for _, baseInfo := range md.baseInfoSlice {
		if baseInfo.Index == index {
			md.lock()
			md.baseInfoSlice = append(md.baseInfoSlice[:index], md.baseInfoSlice[index+1:]...)
			md.unlock()
			return true
		}
	}
	return false
}

func (md *MemoryDataBase) Flush() {
	md.baseInfoSlice = []*BaseInfo{}
}

func (md *MemoryDataBase) Update(newObject interface{}, index int) error {
	for _, baseInfo := range md.baseInfoSlice {
		if baseInfo.Index == index {
			md.lock()
			baseInfo.Value = newObject
			md.unlock()
			return nil
		}
	}

	return errors.New("not found the current index of the database")
}

func (md *MemoryDataBase) Size() int {
	return len(md.baseInfoSlice)
}

func (md *MemoryDataBase) IsEmpty() bool {
	return len(md.baseInfoSlice) == 0
}

func (md *MemoryDataBase) generateIndex() int {
	md.lock()
	md.count += 1
	md.unlock()
	return md.count
}

func (md *MemoryDataBase) computeIsTimeout(age time.Time) bool {
	nowTime := time.Now()
	return nowTime.After(age)
}

func (md *MemoryDataBase) filter() {
	for {
		for index, baseInfo := range md.baseInfoSlice {
			if md.computeIsTimeout(baseInfo.Age) {
				md.lock()
				md.baseInfoSlice = append(md.baseInfoSlice[:index], md.baseInfoSlice[index+1:]...)
				md.unlock()
			}
		}
		time.Sleep(time.Second)
	}
}

func (md *MemoryDataBase) lock() {
	if md.useLock {
		md.mdLock.Lock()
	}
}

func (md *MemoryDataBase) unlock() {
	if md.useLock {
		md.mdLock.Unlock()
	}
}

func NewMemoryDataBase() *MemoryDataBase {
	md := MemoryDataBase{}
	go md.filter()
	return &md
}

func NewMemoryDataBaseWithLock() *MemoryDataBase {
	md := MemoryDataBase{
		useLock: true,
	}
	go md.filter()
	return &md
}
