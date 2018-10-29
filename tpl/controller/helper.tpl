package {{.HPkgName}}

import (
	"gopkg.in/mgo.v2/bson"

    "{{.ProjectName}}/{{.SubProjectName}}/{{.MPkgName}}"
)

type {{.ModelName}}JSON struct {
    ID bson.ObjectId `json:"id"` // ID
    {{with .Fields}} {{range .}}
    {{.Name}} {{.Type}} `json:"{{.TagName}}"` // {{.Comment}}{{end}} {{end}}
}

// {{.ModelName}}AsJSON format {{.ModelName}} to JSON
func {{.ModelName}}AsJSON(m *{{.MPkgName}}.{{.ModelName}}) {{.ModelName}}JSON{
    out := {{.ModelName}}JSON{}
    out.ID = m.ID
    {{with .Fields}} {{range .}}
    out.{{.Name}} = m.{{.Name}}{{end}} {{end}}
    return out
}

// {{.ModelName}}sAsJSON format {{.ModelName}}s to JSON
func {{.ModelName}}sAsJSON(mm []*{{.MPkgName}}.{{.ModelName}}) []{{.ModelName}}JSON{
    out := []{{.ModelName}}JSON{}
    for _, m := range mm {
        out = append(out, {{.ModelName}}AsJSON(m))
    }
    return out
}