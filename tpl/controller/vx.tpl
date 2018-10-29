package {{.CPkgName}}

import (
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"

	"{{.ProjectName}}/lib/db"
	"{{.ProjectName}}/lib/logger"
	"{{.ProjectName}}/lib/render"
	"{{.ProjectName}}/{{.SubProjectName}}/{{.MPkgName}}"
	"{{.ProjectName}}/{{.SubProjectName}}/{{.CDir}}/{{.CPkgName}}/{{.HPkgName}}"
)

type {{.ModelLowerName}}Form struct {
	 {{with .Fields}} {{range .}}
     {{.Name}} {{.Type}} // {{.Comment}}{{end}} {{end}}
}

type {{.ModelLowerName}}CreateForm struct {
	{{.ModelLowerName}}Form
}

type {{.ModelLowerName}}UpdateForm struct {
	{{.ModelLowerName}}Form
}

func (f *{{.ModelLowerName}}CreateForm) valid() error {

	return nil
}

func (f *{{.ModelLowerName}}UpdateForm) valid() error {

	return nil
}

// List{{.ModelName}}s list some {{.ModelName}}s
func List{{.ModelName}}s(c *gin.Context) {
	skip, limit, sorts := render.PageSizeSort(c, 0)

	selector := bson.M{}

	{{.ModelLowerName}}s := []*{{.MPkgName}}.{{.ModelName}}{}

	selector[{{.MPkgName}}.BDeletedAt] = bson.M{{"{"}}{{.MPkgName}}.BExists: false}
	total, err := db.QueryAsPage({{.MPkgName}}.{{.ModelName}}CollectionName, &{{.ModelLowerName}}s, selector, nil, skip, limit, sorts...)
	if err != nil {
		logger.Ctx(c).Error("fail to list {{.ModelLowerName}}s", logger.Err(err))
		render.Err500(c, "服务器错误")
		return
	}

	listJSON := []{{.HPkgName}}.{{.ModelName}}JSON{}
	for _, m := range {{.ModelLowerName}}s {
	    listJSON = append(listJSON, {{.HPkgName}}.{{.ModelName}}AsJSON(m))
	}

	render.OK(c, gin.H{
		"total": total,
		"{{.ModelLowerName}}s": listJSON,
	})
}

// Create{{.ModelName}} creates a {{.ModelName}}
func Create{{.ModelName}}(c *gin.Context) {
	form := new({{.ModelLowerName}}CreateForm)
	err := c.BindJSON(form)
	if err != nil {
		logger.Ctx(c).Error("invalid params", logger.Err(err))
		render.Err400(c)
		return
	}
	if err := form.valid(); err != nil {
		logger.Ctx(c).Error("invalid params", logger.Err(err))
		render.Err500(c, err.Error())
		return
	}

	{{.ModelLowerName}} := new({{.MPkgName}}.{{.ModelName}})
	{{with .Fields}} {{range .}}
        {{$.ModelLowerName}}.{{.Name}} = form.{{.Name}}{{end}} {{end}}

	if err := {{.ModelLowerName}}.Insert(); err != nil {
		logger.Ctx(c).Error("fail to create {{.ModelName}}", logger.Err(err))
		render.Err500(c, "创建失败")
		return
	}
	render.OK(c, gin.H{"{{.ModelLowerName}}": {{.HPkgName}}.{{.ModelName}}AsJSON({{.ModelLowerName}})})
}

// Show{{.ModelName}} show the {{.ModelName}} details
func Show{{.ModelName}}(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		logger.Ctx(c).Error("invalid ID")
		render.Err404(c)
		return
	}

	{{.ModelLowerName}}, err := {{.MPkgName}}.Find{{.ModelName}}(bson.M{"_id": bson.ObjectIdHex(id), {{.MPkgName}}.BDeletedAt: bson.M{{"{"}}{{.MPkgName}}.BExists: false}}, nil)
	if err != nil {
		logger.Ctx(c).Error(err.Error())
		render.Err404(c)
		return
	}

	render.OK(c, gin.H{"{{.ModelLowerName}}": {{.HPkgName}}.{{.ModelName}}AsJSON({{.ModelLowerName}})})
}

// Update{{.ModelName}} update the {{.ModelName}} info
func Update{{.ModelName}}(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		logger.Ctx(c).Error("invalid ID")
		render.Err404(c)
		return
	}

	form := new({{.ModelLowerName}}UpdateForm)
	err := c.BindJSON(form)
	if err != nil {
		logger.Ctx(c).Error("invalid params", logger.Err(err))
		render.Err400(c)
		return
	}
	if err := form.valid(); err != nil {
		logger.Ctx(c).Error("invalid params", logger.Err(err))
		render.Err500(c, err.Error())
		return
	}

    update := bson.M{}
    {{with .Fields}} {{range .}}
    update[{{$.MPkgName}}.B{{$.ModelName}}.{{.Name}}] = form.{{.Name}}{{end}} {{end}}

	if len(update) == 0 {
		render.Err500(c, "没有数据需要更改")
		return
	}
	update[{{.MPkgName}}.BUpdatedAt] = time.Now()

	{{.ModelLowerName}}, err := {{.MPkgName}}.FindAndModify{{.ModelName}}(bson.M{{"{"}}{{.MPkgName}}.Bid: bson.ObjectIdHex(id)}, mgo.Change{
	    Update: bson.M{{"{"}}{{.MPkgName}}.BSet: update},
	    ReturnNew: true,
	    Upsert: false,
	})

	if err != nil {
		logger.Ctx(c).Error("fail to update {{.ModelName}}", logger.Err(err))
		render.Err500(c, "更改失败")
		return
	}

	render.OK(c, gin.H{"{{.ModelLowerName}}": {{.HPkgName}}.{{.ModelName}}AsJSON({{.ModelLowerName}})})
}
