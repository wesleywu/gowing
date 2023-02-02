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

package gwhttpclient

import "net/http"

type HttpRequest struct {
	Method     string         // GET|POST|PUT|DELETE|OPTION
	Url        string         // 要请求的完整URL地址（包括协议、主机名、端口、Path路径、Query查询）
	Headers    *http.Header   // 请求使用的 Headers
	Cookies    []*http.Cookie // 请求使用的 Cookies
	Body       string         // 请求体字符串（GET|OPTION|DELETE请求传空字符串）
	RetryCount int            // 在成功之前，需要反复重试多少次，-1为重复1000次
}

type HttpResponse struct {
	Body       string       // 结果的Body字符串
	StatusCode int          // 结果的Http状态码
	Header     *http.Header // 结果的所有Http Header
}
