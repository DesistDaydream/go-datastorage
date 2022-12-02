package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	// gorm.Model 是一个包含了ID, CreatedAt, UpdatedAt, DeletedAt四个字段的结构体。
	// gorm.Model
	Code  string
	Price uint
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}

	// AutoMigrate 用来刷新数据表，不存在则创建，表名默认为结构体名称的复数，e.g.这里会创建一个名为 products 的表，假如 Product 为 ProductTest，则会创建出一个名为 product_test 的表
	// 结构体中的每个字段都是该表的列，字段名称即是表中列的名称，如果字段名中有多个大写字母，则列名使用下划线分隔，e.g.CreatedAt 字段的列名为 cretaed_at
	// 当结构体中增加字段时，会自动在表中增加列；但是删除结构体中的属性时，并不会删除列
	db.AutoMigrate(&Product{})

	return db
}

// 插入数据
func Insert(db *gorm.DB) {
	db.Create(&Product{Code: "D42", Price: 100})
	db.Create(&Product{Code: "D42", Price: 200})
	db.Create(&Product{Code: "D43", Price: 200})
}

// 查询数据
func Query(db *gorm.DB) {
	// 声明一个 Product 数组，用来存放查询结果
	var products []Product
	// 查询数据，查找 code 列的值为 D42 的所有记录，并将查询结果存放到 products 中
	db.Find(&products, "code = ?", "D42")
	for _, product := range products {
		log.Println(product)
	}
}

// 更新数据
func Update(db *gorm.DB) {
	// 根据条件更新数据
	// UPDATE products SET price=300 WHERE code="D42";
	db.Model(&Product{}).Where("code = ?", "D42").Update("price", 300)
}

// 删除数据
func Delete(db *gorm.DB) {
	// 根据 id 删除数据
	// DELETE FROM products WHERE id=10;
	db.Delete(&Product{}, 1)
}

func main() {
	db := InitDB()
	Insert(db)
	Query(db)
	Update(db)
	// Delete(db)
}
