package controller

import (
	"context"
	"net/http"
	"teach-tech-ai/common"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/service"
	"time"

	"github.com/gin-gonic/gin"
)

type OTPTwilioController interface {
	SendSMS(ctx *gin.Context)
	VerifySMS(ctx *gin.Context)
}

type otpTwilioController struct {
	otpTwilioService service.OTPTwilioService
}

func NewOTPTwilioController(otp service.OTPTwilioService) OTPTwilioController {
	return &otpTwilioController{
		otpTwilioService: otp,
	}
}

const appTimeout = time.Second * 10

func (o *otpTwilioController) SendSMS(ctx *gin.Context) {
	var smsData dto.GenerateOTPRequest
	if err := ctx.ShouldBind(&smsData); err != nil {
		response := common.BuildErrorResponse("OTP Gagal", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	_, cancel := context.WithTimeout(context.Background(), appTimeout)
	defer cancel()

	newOTP := entity.OTPData{
		PhoneNumber: smsData.PhoneNumber,
	}

	_, err := o.otpTwilioService.TwilioSendOTP(newOTP.PhoneNumber)
	if err != nil {
		response := common.BuildErrorResponse("OTP Gagal", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := common.BuildResponse(true, "OTP berhasil terkirim", common.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (o *otpTwilioController) VerifySMS(ctx *gin.Context) {
	var verifyData dto.VerifyOTPRequest
	if err := ctx.ShouldBind(&verifyData); err != nil {
		response := common.BuildErrorResponse("Verifikasi OTP Gagal", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	_, cancel := context.WithTimeout(context.Background(), appTimeout)
	defer cancel()

	newVerifyOTP := entity.VerifyData{
		PhoneNumber: verifyData.PhoneNumber,
		Code: verifyData.Code,
	}

	err := o.otpTwilioService.TwilioVerifyOTP(newVerifyOTP.PhoneNumber, newVerifyOTP.Code)
	if err != nil {
		response := common.BuildErrorResponse("Verifikasi OTP Gagal", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := common.BuildResponse(true, "Verifikasi OTP berhasil", common.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}