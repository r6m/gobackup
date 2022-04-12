package database

import (
	"fmt"

	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
)

// MongoDB database
//
// type: mongodb
// host: 127.0.0.1
// port: 27017
// database:
// username:
// password:
// authdb:
// oplog: false
type MongoDB struct {
	Base
	host     string
	port     string
	database string
	username string
	password string
	authdb   string
	oplog    bool
}

var (
	mongodumpCli = "mongodump"
)

func (ctx *MongoDB) perform() (err error) {
	viper := ctx.viper
	viper.SetDefault("oplog", false)
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("username", "root")
	viper.SetDefault("port", 27017)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.database = viper.GetString("database")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.oplog = viper.GetBool("oplog")
	ctx.authdb = viper.GetString("authdb")

	err = ctx.dump()
	if err != nil {
		return err
	}
	return nil
}

func (ctx *MongoDB) args() []string {
	args := []string{}

	// connection args
	if ctx.host != "" {
		args = append(args, "--host="+ctx.host)
	}
	if ctx.port != "" {
		args = append(args, "--port="+ctx.port)
	}

	// credential args
	if ctx.username != "" {
		args = append(args, "--username="+ctx.username)
	}
	if ctx.password != "" {
		args = append(args, "--password="+ctx.password)
	}
	if ctx.authdb != "" {
		args = append(args, "--authenticationDatabase="+ctx.authdb)
	}

	// will dump all databases if not specified
	if ctx.database != "" {
		args = append(args, "--db="+ctx.database)
	}

	if ctx.oplog {
		args = append(args, "--oplog")
	}

	args = append(args, "--out="+ctx.dumpPath)

	return args
}

func (ctx *MongoDB) dump() error {
	out, err := helper.Exec(mongodumpCli, ctx.args()...)
	if err != nil {
		return fmt.Errorf("-> Dump error: %s", err)
	}
	logger.Info(out)
	logger.Info("dump path:", ctx.dumpPath)
	return nil
}
