package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/galexrt/arpanet/query"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"gorm.io/gen/field"
)

var testCmd = &cobra.Command{
	Use: "test",
	RunE: func(cmd *cobra.Command, args []string) error {
		u := query.User
		user, err := u.Preload(field.Associations).Where(u.Job.In("ambulance"), u.Identifier.Eq("char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57")).Find()
		if err != nil {
			return err
		}

		out, _ := json.MarshalIndent(user, "", "    ")
		fmt.Printf("USERS:\n%s\n", string(out))

		return err
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return query.SetupDB(logger)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
