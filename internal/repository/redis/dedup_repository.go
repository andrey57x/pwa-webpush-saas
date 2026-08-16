package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DedupRepository struct {
	client redis.Cmdable // Поддерживает redis.Client, redis.ClusterClient и моки!
}

func NewDedupRepository(client redis.Cmdable) *DedupRepository {
	return &DedupRepository{client: client}
}

func (r *DedupRepository) FilterDuplicates(ctx context.Context, campaignID uuid.UUID, subIDs []uuid.UUID, ttl time.Duration) ([]uuid.UUID, error) {
	if len(subIDs) == 0 {
		return nil, nil
	}

	pipe := r.client.Pipeline()
	cmds := make([]*redis.BoolCmd, len(subIDs))

	for i, subID := range subIDs {
		key := fmt.Sprintf("push_dedup:%s:%s", campaignID.String(), subID.String())
		cmds[i] = pipe.SetNX(ctx, key, 1, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute redis dedup pipeline: %w", err)
	}

	var uniqueSubIDs []uuid.UUID
	for i, cmd := range cmds {
		isNew, err := cmd.Result()
		if err == nil && isNew {
			uniqueSubIDs = append(uniqueSubIDs, subIDs[i])
		}
	}

	return uniqueSubIDs, nil
}
