package cmd

import (
	"fmt"
	"strings"

	_ "github.com/rs/zerolog"
	_ "github.com/rs/zerolog/log"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		priv, pub := helpers.GenerateRsaKeyPair()
		fmt.Println("Private")
		fmt.Println(strings.ReplaceAll(helpers.ExportPrivateKey(priv), "\n", "\\n"))
		fmt.Println("Public")
		publicStr, _ := helpers.ExportPublicKey(pub)
		fmt.Println(strings.ReplaceAll(publicStr, "\n", "\\n"))
		return
	},
}
