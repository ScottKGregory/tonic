package cmd

import (
	"fmt"
	"strings"

	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	"github.com/scottkgregory/mamba"
	"github.com/scottkgregory/tonic"
	"github.com/scottkgregory/tonic/internal/helpers"
	"github.com/scottkgregory/tonic/internal/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ConfigFile string         `config:"./examples/webapi/config.yaml, The yaml config file to read, true, c"`
	Port       int            `config:"8080, The port to host the site on, false, p"`
	CertGen    bool           `config:"false, Generate new JWT certificates, false, g"`
	Tonic      models.Options `config:""`
}

var cfg AppConfig

var rootCmd = &cobra.Command{
	Use:   "webapi",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if cfg.CertGen {
			priv, pub := helpers.GenerateRsaKeyPair()
			fmt.Println("Private")
			fmt.Println(strings.ReplaceAll(helpers.ExportPrivateKey(priv), "\n", "\\n"))
			fmt.Println("Public")
			publicStr, _ := helpers.ExportPublicKey(pub)
			fmt.Println(strings.ReplaceAll(publicStr, "\n", "\\n"))
			return
		}

		cfg.Tonic.Log.IgnoreRoutes = []string{"/health", "/liveliness", "/readiness"}
		r, _, err := tonic.Init(cfg.Tonic)
		if err != nil {
			panic(err)
		}

		logger := tonic.GetLogger()
		logger.Info().Int("port", cfg.Port).Msg("Starting listener")
		err = r.Run(fmt.Sprintf(":%d", cfg.Port))
		if err != nil {
			logger.Fatal().Err(err).Msg("Error starting listener")
		}
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
