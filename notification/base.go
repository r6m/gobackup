package notification

import (
	"github.com/huacnlee/gobackup/config"
	"github.com/huacnlee/gobackup/logger"
	"github.com/spf13/viper"
)

type Base struct {
	name   string
	model  config.ModelConfig
	viper  *viper.Viper
	report *config.Report
}

type Context interface {
	perform() error
}

func newBase(model config.ModelConfig, report *config.Report) (base Base) {
	return Base{
		name:   model.Name,
		model:  model,
		viper:  model.NotifyBy.Viper,
		report: report,
	}

}

func Run(model config.ModelConfig, report *config.Report) error {
	base := newBase(model, report)

	logger.Info("------------ Notification -------------")
	defer logger.Info("------------ -------------\n")

	var ctx Context
	switch model.NotifyBy.Type {
	case "http":
		ctx = &HTTP{Base: base}
	default:
		logger.Info("No Notification Set")
		return nil
	}

	if err := ctx.perform(); err != nil {
		return err
	}

	return nil
}
