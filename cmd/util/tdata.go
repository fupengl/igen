package util

// TData 模板数据
type TData struct {
	ProjectName    string
	SubProjectName string
	Fields         []*Field // 字段
	ModelName      string   // model 首字母大写
	ModelLowerName string   // model 首字母小写
	MPkgName       string   // model package name
	ModelComment   string   // model 说明
	CollectionName string   // model 全部小写

	CDir     string // controller folder
	CPkgName string // controller package name
	HPkgName string // v1 helper folder

	CRUD bool
}
