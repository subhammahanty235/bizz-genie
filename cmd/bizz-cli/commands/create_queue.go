package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/subhammahanty235/bizzgenie/internal/redisclient"
)

type RedisClients struct {
	externalInstance *redisclient.ExternalRedisClient
	// internalInstance *redisclient.InternalRedisClient
}

// flags --> --setdlq, --maxretries
func (r *RedisClients) CreateQueue(ctx context.Context, args []string) (bool, error) {
	if len(args) == 0 {
		return false, fmt.Errorf("queue name is required")
	}

	externalRedisClient := r.externalInstance.GetExternalRedisClient()
	queueName := args[0]
	setDlq := false
	maxRetries := 3

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--setdlq":
			if i+1 < len(args) {
				boolVal, err := strconv.ParseBool(args[i+1])
				if err != nil {
					return false, fmt.Errorf("invalid value for --setdlq: %s", args[i+1])
				}
				setDlq = boolVal
				i++ // skip next argument
			} else {
				return false, fmt.Errorf("--setdlq requires a value")
			}
		case "--maxretries":
			if i+1 < len(args) {
				intVal, err := strconv.ParseInt(args[i+1], 10, 64)
				if err != nil {
					return false, fmt.Errorf("invalid value for --maxretries: %s", args[i+1])
				}
				maxRetries = int(intVal)
				i++ // skip next argument
			} else {
				return false, fmt.Errorf("--maxretries requires a value")
			}
		default:
			fmt.Printf("Unable to process the paramenter %s \n", args[i])
		}
	}

	queueName = strings.ToLower(strings.TrimSpace(queueName))
	queueMetaKey := fmt.Sprintf("queue_meta:%s", queueName)

	//checking if the queue already exists, if exists, we will throw error
	exists, err := externalRedisClient.Exists(ctx, queueMetaKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if queue exists: %w", err)
	}
	if exists == 1 {
		return false, fmt.Errorf("✅ Queue \"%s\" already exists", queueName)

	}

	//set all the queue data
	queueData := map[string]interface{}{
		"createdAt":                time.Now().UnixMilli(),
		"config_dead_letter_queue": setDlq,
		"maxRetries":               maxRetries,
	}

	//create the queue
	if err := externalRedisClient.HSet(ctx, queueMetaKey, queueData).Err(); err != nil {
		return false, fmt.Errorf("failed to create queue: %w", err)
	}

	fmt.Printf("📌 Queue \"%s\" created successfully.\n", queueName)

	//configure the dead letter queue
	if setDlq {
		fmt.Println("Configuring the Dead letter queue for .\n", queueName, ".......")
		dlqName := fmt.Sprintf("%s_dlq", queueName)
		dlqMetaKey := fmt.Sprintf("queue_meta:%s", dlqName)
		exists, err := externalRedisClient.Exists(ctx, dlqMetaKey).Result()
		if err != nil {
			return false, fmt.Errorf("failed to check if dlq queue exists: %w", err)
		}
		if exists == 1 {
			return false, fmt.Errorf("dead Letter Queue \"%s\" already exists", queueName)
		}

		optionsMap := make(map[string]interface{})

		optionsMap["createdAt"] = time.Now().UnixMilli()
		optionsMap["dead_letter_queue"] = true
		optionsMap["parentQueue"] = queueName

		if maxRetries > 0 {
			optionsMap["maxRetries"] = maxRetries
		} else {
			optionsMap["maxRetries"] = 0
		}

		if err := externalRedisClient.HSet(ctx, dlqMetaKey, optionsMap).Err(); err != nil {
			return false, fmt.Errorf("failed to Create DLQ")
		}

		fmt.Println("✅ Configured the Dead letter queue for .\n", queueName, ".......")

	}

	return true, nil
}
