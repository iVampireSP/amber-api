package cmd

import (
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"rag-new/internal/entity"
)

// Dynamic SQL
//type Querier interface {
//	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
//	FilterWithNameAndRole(name, role string) ([]gen.T, error)
//}

func init() {
	RootCmd.AddCommand(gormGenCmd)
}

var gormGenCmd = &cobra.Command{
	Use: "gorm-gen",
	Run: func(cmd *cobra.Command, args []string) {
		gormGen()
	},
}

func gormGen() {
	app, err := CreateApp()
	if err != nil {
		panic(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/dao",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(app.GORM)

	g.ApplyBasic(
		entity.Chat{},
		entity.ChatMessage{},
		entity.Assistant{},
		entity.AssistantShare{},
		entity.File{},
		entity.Tool{},
		entity.AssistantTool{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	//g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	g.Execute()
}
