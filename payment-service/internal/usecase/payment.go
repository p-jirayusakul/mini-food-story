package usecase

import (
	"context"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/howeyc/crc16"
)

func (i *PaymentImplement) CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, customError *exceptions.CustomError) {

	customError = i.repository.IsOrderExist(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	customError = i.repository.IsOrderItemsNotFinal(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	transactionID, customError = i.repository.CreatePaymentTransaction(ctx, payload)
	if customError != nil {
		return "", customError
	}

	return transactionID, nil
}

func (i *PaymentImplement) GetPaymentLastStatusCodeByTransaction(ctx context.Context, transactionID string) (result string, customError *exceptions.CustomError) {
	return i.repository.GetPaymentLastStatusCodeByTransaction(ctx, transactionID)
}

func (i *PaymentImplement) PaymentTransactionQR(ctx context.Context, transactionID string) (result domain.TransactionQR, customError *exceptions.CustomError) {

	amountOrder, customError := i.repository.GetPaymentAmountByTransaction(ctx, transactionID)
	if customError != nil {
		return domain.TransactionQR{}, customError
	}

	amount := amountOrder
	merchantName := "FOOD STORE"
	merchantCity := "BANGKOK"
	merchantPhone := "0891234567" // หรือ citizenID := "1234567890123"

	// แปลงเป็น string 2 ทศนิยม
	amountStr := fmt.Sprintf("%.2f", amount) // "520.00"

	qrText, err := BuildPromptPayQR(BuildOpts{
		Mobile:   merchantPhone, // หรือ CitizenID: citizenID,
		Amount:   amountStr,
		Ref:      transactionID, // ใส่ ref จะ recon ได้ง่าย
		Merchant: merchantName,
		City:     merchantCity,
		Dynamic:  true,
	})
	if err != nil {
		return domain.TransactionQR{}, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: err,
		}
	}

	expires := time.Now().Add(2 * time.Minute).UTC().Format(time.RFC3339)

	return domain.TransactionQR{
		MethodCode: "PROMPTPAY",
		QrText:     qrText,
		ExpiresAt:  expires,
		Amount:     amount,
	}, nil
}

func (i *PaymentImplement) CallbackPaymentTransaction(ctx context.Context, transactionID string, statusCode string) (customError *exceptions.CustomError) {

	sessionID, customError := i.repository.CallbackPaymentTransaction(ctx, transactionID, statusCode)
	if customError != nil {
		return customError
	}

	if sessionID != uuid.Nil {
		customError = i.cache.DeleteCachedTable(sessionID)
		if customError != nil {
			return customError
		}
	}

	return nil
}

func tlv(tag, val string) string {
	return fmt.Sprintf("%s%02d%s", tag, len(val), val)
}

// 0891234567 -> 0066891234567
func phoneToPromptPayProxy(thMobile string) string {
	m := strings.TrimSpace(thMobile)
	if strings.HasPrefix(m, "0") {
		m = m[1:]
	}
	return "0066" + m
}

type BuildOpts struct {
	// เลือกใส่ “อย่างใดอย่างหนึ่ง”
	Mobile    string // เบอร์มือถือ (ถ้ามี)
	CitizenID string // เลขบัตรประชาชน (13 หลัก) ถ้าจะใช้แทนมือถือ

	Amount   string // เช่น "520.00" (แนะนำให้กำหนดเสมอสำหรับ Dynamic)
	Ref      string // reference/order/transaction id (optional แต่ดีมาก)
	Merchant string // เช่น "FOOD STORE"
	City     string // เช่น "BANGKOK"
	Dynamic  bool   // true=12 (Dynamic), false=11 (Static)
}

// สร้าง Dynamic PromptPay QR Text (Thai QR Payment)
func BuildPromptPayQR(o BuildOpts) (string, error) {
	if o.Mobile == "" && o.CitizenID == "" {
		return "", fmt.Errorf("missing proxy (Mobile or CitizenID)")
	}
	if o.Merchant == "" {
		o.Merchant = "MERCHANT"
	}
	if o.City == "" {
		o.City = "BANGKOK"
	}

	// 00: Payload Format Indicator
	pfi := tlv("00", "01")

	// 01: Point of Initiation Method
	pim := tlv("01", map[bool]string{true: "12", false: "11"}[o.Dynamic])

	// 29: Merchant Account Information (PromptPay)
	//   00 = AID
	//   01/02 = proxy (01=mobile, 02=citizen)
	aid := tlv("00", "A000000677010111")

	var proxyTLV string
	if o.Mobile != "" {
		proxyTLV = tlv("01", phoneToPromptPayProxy(o.Mobile))
	} else {
		// Citizen ID ต้องมี 13 หลัก (เฉพาะตัวเลข)
		proxyTLV = tlv("02", o.CitizenID)
	}
	mai := tlv("29", aid+proxyTLV)

	// 52: MCC, 53: Currency(764), 54: Amount, 58: Country(TH), 59: Name, 60: City
	mcc := tlv("52", "0000")
	cur := tlv("53", "764")
	amt := ""
	if strings.TrimSpace(o.Amount) != "" {
		amt = tlv("54", o.Amount)
	}
	cc := tlv("58", "TH")
	name := tlv("59", o.Merchant)
	city := tlv("60", o.City)

	// 62: Additional Data Field Template (01 = Reference)
	additional := ""
	if o.Ref != "" {
		additional = tlv("62", tlv("01", o.Ref))
	}

	// รวมก่อนคำนวณ CRC
	payload := pfi + pim + mai + mcc + cur + amt + cc + name + city + additional
	payloadForCRC := payload + "6304"

	crc := crc16.ChecksumCCITTFalse([]byte(payloadForCRC))
	crcStr := strings.ToUpper(fmt.Sprintf("%04x", crc))

	return payload + tlv("63", crcStr), nil
}
