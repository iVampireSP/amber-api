package main

import (
	"github.com/spf13/cobra"
	"rag-new/internal/base/conf"
	"strconv"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		conf.CreateConfigIfNotExists()
		initHttpServer()
	},
}

func initHttpServer() {
	app, err := CreateApp()
	if err != nil {
		panic(err)
		return
	}

	if app.Config.Http.Host == "" {
		app.Config.Http.Host = "0.0.0.0"
	}

	if app.Config.Http.Port == 0 {
		app.Config.Http.Port = 8000
	}

	// refresh
	app.Service.Jwks.SetupAuthRefresh()

	var addr = app.Config.Http.Host + ":" + strconv.Itoa(app.Config.Http.Port)
	app.Logger.Sugar.Info("Listening and serving HTTP on ", addr)

	err = app.Gin.Run(addr)
	if err != nil {
		panic(err)
		return
	}

}
