package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/subhammahanty235/bizzgenie/internal/utils"
)

// List queues - will show parameters such as, queue name, created, active (based on the heartbeat), total stored messages
// extra parameters - --showdlq , -n , -a
func (r *RedisClients) ListQueues(ctx context.Context, args []string) (bool, error) {
	externalInstance := r.externalInstance.GetExternalRedisClient()
	var showOnlyName bool = false
	var showOnlyActives bool = false
	// var showDLQs bool = false

	// get the extra parameters (if exists)
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--showdlq":
			// showDLQs = true
		case "-n":
			showOnlyName = true
		case "-a":
			showOnlyActives = true
		}
	}

	// fetch all the queues
	var (
		cursor uint64
		prefix = "queue_meta:"
		result = make(map[string]map[string]string)
	)
	var names []string

	for {
		keys, newCursor, err := externalInstance.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return false, fmt.Errorf("error happed while finding the queues")
		}

		for _, key := range keys {
			values, err := externalInstance.HGetAll(ctx, key).Result()
			if err != nil {
				return false, fmt.Errorf("error happed while finding the queues")
			}

			name := strings.TrimPrefix(key, prefix)
			names = append(names, name)
			activeValue := "true"
			values["active"] = activeValue

			if showOnlyActives {
				if activeValue == "true" {
					result[key] = values
				}
			} else {
				result[key] = values
			}

		}

		cursor = newCursor
		if cursor == 0 {
			break
		}

	}

	//process the result, before showing in a tabular form
	//if only name required than we will only show the names from the names array,
	if len(names) == 0 {
		fmt.Println("no queues found, please add some queues")
		return true, nil
	}

	if showOnlyName {
		for _, name := range names {
			fmt.Println(name)
		}
		fmt.Println("All Queues printed......")
		return true, nil
	}

	//print the data in tabular format
	utils.PrintQueueTable(result)
	return true, nil
}

// query : show queue "name" , details we will show, name , created, fileds, active status , all instances
func (r *RedisClients) ShowQueueDetails(ctx context.Context, queueName string) (bool, error) {
	externalInstance := r.externalInstance.GetExternalRedisClient()
	queueMetaKey := fmt.Sprintf("queue_meta:%s", queueName)

	//check if the queue exists
	exists, err := externalInstance.Exists(ctx, queueMetaKey).Result()
	if err != nil {
		return false, fmt.Errorf("error while finding the queue, please try again later")
	}
	if exists == 0 {
		fmt.Println("Queue does not exists. please check the name")
		return false, nil
	}

	//fetch the details.
	queueDeatils, err := externalInstance.HGetAll(ctx, queueMetaKey).Result()
	if err != nil {
		return false, fmt.Errorf("error while finding the queue properties, please try again later")
	}

	queueDataComplete := make(map[string]string)
	queueDataComplete["createdAt"] = queueDeatils["createdAt"]
	if queueDeatils["maxRetries"] != "" {
		queueDataComplete["maxRetreis"] = queueDeatils["maxRetries"]
	}
	if queueDeatils["config_dead_letter_queue"] != "" {
		queueDataComplete["config_dead_letter_queue"] = queueDeatils["config_dead_letter_queue"]
	}

	// will fetch the other parameters from the internal instance with the heartbeat

	return true, nil

}
