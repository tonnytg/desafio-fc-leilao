package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
)

func (au *AuctionUseCase) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	return au.auctionRepositoryInterface.CloseExpiredAuctions(ctx)
}

func (au *AuctionUseCase) CloseExpiredAuctionById(ctx context.Context, id string) *internal_error.InternalError {
	auction, err := au.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return err
	}

	auction.Status = auction_entity.Closed
	return au.auctionRepositoryInterface.UpdateAuction(ctx, auction)
}
