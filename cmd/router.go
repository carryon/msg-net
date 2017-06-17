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

package cmd

import (
	"strings"

	"github.com/bocheninc/msg-net/config"
	"github.com/bocheninc/msg-net/logger"
	"github.com/bocheninc/msg-net/net/common"
	"github.com/bocheninc/msg-net/router"
	"github.com/spf13/cobra"
)

var routerID string
var routerAddress string
var routerDiscovery string

// routerCmd represents the router command
var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "start rounter service",
	Long:  `supply rounter service for communication of crossed blockchain.`,
	Run:   runRouter,
}

func init() {
	RootCmd.AddCommand(routerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routerCmd.PersistentFlags().String("foo", "", "A help for foo")
	routerCmd.PersistentFlags().StringVar(&routerID, "id", "", "id name")
	routerCmd.PersistentFlags().StringVar(&routerAddress, "address", "", "server listen address")
	routerCmd.PersistentFlags().StringVar(&routerDiscovery, "discovery", "", "discovery addresses")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// routerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	config.BindPFlag("router.id", routerCmd.PersistentFlags().Lookup("id"))
	config.BindPFlag("router.address", routerCmd.PersistentFlags().Lookup("address"))
	config.BindPFlag("router.discovery", routerCmd.PersistentFlags().Lookup("discovery"))
}

func runRouter(cmd *cobra.Command, args []string) {
	logger.SetOut()
	address := config.GetString("router.address")
	if config.GetBool("router.addressAutoDetect") {
		strs := strings.Split(address, ":")
		if ip := common.ChooseInterface(); ip != "" {
			address = strings.Replace(address, strs[0], ip, -1)
		}
	}
	r := router.NewRouter(config.GetString("router.id"), address)
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln(err)
			r.Stop()
		}
	}()
	//go util.SysSignal(func() { r.Stop() })
	r.Start()
}
