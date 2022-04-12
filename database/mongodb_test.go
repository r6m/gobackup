package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongoDB_mongodump(t *testing.T) {
	base := Base{
		dumpPath: "/tmp/gobackup/test",
	}
	ctx := &MongoDB{
		Base:     base,
		host:     "127.0.0.1",
		port:     "4567",
		database: "hello",
		username: "foo",
		password: "bar",
		authdb:   "sssbbb",
		oplog:    true,
	}
	expect := "mongodump --host=127.0.0.1 --port=4567 --username=foo --password=bar --authenticationDatabase=sssbbb --db=hello --oplog --out=/tmp/gobackup/test"

	args := append([]string{mongodumpCli}, ctx.args()...)
	output := strings.Join(args, " ")

	assert.Equal(t, output, expect)
}
