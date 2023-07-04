package alidns

import "github.com/spf13/cobra"

func alidnsUpdateCommand() *cobra.Command {
	alidnsUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "",
		Run:   runAlidnsUpdate,
	}

	return alidnsUpdateCmd
}

func runAlidnsUpdate(cmd *cobra.Command, args []string) {

}
