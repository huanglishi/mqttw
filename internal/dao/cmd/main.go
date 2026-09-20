// =========================
// Gen框架生成代码
// =========================
package main

import (
	"fmt"
	"gofly/internal/dao/conf"
	"gofly/internal/dao/drivers"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gofly/internal/dao/drivers/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
)

// snakeToUpperCamel 下划线转大驼峰，保证首字母大写
func snakeToUpperCamel(s string) string {
	var sb strings.Builder
	needUpper := true
	for _, r := range s {
		if r == '_' {
			needUpper = true
			continue
		}
		if needUpper {
			sb.WriteRune(unicode.ToUpper(r))
			needUpper = false
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",                                // gconv.String(dbConf_arr["prefix"]), // table name prefix, table for `User` would be `gf_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true,                              // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
	})
	if err != nil {
		panic(fmt.Errorf("connect db fail: %v", err.Error()))
	}
	return db
}

func main() {
	// 数据库配置-兼容运行命令目录位置
	pathUrl, _ := os.Getwd()
	var dbConf_arr map[string]any
	var outPath string
	if strings.Contains(pathUrl, "dao") {
		outPath = "../query"
		dbConf_arr = map[string]any{"type": conf.Type, "dbFile": filepath.Join(pathUrl, "../../../resource", "db", "data.db")}
	} else {
		outPath = "dao/query"
		dbConf_arr = map[string]any{"type": conf.Type, "dbFile": filepath.Join(pathUrl, "resource", "db", "data.db")}
	}
	// 指定生成代码的具体相对目录(相对当前文件)，默认为：./query
	// 默认生成需要使用WithContext之后才可以查询的代码，但可以通过设置gen.WithoutContext禁用该模式
	g := gen.NewGenerator(gen.Config{
		// 默认会在 OutPath 目录生成CRUD代码，并且同目录下生成 model 包
		// 所以OutPath最终package不能设置为model，在有数据库表同步的情况下会产生冲突
		// 若一定要使用可以通过ModelPkgPath单独指定model package的名称
		OutPath: outPath,
		// OutPath: "../../dal/query",
		/* ModelPkgPath: "dal/model"*/

		// gen.WithoutContext：禁用WithContext模式
		// gen.WithDefaultQuery：生成一个全局Query对象Q
		// gen.WithQueryInterface：生成Query接口
		Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// Model名称去除prefix配置前缀，保留其他前缀，裁剪tablePrefix后转大驼峰，解决小写开头报错
	g.WithModelNameStrategy(func(tableName string) string {
		short := strings.TrimPrefix(tableName, conf.Prefix)
		return snakeToUpperCamel(short)
	})

	// 文件名命名去掉表前缀配置
	g.WithFileNameStrategy(func(tableName string) string {
		return strings.TrimPrefix(tableName, conf.Prefix)
	})

	// 自定时间数据类型映射
	g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		// 将datetime、timestamp类型映射到自定义的LocalTime
		// "datetime": func(columnType gorm.ColumnType) string {
		// 	return "gtime.GormTime"
		// },
		// "timestamp": func(columnType gorm.ColumnType) string {
		// 	return "gtime.GormTime"
		// },
		// // 可添加其他类型映射（如date）
		// "date": func(columnType gorm.ColumnType) string {
		// 	return "gtime.GormTime"
		// },
		"tinyint": func(columnType gorm.ColumnType) string {
			return "int8"
		},
	})

	//配置软删除,默认 gorm.DeletedAt
	g.WithOpts(
		gen.FieldType("deleted_at", "gorm.DeletedAt"),
	)

	// 通常复用项目中已有的SQL连接配置db(*gorm.DB)
	// 非必需，但如果需要复用连接时的gorm.Config或需要连接数据库同步表信息则必须设置
	MySQLDSN := drivers.DbDSN(dbConf_arr) //从app\dal\model获取dsn链接
	g.UseDB(connectDB(MySQLDSN))

	// 从连接的数据库为所有表生成Model结构体和CRUD代码
	// 也可以手动指定需要生成代码的数据表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行并生成代码
	g.Execute()
}
