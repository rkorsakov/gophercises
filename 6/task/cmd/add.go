package cmd

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new task to your list",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tasks.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		var task strings.Builder
		for _, val := range args {
			task.WriteString(val + " ")
		}
		taskStr := strings.TrimSpace(task.String())
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("tasks"))
			if err != nil {
				return err
			}
			id, _ := b.NextSequence()
			strId := strconv.FormatUint(id, 10)
			err = b.Put([]byte(strId), []byte(taskStr))
			return nil
		})
		fmt.Println("Added task:", taskStr)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
