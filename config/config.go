package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
	"strconv"
)

type Configuration struct {
	Viper *viper.Viper
}

var conf *Configuration
var rootCmd *cobra.Command

func Get() *Configuration {
	return conf
}

func Setup(cmd *cobra.Command) (*Configuration, error) {
	rootCmd = cmd
	v := viper.New()
	configuration := &Configuration{
		Viper: v,
	}
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configuration.GetAbsConfigDir())
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	err := v.Unmarshal(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	v.SetDefault("database.name", "data.json")
	v.SetDefault("server.log.store", false)
	v.SetDefault("gin.log.store", false)

	conf = configuration
	return configuration, nil
}

func (conf *Configuration) Prod() bool {
	return rootCmd.PersistentFlags().Lookup("prod").Value.String() == "true"
}

func (conf *Configuration) Auth() bool {
	return rootCmd.PersistentFlags().Lookup("auth").Value.String() == "true"
}

func (conf *Configuration) GetPort() int {
	s := rootCmd.PersistentFlags().Lookup("port").Value.String()
	i, _ := strconv.Atoi(s)
	return i
}

func (conf *Configuration) GetGRPCPort() int {
	s := rootCmd.PersistentFlags().Lookup("grpc").Value.String()
	i, _ := strconv.Atoi(s)
	return i
}

func (conf *Configuration) GetGlobalID() string {
	return rootCmd.PersistentFlags().Lookup("id").Value.String()
}

func (conf *Configuration) GetArgOne() string {
	return rootCmd.PersistentFlags().Lookup("arg1").Value.String()
}

func (conf *Configuration) GetBrokerIP() string {
	return rootCmd.PersistentFlags().Lookup("broker").Value.String()
}

func (conf *Configuration) GetAbsPluginsDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.getPluginsDir())
}

func (conf *Configuration) GetAbsPluginsDataDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.getPluginsDataDir())
}

func (conf *Configuration) GetAbsDataDir() string {
	return path.Join(conf.getGlobalDir(), conf.getDataDir())
}

func (conf *Configuration) GetAbsConfigDir() string {
	return path.Join(conf.getGlobalDir(), conf.getConfigDir())
}

func (conf *Configuration) GetRootDir() string {
	return rootCmd.PersistentFlags().Lookup("root-dir").Value.String()
}

func (conf *Configuration) getGlobalDir() string {
	rootDir := rootCmd.PersistentFlags().Lookup("root-dir").Value.String()
	appDir := rootCmd.PersistentFlags().Lookup("app-dir").Value.String()
	return path.Join(rootDir, appDir)
}

func (conf *Configuration) getDataDir() string {
	return rootCmd.PersistentFlags().Lookup("data-dir").Value.String()
}

func (conf *Configuration) getConfigDir() string {
	return rootCmd.PersistentFlags().Lookup("config-dir").Value.String()
}

func (conf *Configuration) GetAbsDatabaseFile() string {
	return path.Join(conf.GetAbsDataDir(), viper.GetString("database.name"))
}

func (conf *Configuration) getPluginsDataDir() string {
	return rootCmd.PersistentFlags().Lookup("plugins-data-dir").Value.String()
}

func (conf *Configuration) getPluginsDir() string {
	return rootCmd.PersistentFlags().Lookup("plugins-dir").Value.String()
}
