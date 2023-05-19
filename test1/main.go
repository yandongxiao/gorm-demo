package main

import (
    "context"
    "fmt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
    "log"
    "os"
    "time"
)

type ClusterTemplateSchemaTag struct {
    ID        uint   `gorm:"primarykey"`
    ClusterID uint   `gorm:"uniqueIndex:idx_cluster_id_key"`
    Key       string `gorm:"uniqueIndex:idx_cluster_id_key;column:tag_key"`
    Value     string `gorm:"column:tag_value"`
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy uint
    UpdatedBy uint
}

// NewSqliteDB
// file 参数是 sqlite 数据库文件路径，比如 /tmp/test.db，如果文件不存在，会自动创建
// 你可以通过sqlite3 /tmp/test.db 命令查看数据库内容
func NewSqliteDB(file string) (*gorm.DB, error) {
    orm, err := gorm.Open(sqlite.Open(file), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "tb_",
            SingularTable: true,
        },
        Logger: logger.New(
            log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
            logger.Config{
                SlowThreshold:             0, // print all logs
                LogLevel:                  logger.Error,
                IgnoreRecordNotFoundError: true,
                Colorful:                  true,
            },
        ),
    })

    return orm, err
}

func main() {
    db, err := NewSqliteDB("")
    if err != nil {
        panic(err)
    }

    // 自动创建表
    if err = db.AutoMigrate(&ClusterTemplateSchemaTag{}); err != nil {
        panic(err)
    }

    object := ClusterTemplateSchemaTag{
        ClusterID: 1,
    }

    ctx := context.Background()

    // 创建数据
    db.WithContext(ctx).Create(&object)

    // 获取数据
    // 对应的 sql 语句是：SELECT * FROM `tb_cluster_template_schema_tags` WHERE `tb_cluster_template_schema_tags`.`id` = 1 LIMIT 1
    get := ClusterTemplateSchemaTag{}
    db.WithContext(ctx).First(&get)

    // 更新数据
    // 根据Save方法的注释，它主要是用来更新数据的，如果你想要创建数据，可以使用Create方法
    get.ClusterID = 100
    res := db.WithContext(context.Background()).Save(&get)
    if res.Error != nil {
        panic(err)
    }

    get2 := ClusterTemplateSchemaTag{}
    db.WithContext(ctx).First(&get2)

    if get2.ClusterID != 100 {
        panic(fmt.Sprintf("update failed: got value %v", get2))
    }

    // 删除数据
    // 对应的 sql 语句是：DELETE FROM `tb_cluster_template_schema_tags` WHERE `tb_cluster_template_schema_tags`.`id` = 1
    db.WithContext(ctx).Delete(&get)
    get3 := ClusterTemplateSchemaTag{}
    res = db.WithContext(ctx).First(&get3)
    if res == nil || res.Error != gorm.ErrRecordNotFound {
        panic(fmt.Sprintf("delete failed: got value %v", get3))
    }
}
