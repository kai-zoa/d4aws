package main

import (
	"github.com/kai-zoa/d4aws/service/ecr"
	"github.com/kai-zoa/d4aws/service/leader"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "d4aws",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Usage()
			if err != nil {
				panic(err)
			}
		},
	}
	cobra.OnInitialize()
	root.AddCommand(ecr.Command)
	root.AddCommand(leader.Command)
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
