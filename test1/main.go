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

    // 关于获取数据的进一步解释：
    // db.WithContext(ctx).Raw(sql string, values ...interface{}) 用来执行原生 sql 语句。
    // db.WithContext(ctx).Select(query interface{}, args ...interface{}) specify fields that you want when querying.
    // db.Where("cluster_id = ?", 1)
    //
    // db.Scan() VS db.Find() VS db.First()
    // db.Scan() 用来将查询结果扫描到一个结构体中，如果查询结果有多条，只会扫描第一条。如果查询结果为空，会返回 ErrRecordNotFound 错误。
    // db.Find() 用来将查询结果扫描到一个结构体切片中，如果查询结果为空，会返回空切片
    // db.First() 用来将查询结果扫描到一个结构体中，如果查询结果为空，会返回 ErrRecordNotFound 错误

    // offset 和 limit
    // 为什么会有两个offset和limit? 需要查询出总数和查询数据
    // 相当于 select * from tb_token where user_id = ? and code like ? limit ? offset ?
    // 相当于 select count(*) from tb_token where user_id = ? and code like ?
    // *gorm.DB 是可以被复用的，Find方法会执行一次查询，Count方法也会执行一次查询，所以需要将offset和limit重置。
    // result := d.db.WithContext(ctx).Table("tb_token").
    //     Where("user_id = ?", currentUser.GetID()).
    //     Where("code like ?", fmt.Sprintf("%s%%", generator.AccessTokenPrefix)).
    //     Offset(offset).Limit(limit).
    //     Find(&tokens).Offset(0).Limit(-1).Count(&total)

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

    // db.Exec VS db.Raw
    // db.Exec() 一般用来执行更新语句，不关心返回结果。

    // 3. add new tags
    // 之所以使用 Clauses 方法是因为 gorm 的 Create 方法不支持 ON DUPLICATE KEY UPDATE
    // ON DUPLICATE KEY UPDATE is a MariaDB/MySQL extension to the INSERT statement that, if it finds
    // a duplicate unique or primary key, will instead perform an UPDATE.
    // 下面的语句表示：如果 resource_type, resource_id, tag_key 三个字段的值在数据库中已经存在，则更新 tag_value 字段的值
    // result := d.db.WithContext(ctx).Clauses(clause.OnConflict{
    //     Columns: []clause.Column{
    //         {
    //             Name: "resource_type",
    //         }, {
    //             Name: "resource_id",
    //         }, {
    //             Name: "tag_key",
    //         },
    //     },
    //     DoUpdates: clause.AssignmentColumns([]string{"tag_value"}),
    // }).Create(tags)

    // 删除数据
    // 对应的 sql 语句是：DELETE FROM `tb_cluster_template_schema_tags` WHERE `tb_cluster_template_schema_tags`.`id` = 1
    db.WithContext(ctx).Delete(&get)
    get3 := ClusterTemplateSchemaTag{}
    res = db.WithContext(ctx).First(&get3)
    if res == nil || res.Error != gorm.ErrRecordNotFound {
        panic(fmt.Sprintf("delete failed: got value %v", get3))
    }

    // d.db.WithContext(ctx).Transaction
    // Transaction start a transaction as a block, return error will rollback, otherwise to commit.
}
