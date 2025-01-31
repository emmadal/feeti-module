package models

import (
	"time"
)

// Otp is the struct for OTP in the database
type Otp struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string    `json:"code" gorm:"type:varchar(10);not null;index:idx_code,priority:1" binding:"required,min=5,max=5,numeric"`
	IsUsed      bool      `json:"is_used" gorm:"type:boolean;not null;default:false;index:idx_is_used,priority:2"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(15);not null;index:idx_phone_number,priority:1" binding:"required,e164,min=11,max=14"`
	KeyUID      string    `json:"key_uid" gorm:"type:varchar(100);not null;index:idx_key_uid,priority:2" binding:"required,uuid"`
	ExpiryAt    time.Time `json:"expiry_at" gorm:"index:idx_expiry_at,priority:3" binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// CheckOtp is the struct for checking the OTP
type CheckOtp struct {
	Code        string    `json:"code" binding:"required,min=5,max=5,numeric"`
	PhoneNumber string    `json:"phone_number" binding:"required,e164,min=11,max=14"`
	KeyUID      string    `json:"key_uid" binding:"required,uuid"`
	ExpiryAt    time.Time `json:"expiry_at"`
}

// NewOtp is the struct for creating a new OTP
type NewOtp struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=11,max=14"`
}
