package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func iteratingAll(bucket *bolt.Bucket, space string) {
	space = space + "  "
	bucket.ForEach(func(k, v []byte) error {
		if v == nil {
			fmt.Printf("%s[%s]: \n", space, k)
			// 嵌套迭代
			iteratingAll(bucket.Bucket([]byte(k)), space)
		} else {
			fmt.Printf("%s%s=%s\n", space, k, v)
		}
		return nil
	})
}

func iteratingBucket(bucket *bolt.Bucket, space string) {
	space = space + "  "
	bucket.ForEach(func(k, v []byte) error {
		if v == nil {
			fmt.Printf("%s[%s]\n", space, k)
			// 嵌套迭代
			iteratingBucket(bucket.Bucket([]byte(k)), space)
		}
		return nil
	})
}

func r_transactions(tx *bolt.Tx) {
	// 迭代根中的所有 Bucket
	tx.ForEach(func(bucket_name []byte, bucket *bolt.Bucket) error {
		fmt.Printf("root_bucket=%v\n", string(bucket_name))

		space := ""
		// 迭代所有 Bucket
		iteratingBucket(bucket, space)

		// 迭代所有 Bucket 及其 K/V。
		// 可以指定名称空间，以跌点特定名称空间写下的 Bucket 及其 K/V
		// bucket = tx.Bucket(bucket_name).Bucket([]byte("default"))
		iteratingAll(bucket, space)
		return nil
	})
}

func main() {
	db, err := bolt.Open("metadata.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		r_transactions(tx)
		return nil
	})
}
