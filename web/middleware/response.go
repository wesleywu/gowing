package middleware

import (
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/wesleywu/gowing/errors/gwerror"
)

func ResponseJsonWrapper(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.Status >= http.StatusInternalServerError {
		r.Response.ClearBuffer()
		r.Response.Writeln("服务器开小差了，请稍后再试吧！")
	} else {
		if err := r.GetError(); err != nil {
			r.Response.ClearBuffer()
			_, ok := err.(gvalid.Error)
			if ok {
				validationError := gwerror.NewBadRequestErrorf(r.GetBodyString(), err.Error())
				r.Response.Status = validationError.Code
				r.Response.WriteJsonExit(validationError)
			}
			validationError, ok := err.(gwerror.RequestError)
			if ok {
				r.Response.Status = validationError.Code
				r.Response.WriteJsonExit(validationError)
			}
			serviceError, ok := err.(gwerror.ServiceError)
			if ok {
				r.Response.Status = serviceError.Code
				r.Response.WriteJsonExit(gjson.MustEncodeString(serviceError))
			}
			r.Response.Status = 500
			r.Response.WriteJsonExit(err)
		} else {
			handlerRes := r.GetHandlerResponse()
			r.Response.WriteJsonExit(handlerRes)
		}
	}
}
