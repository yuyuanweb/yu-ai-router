package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v84"
	"github.com/yupi/airouter/go-backend/internal/service"
)

type StripeWebhookController struct {
	stripePaymentService *service.StripePaymentService
}

func NewStripeWebhookController(stripePaymentService *service.StripePaymentService) *StripeWebhookController {
	return &StripeWebhookController{
		stripePaymentService: stripePaymentService,
	}
}

func (s *StripeWebhookController) HandleStripeWebhook(c *gin.Context) {
	payloadBytes, err := c.GetRawData()
	if err != nil {
		log.Printf("read stripe webhook payload failed: err=%v", err)
		c.String(200, "error")
		return
	}
	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := s.stripePaymentService.ConstructWebhookEvent(string(payloadBytes), sigHeader)
	if err != nil {
		log.Printf("verify stripe webhook signature failed: err=%v", err)
		c.String(200, "error")
		return
	}
	switch event.Type {
	case "checkout.session.completed":
		sess, parseErr := service.ParseCheckoutSessionFromEvent(event)
		if parseErr != nil {
			log.Printf("parse stripe checkout session failed: err=%v", parseErr)
			c.String(200, "error")
			return
		}
		if _, handleErr := s.stripePaymentService.HandlePaymentSuccess(sess.ID); handleErr != nil {
			log.Printf("handle stripe payment success failed: sessionId=%s err=%v", sess.ID, handleErr)
			c.String(200, "error")
			return
		}
	case "checkout.session.expired", "charge.refunded":
		// 与 Java 保持一致：记录日志，不做额外业务处理
		log.Printf("stripe event received: type=%s", event.Type)
	default:
		log.Printf("stripe event ignored: type=%s", event.Type)
	}
	c.String(200, "success")
}

func IsStripeEvent(event *stripe.Event, eventType string) bool {
	return event != nil && event.Type == stripe.EventType(eventType)
}
