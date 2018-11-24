// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rjacobs31/trees-against-humanity-server/internal"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a Trees Against Humanity server instance",
	Run: func(cmd *cobra.Command, args []string) {
		addr := ":" + strconv.Itoa(viper.GetInt("port"))
		origins := viper.GetStringSlice("allowed-origins")
		config := internal.ServeConfig{
			Address:        addr,
			AllowedOrigins: origins,
		}
		internal.Serve(config)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntP("port", "p", 8000, "Port of the TAH server")
	serveCmd.Flags().StringArray("allowed-origins", []string{"*"}, "Allowed origins according to CORS standard")
	serveCmd.Flags().StringP("secret", "s", "secret-key", "Key used for encrypting session data")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("allowed-origins", serveCmd.Flags().Lookup("allowed-origins"))
	viper.BindPFlag("secret", serveCmd.Flags().Lookup("secret"))
}
