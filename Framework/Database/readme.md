

    import (
        "github.com/goal-web/contracts"
        "github.com/goal-web/database"
    )
    func init() {
        configs["database"] = func(env contracts.Env) interface{} {
            return database.Config{
                Default: env.StringOption("db.connection", "mysql"),
                Connections: map[string]contracts.Fields{
                    "sqlite": {
                        "driver":   "sqlite",
                        "database": env.GetString("sqlite.database"),
                    },
                    "mysql": {
                        "driver":          "mysql",
                        "host":            env.GetString("db.host"),
                        "port":            env.GetString("db.port"),
                        "database":        env.GetString("db.database"),
                        "username":        env.GetString("db.username"),
                        "password":        env.GetString("db.password"),
                        "unix_socket":     env.GetString("db.unix_socket"),
                        "charset":         env.StringOption("db.charset", "utf8mb4"),
                        "collation":       env.StringOption("db.collation", "utf8mb4_unicode_ci"),
                        "prefix":          env.GetString("db.prefix"),
                        "strict":          env.GetBool("db.struct"),
                        "max_connections": env.GetInt("db.max_connections"),
                        "max_idles":       env.GetInt("db.max_idles"),
                    },
                    "pgsql": {
                        "driver":          "postgres",
                        "host":            env.GetString("db.pgsql.host"),
                        "port":            env.GetString("db.pgsql.port"),
                        "database":        env.GetString("db.pgsql.database"),
                        "username":        env.GetString("db.pgsql.username"),
                        "password":        env.GetString("db.pgsql.password"),
                        "charset":         env.StringOption("db.pgsql.charset", "utf8mb4"),
                        "prefix":          env.GetString("db.pgsql.prefix"),
                        "schema":          env.StringOption("db.pgsql.schema", "public"),
                        "sslmode":         env.StringOption("db.pgsql.sslmode", "disable"),
                        "max_connections": env.GetInt("db.pgsql.max_connections"),
                        "max_idles":       env.GetInt("db.pgsql.max_idles"),
                    },
                },
            }
        }
    }
```

`.env` 的数据库相关配置

```bash
# 默认连接
db.connection=sqlite

sqlite.database=/Users/qbhy/project/go/goal-web/goal/example/database/db.sqlite

db.host=localhost
db.port=3306
db.database=goal
db.username=root
db.password=password

db.pgsql.host=localhost
db.pgsql.port=55433
db.pgsql.database=postgres
db.pgsql.username=postgres
db.pgsql.password=123456
```

### 定义模型 - define a model
`app/models/user.go` 文件

```go
package models

import (
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/class"
)

// UserClass 这个类变量，以后大有用处
var UserClass = class.Make(new(User))

// UserModel 返回 table 实例，继承自查询构造器并且实现了所有 future
func UserModel() *table.Table {
	return table.Model(UserClass, "users")
}

// User 模型结构体
type User struct {
	Id       int64  `json:"id"`
	NickName string `json:"name"`
}
```

### 用法 - method of use
```go
package tests

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/example/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getQuery(name string) contracts.QueryBuilder {
	// 测试用例环境下的简易 goal 应用启动
	app := initApp("/Users/qbhy/project/go/goal-web/goal/tests")

	//return  table.Query("users") 返回 table 实例，使用默认连接
	//tx, _ := app.Get("db").(contracts.DBConnection).Begin()
	//return table.WithTX("users", tx) // 事物环境下执行

	//return table.WithConnection(name, "sqlite") // 返回指定连接的 table 实例，使用连接名
	return table.WithConnection(name, app.Get("db").(contracts.DBConnection)) // 也可以指定连接实例
}

// TestTableQuery 测试不带模型的 table 查询，类似 laravel 的 DB::table()
func TestTableQuery(t *testing.T) {

	getQuery("users").Delete()

	// 不设置模型的情况下，返回 contracts.Fields
	user := getQuery("users").Create(contracts.Fields{
		"name": "qbhy",
	}).(contracts.Fields)

	fmt.Println(user)
	userId := user["id"].(int64)
	// 判断插入是否成功
	assert.True(t, userId > 0)

	// 获取数据总量
	assert.True(t, getQuery("users").Count() == 1)

	// 修改数据
	num := getQuery("users").Where("name", "qbhy").Update(contracts.Fields{
		"name": "goal",
	})
	assert.True(t, num == 1)
	// 判断修改后的数据
	user = getQuery("users").Where("name", "goal").First().(contracts.Fields)

	err := getQuery("users").Chunk(10, func(collection contracts.Collection, page int) error {
		assert.True(t, collection.Len() == 1)
		fmt.Println(collection.ToJson())
		return nil
	})

	assert.Nil(t, err)

	assert.True(t, user["id"] == userId)
	assert.True(t, user["name"] == "goal")
	assert.True(t, getQuery("users").Find(userId).(contracts.Fields)["id"] == userId)
	assert.True(t, getQuery("users").Where("id", userId).Delete() == 1)
	assert.Nil(t, getQuery("users").Find(userId))
}

func TestModel(t *testing.T) {
	initApp("/Users/qbhy/project/go/goal-web/goal/tests")

	fmt.Println("用table查询：",
		getQuery("users").Get().Map(func(user contracts.Fields) {
			fmt.Println("用table查询", user)
		}).ToJson()) // query 返回 Collection<contracts.Fields>

	user := models.UserModel().Create(contracts.Fields{
		"name": "qbhy",
	}).(models.User)

	fmt.Println("创建后返回模型", user)

	fmt.Println("用table查询：",
		getQuery("users").Get().Map(func(user contracts.Fields) {
			fmt.Println("用table查询", user)
		}).ToJson()) // query 返回 Collection<contracts.Fields>

		// 用模型查询
	fmt.Println(models.UserModel(). // model 返回 Collection<models.User>
					Get().
					Map(func(user models.User) {
			fmt.Println("id:", user.Id)
		}).ToJson())

	fmt.Println(models.UserModel().Where("id", ">", 0).Delete())
}
