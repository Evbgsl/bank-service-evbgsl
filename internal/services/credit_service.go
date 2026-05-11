package services

import (
	"errors"
	"math"
	"time"

	"github.com/evbgsl/bank-service-evbgsl/internal/models"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
)

var (
	ErrInvalidCreditData = errors.New("invalid credit data")
	ErrCreditNotFound    = errors.New("credit not found")
)

type CreditService struct {
	creditRepo   *repositories.CreditRepository
	accountRepo  *repositories.AccountRepository
	emailService *EmailService
}

func NewCreditService(
	creditRepo *repositories.CreditRepository,
	accountRepo *repositories.AccountRepository,
	emailService *EmailService,
) *CreditService {
	return &CreditService{
		creditRepo:   creditRepo,
		accountRepo:  accountRepo,
		emailService: emailService,
	}
}

func (s *CreditService) CreateCredit(
	userID int64,
	req models.CreateCreditRequest,
) (*models.CreateCreditResponse, error) {
	if req.AccountID <= 0 {
		return nil, ErrInvalidCreditData
	}

	if req.PrincipalAmount <= 0 {
		return nil, ErrInvalidCreditData
	}

	if req.InterestRate < 0 {
		return nil, ErrInvalidCreditData
	}

	if req.TermMonths <= 0 || req.TermMonths > 360 {
		return nil, ErrInvalidCreditData
	}

	account, err := s.accountRepo.FindByIDAndUserID(req.AccountID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	monthlyPayment := calculateAnnuityPayment(
		req.PrincipalAmount,
		req.InterestRate,
		req.TermMonths,
	)

	schedule := generatePaymentSchedule(
		req.PrincipalAmount,
		req.InterestRate,
		req.TermMonths,
		monthlyPayment,
	)

	credit := &models.Credit{
		UserID:          userID,
		AccountID:       account.ID,
		PrincipalAmount: roundMoney(req.PrincipalAmount),
		InterestRate:    roundMoney(req.InterestRate),
		TermMonths:      req.TermMonths,
		MonthlyPayment:  monthlyPayment,
		RemainingAmount: roundMoney(req.PrincipalAmount),
		Status:          "ACTIVE",
	}

	if err := s.creditRepo.CreateCreditWithSchedule(credit, schedule); err != nil {
		return nil, err
	}

	return &models.CreateCreditResponse{
		ID:              credit.ID,
		AccountID:       credit.AccountID,
		PrincipalAmount: credit.PrincipalAmount,
		InterestRate:    credit.InterestRate,
		TermMonths:      credit.TermMonths,
		MonthlyPayment:  credit.MonthlyPayment,
		RemainingAmount: credit.RemainingAmount,
		Status:          credit.Status,
		Message:         "credit created successfully",
	}, nil
}

func (s *CreditService) GetUserCredits(userID int64) ([]models.Credit, error) {
	return s.creditRepo.FindByUserID(userID)
}

func (s *CreditService) GetCreditSchedule(
	userID int64,
	creditID int64,
) ([]models.PaymentSchedule, error) {
	if creditID <= 0 {
		return nil, ErrInvalidCreditData
	}

	schedule, err := s.creditRepo.FindScheduleByCreditIDAndUserID(creditID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrCreditNotFound) {
			return nil, ErrCreditNotFound
		}

		return nil, err
	}

	return schedule, nil
}

func (s *CreditService) ProcessDuePayments() (*models.PaymentProcessingResult, error) {
	result, err := s.creditRepo.ProcessDuePayments()
	if err != nil {
		return nil, err
	}

	if s.emailService == nil {
		return result, nil
	}

	for _, notification := range result.Notifications {
		_ = s.emailService.SendCreditPaymentEmail(
			notification.UserEmail,
			notification.Amount,
			notification.Status,
		)
	}

	return result, nil
}

func calculateAnnuityPayment(
	principal float64,
	annualRate float64,
	termMonths int,
) float64 {
	if annualRate == 0 {
		return roundMoney(principal / float64(termMonths))
	}

	monthlyRate := annualRate / 100 / 12

	payment := principal * monthlyRate * math.Pow(1+monthlyRate, float64(termMonths)) /
		(math.Pow(1+monthlyRate, float64(termMonths)) - 1)

	return roundMoney(payment)
}

func generatePaymentSchedule(
	principal float64,
	annualRate float64,
	termMonths int,
	monthlyPayment float64,
) []models.PaymentSchedule {
	schedule := make([]models.PaymentSchedule, 0, termMonths)

	remaining := principal
	monthlyRate := annualRate / 100 / 12

	for i := 1; i <= termMonths; i++ {
		interestPart := roundMoney(remaining * monthlyRate)
		principalPart := roundMoney(monthlyPayment - interestPart)

		if i == termMonths {
			principalPart = roundMoney(remaining)
			monthlyPayment = roundMoney(principalPart + interestPart)
		}

		remaining = roundMoney(remaining - principalPart)

		schedule = append(schedule, models.PaymentSchedule{
			PaymentNumber: i,
			PaymentDate:   time.Now().AddDate(0, i, 0),
			Amount:        monthlyPayment,
			PrincipalPart: principalPart,
			InterestPart:  interestPart,
			Status:        "PLANNED",
		})
	}

	return schedule
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
