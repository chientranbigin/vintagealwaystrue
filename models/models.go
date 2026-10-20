package models

import (
	"time"
)

type Product struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Type               *string        `json:"type"`
	Name               *string        `json:"name"`
	ImportPrice        *float64       `json:"import_price"`
	Price              *int64         `json:"price"`
	URLDetail          *string        `json:"url_detail"`
	PathThumb          *string        `json:"path_thumb"`
	ImageThumbScale    *string        `json:"image_thumb_scale"`
	FBPostID           *string        `json:"fb_post_id"`
	Status             *string        `json:"status"`
	CreatedAt          *time.Time     `json:"created_at"`
	UpdatedAt          *time.Time     `json:"updated_at"`
	Description        *string        `json:"description"`
	Color              *string        `json:"color"`
	Fabric             *string        `json:"fabric"`
	DetailFabric       *string        `json:"detail_fabric"`
	Title              *string        `json:"title"`
	Images             []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
	Sizes              []Size         `gorm:"foreignKey:ProductID" json:"sizes"`
	LatestUpload       *string        `gorm:"->;-:migration" json:"latest_upload,omitempty"`
	ImageThumbScaleURL *string        `gorm:"-" json:"image_thumb_scale_url"`
}

type ProductImage struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint       `json:"product_id"`
	FilePath  string     `json:"file_path"`
	IsMain    bool       `json:"is_main"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Size struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID *uint      `json:"product_id"`
	Name      *string    `json:"name"`
	Value     *float64   `json:"value"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Order struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             *string    `json:"name"`
	Phone            *string    `json:"phone"`
	Address          *string    `json:"address"`
	Note             *string    `json:"note"`
	Desc             *string    `json:"desc"`
	AdditionalAmount *int64     `json:"additional_amount"`
	FreeshipAmount   *float64   `json:"freeship_amount"`
	DepositAmount    *int64     `json:"deposit_amount"`
	TotalAmount      *int64     `json:"total_amount"`
	FinalAmount      *int64     `json:"final_amount"`
	IsFreeship       bool       `json:"is_freeship"`
	IsPaid           bool       `json:"is_paid"`
	Status           *string    `json:"status"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	ConsumerID       *int       `json:"consumer_id"`
	Expenses         *string    `gorm:"type:longtext" json:"expenses"`
	Products         []Product  `gorm:"many2many:order_products;" json:"products"`
	ViettelAmount    int64      `gorm:"-" json:"viettel_amount"`
	IsOld            bool       `gorm:"-" json:"is_old"`
	CreatedAtHuman   string     `gorm:"-" json:"created_at_human"`
}

type OrderProduct struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   *uint      `json:"order_id"`
	ProductID *uint      `json:"product_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Consumer struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      *string    `json:"name"`
	Phone     *string    `json:"phone"`
	Address   *string    `json:"address"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Expense struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	EntryDate   string     `gorm:"type:date" json:"entry_date"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProductUploadLog struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID    string     `json:"session_id"`
	FileName     string     `json:"file_name"`
	FilePath     string     `json:"file_path"`
	ProductCode  *string    `json:"product_code"`
	ProductID    *int       `json:"product_id"`
	Amount       *float64   `json:"amount"`
	Status       string     `json:"status"`
	Message      *string    `json:"message"`
	DetectedText *string    `json:"detected_text"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type User struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  *string    `json:"-"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

