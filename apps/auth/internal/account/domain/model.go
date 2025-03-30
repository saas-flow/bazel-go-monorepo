package domain

import (
	"regexp"
	"strconv"
	"time"

	"github.com/saas-flow/monorepo/libs/security"
	"gorm.io/gorm"
)

type AuthProvider struct {
	ID           string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string    `gorm:"column:name;type:text;unique;not null" json:"name"`
	ClientID     string    `gorm:"column:client_id;type:text;not null" json:"client_id"`
	ClientSecret string    `gorm:"column:client_secret;type:text;not null" json:"client_secret"`
	AuthURL      string    `gorm:"column:auth_url;type:text;not null" json:"auth_url"`
	TokenURL     string    `gorm:"column:token_url;type:text;not null" json:"token_url"`
	UserInfoURL  string    `gorm:"column:user_info_url;type:text" json:"user_info_url"`
	Scopes       string    `gorm:"column:scopes;type:text" json:"scopes"`
	RedirectURI  string    `gorm:"column:redirect_uri;type:text" json:"redirect_uri"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type PasswordPattern string

var (
	MinLength                        = func(n int) PasswordPattern { return PasswordPattern("^.{" + strconv.Itoa(n) + ",}$") } // Minimal N karakter
	LowerCase        PasswordPattern = "[a-z]"                                                                                 // Setidaknya 1 huruf kecil
	UpperCase        PasswordPattern = "[A-Z]"                                                                                 // Setidaknya 1 huruf besar
	Number           PasswordPattern = "[0-9]"                                                                                 // Setidaknya 1 angka
	SpecialCharacter PasswordPattern = "[@$!%*?&]"                                                                             // Setidaknya 1 karakter spesial
)

type PasswordVersion struct {
	ID      string        `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	Version int           `gorm:"column:version;unique" json:"version"`
	Rules   PasswordRules `gorm:"foreignKey:VersionID" json:"rules"`
}

func (m *PasswordVersion) ValidatePassword(password string) bool {
	for _, rule := range m.Rules {
		matched, _ := regexp.MatchString(string(rule.Pattern), password)
		if !matched {
			return false
		}
	}
	return true
}

type PasswordRule struct {
	ID          string          `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey" json:"-"`
	VersionID   string          `gorm:"column:version_id;type:uuid;index" json:"version_id"`
	Name        string          `gorm:"column:name" json:"name"`
	DisplayName string          `gorm:"column:display_name" json:"display_name"`
	Pattern     PasswordPattern `gorm:"column:pattern" json:"pattern"`
}

type PasswordRules []PasswordRule

type UserType string

var (
	User           UserType = "USER"
	ServiceAccount UserType = "SERVICE_ACCOUNT"
)

type Account struct {
	ID                 string          `gorm:"column:id;default:uuid_generate_v4()" json:"id"`
	Type               UserType        `gorm:"column:type" json:"type"`
	Username           string          `gorm:"column:username" json:"username"`
	Password           string          `gorm:"column:password" json:"-"`
	FirstName          string          `gorm:"column:first_name" json:"first_name"`
	LastName           string          `gorm:"column:last_name" json:"last_name"`
	PasswordVersionID  string          `gorm:"column:password_version_id" json:"password_version"`
	PasswordVersion    PasswordVersion `gorm:"foreignKey:PasswordVersionID" json:"-"`
	PasswordChangeAt   time.Time       `gorm:"column:password_change_at" json:"password_change_at"`
	SkipPasswordUpdate bool            `gorm:"column:skip_password_update" json:"skipp_password_update"`
	CreatedAt          time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`
	AuthProviderID     string          `gorm:"column:auth_provider_id" json:"auth_provider_id"`
	ProviderUserID     string          `gorm:"column:provider_user_id" json:"provider_user_id"`
}

func (m *Account) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if m.Type == User {
		hash, err := security.HashArgon2(m.Password)
		if err != nil {
			return err
		}

		m.Password = hash
		m.PasswordChangeAt = now
	}

	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func (m *Account) BeforeSave(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

func (m *Account) ComparePassword(password string) bool {
	return security.VerifyHashArgon2(password, m.Password)
}
