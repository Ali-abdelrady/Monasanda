package models

import (
	"fmt"
	"time"

	"mosanda-vacation-attendace/pkg/osioutil"

	"github.com/google/uuid"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticatedUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (user *User) Validate() error {
	if user.Email == "" {
		return fmt.Errorf("email field is required")
	}

	return nil
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type VacationBatch struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	CreatedBy int       `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

const (
	AttendanceStatusPresent = "present"
	AttendanceStatusAbsent  = "absent"
)

type Attendance struct {
	ID             int                `json:"id"`
	BatchID        int                `json:"batchId"`
	Name           *string            `json:"name"`
	Location       *string            `json:"location"`
	PhoneNumber    *string            `json:"phoneNumber"`
	MilitaryNumber *string            `json:"militaryNumber"`
	Status         string             `json:"status"`
	Notes          *string            `json:"notes"`
	StartDate      *time.Time         `json:"startDate"`
	EndDate        *time.Time         `json:"endDate"`
	CheckedInAt    *time.Time         `json:"checkedInAt"`
	CheckedInBy    *AuthenticatedUser `json:"checkedInBy"`
	CheckedInByID  *int               `json:"-"`
	CreatedBy      int                `json:"createdBy"`
	CreatedAt      *time.Time         `json:"createdAt"`
	UpdatedAt      *time.Time         `json:"updatedAt"`
}

func (attendance *Attendance) Validate() error {
	if attendance == nil {
		return fmt.Errorf("attendance record is required")
	}
	if attendance.BatchID <= 0 {
		return fmt.Errorf("batchId must be a positive integer")
	}
	if attendance.Status != "" &&
		attendance.Status != AttendanceStatusPresent &&
		attendance.Status != AttendanceStatusAbsent {
		return fmt.Errorf("status must be either present or absent")
	}
	if attendance.PhoneNumber != nil && len(*attendance.PhoneNumber) > 12 {
		return fmt.Errorf("phoneNumber must not exceed 12 characters")
	}
	if attendance.MilitaryNumber != nil && len(*attendance.MilitaryNumber) > 13 {
		return fmt.Errorf("militaryNumber must not exceed 12 characters")
	}
	if attendance.StartDate != nil && attendance.EndDate != nil && attendance.EndDate.Before(*attendance.StartDate) {
		return fmt.Errorf("endDate must not be before startDate")
	}
	return nil
}

func (batch *VacationBatch) Validate() error {
	if batch == nil {
		return fmt.Errorf("vacation batch is required")
	}

	if batch.Name == "" {
		return osioutil.GetRefErrorTxt("MDL0201", "name field is requeird")
	}

	// if batch.CreatedBy <= 0 {
	// 	return fmt.Errorf("createdBy must be a positive integer")
	// }
	return nil
}
