package main

import (
	"fmt"
	"os"

	"model/mysqlrepo/initdb"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	var db gorm.DB
	var err error

	db, err = gorm.Open("mysql", "root:yunjikeji@tcp(10.0.1.31:3306)/osinstall?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	//db.SingularTable(true)
	err = initdb.DropOsInstallTables(&db, nil)
	err = initdb.InitOsInstallTables(&db, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("done")
}
