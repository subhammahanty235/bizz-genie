package utils

import (
	"fmt"
	"strings"
)

func PrintQueueTable(result map[string]map[string]string) {
	// Print table header
	fmt.Printf("%-20s %-15s %-10s\n", "Queue Name", "Created", "Active")
	fmt.Println(strings.Repeat("-", 50))

	// Print each row
	for fullKey, fields := range result {
		// Extract queue name
		queueName := strings.TrimPrefix(fullKey, "queue_meta:")

		// Get field values
		created := fields["created"]
		active := fields["active"]

		// Print row
		fmt.Printf("%-20s %-15s %-10s\n", queueName, created, active)
	}
}
