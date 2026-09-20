package dao

import (
	"fmt"
	"gofly/internal/dao/conf"
	"gofly/internal/dao/drivers"
	"gofly/internal/dao/query"
	"gofly/internal/utils/tools/instance"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type (
	// Condition query condition
	// field.Expr and subquery are expect value
	Condition = gen.Condition
)

// Create a writer for logging entries/创建写日志的writer
type logWriter struct {
	stdout bool
}

// Default database configuration/默认数据库配置
const (
	DefaultGroupName = "default"            // Default group name.
	keyDatabase      = "component.database" // Database of Component Names
)

// Return the query object. The parameter 'name' is the name of the database configuration, with the default value being 'default'.
// 返回query对象，参数name是数据库配置名称，默认default。
func Query(name ...string) *query.Query {
	return query.Use(DB(name...))
}

// Instantiate the data object/实例化数据对象
func DB(name ...string) *gorm.DB {
	return database(name...)
}

// Database returns an instance of database ORM object with specified configuration group name.
// Note that it panics if any error occurs duration instance creating.
func database(name ...string) *gorm.DB {
	var group = DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", keyDatabase, group)
	db := instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		// Create a new ORM object with given configurations.
		// fmt.Println("1.执行数据对象方法", group)
		if db, err := NewConnectDB(group); err == nil {
			return db
		} else {
			// If panics, often because it does not find its configuration for given group.
			panic(err)
		}
	})

	if db != nil {
		return db.(*gorm.DB)
	}
	return nil
}

// 清除DB实例对象
func Clear(name ...string) {
	var group = DefaultGroupName
	if len(name) > 0 {
		group = name[0]
	}
	instance.ClearOne(fmt.Sprintf("%s.%s", keyDatabase, group))
}

// NewConnectDB creates and returns an ORM object with global configurations.
func NewConnectDB(groupName string) (*gorm.DB, error) {
	// Retrieve the database configuration and determine whether the database configuration exists.
	// 1.仅在终端输出
	newLogger := logger.Discard //日志级别：Silent、Error、Warn、Info
	//2.在debug模式下才判处理文件日志（同时输出到文件和控制台）
	newLogger = logger.New(
		logWriter{stdout: true}, //Set log type
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level,eg: Silent、Error、Warn、Info
			IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)
	pathUrl, err := os.Getwd()
	if err != nil {
		pathUrl = ""
	}
	// ails runtime.AppDataDirectory，应用数据目录，数据库放这里
	dbConf_arr := map[string]any{"type": conf.Type, "dbFile": filepath.Join(pathUrl, "resource", "db", "data.db")}
	db, err := gorm.Open(drivers.DbOpen(dbConf_arr), &gorm.Config{
		Logger: newLogger,
		DryRun: false, //生成 SQL 但不执行，可以用于准备或测试生成的 SQL
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Prefix,                       // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   true,                              // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect db fail: %w", err)
	}

	//Set up the connection pool/设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Set up the connection pool fail: %w", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(conf.MaxOpen)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Hour)
	// SetConnMaxIdleTime 设置连接最大空闲时间（空闲30秒就关闭）
	sqlDB.SetConnMaxIdleTime(time.Duration(conf.MaxIdletime) * time.Second)

	// Enable Debug mode/开启Debug模式
	if conf.RunEnv == "debug" {
		return db.Debug(), nil
	}
	return db, nil
}

// Customized SQL log output (to file)/自定义的sql日志输出（到文件）
func (w logWriter) Printf(format string, args ...interface{}) {
	msg := "\n" + fmt.Sprintf(format, args...)
	if conf.RunEnv == "debug" && w.stdout {
		fmt.Println(msg) // 打印到控制台
	}
	//处理文件名:
	var filePath string = ""
	path, err := os.Getwd() //获取当前路径
	if err != nil {
		filePath = "./runtime/log/gorm/"
	} else {
		filePath = filepath.Join(path, "/runtime/log/gorm/")
	}
	os.MkdirAll(filePath, 0755)
	filePathName := filepath.Join(filePath, time.Now().Format("2006-01-02")+".log")
	// 替换掉彩色打印符号
	msg = strings.ReplaceAll(msg, logger.Reset, "")
	msg = strings.ReplaceAll(msg, logger.Red, "")
	msg = strings.ReplaceAll(msg, logger.Green, "")
	msg = strings.ReplaceAll(msg, logger.Yellow, "")
	msg = strings.ReplaceAll(msg, logger.Blue, "")
	msg = strings.ReplaceAll(msg, logger.Magenta, "")
	msg = strings.ReplaceAll(msg, logger.Cyan, "")
	msg = strings.ReplaceAll(msg, logger.White, "")
	msg = strings.ReplaceAll(msg, logger.BlueBold, "")
	msg = strings.ReplaceAll(msg, logger.MagentaBold, "")
	msg = strings.ReplaceAll(msg, logger.RedBold, "")
	msg = strings.ReplaceAll(msg, logger.YellowBold, "")

	//写入文件
	file, err := os.OpenFile(filePathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("日志文件的打开错误 :", err)
	}
	defer file.Close()
	if _, err := file.WriteString(msg); err != nil {
		fmt.Println("写入日志文件错误 :", err)
	}
}

// Offset 处理分页偏移量，页数 (pageNo) 转 偏移量 (offset)
func Offset(PageNo, PageSize int) int {
	if PageNo <= 0 {
		PageNo = 1
	}
	if PageSize <= 0 {
		PageSize = 10 // 默认每页10条
	}
	return (PageNo - 1) * PageSize
}
