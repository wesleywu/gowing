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

type MultiType uint8

const (
	Exact MultiType = iota
	Between
	NotBetween
	In
	NotIn
)

var (
	MultiTypeNames = map[MultiType]string{
		Exact:      "exact",
		Between:    "between",
		NotBetween: "notbetween",
		In:         "in",
		NotIn:      "notin",
	}
	MultiTypeValues = map[string]MultiType{
		"exact":      Exact,
		"between":    Between,
		"notbetween": NotBetween,
		"in":         In,
		"notin":      NotIn,
	}
)

func (s *MultiType) String() string {
	return MultiTypeNames[*s]
}

func (s *MultiType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *MultiType) UnmarshalJSON(data []byte) (err error) {
	var typeString string
	if err := json.Unmarshal(data, &typeString); err != nil {
		return err
	}
	if *s, err = ParseMultiType(typeString); err != nil {
		return err
	}
	return nil
}

func ParseMultiType(s string) (MultiType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := MultiTypeValues[s]
	if !ok {
		return Exact, gerror.Newf("%q 不是合法的 MultiType", s)
	}
	return value, nil
}
