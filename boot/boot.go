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

package boot

import (
	"fmt"
	"github.com/WesleyWu/gowing"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
)

func init() {
	_ = gtime.SetTimeZone("Asia/Shanghai") //设置系统时区
	showLogo()
	g.Log().SetFlags(glog.F_ASYNC | glog.F_TIME_DATE | glog.F_TIME_TIME | glog.F_TIME_MILLI | glog.F_FILE_SHORT)
}

func showLogo() {
	fmt.Println(strings.Join([]string{
		"   ______      _       __ _",
		"  / ____/____ | |     / /(_)____   ____ _",
		" / / __ / __ \\| | /| / // // __ \\ / __ `/",
		"/ /_/ // /_/ /| |/ |/ // // / / // /_/ /",
		"\\____/ \\____/ |__/|__//_//_/ /_/ \\__, /  Version: " + gowing.Version,
		"                                /____/",
	}, "\n"))
}
