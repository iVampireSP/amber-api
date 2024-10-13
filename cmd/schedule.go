package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"rag-new/internal/base"
	"rag-new/internal/batch"
	"sync"
	"time"
)

func init() {
	RootCmd.AddCommand(scheduleCmd)
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

	var ctx = context.Background()

	wg.Add(1)
	// 启动一个定时器
	go func() {
		app.Logger.Sugar.Info("Batch DeleteExpiredChats is ready.")
		for {
			var chatDeleteExpired = &batch.ChatDeleteExpired{
				BeforeTime:  time.Now(),
				ChatService: app.Service.Chat,
			}
			err := app.Batch.DeleteExpiredChats(ctx, chatDeleteExpired)
			if err != nil {
				app.Logger.Sugar.Error(err)
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	wg.Add(1)
	// 向量 chunk
	go func() {
		app.Logger.Sugar.Info("Vector Chunk is ready.")
		for {
			app.Logger.Sugar.Info("Vector Chunk running.")
			var vectorChunk = &batch.ChunkVectorBatch{
				LibraryService: app.Service.Library,
				DAO:            app.DAO,
			}
			err := app.Batch.ChunkVector(ctx, vectorChunk)
			if err != nil {
				app.Logger.Sugar.Error(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	wg.Add(1)
	// 未结算的 Token 计费
	go func() {
		app.Logger.Sugar.Info("Token billing is ready.")
		for {
			app.Logger.Sugar.Info("Token billing running.")
			var tokenBilling = &batch.UnsettedTokenBilling{
				AccountService:       app.Service.Account,
				UnsettedTokenService: app.Service.UnsettledToken,
				Config:               app.Config,
				DAO:                  app.DAO,
			}
			err := app.Batch.UnsettedTokenBilling(tokenBilling)
			if err != nil {
				app.Logger.Sugar.Error(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	wg.Wait()
}
