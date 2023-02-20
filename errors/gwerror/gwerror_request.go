/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gwerror

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/lib/pq"
)

type RequestError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	ReqBody string `json:"reqBody,omitempty"`
}

func (e RequestError) Error() string {
	return e.Message
}

func NewRequestErrorf(code int, req interface{}, format string, v ...any) RequestError {
	var (
		reqBody string
		ok      bool
	)
	reqBody = ""
	if req != nil {
		if reqBody, ok = req.(string); !ok {
			reqBody = gjson.MustEncodeString(req)
		}
	}
	return RequestError{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
		ReqBody: reqBody,
	}
}

func NewBadRequestErrorf(req interface{}, format string, v ...any) RequestError {
	return NewRequestErrorf(400, req, format, v...)
}

func NewNotFoundErrorf(req interface{}, format string, v ...any) RequestError {
	return NewRequestErrorf(404, req, format, v...)
}

func NewPkConflictErrorf(req interface{}, format string, v ...any) RequestError {
	return NewRequestErrorf(409, req, format, v...)
}

func NewDataTooLongErrorf(req interface{}, format string, v ...any) RequestError {
	return NewRequestErrorf(413, req, format, v...)
}

func DbErrorToRequestError(req interface{}, err error, dbType string) (error, bool) {
	if err == nil {
		return nil, false
	}
	underlyingError := gerror.Unwrap(err)
	if underlyingError == nil {
		return nil, false
	}
	switch dbType {
	case "mysql":
		if driverError, ok := underlyingError.(*mysql.MySQLError); ok {
			switch driverError.Number {
			case 1062:
				return NewPkConflictErrorf(req, "%v: %v", "主键冲突", driverError.Error()), true
			case 1406:
				return NewDataTooLongErrorf(req, "%v: %v", "数据过长", driverError.Error()), true
			}
		}
	case "pgsql":
		if driverError, ok := underlyingError.(*pq.Error); ok {
			switch driverError.Code {
			case "23505":
				return NewPkConflictErrorf(req, "%v: %v", "主键冲突", driverError.Error()), true
			case "22001":
				return NewDataTooLongErrorf(req, "%v: %v", "数据过长", driverError.Error()), true
			}
		}
	}
	return err, false
}
