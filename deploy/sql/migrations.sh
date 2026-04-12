# 创建迁移文件
migrate create -ext sql -dir deploy/sql/migrations -seq modify_user_pk

## 编写SQL迁移文件

## 应用迁移
migrate -database "mysql://root:go-chat@tcp(124.221.43.166:13306)/DiceTales" -path deploy/sql/migrations up

# 假设从直接连接数据库生成最新模型代码，二选一
goctl model mysql datasource -url="root:go-chat@tcp(124.221.43.166:13306)/DiceTales" -table="user" -dir="./apps/user/model"
goctl model mysql ddl -src deploy/sql/user.sql -dir ./apps/user/model