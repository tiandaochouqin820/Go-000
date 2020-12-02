package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

var (
	userName  string = "root"
	password  string = "root"
	ipAddrees string = "127.0.0.1"
	port      int    = 3306
	dbName    string = "gk_go"
	charset   string = "utf8"
)

type User struct {
	Id   int
	Name string
}

//数据库连接失败直接panic
func connectMysql() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

//model层查询数据
func model() (*sql.Rows, error) {
	db := connectMysql()
	rows, err := db.Query("select * from User")
	if err != nil {
		return nil, errors.Wrap(err, "model:data not query")
	}
	return rows, nil
}

//dao层调用model函数获取数据
func dao() (*sql.Rows, error) {
	data, err := model()
	if err != nil {
		errors.WithMessage(err, "dao:data not get")
	}
	return data, nil
}

//biz层组装业务逻辑
func biz() (*sql.Rows, error) {
	data, err := dao()
	if err != nil {
		errors.WithMessage(err, "biz:data assembly fail")
	}
	return data, nil
}

//api gateway层记录日志
func route() {
	data, err := biz()
	if err != nil {
		log.Println("route: %v\n", err)
	}
	log.Println(data)
}

func main() {
	route()
}
