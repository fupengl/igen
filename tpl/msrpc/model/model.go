package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"igen/lib"
)

// Mongo通用字段
const (
	Bid        string = "_id"
	BCreatedAt string = "createdAt"
	BUpdatedAt string = "updatedAt"
	BDeletedAt string = "deletedAt"
	BSet       string = "$set"
	BPush      string = "$push"
	BExists    string = "$exists"
)

// Model 公共字段
type Model struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`                                  // 唯一ID
	CreatedAt lib.Time      `bson:"createdAt" json:"createdAt"`                     // 创建时间
	UpdatedAt lib.Time      `bson:"updatedAt" json:"updatedAt"`                     // 修改时间
	DeletedAt *lib.Time     `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"` // 删除时间
}

// Init 初始化
func (m *Model) Init() {
	t := lib.Time{Time: time.Now()}

	if !m.ID.Valid() {
		m.ID = bson.NewObjectId()
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = t
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = t
	}
}
