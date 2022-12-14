# Gorm 库介绍

> 参考：
>
> * [GitHub 项目，jinzhu/gorm](https://github.com/jinzhu/gorm)
>   * 新版已移动到 [go-gorm/gorm](https://github.com/go-gorm/gorm)
> * [官网](https://gorm.io/)
> * 文档：https://gorm.io/zh_CN/docs/index.html

# 约定,一些关于 gorm 中的特殊说明

## gorm.Model 结构体

gorm.Model 是一个包含了ID, CreatedAt, UpdatedAt, DeletedAt四个字段的GoLang结构体。

你可以将它嵌入到你自己的 Model 中，也可以完全使用自己的 Model。

```go
// gorm.Model 定义
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}
```

```go
// Inject fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` into model `User`
// 将 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`字段注入到`User`模型中
type User struct {
  gorm.Model
  Name string
}
```

```go
// 不使用gorm.Model定义模型
type User struct {
  ID   int
  Name string
}
ID 作为主键
GORM 默认会使用名为ID的字段作为表的主键。

type User struct {
  ID   string // 名为`ID`的字段会默认作为表的主键
  Name string
}

// 使用`AnimalID`作为主键
type Animal struct {
  AnimalID int64 `gorm:"primary_key"`
  Name     string
  Age      int64
}
```

## 表名（Table Name）

表名默认就是结构体名称的复数，例如：

```go
type User struct {} // 默认表名是 `users`

// 将 User 的表名设置为 `profiles`
func (User) TableName() string {
  return "profiles"
}

func (u User) TableName() string {
  if u.Role == "admin" {
    return "admin_users"
  } else {
    return "users"
  }
}


// 禁用默认表名的复数形式，如果置为 true，则 `User` 的默认表名是 `user`
db.SingularTable(true)
```

## 指定表名称

```go
// 使用User结构体创建名为`deleted_users`的表
db.Table("deleted_users").CreateTable(&User{})

var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
//// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "jinzhu").Delete()
//// DELETE FROM deleted_users WHERE name = 'jinzhu';
```

## 更改默认表名称（table name）

可以通过定义DefaultTableNameHandler来设置默认表名的命名规则

```go
gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
  return "prefix_" + defaultTableName;
}
```

## 下划线分割命名（Snake Case）的列名

列名根据结构体中的字段自动生成，若字段中的名字有多个大写，则通过下划线进行分割

```go
type User struct {
  ID        uint      // column name is `id`
  Name      string    // column name is `name`
  Birthday  time.Time // column name is `birthday`
  CreatedAt time.Time // column name is `created_at`
}

// 手动设置列名
type Animal struct {
  AnimalId    int64     `gorm:"column:beast_id"`         // set column name to `beast_id`
  Birthday    time.Time `gorm:"column:day_of_the_beast"` // set column name to `day_of_the_beast`
  Age         int64     `gorm:"column:age_of_the_beast"` // set column name to `age_of_the_beast`
}
```

## 时间点（Timestamp）跟踪

CreatedAt

如果模型有 CreatedAt字段，该字段的值将会是初次创建记录的时间。

```go
db.Create(&user) // `CreatedAt`将会是当前时间

// 可以使用`Update`方法来改变`CreateAt`的值
db.Model(&user).Update("CreatedAt", time.Now())
```

UpdatedAt

如果模型有UpdatedAt字段，该字段的值将会是每次更新记录的时间。

```go
db.Save(&user) // `UpdatedAt`将会是当前时间

db.Model(&user).Update("name", "jinzhu") // `UpdatedAt`将会是当前时间
```

DeletedAt

如果模型有DeletedAt字段，调用Delete删除该记录时，将会设置DeletedAt字段为当前时间，而不是直接将记录从数据库中删除。 了解什么是 [软删除](https://gorm.io/zh_CN/docs/delete.html#Soft-Delete)
