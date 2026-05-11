package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/evbgsl/bank-service-evbgsl/internal/models"
)

const cbrDailyInfoURL = "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"

var ErrKeyRateNotFound = errors.New("key rate not found")

type CBRService struct {
	bankMargin float64
	httpClient *http.Client
}

func NewCBRService(bankMargin float64) *CBRService {
	return &CBRService{
		bankMargin: bankMargin,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *CBRService) GetKeyRate() (*models.KeyRateResponse, error) {
	soapRequest := buildKeyRateSOAPRequest()

	rawBody, err := s.sendSOAPRequest(soapRequest)
	if err != nil {
		return nil, err
	}

	rate, err := parseKeyRateXML(rawBody)
	if err != nil {
		return nil, err
	}

	rate = roundIntegrationMoney(rate)
	margin := roundIntegrationMoney(s.bankMargin)
	finalRate := roundIntegrationMoney(rate + margin)

	return &models.KeyRateResponse{
		CentralBankRate: rate,
		BankMargin:      margin,
		FinalRate:       finalRate,
		Message:         "key rate received from Central Bank of Russia SOAP service",
	}, nil
}

func buildKeyRateSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                 xmlns:xsd="http://www.w3.org/2001/XMLSchema"
                 xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
  <soap12:Body>
    <KeyRate xmlns="http://web.cbr.ru/">
      <fromDate>%s</fromDate>
      <ToDate>%s</ToDate>
    </KeyRate>
  </soap12:Body>
</soap12:Envelope>`, fromDate, toDate)
}

func (s *CBRService) sendSOAPRequest(soapRequest string) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		cbrDailyInfoURL,
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cbr request failed: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cbr response read failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cbr returned status %d: %s", resp.StatusCode, string(rawBody))
	}

	return rawBody, nil
}

func parseKeyRateXML(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, fmt.Errorf("cbr xml parse failed: %w", err)
	}

	rateElements := doc.FindElements("//diffgram/KeyRate/KR/Rate")
	if len(rateElements) == 0 {
		rateElements = doc.FindElements("//KeyRate/KR/Rate")
	}

	if len(rateElements) == 0 {
		return 0, ErrKeyRateNotFound
	}

	rateText := strings.TrimSpace(rateElements[0].Text())
	rateText = strings.ReplaceAll(rateText, ",", ".")

	var rate float64
	if _, err := fmt.Sscanf(rateText, "%f", &rate); err != nil {
		return 0, fmt.Errorf("cbr key rate convert failed: %w", err)
	}

	return rate, nil
}

func roundIntegrationMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
