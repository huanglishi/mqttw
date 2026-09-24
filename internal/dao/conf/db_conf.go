package conf

// 数据库配置
const (
	RunEnv      = "release" //环境状态：生成=release，开发=debug
	Type        = "sqlite"  //数据库类型
	Prefix      = ""        //table name prefix
	MaxIdle     = 10        //连接池最大闲置的连接数
	MaxOpen     = 32        //连接池最大打开的连接数
	MaxLifetime = 1         //连接对象可重复使用的时间长度（单位：小时）
	MaxIdletime = 30        //设置连接最大空闲时间（单位：秒）
)
