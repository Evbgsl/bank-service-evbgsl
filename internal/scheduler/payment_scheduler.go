package scheduler

import (
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/services"
	"github.com/sirupsen/logrus"
)

func StartPaymentScheduler(
	creditService *services.CreditService,
	intervalHours int,
	log *logrus.Logger,
) {
	if intervalHours <= 0 {
		intervalHours = 12
	}

	interval := time.Duration(intervalHours) * time.Hour

	log.Infof("payment scheduler started with interval %s", interval)

	go func() {
		time.Sleep(5 * time.Second)
		runPaymentProcessing(creditService, log)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			runPaymentProcessing(creditService, log)
		}
	}()
}

func runPaymentProcessing(
	creditService *services.CreditService,
	log *logrus.Logger,
) {
	log.Info("payment scheduler: processing due credit payments")

	result, err := creditService.ProcessDuePayments()
	if err != nil {
		log.Errorf("payment scheduler failed: %v", err)
		return
	}

	log.Infof(
		"payment scheduler finished: processed=%d paid=%d overdue=%d skipped=%d",
		result.Processed,
		result.Paid,
		result.Overdue,
		result.Skipped,
	)
}
