package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	ar.triggerCreateRoutineCheckExpire(context.Background(), auctionEntity)

	return nil
}

func (ar *AuctionRepository) triggerCreateRoutineCheckExpire(ctx context.Context, auctionEntity *auction_entity.Auction) {
	go func() {

		auctionInterval := os.Getenv("AUCTION_INTERVAL")
		auctionTimeLoopInterval, err := time.ParseDuration(auctionInterval)
		if err != nil {
			auctionTimeLoopInterval = time.Minute * 1
		}

		auctionExpire := os.Getenv("AUCTION_EXPIRE")
		auctionLimitTimeToExpire, err := time.ParseDuration(auctionExpire)
		if err != nil {
			auctionLimitTimeToExpire = time.Minute * 3
		}

		for {
			logger.Info("Checking: %s", zap.String("Auction ID", auctionEntity.Id))
			log.Println("Auction is active:", auctionEntity.Id, auctionEntity.ProductName)

			durationAuctionCreated := time.Since(auctionEntity.Timestamp)
			if durationAuctionCreated > auctionLimitTimeToExpire {
				logger.Info("Auction expired: %s", zap.String("Auction ID", auctionEntity.Id))
				log.Printf("Auction is expired: %s - changing to completed", auctionEntity.Id)

				filter := bson.M{"_id": auctionEntity.Id, "status": auction_entity.Active}
				update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

				_, err := ar.Collection.UpdateOne(ctx, filter, update)
				if err != nil {
					logger.Error("Error trying to change status auction from Active to Complete", err)
					log.Println("Error trying to change status auction from Active to Complete", err)
				}

				break
			}

			time.Sleep(auctionTimeLoopInterval)
		}
	}()
}
