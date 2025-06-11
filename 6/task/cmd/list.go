package cmd

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"log"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows all your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("tasks"))
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("%s. %s\n", k, v)
			}

			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
