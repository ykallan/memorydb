# memorydb
一款基于内存的，无需第三方数据库配置的，带有过期时间的数据库 db-base-memory-with-expiretime


### 1.安装

```shell
go get github.com/ykallan/memorydb
```

### 2.使用
```go
ms := memorydb.New()

setValue := 1
expire := 10

index := ms.Set(setValue, expire)

getValue := ms.Get(index)
```

### 3.目前支持的方法

| 方法       | 描述          |
|:---------|:------------|
| Set      | 新增数据        |
| SetBatch | 批量新增数据      |
| Get      | 通过id获取数据    |
| GetAll   | 获取所有数据      |
| Remove   | 通过id删除某一条数据 |
| Flush    | 清空所有数据      |
| Update   | 通过id更新某一条数据 |
| Size     | 当前库中存有的数据量  |
| Empty    | 当前数据库是否为空   |


### 4.项目中的一点使用须知

- 目前在Set、Remove、Update、生成数据对应id、筛选过期value的时候，默认都没有使用锁，可能会出现线程安全的问题，如果对线程安全有要求，可以使用`memorydb.NewWithLock()`创建数据库对象。


- 对于数据库加锁与否的简单性能测试，可能有十足的偶然性：
```go
ms := memorydb.NewWithLock()
//2407669  lock
//2118867  unlock
time.Sleep(time.Second)
start := time.Now()
for i := 0; i < 10000000; i++ {
    ms.Set(i, 100)
}
end := time.Now()
fmt.Println(end.UnixMicro() - start.UnixMicro())


```

- 目前使用的是浅拷贝，可能会出现数据安全问题，插入数据的时候，需要注意
- 