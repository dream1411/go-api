package models

import "time"

type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"-"`
	Email            *string   `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	UserTypeID       int       `json:"user_type_id"`
	BranchID         int       `json:"branch_id"`
	ProfileImagePath *string   `json:"profile_image_path"`
	CreateDate       time.Time `json:"createDate"`
	UpdateDate       time.Time `json:"updateDate"`
	StatusID         int       `json:"status_id"`
	RoleID           int       `json:"role_id"`
	Nickname         string    `json:"nickname"`
	PhoneNumber      string    `json:"phone_number"`
	EditBy           *string   `json:"edit_by"`
	StartDate        time.Time `json:"start_date"`
	UIndex           int       `json:"u_index"`
	UCode            string    `json:"u_code"`
	Permission       *string   `json:"permission"`
}
