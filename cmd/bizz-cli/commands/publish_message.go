package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	datamodels "github.com/subhammahanty235/bizzgenie/internal/dataModels"
)

// required parameters , queuename , message. extra parameters: priority ,
func (r *RedisClients) PublishJob(ctx context.Context, queueName string, args []string, message string) (bool, error) {
	externalInstance := r.externalInstance.GetExternalRedisClient()
	queueMetaKey := fmt.Sprintf("queue_meta:%s", queueName)
	queueKey := fmt.Sprintf("queue:%s", queueName)

	exists, err := externalInstance.Exists(ctx, queueMetaKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if queue exists: %w", err)
	}
	if exists == 0 {
		return false, fmt.Errorf("❌ Queue \"%s\" does not exist. Create it first", queueName)
	}

	messageId := fmt.Sprintf("message:%d", time.Now().UnixMilli())

	var setPriority int64 = 3
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-p":
			priority, _ := strconv.ParseInt(args[i+1], 10, 64)
			setPriority = priority
		}
	}

	timestamp := time.Now().UnixMilli()
	score := float64(setPriority)*1e10 + float64(timestamp)
	messageObj := datamodels.NewMessage(queueName, messageId, message, datamodels.MessageOptions{Priority: setPriority})
	messageJSON, _ := json.Marshal(messageObj.ToJSON())

	if err := externalInstance.ZAdd(ctx, queueKey, &redis.Z{
		Score:  score,
		Member: messageJSON,
	}).Err(); err != nil {
		return false, fmt.Errorf("failed to publish notification: %w", err)
	}

	return true, nil

}
