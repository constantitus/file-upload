package db

import (
	"fmt"
	"main/cache"
	"main/db"
	"strings"

	"testing"
)

// Prints the contents of the database
func TestDB(t *testing.T) {
    err := db.Initialize()
    if err != nil {
        t.Fatal(err)
    }

    var message strings.Builder
    message.WriteString("\n")
    for _, uuid := range cache.Keys() {
        data, got := cache.Get(uuid)
        if got {
            exp, _ := cache.GetExp(uuid)
            message.WriteString(uuid)
            message.WriteString(" ")
            message.WriteString(data.Name)
            message.WriteString(" Admin:")
            message.WriteString(fmt.Sprint(data.Rank))
            message.WriteString(" ")
            message.WriteString(fmt.Sprintln(exp))
        }
    }
    t.Log(message.String())
}
