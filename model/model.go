package model

import (
	"os"
	"path"
	"time"

	"github.com/huacnlee/gobackup/archive"
	"github.com/huacnlee/gobackup/compressor"
	"github.com/huacnlee/gobackup/config"
	"github.com/huacnlee/gobackup/database"
	"github.com/huacnlee/gobackup/encryptor"
	"github.com/huacnlee/gobackup/logger"
	"github.com/huacnlee/gobackup/notification"
	"github.com/huacnlee/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
	Report *config.Report
}

// Perform model
func (ctx Model) Perform() {
	ctx.Report.StartTime = time.Now()
	logger.Info("======== " + ctx.Config.Name + " ========")
	logger.Info("WorkDir:", ctx.Config.DumpPath+"\n")

	defer func() {
		if r := recover(); r != nil {
			ctx.cleanup()
		}

		ctx.cleanup()
	}()

	defer func() {
		ctx.Report.Duration = time.Since(ctx.Report.StartTime)
		ctx.Report.EndTime = ctx.Report.StartTime.Add(ctx.Report.Duration)
		if err := notification.Run(ctx.Config, ctx.Report); err != nil {
			logger.Error(err)
		}
	}()

	err := database.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		ctx.Report.Status = "error"
		ctx.Report.Message = "Backup Error\n" + ctx.Report.Message + err.Error()
		return
	}

	if ctx.Config.Archive != nil {
		err = archive.Run(ctx.Config)
		if err != nil {
			logger.Error(err)
			ctx.Report.Status = "error"
			ctx.Report.Message = "Archiving Error\n" + ctx.Report.Message + err.Error()
			return
		}
		ctx.Report.Status = "success"
	}

	archivePath, err := compressor.Run(ctx.Config)
	if err != nil {
		logger.Error(err)
		ctx.Report.Status = "error"
		ctx.Report.Message = "Compressing Error\n" + ctx.Report.Message + err.Error()
		return
	}
	ctx.Report.Filename = path.Base(archivePath)
	ctx.Report.Status = "success"

	archivePath, err = encryptor.Run(archivePath, ctx.Config)
	if err != nil {
		logger.Error(err)
		ctx.Report.Status = "error"
		ctx.Report.Message = "Encrypting Error\n" + ctx.Report.Message + err.Error()
		return
	}
	ctx.Report.Filename = path.Base(archivePath)
	ctx.Report.Status = "success"

	err = storage.Run(ctx.Config, archivePath)
	if err != nil {
		logger.Error(err)
		ctx.Report.Status = "error"
		ctx.Report.Message = "Storing Error\n" + ctx.Report.Message + err.Error()
		return
	}
	ctx.Report.Status = "success"
	ctx.Report.Message = "Backup performed successfuly"

}

// Cleanup model temp files
func (ctx Model) cleanup() {
	logger.Info("Cleanup temp: " + ctx.Config.TempPath + "/\n")
	err := os.RemoveAll(ctx.Config.TempPath)
	if err != nil {
		logger.Error("Cleanup temp dir "+ctx.Config.TempPath+" error:", err)
	}
	logger.Info("======= End " + ctx.Config.Name + " =======\n\n")
}
