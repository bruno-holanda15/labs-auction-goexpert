package auction_usecase

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"os"
	"strconv"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition" binding:"oneof=0 1 2"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

func NewAuctionUseCase(
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface,
	bidRepositoryInterface bid_entity.BidEntityRepository) AuctionUseCaseInterface {

	interval := getCompleteBatchAuctionInterval()
	batchSize := getMaxBatchSize()

	auctionusecase := &AuctionUseCase{
		auctionRepositoryInterface: auctionRepositoryInterface,
		bidRepositoryInterface:     bidRepositoryInterface,
		completeTimeInterval:       interval,
		timer:                      time.NewTimer(interval),
		maxAuctionBatchSize: batchSize,
		auctionChannel:             make(chan auction_entity.Auction),
	}

	auctionusecase.triggerCompleteAuctionsRoutine(context.Background())
	return auctionusecase
}

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionInput AuctionInputDTO) *internal_error.InternalError

	FindAuctionById(
		ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)

	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository
	completeTimeInterval       time.Duration
	timer                      *time.Timer
	maxAuctionBatchSize        int
	auctionChannel             chan auction_entity.Auction
}

var auctionBatch []auction_entity.Auction

const maxAuctionBatchSize = 5

func (au *AuctionUseCase) triggerCompleteAuctionsRoutine(ctx context.Context) {
	go func() {
		defer close(au.auctionChannel)

		for {
			select {
			case auctionEntity, ok := <-au.auctionChannel:
				if !ok {
					if len(auctionBatch) > 0 {
						if err := au.auctionRepositoryInterface.CompleteAuctions(ctx, auctionBatch); err != nil {
							logger.Error("error trying to complete auctions batch", err)
						}
					}
					return
				}

				auctionBatch = append(auctionBatch, auctionEntity)

				if len(auctionBatch) >= maxAuctionBatchSize {
					if err := au.auctionRepositoryInterface.CompleteAuctions(ctx, auctionBatch); err != nil {
						logger.Error("error trying to complete auctions batch", err)
					}

					auctionBatch = nil
					au.timer.Reset(au.completeTimeInterval)
				}
			case <-au.timer.C:
				if err := au.auctionRepositoryInterface.CompleteAuctions(ctx, auctionBatch); err != nil {
					logger.Error("error trying to complete auctions batch", err)
				}
				auctionBatch = nil
				au.timer.Reset(au.completeTimeInterval)
			}

		}
	}()
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(
		ctx, auction); err != nil {
		return err
	}

	au.auctionChannel <- *auction
	return nil
}

func getCompleteBatchAuctionInterval() time.Duration {
	interval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	batchSize, err := strconv.Atoi(os.Getenv("COMPLETE_AUCTION_BATCH_SIZE"))
	if err != nil {
		return 10
	}

	return batchSize
}
