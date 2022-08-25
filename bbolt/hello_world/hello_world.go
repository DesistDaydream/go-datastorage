package main

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

var root_bucket_1 = []byte("hello_world_bucket_1")
var root_bucket_2 = []byte("hello_world_bucket_2")
var key = []byte("hello_world")
var value = []byte("DesistDaydream")

func crud(bucket *bbolt.Bucket) {
	// 增。将 键值对 写入到 bucket 中
	bucket.Put(key, value)

	// 查
	v := bucket.Get(key)
	fmt.Printf("查：\nkey=%s,value=%s\n", key, v)

	// 改
	bucket.Put(key, []byte("DesistDaydream_2"))
	fmt.Printf("改：\nkey=%s,value=%s\n", key, bucket.Get(key))

	// 删
	bucket.Delete(key)
	fmt.Printf("删：\nkey=%s,value=%s\n", key, bucket.Get(key))

	// 再写入键，供后续读事务使用
	bucket.Put(key, value)
}

// 读写事务
func rw_transactions(tx *bbolt.Tx) {
	var buckets []*bbolt.Bucket
	// 创建多个顶层 Bucket
	bucket_1, _ := tx.CreateBucketIfNotExists(root_bucket_1)
	bucket_2, _ := tx.CreateBucketIfNotExists(root_bucket_2)

	buckets = append(buckets, bucket_1, bucket_2)

	for _, bucket := range buckets {
		// 增删改查。对 Bucket 执行增删改查操作。
		// 注意：
		// 对 Bucket 执行增删改查操作,必须在一个 Transactions(事务) 中。
		// 如果把代码从 db.Update() 中移出，将会 Panic: assertion failed: tx closed
		crud(bucket)
	}
}

// 只读事务
func r_transactions(tx *bbolt.Tx) {
	fmt.Printf("迭代根下的所有 Bucket:\n")
	// 迭代所有根下的 Bucket。将桶名称赋值给 name 变量；并实例化桶为 bucket 变量。
	tx.ForEach(func(bucket_name []byte, bucket *bbolt.Bucket) error {
		fmt.Printf("root_bucket=%s\n", bucket_name)

		// 使用 ForEach() 迭代 Bucket
		bucket.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s,value=%s\n", k, v)
			return nil
		})

		// 只读事务中可以创建一个 Cursor
		c := bucket.Cursor()
		// 使用 Cursor 迭代 Bucket
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s,value=%s\n", k, v)
		}

		return nil
	})

}

func main() {
	// 打开数据库，若不存在则自动创建
	db, err := bbolt.Open("pkg/bbolt/hello_world/hello_world.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// bbolt.Tx 结构体表示数据库上的一个 只读 或者 读写 的事务。
	// 重要提示: 完成事务后，必须提交或回滚事务。除非没有更多的事务正在使用页面，否则作者无法收回页面。长时间运行的读取事务会导致数据库快速增长。

	// 启动一个 读写 Transactions(事务)。可用于创建和移除 Bucket、创建和移除 Key
	db.Update(func(tx *bbolt.Tx) error {
		rw_transactions(tx)
		return nil
	})

	// 启动一个 只读 Transactions(事务)。可用于检索 Key 和创建 Cursors
	db.View(func(tx *bbolt.Tx) error {
		r_transactions(tx)
		return nil
	})
}
