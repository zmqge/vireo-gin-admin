package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zmqge/vireo-gin-admin/app/admin/services"
)

type MFA struct {
	service *services.MFA
}

func NewMFA(service *services.MFA) *MFA {
	return &MFA{service: service}
}

// 启用 MFA
func (c *MFA) Enable(ctx *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secretKey, qrCodeURL, err := c.service.Enable(request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"secret_key":  secretKey,
		"qr_code_url": qrCodeURL,
	})
}

// 验证 OTP
func (c *MFA) VerifyOTP(ctx *gin.Context) {
	var request struct {
		UserID  uint   `json:"user_id" binding:"required"`
		OTPCode string `json:"otp_code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValid, err := c.service.VerifyOTP(request.UserID, request.OTPCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"is_valid": isValid})
}

// 禁用 MFA
func (c *MFA) DisableMFA(ctx *gin.Context) {
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.DisableMFA(request.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "MFA 已禁用"})
}
