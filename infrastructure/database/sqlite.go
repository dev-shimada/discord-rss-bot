package database

import (
	"fmt"
	"log"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDB() *gorm.DB {
	if err := RetryConnectDB(sqlite.Open("sqlite/rss_subscriptions.db"), &gorm.Config{}, 100); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected")
	if err := db.AutoMigrate(&model.Subscription{}, &model.RssEntry{}); err != nil {
		log.Fatalln(err)
	}
	return db
}

func RetryConnectDB(dialector gorm.Dialector, opt gorm.Option, count uint) error {
	var err error
	for count > 1 {
		if db, err = gorm.Open(dialector, opt); err != nil {
			time.Sleep(time.Second * 2)
			count--
			fmt.Printf("retry... coutn:%v\n", count)
			continue
		}
		break
	}
	return err
}

func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
