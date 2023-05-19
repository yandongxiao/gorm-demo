该测试场景包括：

1. 使用 gorm 插入数据、查询数据、更新数据、删除数据。
2. 使用内存型数据库来作为 gorm 的底层存储。方便测试。
3. 如果对 gorm 的操作进行单元测试。

# 安装 sqlite3

```bash
# 在mac上安装 sqlite3
brew install sqlite3

# 启动 sqlite3
sqlite3

# 查看所有数据库表
sqlite> .tables

# 创建数据库表
sqlite> create table users (id integer primary key autoincrement, name text, age integer);

# 查看数据库表结构
sqlite> .schema users

# 对数据库表进行增删改查
sqlite> insert into users (name, age) values ('张三', 18);
sqlite> select * from users;
sqlite> update users set age = 20 where name = '张三';
sqlite> delete from users where name = '张三';
```

# 
