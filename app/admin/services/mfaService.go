package services

import (
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/zmqge/vireo-gin-admin/app/admin/models"
	"gorm.io/gorm"
)

type MFA struct {
	db *gorm.DB
}

func NewMFA(db *gorm.DB) *MFA {
	return &MFA{db: db}
}

// 启用 MFA：生成密钥并保存到数据库
func (s *MFA) Enable(userID uint) (secretKey, qrCodeURL string, err error) {
	// 生成 TOTP 密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "YourApp",
		AccountName: fmt.Sprintf("user_%d", userID),
	})
	if err != nil {
		return "", "", err
	}

	// 保存到数据库
	mfa := models.UserMFA{
		UserID:    userID,
		SecretKey: key.Secret(),
		IsEnabled: true,
	}
	if err := s.db.Save(&mfa).Error; err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// 验证 OTP
func (s *MFA) VerifyOTP(userID uint, otpCode string) (bool, error) {
	var mfa models.UserMFA
	if err := s.db.Where("user_id = ?", userID).First(&mfa).Error; err != nil {
		return false, err
	}
	return totp.Validate(otpCode, mfa.SecretKey), nil
}

// 禁用 MFA
func (s *MFA) DisableMFA(userID uint) error {
	return s.db.Model(&models.UserMFA{}).Where("user_id = ?", userID).Update("is_enabled", false).Error
}
