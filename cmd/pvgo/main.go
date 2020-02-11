package main

import (
	"github.com/metrumresearchgroup/projectValidator/cmd"
	"github.com/spf13/viper"
)

func main(){
	viper.SetEnvPrefix("pvgo")
	viper.AutomaticEnv()
	cmd.Execute()
}
