package models

import "time"

// User represents the structure for the user table in the database
type User struct {
	ID               int        `json:"id"`
	Username         string     `json:"username"`
	FirstName        string     `json:"first_name,omitempty"`
	LastName         string     `json:"last_name,omitempty"`
	Password         string     `json:"-"` // ไม่แนะนำให้ส่ง password ใน JSON response
	UserTypeID       *int       `json:"user_type_id,omitempty"`
	BranchID         *int       `json:"branch_id,omitempty"`
	ProfileImagePath *string    `json:"profile_image_path,omitempty"`
	CreateDate       time.Time  `json:"create_date"`
	UpdateDate       time.Time  `json:"update_date"`
	StatusID         *int       `json:"status_id,omitempty"`
	RoleID           *int       `json:"role_id,omitempty"`
	Nickname         *string    `json:"nickname,omitempty"`
	PhoneNumber      *string    `json:"phone_number,omitempty"`
	Email            *string    `json:"email,omitempty"`
	EditBy           *string    `json:"edit_by,omitempty"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	UIndex           *int       `json:"u_index,omitempty"`
	UCode            *string    `json:"u_code,omitempty"` // คำนวณเป็น Virtual column ใน DB
	Permission       *string    `json:"permission,omitempty"`
}
