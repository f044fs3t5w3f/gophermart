package accrual

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/f044fs3t5w3f/gophermart/internal/accrual/client"
	"github.com/f044fs3t5w3f/gophermart/internal/logger"
	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"go.uber.org/zap"
)

var ErrServiceShuttedDown = errors.New("service is shutted down")

type repo interface {
	// transactionable
	ListOrdersForUpdate(ctx context.Context) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status models.OrderStatus) error
	UpdateOrderStatusAndAccrual(ctx context.Context, orderID int64, status models.OrderStatus, accrual float64) error
}

type accrualClient interface {
	GetInfo(ctx context.Context, orderID string) (*client.AccrualServiceResponse, error)
}

func NewAccuralService(ctx context.Context, repo repo, log *zap.Logger, accrualServiceURL string) *AccuralService {
	client := client.NewAccrualClient(accrualServiceURL)

	log.Info("Starting accural servuce ", zap.String("url", accrualServiceURL))
	accuralService := &AccuralService{
		ctx:        ctx,
		repository: repo,
		client:     client,
		log:        log,
		wg:         &sync.WaitGroup{},
	}
	go func() {
		<-ctx.Done()
		accuralService.shuttedDown.Store(true)
	}()
	return accuralService
}

type AccuralService struct {
	ctx         context.Context
	repository  repo
	client      accrualClient
	log         *zap.Logger
	wg          *sync.WaitGroup
	shuttedDown atomic.Bool
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

func (p *AccuralService) AddToFetchList(order *models.Order) error {
	if p.shuttedDown.Load() {
		return ErrServiceShuttedDown
	}
	p.log.Info("accural service: start to proccess", zap.String("order", order.Number))
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		firstAttempt := true
		sleep := 1 * time.Second
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
			}
			if firstAttempt {
				firstAttempt = false
			} else {
				timer := time.NewTimer(sleep * time.Second)
				select {
				case <-timer.C:
				case <-p.ctx.Done():
					return
				}
				sleep = min(sleep*2, 300*time.Second)
			}
			accrual, err := p.client.GetInfo(p.ctx, order.Number)
			if err != nil {
				p.log.Error("failed to get accrual info", zap.Error(err))
				continue
			}
			p.log.Info("Got accural service response", zap.String("accural", fmt.Sprintf("%+v\n", accrual)))
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
				err := p.repository.UpdateOrderStatus(p.ctx, order.ID, models.OrderStatusProcessing)
				if err != nil {
					p.log.Error("failed to update order status", zap.String("orderId", order.Number))
				}
			case client.StatusProcessed:
				status := models.OrderStatusProcessed

				p.log.Info("To update order status", zap.String("accural", fmt.Sprintf("%+v\n", accrual)))
				err := p.repository.UpdateOrderStatusAndAccrual(p.ctx, order.ID, status, accrual.Accrual)
				if err != nil {
					p.log.Error("failed to update order status and accrual", zap.String("orderId", order.Number))
				} else {
					return
				}

			case client.StatusInvalid:
				status := models.OrderStatusInvalid
				err := p.repository.UpdateOrderStatus(p.ctx, order.ID, status)
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
	return nil
}

func (p *AccuralService) Wait() {
	p.log.Info("accuralService: wait for current responses...")
	p.wg.Wait()
	p.log.Info("accuralService: closed")
}
