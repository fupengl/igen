package model

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2"

	"igen/lib/db"
	"igen/lib/logger"
)

// EnsureIndexes 初始化所有表单索引
func EnsureIndexes() {

}

// ensures an index to the collection
func ensureIndex(collectionName string, indexes ...mgo.Index) error {
	if len(indexes) == 0 {
		return nil
	}

	errs := []string{}
	var err error
	db.Exec(collectionName, func(c *mgo.Collection) {
		for _, index := range indexes {
			err = c.EnsureIndex(index)
			if err != nil {
				logger.Errorf("Ensure %s Index error: %s", collectionName, err.Error())
				errs = append(errs, err.Error())
			}

		}
	})
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "; "))
}
