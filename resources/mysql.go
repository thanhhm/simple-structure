package resources

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("Connect to mysql error: ", err.Error())
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		break
	}
	return db, err
}
