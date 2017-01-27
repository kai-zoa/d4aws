package ecr

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Command provides a ECR command root
var Command = &cobra.Command{
	Use:   "ecr",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		if err != nil {
			panic(err)
		}
	},
}

var (
	getLoginCommand = &cobra.Command{
		Use:   "get-login",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ecr, err := New()
			if err != nil {
				panic(err)
			}
			input := GetLoginCommandInput{
				RegistryIDs: registryIDs,
			}
			s, err := ecr.GetLoginCommand(&input)
			fmt.Println(s)
		},
	}
	registryIDs = []string{}
)

func init() {
	Command.AddCommand(getLoginCommand)
	getLoginCommand.
		PersistentFlags().
		StringArrayVar(&registryIDs, "registry-ids", registryIDs, "")
}
