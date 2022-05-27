# Go 基础知识与框架体系 part 9: hashmap


这篇文章总结了 Go 的知识体系: hashmap，包括其中的底层实现等等。

<!--more-->

## hashmap

{{< admonition type=question open=true >}}
Hashmap 的内部结构是如何实现的呢？
{{< /admonition >}}

### 1. 内部结构

Go的 map 是 unordered map，即无法对 key 值排序遍历。跟传统的 hashmap 的实现方法一样，它通过一个 buckets 数组实现，所有元素被 hash 到数组的 bucket 中，buckets 就是指向了这个内存连续分配的数组。`B`字段说明 hash 表大小是2的指数，即$2^B$。每次扩容会增加到上次大小的两倍，即 $2^{B+1}$。当 bucket 填满后，将通过 `overflow` 指针来 `mallocgc` 一个bucket出来形成链表，也就是为哈希表解决冲突问题。

```go
// A header for a Go map.
type hmap struct {
	count int // len()返回的map的大小 即有多少kv对
	flags uint8
	B     uint8  // 表示hash table总共有2^B个buckets 
	hash0 uint32 // hash seed
	buckets    unsafe.Pointer // 按照low hash值可查找的连续分配的数组，初始时为16个Buckets.
	oldbuckets unsafe.Pointer 
	nevacuate  uintptr      
	overflow *[2]*[]*bmap //溢出链 当初始buckets都满了之后会使用overflow
}

// A bucket for a Go map.
type bmap struct {
    tophash [bucketCnt]uint8
// Followed by bucketCnt keys and then bucketCnt values.
// NOTE: packing all the keys together and then all the values together makes the
// code a bit more complicated than alternating key/value/key/value/... but it allows
// us to eliminate padding which would be needed for, e.g., map[int64]int8.
// Followed by an overflow pointer.
}
```
bucket map 的数据结构中，`tophash` 是个大小为 8(bucketCnt) 的数组，存储了8个 key 的 hash 值的高八位值。

{{< admonition type=tip open=true >}}
在对 key/value 对增删查的时候，先比较 key 的 hash 值高八位是否相等，然后再比较具体的key值。
{{< /admonition  >}}

### 2. map 初始化

B的初始大小是0，若指定了map的大小hint且hint大于8，那么buckets会在make时就通过newarray分配好，否则buckets会在第一次put的时候分配。随着hashmap中key/value对的增多，buckets需要重新分配，每一次都要**重新hash并进行元素拷贝**，所以最好在初始化时就给map指定一个合适的大小。

`makemap` 有`h`和`bucket`这两个参数，是留给编译器的。如果编译器决定`hmap`结构体和第一个`bucket`可以在栈上创建，这两个入参可能不是`nil`的。

```go
// makemap implemments a Go map creation make(map[k]v, hint)
func makemap(t *maptype, hint int64, h *hmap, bucket unsafe.Pointer) *hmap{
  B := uint8(0)
  for ; hint > bucketCnt && float32(hint) > loadFactor*float32(uintptr(1)<<B); B++ {
  }
  // 确定初始B的初始值 这里hint是指kv对的数目 而每个buckets中可以保存8个kv对
  // 因此上式是要找到满足不等式 hint > loadFactor*(2^B) 最小的B
  if B != 0 {
    buckets = newarray(t.bucket, uintptr(1)<<B)
  }
  h = (*hmap)(newobject(t.hmap))
  return h
}
```
### 3. map 存值

存储的步骤和第一部分的分析一致。首先用`key`的`hash`值低8位找到`bucket`，然后在`bucket`内部比对`tophash`和高8位与其对应的`key`值与入参`key`是否相等，若找到则更新这个值。若`key`不存在，则`key`优先存入在查找的过程中遇到的空的`tophash`数组位置。若当前的`bucket`已满则需要另外分配空间给这个`key`，新分配的`bucket`将挂在`overflow`链表后。

```go
func mapassign1(t *maptype, h *hmap, key unsafe.Pointer, val unsafe.Pointer) {
  hash := alg.hash(key, uintptr(h.hash0))
  if h.buckets == nil {
    h.buckets = newarray(t.bucket, 1)
  }
again:
  //根据低8位hash值找到对应的buckets
  bucket := hash & (uintptr(1)<> (sys.PtrSize*8 - 8))
  for {
    //遍历每一个bucket 对比所有tophash是否与top相等
    //若找到空tophash位置则标记为可插入位置
    for i := uintptr(0); i < bucketCnt; i++ {
      if b.tophash[i] != top {
        if b.tophash[i] == empty && inserti == nil {
          inserti = &b.tophash[i]
        } 
        continue
      }
      //当前tophash对应的key位置可以根据bucket的偏移量找到
      k2 := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
      if !alg.equal(key, k2) {
        continue
      }
      //找到符合tophash对应的key位置
      typedmemmove(t.elem, v2, val)
      goto done
    }
    //若overflow为空则break
    ovf := b.overflow(t)
  }
  // did not find mapping for key.  Allocate new cell & add entry.
  if float32(h.count) >= loadFactor*float32((uintptr(1)<= bucketCnt {
    hashGrow(t, h)
    goto again // Growing the table invalidates everything, so try again
  }
  // all current buckets are full, allocate a new one.
  if inserti == nil {
    newb := (*bmap)(newobject(t.bucket))
    h.setoverflow(t, b, newb)
    inserti = &newb.tophash[0]
  }
  // store new key/value at insert position
  kmem := newobject(t.key)
  vmem := newobject(t.elem)
  typedmemmove(t.key, insertk, key) 
  typedmemmove(t.elem, insertv, val)
  *inserti = top
  h.count++
}
```

### 4. hash grow 扩容和迁移

在往map中存值时若所有的bucket已满，需要在堆中new新的空间时需要计算是否需要扩容。扩容的时机是`count > loadFactor(2^B)`。这里的`loadfactor`选择为6.5。

{{< admonition type=question >}}
为什么选取`loadfactor`为6.5呢？
{{< /admonition >}}

这是因为`loadfactor`和`overflow`溢出率、`bytes/entry`、`hitprobe`、`missprobe`相关。

- `overflow`溢出率是指平均一个bucket有多少个kv的时候会溢出。
- `bytes/entry`是指平均存一个kv需要额外存储多少字节的数据。
- `hitprobe`是指找到一个存在的key平均需要找多少次。
- `missprobe`是指找到一个不存在的key平均需要找多少次。选取6.5是为了平衡这组数据。

在没有溢出时hashmap总共可以存储`8(2^B)`个KV对，当hashmap已经存储到`6.5(2^B)`个KV对时表示hashmap已经趋于溢出，即很有可能在存值时用到`overflow`链表，这样会增加`hitprobe`和`missprobe`。

| **loadfactor** | **%overflow** | **bytes/entry** | **hitprobe** | **missprobe** |
|----------------|---------------|-----------------|--------------|---------------|
| 4.00           | 2.13          | 20.77           | 3.00         | 4.00          |
| 4.50           | 4.05          | 17.30           | 3.25         | 4.50          |
| 5.00           | 6.85          | 14.77           | 3.50         | 5.00          |
| 5.50           | 10.55         | 12.94           | 3.75         | 5.50          |
| 6.00           | 15.27         | 11.67           | 4.00         | 6.00          |
| 6.50           | 20.90         | 10.79           | 4.25         | 6.50          |
| 7.00           | 27.14         | 10.15           | 4.50         | 7.00          |
| 7.50           | 34.03         | 9.73            | 4.50         | 7.50          |
| 8.00           | 41.10         | 9.40            | 5.00         | 8.00          |

但这个迁移并没有在扩容之后一次性完成，而是逐步完成的，每一次insert或remove时迁移1到2个pair，即**增量扩容**。

**增量扩容的原因**:主要是缩短map容器的响应时间。若hashmap很大扩容时很容易导致系统停顿无响应。

**增量扩容本质上**就是将总的扩容时间分摊到了每一次hash操作上。由于这个工作是逐渐完成的，导致数据一部分在old table中一部分在new table中。old的bucket不会删除，只是加上一个已删除的标记。只有当所有的bucket都从old table里迁移后才会将其释放掉。(摘录于[Golang中文社区](https://studygolang.com/articles/11979))

