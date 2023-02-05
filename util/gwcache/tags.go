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

package gwcache

import (
	"context"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

const (
	enabledDefault = false
	adapterDefault = "redis"
	ttlDefault     = uint32(600)
)

type Setting struct {
	Enabled bool
	Name    string // todo implements custom specified name
	Key     string // todo implements custom specified key
	Adapter string // todo implements memory adapter
	Ttl     uint32
}

func ParseCacheTag(ctx context.Context, cacheTag string) Setting {
	ctx, span := gtrace.NewSpan(ctx, "parseCacheTag")
	defer span.End()
	if cacheTag == "" {
		return Setting{
			Enabled: enabledDefault,
			Adapter: adapterDefault,
			Ttl:     ttlDefault,
		}
	}
	enabled := enabledDefault
	name := ""
	key := ""
	adapter := adapterDefault
	ttl := ttlDefault
	for _, s := range strings.Split(cacheTag, ",") {
		s = strings.Trim(s, " ")
		attrib := strings.Split(strings.Trim(s, " "), "=")
		if len(attrib) == 1 && attrib[0] == "true" {
			enabled = true
			continue
		}
		if len(attrib) == 2 {
			attribName := strings.Trim(attrib[0], " ")
			attribValue := strings.Trim(attrib[1], " ")
			switch attribName {
			case "name":
				name = attribValue
			case "key":
				key = attribValue
			case "adapter":
				adapter = attribValue
			case "ttl":
				ttl = gconv.Uint32(attribValue)
			}
		}
	}
	return Setting{
		Enabled: enabled,
		Name:    name,
		Key:     key,
		Adapter: adapter,
		Ttl:     ttl,
	}
}
