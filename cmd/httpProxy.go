/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/viper"
	"github.com/spf13/cobra"

	"github.com/weiqiang333/infra-socks-http/internal/httpproxy"
)

// httpProxyCmd represents the httpProxy command
var httpProxyCmd = &cobra.Command{
	Use:   "httpProxy",
	Short: "Proxy converts socks to http",
	Long: `Proxy between socks and http:
	Proxy converts socks to http.`,
	Run: func(cmd *cobra.Command, args []string) {
		httpproxy.HttpProxy()
	},
}

func init() {
	rootCmd.AddCommand(httpProxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpProxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpProxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	httpProxyCmd.Flags().StringP("ListenAddress", "l", ":8080", "HTTP proxy local listening address")
	httpProxyCmd.Flags().StringP("SocksProxy","s", "127.0.0.1:9999", "Socks proxy")
	viper.BindPFlags(httpProxyCmd.Flags())
}
