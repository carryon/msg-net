// Copyright (C) 2017, Beijing Bochen Technology Co.,Ltd.  All rights reserved.
//
// This file is part of msg-net
//
// The msg-net is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The msg-net is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/bocheninc/msg-net/cmd"
	"github.com/bocheninc/msg-net/config"
)

func main() {
	if port := config.GetString("profiler.port"); port != "" {
		go func() {
			http.ListenAndServe("0.0.0.0:"+port, nil)
		}()
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
