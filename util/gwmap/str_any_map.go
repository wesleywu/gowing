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

package gwmap

type StrStrMap struct {
	InnerMap map[string]interface{}
}

func (s *StrStrMap) Get(key string) string {
	if s.InnerMap == nil {
		return ""
	}
	item, ok := s.InnerMap[key].([]string)
	if !ok {
		return ""
	}
	if len(item) == 0 {
		return ""
	}
	return item[0]
}

func (s *StrStrMap) Set(key string, value string) {
	if s.InnerMap == nil {
		s.InnerMap = map[string]interface{}{}
	}
	s.InnerMap[key] = value
}

func (s *StrStrMap) Keys() []string {
	out := make([]string, 0, len(s.InnerMap))
	for key := range s.InnerMap {
		out = append(out, key)
	}
	return out
}
