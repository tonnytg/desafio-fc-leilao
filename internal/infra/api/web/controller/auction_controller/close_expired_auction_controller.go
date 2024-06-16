package auction_controller

import (
	"context"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func (u *AuctionController) CloseExpiredAuctions(ctx context.Context) {
	log.Println("AuctionController - CloseExpiredAuctions")

	intervalCheckAuctionExpire, _ := time.ParseDuration(os.Getenv("AUCTION_INTERVAL"))
	log.Println("intervalCheckAuctionExpire:", intervalCheckAuctionExpire)
	ticker := time.NewTicker(intervalCheckAuctionExpire * time.Second)
	//defer ticker.Stop()

	log.Println("Start Loop")

	for {
		select {
		case <-ticker.C:

			var status auction_usecase.AuctionStatus = 0

			log.Println("getting values")
			values, err := u.auctionUseCase.FindAuctions(ctx, status, "", "")
			if err != nil {
				u.logger.Error("Error retrieving auctions", zap.Error(err))
				continue
			}

			auctionExpireTimeout := os.Getenv("AUCTION_EXPIRE_TIMEOUT")
			log.Println("auctionExpireTimeout:", auctionExpireTimeout)
			timeToLive, errParse := time.ParseDuration(auctionExpireTimeout)
			if errParse != nil {
				timeToLive = 60 * 5
			}

			log.Println("checking values")
			for i, v := range values {
				log.Printf("[%d] Auction: %v - Status: %v", i, v.Id, v.Status)

				// Check if Auction has expired
				timeNow := time.Now()
				duration := timeNow.Sub(v.Timestamp)

				if duration.Seconds() > float64(timeToLive) {
					log.Println("Expired:", v.Id)

					// Update Auction to new status 0 -> Active to 3 -> Closed
					err := u.auctionUseCase.CloseExpiredAuctionById(ctx, v.Id)
					if err != nil {
						u.logger.Error("Error closing expired auction", zap.Error(err))
					} else {
						u.logger.Info("Closed expired auction", zap.String("Auction ID", v.Id))
					}
				}
			}
		case <-ctx.Done():
			u.logger.Info("Stopping CloseExpiredAuctions routine")
			return

		default:
			log.Println("Default")
		}
	}
}
