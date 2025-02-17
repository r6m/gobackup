package database

import (
	"fmt"
	"path"
	"strings"

	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
)

// MySQL database
//
// type: mysql
// host: 127.0.0.1
// port: 3306
// database:
// username: root
// password:
// additional_options:
type MySQL struct {
	Base
	host              string
	port              string
	database          string
	username          string
	password          string
	additionalOptions []string
}

func (ctx *MySQL) perform() (err error) {
	viper := ctx.viper
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("username", "root")
	viper.SetDefault("port", 3306)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.database = viper.GetString("database")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	addOpts := viper.GetString("additional_options")
	if len(addOpts) > 0 {
		ctx.additionalOptions = strings.Split(addOpts, " ")
	}

	// mysqldump command
	if ctx.database == "" {
		logger.Warn("mysql database is not specified, passing --all-databases")
	}

	err = ctx.dump()
	return
}

func (ctx *MySQL) dumpArgs() []string {
	dumpArgs := []string{}
	if len(ctx.host) > 0 {
		dumpArgs = append(dumpArgs, "--host", ctx.host)
	}
	if len(ctx.port) > 0 {
		dumpArgs = append(dumpArgs, "--port", ctx.port)
	}
	if len(ctx.username) > 0 {
		dumpArgs = append(dumpArgs, "-u", ctx.username)
	}
	if len(ctx.password) > 0 {
		dumpArgs = append(dumpArgs, `-p`+ctx.password)
	}
	if len(ctx.additionalOptions) > 0 {
		dumpArgs = append(dumpArgs, ctx.additionalOptions...)
	}

	if ctx.database == "" {
		dumpArgs = append(dumpArgs, "--all-databases")
	} else {
		dumpArgs = append(dumpArgs, ctx.database)
	}

	filename := ctx.database + ".sql"
	if ctx.database == "" {
		filename = "all-databases.sql"
	}

	dumpFilePath := path.Join(ctx.dumpPath, filename)
	dumpArgs = append(dumpArgs, "--result-file="+dumpFilePath)
	return dumpArgs
}

func (ctx *MySQL) dump() error {
	logger.Info("-> Dumping MySQL...")
	_, err := helper.Exec("mysqldump", ctx.dumpArgs()...)
	if err != nil {
		return fmt.Errorf("-> Dump error: %s", err)
	}
	logger.Info("dump path:", ctx.dumpPath)
	return nil
}
