package models

import (
	"time"
)

// OTP is the struct for OTP in the database
type OTP struct {
	ID          int64     `json:"id" gorm:"primaryKey;unique;not null;autoIncrement"`
	Code        string    `json:"code" gorm:"type:varchar(6);not null" binding:"required,min=6,max=6,numeric"`
	IsUsed      bool      `json:"is_used" gorm:"type:boolean;not null;default:false"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(15);not null" binding:"required,e164,min=11,max=14"`
	KeyUID      string    `json:"key_uid" gorm:"type:varchar(100);not null" binding:"required,uuid"`
	ExpiryAt    time.Time `json:"expiry_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CheckOTP is the struct for checking the OTP
type CheckOTP struct {
	Code        string    `json:"code" binding:"required,min=6,max=6,numeric"`
	PhoneNumber string    `json:"phone_number" binding:"required,e164,min=11,max=14"`
	KeyUID      string    `json:"key_uid" binding:"required,uuid"`
	ExpiryAt    time.Time `json:"expiry_at"`
	CreatedAt   time.Time `json:"created_at"`
}
