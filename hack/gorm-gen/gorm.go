package main

import (
	"gorm.io/gen"
	"rag-new/internal/entity"
)

// Dynamic SQL
//type Querier interface {
//	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
//	FilterWithNameAndRole(name, role string) ([]gen.T, error)
//}

func main() {
	//app, err := cmd.CreateApp()
	//if err != nil {
	//	panic(err)
	//}
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../internal/dao",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	//g.UseDB(app.GORM)

	g.ApplyBasic(
		entity.Chat{},
		entity.ChatMessage{},
		entity.Assistant{},
		entity.AssistantKey{},
		entity.File{},
		entity.Tool{},
		entity.AssistantTool{},
		entity.Embedding{},
		entity.Memory{},
		entity.Library{},
		entity.Document{},
		entity.DocumentChunk{},
		//entity.UserFile{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	g.Execute()
}
