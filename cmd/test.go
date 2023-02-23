package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/galexrt/arpanet/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
	RunE: func(cmd *cobra.Command, args []string) error {
		u := query.Citizen
		user, err := u.Where(u.Firstname.Like("Prof. Dr.%"), u.Lastname.Like("Scott")).First()

		out, _ := json.Marshal(user)
		fmt.Printf("USER: %s\n", string(out))

		return err
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return query.SetupDB(logger)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
