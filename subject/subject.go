package subject

const (
	SubjectWalletCreate       = "wallet.create"
	SubjectWalletGet          = "wallet.get"
	SubjectWalletLock         = "wallet.lock"
	SubjectWalletUnlock       = "wallet.unlock"
	SubjectWalletBalance      = "wallet.balance"
	SubjectWalletDeposit      = "wallet.deposit"
	SubjectWalletWithdraw     = "wallet.withdraw"
	SubjectWalletDisable      = "wallet.disable"
	SubjectWalletEnable       = "wallet.enable"
	SubjectWalletCheckBalance = "wallet.check_balance"
)

const (
	SubjectUserCreate  = "user.create"
	SubjectUserGet     = "user.get"
	SubjectUserLock    = "user.lock"
	SubjectUserUnlock  = "user.unlock"
	SubjectUserDisable = "user.disable"
	SubjectUserEnable  = "user.enable"
)

const (
	SubjectNotificationCreate = "notification.create"
	SubjectNotificationGet    = "notification.get"
	SubjectNotificationBulk   = "notification.bulk"
	SubjectNotificationDelete = "notification.delete"
	SubjectOTPCreate          = "otp.create"
	SubjectOTPCheck           = "otp.check"
)

const (
	SubjectTransactionCreate = "transaction.create"
	SubjectTransactionGet    = "transaction.get"
)
