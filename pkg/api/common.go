package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Unknwon/macaron"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/metrics"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/setting"
)

var (
	NotFound    = ApiError(404, "Not found", nil)
	ServerError = ApiError(500, "Server error", nil)
)

type Response interface {
	Handle(c *middleware.Context)
}

func wrap(action interface{}) macaron.Handler {
	return func(c *middleware.Context) {
		var res Response
		val, err := c.Invoke(action)
		if err == nil && val != nil && len(val) > 0 {
			res = val[0].Interface().(Response)
		} else {
			res = ServerError
		}

		res.Handle(c)
	}
}

type NormalResponse struct {
	status int
	body   []byte
	header http.Header
}

func (r *NormalResponse) Handle(c *middleware.Context) {
	out := c.Resp
	header := out.Header()
	for k, v := range r.header {
		header[k] = v
	}
	out.WriteHeader(r.status)
	out.Write(r.body)
}

func (r *NormalResponse) Cache(ttl string) *NormalResponse {
	return r.Header("Cache-Control", "public,max-age="+ttl)
}

func (r *NormalResponse) Header(key, value string) *NormalResponse {
	r.header.Set(key, value)
	return r
}

// functions to create responses

func Empty(status int) *NormalResponse {
	return Respond(status, nil)
}

func Json(status int, body interface{}) *NormalResponse {
	return Respond(status, body).Header("Content-Type", "application/json")
}

func ApiSuccess(message string) *NormalResponse {
	resp := make(map[string]interface{})
	resp["message"] = message
	return Respond(200, resp)
}

type HtmlResponse struct {
	status int
	view   string
	data   map[string]interface{}
}

func (r *HtmlResponse) Handle(c *middleware.Context) {
	c.Data = r.data
	c.HTML(r.status, r.view)
}

func HtmlErrorView(status int, message string, err error) *HtmlResponse {
	r := &HtmlResponse{status: status, data: make(map[string]interface{})}
	r.data["Title"] = message
	r.view = "500"

	if err != nil {
		log.Error(4, "%s: %v", message, err)
		if setting.Env != setting.PROD {
			r.data["ErrorMsg"] = err
		}
	}

	return r
}

type RedirectResponse struct {
	url string
}

func (r *RedirectResponse) Handle(c *middleware.Context) {
	c.Redirect(r.url)
}

func Redirect(url string) *RedirectResponse {
	if strings.HasPrefix(url, "/") {
		url = setting.AppSubUrl + url
	}
	return &RedirectResponse{url: url}
}

func ApiError(status int, message string, err error) *NormalResponse {
	resp := make(map[string]interface{})

	if err != nil {
		log.Error(4, "%s: %v", message, err)
		if setting.Env != setting.PROD {
			resp["error"] = err.Error()
		}
	}

	switch status {
	case 404:
		metrics.M_Api_Status_404.Inc(1)
		resp["message"] = "Not Found"
	case 500:
		metrics.M_Api_Status_500.Inc(1)
		resp["message"] = "Internal Server Error"
	}

	if message != "" {
		resp["message"] = message
	}

	return Json(status, resp)
}

func Respond(status int, body interface{}) *NormalResponse {
	var b []byte
	var err error
	switch t := body.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			return ApiError(500, "body json marshal", err)
		}
	}
	return &NormalResponse{
		body:   b,
		status: status,
		header: make(http.Header),
	}
}

func CheckQuota(c *middleware.Context, quotaDef *middleware.QuotaDef) Response {
	limitReached, err := middleware.QuotaReached(c, quotaDef)
	if err != nil {
		return ApiError(500, "Failed to check quota", err)
	}
	if limitReached {
		return ApiError(403, fmt.Sprintf("%s Quota reached", quotaDef.Name), err)
	}
	return nil
}
