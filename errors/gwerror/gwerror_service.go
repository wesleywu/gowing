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

	"github.com/gogf/gf/v2/encoding/gjson"
)

type ServiceError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	ReqBody string `json:"reqBody,omitempty"`
	Err     error  `json:"err"`
}

func (e ServiceError) Error() string {
	return e.Message + e.Err.Error()
}
func WrapServiceErrorf(err error, req interface{}, format string, v ...any) ServiceError {
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
	code := 500
	if rerr, ok := err.(RequestError); ok {
		code = rerr.Code
	}
	return ServiceError{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
		ReqBody: reqBody,
		Err:     err,
	}
}
