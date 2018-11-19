## 依赖

    go get github.com/astaxie/beego
    go get -u github.com/beego/bee
    go get github.com/garyburd/redigo/redis
    go get github.com/golang/freetype
    go get github.com/dgrijalva/jwt-go
    go get github.com/spf13/cast
    go get github.com/ethereum/go-ethereum
    go get github.com/go-sql-driver/mysql
    go get github.com/satori/go.uuid

## 构建后台管理页面

    cd front
    npm install
    npm run build

## 创建日志文件夹和上传文件夹
    mkdir -p /var/upload/uphp/gcexserver
    mkdir /var/log/kuangchi

## 创建数据库
    create database gcexserver default character set utf8;

## 启动

    bee run

## 创建超级管理员

初始化超级管理员:
    insert into admin(username, password, salt, super, create_time, last_time) value("admin", "3eee0d5225c27509342f9a1e3a11de59237827c15befcfc8085ac61abd3ba3e6bef2635ad31c0cbd822f328626d2ad359ed8", "iFkzSESAgZ", 1, NOW(), NOW());

username: admin

password: 123456

## 后台访问地址

    /kc/admin-front/
