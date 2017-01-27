package leader

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Command provides a PrimaryManager command root.
var Command = &cobra.Command{
	Use:   "leader",
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
	getIPCommand = &cobra.Command{
		Use:   "ip",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				panic(errors.New("ip command takes an argument"))
			}
			pm, err := New(args[0])
			if err != nil {
				panic(err)
			}
			if flagPublic {
				ip, err := pm.GetPublicIPAddress()
				if err != nil {
					panic(err)
				}
				fmt.Println(ip)
			} else {
				ip, err := pm.GetPrivateIPAddress()
				if err != nil {
					panic(err)
				}
				fmt.Println(ip)
			}
		},
	}
	flagPublic = false
)

func init() {
	Command.AddCommand(getIPCommand)
	getIPCommand.PersistentFlags().
		BoolVarP(&flagPublic, "public", "p", false, "")
}
