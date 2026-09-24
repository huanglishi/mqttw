// =====================================
// MySQL database driver operation
// Other databases are configured based on the GORM files and this MySQL database configuration.
// GORM Doc: https://gorm.io/zh_CN/docs/connecting_to_the_database.html
// ======================================
package drivers

import (
	"fmt"
	"gofly/internal/dao/drivers/sqlite"

	"gorm.io/gorm"
)

// DSN data source name
func GormSqliteDSN(dbConf_arr map[string]any) string {
	dsn := fmt.Sprintf("file:%v?_pragma=foreign_keys(1)&_pragma=journal_mode(DELETE)&_pragma=busy_timeout(5000)", dbConf_arr["dbFile"])
	return dsn
}

// GormSqlite is used to initialize the  driver adapter.
func GormSqlite(dbConf_arr map[string]any) gorm.Dialector {
	return sqlite.Open(GormSqliteDSN(dbConf_arr))
}
