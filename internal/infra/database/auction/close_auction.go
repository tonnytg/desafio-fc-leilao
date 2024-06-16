package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"time"
)

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			filter := bson.M{"end_time": bson.M{"$lt": now}, "status": auction_entity.Open}
			update := bson.M{"$set": bson.M{"status": auction_entity.Closed}}

			result, err := ar.Collection.UpdateMany(ctx, filter, update)
			if err != nil {
				logger.Error("Error closing expired auctions", err)
			} else {
				logger.Info(
					"Closed expired auctions",
					zap.Int64("count", result.ModifiedCount),
				)
			}
		}
	}
}
