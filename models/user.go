package models

import (
	"time"
)

// User is the struct for a user
type User struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;unique;not null"`
	FirstName   string    `json:"first_name" gorm:"type:varchar(100);not null" binding:"required,alpha"`
	LastName    string    `json:"last_name" gorm:"type:varchar(100);not null" binding:"required,alpha"`
	Email       string    `json:"email" gorm:"type:varchar(200)"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(15);uniqueIndex;unique;not null" binding:"required,e164,min=11,max=14"`
	DeviceToken string    `json:"device_token" gorm:"type:varchar(200);not null" binding:"required;min=10,max=200"`
	Pin         string    `json:"pin" gorm:"type:varchar(100);not null" binding:"required"`
	Quota       uint      `json:"quota" gorm:"type:bigint;default:0;not null"`
	Locked      bool      `json:"locked" gorm:"type:boolean;default:false"`
	Photo       string    `json:"photo" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean;default:true;index"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Wallet is the struct for a wallet
type Wallet struct {
	ID        int64     `json:"id" gorm:"primaryKey;unique"`
	UserID    int64     `json:"user_id" gorm:"type:bigint;not null" binding:"required,number,gt=0"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Balance   uint      `json:"balance" gorm:"type:bigint;default:0;not null" binding:"required,number,min=0,max=50000000"`
	Currency  string    `json:"currency" gorm:"type:varchar(3);default:XOF;not null" binding:"alpha,oneof=XOF GHS XAF GNH EUR USD"`
	IsActive  bool      `json:"is_active" gorm:"type:boolean;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserLogin is the struct for login
type UserLogin struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=11,max=14"`
	Pin         string `json:"pin" binding:"required,len=4,numeric,min=4,max=4"`
}

// UserResetPin is the struct for resetting the pin
type UserResetPin struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=11,max=14"`
	Pin         string `json:"pin" binding:"required,len=4,numeric,min=4,max=4"`
	CodeOTP     string `json:"code_otp" binding:"required,len=6,numeric,min=6,max=6"`
	KeyUID      string `json:"key_uid" binding:"required,uuid"`
}

// UserUpdatePin is the struct for updating the pin
type UserUpdatePin struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164,min=11,max=14"`
	OldPin      string `json:"old_pin" binding:"required,len=4,numeric,min=4,max=4"`
	NewPin      string `json:"new_pin" binding:"required,len=4,numeric,min=4,max=4"`
	CodeOTP     string `json:"code_otp" binding:"required,len=6,numeric,min=6,max=6"`
	KeyUID      string `json:"key_uid" binding:"required,uuid"`
}
