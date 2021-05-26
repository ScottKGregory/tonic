package cmd

import (
	"fmt"

	"github.com/scottkgregory/mamba"
	"github.com/scottkgregory/tonic"
	"github.com/scottkgregory/tonic/pkg/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ConfigFile string         `config:"./examples/webapi/config.yaml, The yaml config file to read, true, c"`
	Port       int            `config:"8080, The port to host the site on, false, p"`
	Tonic      models.Options `config:""`
}

var cfg AppConfig

var rootCmd = &cobra.Command{
	Use:   "webapi",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// priv, pub := helpers.GenerateRsaKeyPair()
		// fmt.Println(helpers.ExportPrivateKey(priv))
		// fmt.Println(helpers.ExportPublicKey(pub))

		cfg.Tonic.Log.IgnoreRoutes = []string{"/health", "/liveliness", "/readiness"}
		r, _, err := tonic.Init(cfg.Tonic)
		if err != nil {
			panic(err)
		}

		// logger := helpers.GetLogger()
		// logger.Info().Msg("Starting listener")
		r.Run(fmt.Sprintf(":%d", cfg.Port))
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	mamba.MustBind(&AppConfig{}, rootCmd, &mamba.Options{PrefixEmbedded: false})
}

func initConfig() {
	cfgFile := viper.GetString("configfile")

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Error unmarshalling config: %s", err))
	}
}
