package dto

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// Failed
	MESSAGE_FAILED_REGISTER_USER           = "gagal membuat pengguna"
	MESSAGE_FAILED_GET_USER_TOKEN          = "gagal mendapatkan token pengguna"
	MESSAGE_FAILED_GET_USER                = "gagal mendapatkan pengguna"
	MESSAGE_FAILED_LOGIN                   = "login gagal"
	MESSAGE_FAILED_WRONG_EMAIL_OR_PASSWORD = "email atau password salah"
	MESSAGE_FAILED_UPDATE_USER             = "gagal memperbarui pengguna"
	MESSAGE_FAILED_CHANGE_PASSWORD         = "gagal mengubah password"
	MESSAGE_FAILED_RESET_PASSWORD          = "gagal mereset password"
	MESSAGE_FAILED_DELETE_USER             = "gagal menghapus pengguna"
	MESSAGE_FAILED_PROCESSING_REQUEST      = "gagal memproses permintaan"
	MESSAGE_FAILED_DENIED_ACCESS           = "akses ditolak"
	MESSAGE_FAILED_SEND_OTP_EMAIL          = "gagal mengirim verifikasi otp ke email"
	MESSAGE_FAILED_VERIFY_EMAIL            = "gagal memverifikasi email"
	MESSAGE_FAILED_REFRESHING_TOKEN        = "gagal memperbarui token"
	MESSAGE_FAILED_LOGOUT                  = "gagal logout"
	MESSAGE_FAILED_UPLOAD_PROFILE_PICTURE  = "gagal mengunggah gambar"
	MESSAGE_FAILED_GET_PROFILE_PICTURE     = "gagal mendapatkan gambar profil"
	MESSAGE_FAILED_DELETE_PROFILE_PICTURE  = "gagal menghapus gambar profil"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER          = "membuat pengguna berhasil"
	MESSAGE_SUCCESS_GET_USER               = "mendapatkan pengguna berhasil"
	MESSAGE_SUCCESS_LOGIN                  = "login berhasil"
	MESSAGE_SUCCESS_UPDATE_USER            = "memperbarui pengguna berhasil"
	MESSAGE_SUCCESS_CHANGE_PASSWORD        = "mengubah password berhasil"
	MESSAGE_SUCCESS_RESET_PASSWORD         = "mereset password berhasil"
	MESSAGE_SUCCESS_DELETE_USER            = "menghapus pengguna berhasil"
	MESSAGE_SEND_OTP_EMAIL_SUCCESS         = "mengirim verifikasi otp ke email berhasil"
	MESSAGE_SUCCESS_VERIFY_EMAIL           = "memverifikasi email berhasil"
	MESSAGE_SUCCESS_REFRESH_TOKEN          = "memperbarui token berhasil"
	MESSAGE_SUCCESS_LOGOUT                 = "logout berhasil"
	MESSAGE_SUCCESS_UPLOAD_PROFILE_PICTURE = "mengunggah foto profil berhasil"
	MESSAGE_SUCCESS_DELETE_PROFILE_PICTURE = "menghapus foto profil berhasil"
)

var (
	ErrCreateUser                     = errors.New("gagal membuat pengguna")
	ErrGetAllUser                     = errors.New("gagal mendapatkan semua pengguna")
	ErrGetUserById                    = errors.New("gagal mendapatkan pengguna berdasarkan ID")
	ErrGetUserByEmail                 = errors.New("gagal mendapatkan pengguna berdasarkan email")
	ErrEmailAlreadyExists             = errors.New("email sudah ada")
	ErrUpdateUser                     = errors.New("gagal memperbarui pengguna")
	ErrUserNotAdmin                   = errors.New("pengguna bukan admin")
	ErrUserNotFound                   = errors.New("pengguna tidak ditemukan")
	ErrEmailNotFound                  = errors.New("email tidak ditemukan")
	ErrDeleteUser                     = errors.New("gagal menghapus pengguna")
	ErrPasswordNotMatch               = errors.New("password tidak cocok")
	ErrPasswordSame                   = errors.New("password baru tidak boleh sama dengan password lama")
	ErrInvalidOldPassword             = errors.New("password lama salah")
	ErrEmailOrPassword                = errors.New("email atau password salah")
	ErrAccountNotVerified             = errors.New("akun belum diverifikasi")
	ErrAccountNotVerifiedWhenRegister = errors.New("akun belum diverifikasi. OTP baru telah dikirim")
	ErrTokenInvalid                   = errors.New("token tidak valid")
	ErrTokenExpired                   = errors.New("token kadaluarsa")
	ErrAccountAlreadyVerified         = errors.New("akun sudah diverifikasi. silakan login")
	ErrProfilePictureNotFound         = errors.New("foto profil tidak ditemukan")
	ErrUserNotFoundGorm               = gorm.ErrRecordNotFound
)

type (
	UserCreateDTO struct {
		ID       uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
		Email    string    `json:"email" form:"email" binding:"required"`
		Name     string    `json:"name" form:"name" binding:"required"`
		Password string    `json:"password" form:"password" binding:"required"`
	}

	UserMeResponseDTO struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		Name         string    `json:"name"`
		AsalInstansi string    `json:"asal_instansi"`
		DateOfBirth  time.Time `json:"date_of_birth"`
		IsVerified   bool      `json:"is_verified"`
		RoleID       string    `json:"role_id"`
	}

	UserUpdateInfoDTO struct {
		ID uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
		// GoogleID     string    `gorm:"type:varchar(255);" json:"google_id"`
		Name         string    `json:"name" form:"name" binding:"required"`
		AsalInstansi string    `json:"asal_instansi" form:"asal_instansi" binding:"required"`
		DateOfBirth  time.Time `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	}

	SendUserOTPByEmail struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	VerifyUserOTPByEmail struct {
		Email string `json:"email" form:"email" binding:"required"`
		OTP   string `json:"otp" form:"otp" binding:"required"`
	}

	UserUpdateEmailDTO struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserUpdatePhoneDTO struct {
		Phone string `json:"phone" form:"phone" binding:"required"`
	}

	UploadFileDTO struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}

	UserChangePassword struct {
		OldPassword string `json:"old_password" form:"old_password" binding:"required"`
		NewPassword string `json:"new_password" form:"new_password" binding:"required"`
	}

	ForgotPassword struct {
		Email string `json:"email" form:"email" binding:"required"`
	}

	UserLoginDTO struct {
		Email    string `json:"email" form:"email" binding:"email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	UserLoginResponseDTO struct {
		SessionToken string `json:"session_token"`
		RefreshToken string `json:"refresh_token"`
		Role         string `json:"role"`
	}

	UserRefreshDTO struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
	}
)
