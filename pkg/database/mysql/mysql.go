package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// databaseInfo 数据库连接信息
type databaseInfo struct {
	UserName string
	Password string
	Protocol string
	Server   string
	Port     int64
	Database string
}

// User 表结构
type User struct {
	ID         int    `json:"id" form:"id"`
	Username   string `json:"username" form:"username"`
	Password   string `json:"password" form:"password"`
	Status     int    `json:"status" form:"status"` // 0 正常状态， 1删除
	Createtime int64  `json:"createtime" form:"createtime"`
}

// ConnDB 连接数据库
func (i *databaseInfo) ConnDB() (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", i.UserName, i.Password, i.Protocol, i.Server, i.Port, i.Database)
	DB, err := sql.Open("mysql", conn)
	fmt.Println(err)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		os.Exit(3)
	}
	return DB, err
}

// CreateTable 创建表
func CreateTable(DB *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS users(
	id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	username VARCHAR(64),
	password VARCHAR(64),
	status INT(4),
	createtime INT(10)
	); `

	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}

// InsertData 插入数据
func InsertData(DB *sql.DB) {
	result, err := DB.Exec("insert INTO users(username,password) values(?,?)", "test", "123456")
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)

	rowsaffected, err := result.RowsAffected() //通过RowsAffected获取受影响的行数
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

// QueryOne 查询单行
func QueryOne(DB *sql.DB) {
	user := new(User) //用new()函数初始化一个结构体对象
	row := DB.QueryRow("select id,username,password from users where id=?", 1)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Println("Single row data:", *user)
}

// QueryMulti 查询多行
func QueryMulti(DB *sql.DB) {
	user := new(User)
	rows, err := DB.Query("select id,username,password from users where id = ?", 2)

	defer func() {
		if rows != nil {
			rows.Close() //关闭掉未scan的sql连接
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Password) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)
			return
		}
		fmt.Println("scan successd:", *user)
	}
}

// UpdateData 更新数据
func UpdateData(DB *sql.DB) {
	result, err := DB.Exec("UPDATE users set password=? where id=?", "111111", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return
	}
	fmt.Println("update data successd:", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

// DeleteData 删除数据
func DeleteData(DB *sql.DB) {
	result, err := DB.Exec("delete from users where id=?", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)
		return
	}
	fmt.Println("delete data successd:", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

func main() {
	// 设置连接数据库的信息
	i := new(databaseInfo)
	i.UserName = "root"
	i.Password = "mysql"
	i.Protocol = "tcp"
	i.Server = "127.0.0.1"
	i.Port = 3306
	i.Database = "test"
	// 连接数据库
	DB, _ := i.ConnDB()

	//最大连接周期，超时的连接就close
	DB.SetConnMaxLifetime(100 * time.Second)
	//设置最大连接数
	DB.SetMaxOpenConns(100)

	// 创建表
	CreateTable(DB)
	// 插入数据
	InsertData(DB)
	// 查询一行
	QueryOne(DB)
	// 查询多行
	QueryMulti(DB)
	// 更新数据
	UpdateData(DB)
	// 删除数据
	DeleteData(DB)
}
