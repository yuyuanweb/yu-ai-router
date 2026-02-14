package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
	"github.com/stripe/stripe-go/v84/webhook"
	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/errno"
)

type StripePaymentService struct {
	cfg            *config.Config
	rechargeService *RechargeService
}

func NewStripePaymentService(cfg *config.Config, rechargeService *RechargeService) *StripePaymentService {
	if cfg.StripeAPIKey != "" {
		stripe.Key = cfg.StripeAPIKey
	}
	return &StripePaymentService{
		cfg:             cfg,
		rechargeService: rechargeService,
	}
}

func (s *StripePaymentService) CreateCheckoutSession(userID int64, amount float64, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	if s.cfg.StripeAPIKey == "" {
		return nil, errno.NewWithMessage(errno.SystemError, "STRIPE_API_KEY 未配置")
	}
	record, err := s.rechargeService.CreateRechargeRecord(userID, amount, "stripe")
	if err != nil {
		return nil, err
	}
	amountInCents := int64(amount * 100)
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyCNY)),
					UnitAmount: stripe.Int64(amountInCents),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Yu AI Router 账户充值"),
						Description: stripe.String(fmt.Sprintf("充值金额：¥%.4f", amount)),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.AddMetadata("userId", strconv.FormatInt(userID, 10))
	params.AddMetadata("recordId", strconv.FormatInt(record.ID, 10))
	params.AddMetadata("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	sess, err := session.New(params)
	if err != nil {
		log.Printf("create stripe checkout session failed: userId=%d amount=%.4f err=%v", userID, amount, err)
		return nil, errno.NewWithMessage(errno.OperationError, "创建支付会话失败")
	}
	return sess, nil
}

func (s *StripePaymentService) HandlePaymentSuccess(sessionID string) (bool, error) {
	if s.cfg.StripeAPIKey == "" {
		return false, errno.NewWithMessage(errno.SystemError, "STRIPE_API_KEY 未配置")
	}
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		log.Printf("retrieve stripe session failed: sessionId=%s err=%v", sessionID, err)
		return false, errno.NewWithMessage(errno.OperationError, "获取支付会话失败")
	}
	if sess.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return false, nil
	}
	recordIDRaw, ok := sess.Metadata["recordId"]
	if !ok || recordIDRaw == "" {
		return false, errno.NewWithMessage(errno.OperationError, "支付会话缺少 recordId")
	}
	recordID, parseErr := strconv.ParseInt(recordIDRaw, 10, 64)
	if parseErr != nil || recordID <= 0 {
		return false, errno.New(errno.ParamsError)
	}
	return s.rechargeService.CompleteRecharge(recordID, sessionID)
}

func (s *StripePaymentService) ConstructWebhookEvent(payload, sigHeader string) (*stripe.Event, error) {
	if s.cfg.StripeWebhookSecret == "" {
		return nil, errno.NewWithMessage(errno.SystemError, "STRIPE_WEBHOOK_SECRET 未配置")
	}
	event, err := webhook.ConstructEvent([]byte(payload), sigHeader, s.cfg.StripeWebhookSecret)
	if err != nil {
		return nil, errno.NewWithMessage(errno.ForbiddenError, "Webhook 签名验证失败")
	}
	return &event, nil
}

func ParseCheckoutSessionFromEvent(event *stripe.Event) (*stripe.CheckoutSession, error) {
	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}
