// ===============================
// Database-driven unified entry point
// If you are not using MySQL, refer to the method for introducing MySQL, and expand the new database as follows. In DbOpen and DbDSN, methods for extending other databases have already been reserved.
// ------------------------------------
// 数据库驱动统一入口
// 如果你使用的不是Mysql，参考Mysql引入方法，扩展新的数据库，如下DbOpen和DbDSN中已经预留扩展其他数据库方法。
// ===============================
package drivers

import (
	"gorm.io/gorm"
)

// Initialize the database driver adapter/初始化 数据库 驱动程序适配器
func DbOpen(dbConf_arr map[string]any) gorm.Dialector {
	switch dbConf_arr["type"] {
	// case "mysql":
	// 	return GormMysql(dbConf_arr)
	// case "pgsql":
	// 	return GormPostgreSQL(dbConf_arr)
	// case "oracle":
	// 	return GormOracle(dbConf_arr)
	// case "mssql":
	// 	return GormMssql(dbConf_arr)
	case "sqlite":
		return GormSqlite(dbConf_arr)
	default:
		return GormSqlite(dbConf_arr)
	}
}

// Database DSN connection/数据库DSN链接
func DbDSN(dbConf_arr map[string]any) string {
	switch dbConf_arr["type"] {
	// case "mysql":
	// 	return GormMysqlDSN(dbConf_arr)
	// case "pgsql":
	// 	return GormPostgreSQLDSN(dbConf_arr)
	// case "oracle":
	// 	return GormOracleDSN(dbConf_arr)
	// case "mssql":
	// 	return GormMssqlDSN(dbConf_arr)
	case "sqlite":
		return GormSqliteDSN(dbConf_arr)
	default:
		return GormSqliteDSN(dbConf_arr)
	}
}
