package accrual

import (
	"context"
	"time"

	"github.com/f044fs3t5w3f/gophermart/internal/accrual/client"
	"github.com/f044fs3t5w3f/gophermart/internal/logger"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"go.uber.org/zap"
)

// type txKeyType struct{}

// var txKey = txKeyType{}

// type transactionable interface {
// 	beginTx() *sql.Tx
// 	commitTx(tx *sql.Tx) error
// }

type repo interface {
	// transactionable
	ListOrdersForUpdate(ctx context.Context) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, status models.OrderStatus) error
	UpdateOrderStatusAndAccrual(ctx context.Context, orderId int64, status models.OrderStatus, accrual float64) error
}

type accrualClient interface {
	GetInfo(ctx context.Context, orderId string) (*client.AccrualServiceResponse, error)
}

func NewAccuralService(ctx context.Context, repo repo, log *zap.Logger, accrualServiceURL string) *AccuralService {
	client := client.NewAccrualClient(accrualServiceURL)
	return &AccuralService{
		ctx:        ctx,
		repository: repo,
		client:     client,
		log:        log,
	}
}

type AccuralService struct {
	ctx        context.Context
	repository repo
	client     accrualClient
	log        *zap.Logger
}

func (p *AccuralService) LoadOld() {
	orders, err := p.repository.ListOrdersForUpdate(p.ctx)
	if err != nil {
		return
	}
	for _, order := range orders {
		p.AddToFetchList(order)
	}
}

func (p *AccuralService) AddToFetchList(order *models.Order) {
	firstAttempt := true
	sleep := 1 * time.Second
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
			}
			if firstAttempt {
				firstAttempt = false
			} else {
				time.Sleep(sleep)
				sleep = min(sleep*2, 1*time.Hour)
			}
			accrual, err := p.client.GetInfo(p.ctx, order.Number)
			if err != nil {
				p.log.Error("failed to get accrual info", zap.Error(err))
				continue
			}

			if accrual.Order != order.Number {
				p.log.Error(
					"Response order number doesn't match request order number",
					zap.String("request order", order.Number),
					zap.String("response order", accrual.Order),
				)
				return
			}

			switch accrual.Status {
			case client.StatusRegistered:
				// do nothing
			case client.StatusProcessing:
				err := p.repository.UpdateOrderStatus(p.ctx, order.Id, models.OrderStatusProcessing)
				if err != nil {
					p.log.Error("failed to update order status", zap.String("orderId", order.Number))
				}
			case client.StatusProcessed:
				status := models.OrderStatusProcessed
				err := p.repository.UpdateOrderStatusAndAccrual(p.ctx, order.Id, status, accrual.Accrual)
				if err != nil {
					p.log.Error("failed to update order status and accrual", zap.String("orderId", order.Number))
				} else {
					return
				}

			case client.StatusInvalid:
				status := models.OrderStatusInvalid
				err := p.repository.UpdateOrderStatus(p.ctx, order.Id, status)
				if err != nil {
					logger.Log.Error("failed to update order status", zap.String("orderId", order.Number))
				} else {
					return
				}
			default:
				p.log.Info("Strange status", zap.String("orderId", order.Number), zap.String("status", accrual.Status))
			}
		}
	}()
}
