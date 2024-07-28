package main

import (
	"github.com/spf13/cobra"
	"rag-new/internal/base/conf"
	"strconv"
	"sync"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动 http
	go func() {
		// refresh
		app.Service.Jwks.SetupAuthRefresh()
		var addr = app.Config.Http.Host + ":" + strconv.Itoa(app.Config.Http.Port)
		app.Logger.Sugar.Info("Listening and serving HTTP on ", addr)

		err := app.HttpServer.BizRouter().Run(app.Config.Http.Host + ":" + strconv.Itoa(app.Config.Http.Port))
		if err != nil {
			panic(err)
			return
		}

		wg.Done()
	}()

	// 启动 metrics
	if app.Config.Metrics.Enabled {
		go func() {
			err := app.HttpServer.MetricRouter().Run(app.Config.Metrics.Host + ":" + strconv.Itoa(app.Config.Metrics.Port))
			if err != nil {
				panic(err)
				return
			}
			wg.Done()

		}()
	}

	wg.Wait()

}
