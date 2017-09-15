// Copyright © 2017 Brian Ketelsen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/atotto/clipboard"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, Alias, Channel, Event string

var baseURL string

var (
	CommitHash string
	BuildTime  string
	Tag        string
)

type Submission struct {
	URL       string `json:"url"`
	ShortCode string `json:"short_code"`
}

type Response struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	Error     string `json:"error"`
}

const track = "?utm"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bjk URL SHORTCODE",
	Short: "a URL Shortening service and corresponding command line tool",
	Long: `bjk is a URL shortening service that automatically appends
the appropriate tracking tags to a URL.  The command line tool can be
used to submit a new link using personalized values.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("url", args[0])

		fmt.Println("short_code", args[1])
		err := submit(args[0], args[1])
		if err != nil {
			fmt.Println("Error: ", err)
		}

	},
}

func submit(url, shortCode string) error {
	fmt.Println("provided ", shortCode, url)
	Alias = viper.GetString("alias")
	Channel = viper.GetString("channel")
	Event = viper.GetString("event")
	if Alias == "" {
		fmt.Println("Alias is required.  Set with -a or in config file.")
		return errors.New("Alias not provided.")
	}

	if Event == "" {
		fmt.Println("Event is required.  Set with -e or in config file.")
		return errors.New("Event not provided.")
	}
	if Channel == "" {
		fmt.Println("Channel is required.  Set with -c or in config file.")
		return errors.New("Channel not provided.")
	}
	reqURL := build(url)
	// submit to server
	fmt.Println("Submitting to ", baseURL)
	jsonValue, err := json.Marshal(Submission{ShortCode: shortCode, URL: reqURL})
	if err != nil {
		return errors.Wrap(err, "creating JSON")
	}
	resp, err := http.Post(baseURL+"/save", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {

		return errors.Wrap(err, "posting to server")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return errors.Wrap(err, "Unmarshaling result")
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	fmt.Println(result.URL)
	err = clipboard.WriteAll(result.URL)
	if err != nil {
		fmt.Println("Unable to insert into clipboard: ", err)
	}
	return nil
}

func build(url string) string {
	return fmt.Sprintf("%s%s=%s-%s-%s", url, track, Event, Channel, Alias)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	if os.Getenv("BASE_URL") == "" {
		baseURL = "bjk.fyi"
	} else {
		baseURL = os.Getenv("BASE_URL")
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bjk.yaml)")
	RootCmd.PersistentFlags().StringVar(&baseURL, "server", "http://bjk.fyi", "URL Shortening server")
	RootCmd.PersistentFlags().StringVarP(&Alias, "alias", "a", "", "Alias")
	viper.BindPFlag("alias", RootCmd.PersistentFlags().Lookup("alias"))
	RootCmd.PersistentFlags().StringVarP(&Event, "event", "e", "", "event")
	viper.BindPFlag("event", RootCmd.PersistentFlags().Lookup("event"))
	RootCmd.PersistentFlags().StringVarP(&Channel, "channel", "c", "", "channel")
	viper.BindPFlag("channel", RootCmd.PersistentFlags().Lookup("channel"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".bjk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bjk")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}