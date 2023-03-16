# 如何通过Gorm连接数据库

```
func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("连接数据库成功")
}
```

# 如何通过Gorm生成表解构

```
func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("连接数据库成功")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalln(err)
	}
}
```

# Gorm中对零值的处理，update和updates都有哪些坑？

在 Golang 定义一个变量且不进行赋值时，初始值就是该数据类型的零值，比如 0、""、nil 或 false 等。
在使用 gorm 更新数据时，只要将变更的字段值进行修改，然后执行 `db.Updates()` 操作即可。

查看 gorm 文档得知，`Updates` 方法支持 `struct` 和 `map[string]interface{}` 参数。当使用 `struct` 更新时，默认情况下，GORM 只会更新非零值的字段，所以可以使用 `map` 更新字段，或者使用 `Select` 指定要更新的字段。

```
// 使用 map 更新
db.Model(&user).Updates(map[string]interface{}{"age": 0})

// 使用 Select 指定字段更新
db.Model(&user).Select("Age").Updates(User{Age: 0})
```

# 添加数据有几种方式？

```
// 单条插入
user := User{Name: "Channer", Age: 30, CreatedAt: time.Now()}
db.Create(&user)

// 多条一次性插入
users := []User{{Name: "Channer1"}, {Name: "Channer2"}}
db.Create(&users)

// 多条分批插入
users := []User{{Name: "Channer1"}, ..., {Name: "Channer1000"}}
// 每次插入 100 条数据
db.CreateInBatches(&users, 100)
// upsert
db.Where(User{Name: "nobody"}).FirstOrCreate(&user)
```

# 查询数据有几种方式？

```
// 获取单条记录
db.First(&user) // 主键升序第一条
db.Take(&user)  // 默认第一条
db.Last(&user)  // 主键降序最后一条

// 获取全部匹配记录
var users []User

db.Where(User{Name: "channer"}).Scan(&users)
db.Where(User{Name: "channer"}).Order("id desc, name").Limit(10).Offset(5).Find(&users)
db.Where("name = ? AND age >= ?", "channer", "18").Find(&users)

// 子查询
db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
```

# Gorm如何实现删除数据的？

```
// 删除指定数据
db.Delete(&User{}, 10) // user.id == 10
db.Where("name = ?", "channer").Delete(&User{}) // user.name == "channer"

// 软删除，标记删除
type User struct {
  ID      int
  Deleted gorm.DeletedAt
  Name    string
}

// 软删除
db.Delete(&user)
db.Where("name = ?", "channer").Delete(&User{})

// 永久删除
db.Unscoped().Delete(&user)
```

# 如何进行一对一，一对多，多对多查询？

```
type Employer struct {
    gorm.Model
    Name string `gorm:"unique;not_null"`
    CompanyId int
    Company Company
    CreditCards []CreditCard
    
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
  gorm.Model
  Name string
  Employers []Employers `gorm:"many2many:user_languages;"`
}

type Company struct {
    gorm.Model
    Name string
}
type CreditCard struct {
    gorm.Model
    Number string
    EmployerId uint
}

// 一对多
var employer Employer
db.Model(&Employer{}).Preload("Company").First(&employer)
fmt.Println(employer.Company.Name)

// 一对多
db.Model(&Employer{}).Preload("CreditCards").First(&employer)
fmt.Println(employer.CreditCards)

// 多对多
db.Model(&Employer{}).Preload("Languages").First(&employer)
fmt.Println(employer.Languages)
db.Model(&Language{}).Preload("Employers").FIrst(&language)
fmt.Println(language.Employers)
```