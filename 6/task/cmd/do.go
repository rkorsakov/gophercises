package cmd

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"log"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("tasks"))
			v := b.Get([]byte(args[0]))
			if v == nil {
				return fmt.Errorf("task '%s' not found", args[0])
			}
			fmt.Println("You've completed the task: " + string(v))
			return b.Delete([]byte(args[0]))
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
