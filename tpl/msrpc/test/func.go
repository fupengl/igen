package tests

import (
	"gopkg.in/mgo.v2"

	"igen/lib/db"
)

// DropCollections 删除数据表
// 同时会删除索引
func DropCollections(collections ...string) {
	for _, collection := range collections {
		db.Exec(collection, func(c *mgo.Collection) {
			c.DropCollection()
		})
	}
}

// RemoveDocuments 清空数据
// 不会删除索引
func RemoveDocuments(collections ...string) {
	for _, collection := range collections {
		db.Exec(collection, func(c *mgo.Collection) {
			c.RemoveAll(nil)
		})
	}
}
