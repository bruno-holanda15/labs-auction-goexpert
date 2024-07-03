package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/mocks"
	"fullcycle-auction_go/internal/internal_error"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuctionUseCase_CreateAuction(t *testing.T) {
	errT := internal_error.NewInternalServerError("teste")
	ea := map[string]*internal_error.InternalError{
		"teste": errT,
	}
	au_repo := mocks.NewMockAuctionRepository(ea)
	bid_repo := mocks.NewMockBidRepository(ea)

	usecase := NewAuctionUseCase(au_repo, bid_repo)

	dto := AuctionInputDTO{
		ProductName: "PS5",
		Category:    "tech",
		Description: "Testando fluxo de create do usecase",
		Condition:   ProductCondition(auction_entity.Refurbished),
	}

	output := usecase.CreateAuction(context.Background(), dto)
	assert.Equal(t, output, (*internal_error.InternalError)(nil))

	ea = map[string]*internal_error.InternalError{
		"create_auction": errT,
	}
	au_repo = mocks.NewMockAuctionRepository(ea)
	usecase = NewAuctionUseCase(au_repo, bid_repo)
	
	output = usecase.CreateAuction(context.Background(), dto)
	assert.Equal(t, output, errT)
}
