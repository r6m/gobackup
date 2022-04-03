package database

import (
	"os"
	"path"

	"github.com/huacnlee/gobackup/helper"
	"github.com/huacnlee/gobackup/logger"
)

// PostgreSQL database
//
// type: postgresql
// host: localhost
// port: 5432
// database: test
// username:
// password:
type PostgreSQL struct {
	Base
	host        string
	port        string
	database    string
	username    string
	password    string
	dumpCommand string
}

func (ctx PostgreSQL) perform() (err error) {
	viper := ctx.viper
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 5432)

	ctx.host = viper.GetString("host")
	ctx.port = viper.GetString("port")
	ctx.database = viper.GetString("database")
	ctx.username = viper.GetString("username")
	ctx.password = viper.GetString("password")
	ctx.dumpCommand = "pg_dump"

	if ctx.database == "" {
		ctx.dumpCommand = "pg_dumpall"
		logger.Warn("postgres database is not specified, using 'pg_dumpall' to backup all databases")
	}

	err = ctx.dump()
	return
}

func (ctx *PostgreSQL) args() []string {
	// mysqldump command
	dumpArgs := []string{}

	if len(ctx.host) > 0 {
		dumpArgs = append(dumpArgs, "--host="+ctx.host)
	}
	if len(ctx.port) > 0 {
		dumpArgs = append(dumpArgs, "--port="+ctx.port)
	}
	if len(ctx.username) > 0 {
		dumpArgs = append(dumpArgs, "--username="+ctx.username)
	}

	if ctx.database != "" {
		dumpArgs = append(dumpArgs, ctx.database)
	}

	dumpFilePath := path.Join(ctx.dumpPath, ctx.database+".sql")
	dumpArgs = append(dumpArgs, "-f", dumpFilePath)

	return dumpArgs
}

func (ctx *PostgreSQL) dump() error {
	logger.Info("-> Dumping PostgreSQL...")
	if len(ctx.password) > 0 {
		os.Setenv("PGPASSWORD", ctx.password)
	}
	_, err := helper.Exec(ctx.dumpCommand, ctx.args()...)
	if err != nil {
		return err
	}
	logger.Info("dump path:", ctx.dumpPath)
	return nil
}
