package {{.MPkgName}}

{{if .CRUD}}
import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"{{.ProjectName}}/lib/db"
)

// {{.ModelName}}CollectionName is the {{.ModelName}} collection name
const {{.ModelName}}CollectionName = "{{.CollectionName}}"
{{end}}

// {{.ModelName}} {{.ModelComment}}
type {{.ModelName}} struct {
{{if .CRUD}}
	Model `bson:",inline"`
{{end}}
	{{with .Fields}} {{range .}}
    {{.Name}} {{.Type}} `bson:"{{.TagName}}" json:"{{.TagName}}"` // {{.Comment}}{{end}} {{end}}
}

// B{{.ModelName}} defines {{.ModelName}} model bson field name
var B{{.ModelName}} struct {
    {{with .Fields}} {{range .}}
    {{.Name}} string {{end}} {{end}}
}

func init(){
    {{with .Fields}} {{range .}}
    B{{$.ModelName}}.{{.Name}} = "{{.TagName}}" {{end}} {{end}}
}

{{if .CRUD}}
// Insert should create a new {{.ModelName}}
func (m *{{.ModelName}}) Insert() (err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		m.Init()
		err = c.Insert(m)
	})
	return
}

func ensure{{.ModelName}}Index() {
	ensureIndex({{.ModelName}}CollectionName,
		mgo.Index{Key: []string{BDeletedAt}},
	)
}
// Find{{.ModelName}}s returns batch of {{.ModelName}}s
func Find{{.ModelName}}s (selector bson.M, fields bson.M, skip int, limit int, sort ...string) (models []*{{.ModelName}}, err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		err = c.Find(selector).Select(fields).Sort(sort...).Skip(skip * limit).Limit(limit).All(&models)
	})
	return
}

// Find{{.ModelName}} returns a {{.ModelName}}
func Find{{.ModelName}}(selector, fields bson.M) (model *{{.ModelName}}, err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		err = c.Find(selector).Select(fields).One(&model)
		return
	})
	return
}

// FindAndModify{{.ModelName}} find and modify a {{.ModelName}}
func FindAndModify{{.ModelName}}(selector bson.M, change mgo.Change) (model *{{.ModelName}}, err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		_, err = c.Find(selector).Apply(change, &model)
	})
	return
}

// Update{{.ModelName}} update a {{.ModelName}}
func Update{{.ModelName}}(selector bson.M, update bson.M) (err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		err = c.Update(selector, update)
	})
	return
}

// Delete{{.ModelName}} remove a {{.ModelName}}
func Delete{{.ModelName}}(selector bson.M) (err error) {
	db.Exec({{.ModelName}}CollectionName, func(c *mgo.Collection) {
		err = c.Remove(selector)
	})
	return
}
{{end}}