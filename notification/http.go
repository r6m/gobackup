package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/huacnlee/gobackup/logger"
)

type HTTP struct {
	Base
	method  string
	url     string
	headers map[string]string
}

type httpPayload struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

const tplMessage = ``

func (ctx *HTTP) perform() error {
	ctx.method = strings.ToUpper(ctx.viper.GetString("method"))
	ctx.url = ctx.viper.GetString("url")
	ctx.headers = ctx.viper.GetStringMapString("headers")

	payload := &httpPayload{
		Name:    ctx.Base.model.Name,
		Status:  ctx.report.Status,
		Message: ctx.report.Message,
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(ctx.method, ctx.url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range ctx.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Error(fmt.Sprintf("%s: %s\n%s", ctx.method, ctx.url, string(body)))
		return errors.New("status code is not 200")
	}
	logger.Info("response body:", resp.Status)

	return nil
}
