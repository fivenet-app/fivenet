package cmd

import (
	"github.com/galexrt/arpanet/pkg/config"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
	RunE: func(cm *cobra.Command, args []string) error {
		db, err := sqlx.Open("mysql", config.C.Database.DSN)
		if err != nil {
			return err
		}

		stmt := `SELECT documentreference.id AS "documentreference.id",
                 documentreference.created_at AS "documentreference.created_at",
                 documentreference.source_document_id AS "documentreference.source_document_id",
                 documentreference.reference AS "documentreference.reference",
                 documentreference.target_document_id AS "documentreference.target_document_id",
                 documentreference.creator_id AS "documentreference.creator_id",
                 source_document.id AS "source_document.id",
                 source_document.created_at AS "source_document.created_at",
                 source_document.updated_at AS "source_document.updated_at",
                 source_document.category_id AS "source_document.category_id",
                 source_document.title AS "source_document.title",
                 source_document.creator_id AS "source_document.creator_id",
                 source_document.state AS "source_document.state",
                 source_document.closed AS "source_document.closed",
                 category.id AS "category.id",
                 category.name AS "category.name",
                 category.description AS "category.description",
                 ref_creator.id AS "ref_creator.id",
                 ref_creator.identifier AS "ref_creator.identifier",
                 ref_creator.job AS "ref_creator.job",
                 ref_creator.job_grade AS "ref_creator.job_grade",
                 ref_creator.firstname AS "ref_creator.firstname",
                 ref_creator.lastname AS "ref_creator.lastname"
            FROM arpanet_documents_references AS documentreference
                 LEFT JOIN arpanet_documents AS source_document ON (documentreference.source_document_id = source_document.id)
                 LEFT JOIN arpanet_documents_categories AS category ON (category.id = source_document.category_id)
                 LEFT JOIN users AS ref_creator ON (documentreference.creator_id = ref_creator.id)
            WHERE (
                      (documentreference.target_document_id = 1)
                          AND documentreference.deleted_at IS NULL
                  )
            LIMIT 25;`

		var dest []*documents.DocumentReference

		udb := db.Unsafe()
		err = udb.Get(&dest, stmt)
		return err
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
