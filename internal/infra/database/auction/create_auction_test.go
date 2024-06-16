package auction_test

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestTriggerCreateRoutineCheckExpire(t *testing.T) {
	ctx := context.Background()

	if err := godotenv.Load("../../../../cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	os.Setenv("MONGODB_URL", "mongodb://admin:admin@localhost:27017/auctions?authSource=admin")

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	ar := auction.NewAuctionRepository(databaseConnection)

	a := auction_entity.Auction{
		Id:          "12345678-a12b-12ab-1234-1a2bc34d56ef",
		ProductName: "test_name_auction",
		Category:    "test_category_auction",
		Description: "a b c d e f g h i j k",
	}

	// Set new value to expire auction to make test more quickly
	os.Setenv("AUCTION_EXPIRE", "2s")

	_ = ar.CreateAuction(ctx, &a)

	values, _ := ar.FindAuctions(ctx, 0, "", "")

	var hasError bool = true

	for _, v := range values {

		log.Println("checking", v.Id, v.ProductName)

		if v.ProductName == a.ProductName {
			log.Printf("expected %s got %s\n", a.ProductName, v.ProductName)
			hasError = false
		}
	}
	if hasError {
		t.Error("error to check auction has created")
	}

	values, _ = ar.FindAuctions(ctx, 0, "", "")

	hasError = true

	for _, v := range values {

		log.Println("checking", v.Id, v.ProductName)

		if v.ProductName == a.ProductName {
			log.Printf("expected %s got %s\n", a.ProductName, v.ProductName)
			log.Printf("expected %v got %v\n", 1, v.Status)
			if v.Status == 1 {
				hasError = false
			}
		}
	}
	if hasError {
		t.Error("error to check auction has correct status")
	}

	return
}
