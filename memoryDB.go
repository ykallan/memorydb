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

type MemoryStorage struct {
	count         int
	useLock       bool
	msLock        sync.Mutex
	baseInfoSlice []*BaseInfo
}

func (ms *MemoryStorage) Set(value interface{}, expire int) int {
	index := ms.generateIndex()
	baseInfo := &BaseInfo{
		Index: index,
		Age:   time.Now().Add(time.Duration(expire) * time.Second),
		Value: value,
	}
	ms.lock()
	ms.baseInfoSlice = append(ms.baseInfoSlice, baseInfo)
	ms.unlock()
	return index
}

func (ms *MemoryStorage) SetBatch(valueSlice []interface{}, expire int) []int {
	var indexSlice []int
	ms.lock()
	for _, value := range valueSlice {
		index := ms.generateIndex()
		baseInfo := &BaseInfo{
			Index: index,
			Age:   time.Now().Add(time.Duration(expire) * time.Second),
			Value: value,
		}
		ms.baseInfoSlice = append(ms.baseInfoSlice, baseInfo)
	}
	ms.unlock()
	return indexSlice
}

func (ms *MemoryStorage) Get(index int) interface{} {
	if len(ms.baseInfoSlice) == 0 {
		return nil
	}
	for _, baseInfo := range ms.baseInfoSlice {
		if baseInfo.Index == index {
			return baseInfo.Value
		}
	}
	return nil
}

func (ms *MemoryStorage) GetAll() []interface{} {
	var result []interface{}

	for _, baseInfo := range ms.baseInfoSlice {
		result = append(result, baseInfo.Value)
	}

	return result
}

func (ms *MemoryStorage) Remove(index int) bool {
	for index, baseInfo := range ms.baseInfoSlice {
		if baseInfo.Index == index {
			ms.lock()
			ms.baseInfoSlice = append(ms.baseInfoSlice[:index], ms.baseInfoSlice[index+1:]...)
			ms.unlock()
			return true
		}
	}
	return false
}

func (ms *MemoryStorage) Flush() {
	ms.baseInfoSlice = []*BaseInfo{}
}

func (ms *MemoryStorage) Update(newObject interface{}, index int) error {
	for _, baseInfo := range ms.baseInfoSlice {
		if baseInfo.Index == index {
			ms.lock()
			baseInfo.Value = newObject
			ms.unlock()
			return nil
		}
	}

	return errors.New("not found the current index of the database")
}

func (ms *MemoryStorage) Size() int {
	return len(ms.baseInfoSlice)
}

func (ms *MemoryStorage) Empty() bool {
	return len(ms.baseInfoSlice) == 0
}

func (ms *MemoryStorage) generateIndex() int {
	ms.lock()
	ms.count += 1
	ms.unlock()
	return ms.count
}

func (ms *MemoryStorage) computeIsTimeout(age time.Time) bool {
	nowTime := time.Now()
	return nowTime.After(age)
}

func (ms *MemoryStorage) filter() {
	for {
		for index, baseInfo := range ms.baseInfoSlice {
			if ms.computeIsTimeout(baseInfo.Age) {
				ms.lock()
				ms.baseInfoSlice = append(ms.baseInfoSlice[:index], ms.baseInfoSlice[index+1:]...)
				ms.unlock()
			}
		}
		time.Sleep(time.Second)
	}
}

func (ms *MemoryStorage) lock() {
	if ms.useLock {
		ms.msLock.Lock()
	}
}

func (ms *MemoryStorage) unlock() {
	if ms.useLock {
		ms.msLock.Unlock()
	}
}

func New() *MemoryStorage {
	ms := MemoryStorage{}
	go ms.filter()
	return &ms
}

func NewWithLock() *MemoryStorage {
	ms := MemoryStorage{
		useLock: true,
	}
	go ms.filter()
	return &ms
}
