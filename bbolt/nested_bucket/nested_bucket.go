package main

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

// 使用字符串 / 作为顶层 bucket 的名称
var root_bucket = []byte("/")
var bucket_1 = []byte("bucket_1")

// var bucket_2 = []byte("bucket_2")
var key = []byte("hello_world")
var value = []byte("DesistDaydream")
var key_1 = []byte("hello_world_1")
var value_1 = []byte("DesistDaydream_1")

func rw_transactions(tx *bbolt.Tx) {
	// 如果 bucket 不存在则，创建一个 bucket
	bucket, _ := tx.CreateBucketIfNotExists(root_bucket)

	// 将 key-value 写入到 bucket 中
	bucket.Put(key, value)

	// 在顶层 bucket 下再创建一个 bucket
	bucket, _ = bucket.CreateBucketIfNotExists(bucket_1)

	// 将 key-value 写入到 bucket_1 中
	bucket.Put(key_1, value_1)
}

func iteratingBucket1(tx *bbolt.Tx) {
	b := tx.Bucket([]byte(root_bucket))
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("key=%s,value=%s\n", k, v)
	}
}

func iteratingBucket2(b *bbolt.Bucket, space string) {
	space = space + "  "
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			fmt.Printf("%sbucket=%s\n", space, k)
			iteratingBucket2(b.Bucket([]byte(k)), space)
		} else {
			fmt.Printf("%skey=%s, value=%s\n", space, k, v)
		}
		return nil
	})
}

func main() {
	// 打开数据库，若不存在则自动创建
	db, err := bbolt.Open("bbolt/nested_bucket/nested_bucket.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 启动读写事务
	db.Update(func(tx *bbolt.Tx) error {
		rw_transactions(tx)
		return nil
	})

	// 启动只读事务
	db.View(func(tx *bbolt.Tx) error {
		iteratingBucket1(tx)

		b := tx.Bucket(root_bucket)
		space := ""
		iteratingBucket2(b, space)
		return nil
	})
}
