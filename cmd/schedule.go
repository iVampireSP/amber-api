package main

import (
	"context"
	"github.com/spf13/cobra"
	"rag-new/internal/base"
	"rag-new/internal/batch"
	"sync"
	"time"
)

func init() {
	rootCmd.AddCommand(scheduleCmd)
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule commands",
	Long:  `Schedule commands`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := CreateApp()
		if err != nil {
			panic(err)
			return
		}

		runSchedule(app)
	},
}

func runSchedule(app *base.Application) {
	var wg sync.WaitGroup

	wg.Add(1)
	// 启动一个定时器
	go func() {
		app.Logger.Sugar.Info("Batch DeleteExpiredChats is ready.")
		for {
			var chatDeleteExpired = &batch.ChatDeleteExpired{
				BeforeTime:  time.Now(),
				ChatService: app.Service.Chat,
			}
			err := app.Batch.DeleteExpiredChats(context.Background(), chatDeleteExpired)
			if err != nil {
				app.Logger.Sugar.Error(err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	wg.Wait()
}
