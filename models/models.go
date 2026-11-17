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
	URLDetail          *string        `gorm:"column:url_detail" json:"url_detail"`
	PathThumb          *string        `json:"path_thumb"`
	ImageThumbScale    *string        `json:"image_thumb_scale"`
	FBPostID           *string        `gorm:"column:fb_post_id" json:"fb_post_id"`
	FBAlbumPhotoID     *string        `gorm:"column:fb_album_photo_id" json:"fb_album_photo_id"`
	FBAlbumApp         int            `gorm:"column:fb_album_app;default:1" json:"fb_album_app"`
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
	DisplayThumbURL    *string        `gorm:"-" json:"display_thumb"`
}

// DisplayThumb returns ImageThumbScale if present, otherwise PathThumb
func (p Product) DisplayThumb() string {
	if p.ImageThumbScale != nil && *p.ImageThumbScale != "" {
		return *p.ImageThumbScale
	}
	if p.PathThumb != nil {
		return *p.PathThumb
	}
	return ""
}

// BestQualityThumb returns the original full-resolution image (PathThumb) for Facebook uploads instead of low-res thumbnail
func (p Product) BestQualityThumb() string {
	if p.PathThumb != nil && *p.PathThumb != "" {
		return *p.PathThumb
	}
	if len(p.Images) > 0 {
		for _, img := range p.Images {
			if img.IsMain && img.FilePath != "" {
				return img.FilePath
			}
		}
		if p.Images[0].FilePath != "" {
			return p.Images[0].FilePath
		}
	}
	if p.ImageThumbScale != nil && *p.ImageThumbScale != "" {
		return *p.ImageThumbScale
	}
	return ""
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
	ShippingFee      *float64   `json:"shipping_fee"`
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
	Product      *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
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

type Category struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code              string     `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	Name              string     `gorm:"type:varchar(100)" json:"name"`
	FBAlbumID         *string    `gorm:"column:fb_album_id;type:varchar(100)" json:"fb_album_id"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	NeedsAlbumReorder bool       `gorm:"column:needs_album_reorder;default:false" json:"needs_album_reorder"`
	CreatedAt         *time.Time `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

type AlbumResyncJob struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryCode       string     `gorm:"type:varchar(50);index" json:"category_code"`
	CategoryName       string     `gorm:"type:varchar(100)" json:"category_name"`
	AlbumID            string     `gorm:"type:varchar(100)" json:"album_id"`
	TriggerType        string     `gorm:"type:varchar(20)" json:"trigger_type"` // "MANUAL" or "NIGHTLY"
	Status             string     `gorm:"type:varchar(30);index" json:"status"` // "PENDING", "DELETING", "UPLOADING", "PAUSED_LIMIT", "COMPLETED", "CANCELLED", "FAILED"
	Phase              string     `gorm:"type:varchar(20)" json:"phase"` // "DELETE" or "UPLOAD"
	TotalDelete        int        `json:"total_delete"`
	DeletedCount       int        `json:"deleted_count"`
	TotalUpload        int        `json:"total_upload"`
	UploadedCount      int        `json:"uploaded_count"`
	RemainingDeleteIDs string     `gorm:"type:longtext" json:"-"`
	RemainingUploadIDs string     `gorm:"type:longtext" json:"-"`
	ResumeAt           *time.Time `json:"resume_at"`
	LastError          string     `gorm:"type:text" json:"last_error"`
	StartedAt          *time.Time `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

