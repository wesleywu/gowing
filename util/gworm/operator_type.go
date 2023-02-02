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

type OperatorType uint8

const (
	EQ OperatorType = iota
	NE
	GT
	GTE
	LT
	LTE
	Like
	NotLike
	Null
	NotNull
)

var (
	OperatorTypeNames = map[OperatorType]string{
		EQ:      "eq",
		NE:      "ne",
		GT:      "gt",
		GTE:     "gte",
		LT:      "lt",
		LTE:     "lte",
		Like:    "like",
		NotLike: "notlike",
		Null:    "null",
		NotNull: "notnull",
	}
	OperatorTypeValues = map[string]OperatorType{
		"eq":      EQ,
		"ne":      NE,
		"gt":      GT,
		"gte":     GTE,
		"lt":      LT,
		"lte":     LTE,
		"like":    Like,
		"notlike": NotLike,
		"null":    Null,
		"notnull": NotNull,
	}
)

func (s *OperatorType) String() string {
	return OperatorTypeNames[*s]
}

func (s *OperatorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *OperatorType) UnmarshalJSON(data []byte) (err error) {
	var typeString string
	if err := json.Unmarshal(data, &typeString); err != nil {
		return err
	}
	if *s, err = ParseOperatorType(typeString); err != nil {
		return err
	}
	return nil
}

func ParseOperatorType(s string) (OperatorType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := OperatorTypeValues[s]
	if !ok {
		return EQ, gerror.Newf("%q 不是合法的 OperatorType", s)
	}
	return value, nil
}
