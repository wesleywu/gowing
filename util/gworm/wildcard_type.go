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

package gworm

import (
	"encoding/json"
	"github.com/gogf/gf/v2/errors/gerror"
	"strings"
)

type WildcardType uint8

const (
	None WildcardType = iota
	Contains
	StartsWith
	EndsWith
)

var (
	WildcardTypeNames = map[WildcardType]string{
		None:       "none",
		Contains:   "contains",
		StartsWith: "startswith",
		EndsWith:   "endswith",
	}
	WildcardTypeValues = map[string]WildcardType{
		"none":       None,
		"contains":   Contains,
		"startswith": StartsWith,
		"endswith":   EndsWith,
	}
)

func (s *WildcardType) String() string {
	return WildcardTypeNames[*s]
}

func (s *WildcardType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *WildcardType) UnmarshalJSON(data []byte) (err error) {
	var typeString string
	if err := json.Unmarshal(data, &typeString); err != nil {
		return err
	}
	if *s, err = ParseWildcardType(typeString); err != nil {
		return err
	}
	return nil
}

func ParseWildcardType(s string) (WildcardType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := WildcardTypeValues[s]
	if !ok {
		return None, gerror.Newf("%q 不是合法的 WildcardType", s)
	}
	return value, nil
}
