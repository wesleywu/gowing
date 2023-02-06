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

package gwtypes

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"
)

func ParseOperatorType(s string) (OperatorType, error) {
	if g.IsEmpty(s) {
		return OperatorType_EQ, nil
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := OperatorType_value[s]
	if !ok {
		return OperatorType_EQ, gerror.Newf("%q 不是合法的 OperatorType", s)
	}
	return OperatorType(value), nil
}

func ParseMultiType(s string) (MultiType, error) {
	if g.IsEmpty(s) {
		return MultiType_In, nil
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := MultiType_value[s]
	if !ok {
		return MultiType_Exact, gerror.Newf("%q 不是合法的 MultiType", s)
	}
	return MultiType(value), nil
}

func ParseWildcardType(s string) (WildcardType, error) {
	if g.IsEmpty(s) {
		return WildcardType_None, nil
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := WildcardType_value[s]
	if !ok {
		return WildcardType_None, gerror.Newf("%q 不是合法的 WildcardType", s)
	}
	return WildcardType(value), nil
}

func MustParseOperatorType(s string, d OperatorType) OperatorType {
	if g.IsEmpty(s) {
		return d
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := OperatorType_value[s]
	if !ok {
		return d
	}
	return OperatorType(value)
}

func MustParseMultiType(s string, d MultiType) MultiType {
	if g.IsEmpty(s) {
		return d
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := MultiType_value[s]
	if !ok {
		return d
	}
	return MultiType(value)
}

func MustParseWildcardType(s string, d WildcardType) WildcardType {
	if g.IsEmpty(s) {
		return d
	}
	s = strings.TrimSpace(gstr.CaseCamel(s))
	value, ok := WildcardType_value[s]
	if !ok {
		return d
	}
	return WildcardType(value)
}
