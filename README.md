# skiplist

## 📖 简介

`skiplist`一个基于`golang`实现的跳跃链表。其添加,查询,删除时间复杂度为`O(logn)`。并且是并发安全的。

## 🚀 功能

- 跳跃链表的实现，其添加,查询,删除时间复杂度为`O(logn)`。
- 并发安全，可在并发环境下使用。
- 支持设置最大层高度和基础概率，可根据数据规模进行调整

## 🧰 安装

```
go get -u github.com/sirodeneko/skiplist
```

## 🛠 使用

```
list = skiplist.New()
list.Set(1,"1")
list.Set(2,"2")
list.Set(3,"3")

list.Get(1)
list.Get(2)
list.Get(3)

list.Remove(1)
```



## 🔍 测试

```
$ go test -bench . -benchmem
Structure sizes: SkipList is 136, Element is 48 bytes
goos: windows
goarch: amd64
pkg: skiplist
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkIncSet-8   2857716      412.9 ns/op   6921730.49 MB/s     61 B/op   3 allocs/op
BenchmarkIncGet-8   5116286      207.2 ns/op   24690854.69 MB/s    0 B/op    0 allocs/op
BenchmarkDecSet-8   4539894      239.9 ns/op   18924755.93 MB/s    61 B/op   3 allocs/op
BenchmarkDecGet-8   5186264      213.8 ns/op   24251998.98 MB/s    0 B/op    0 allocs/op
PASS
ok      skiplist        21.348s

```

