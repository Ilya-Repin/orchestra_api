package model

import (
	"github.com/google/uuid"
	"regexp"
	"time"
)

type MemberStatus string

const (
	StatusPending  MemberStatus = "pending"
	StatusApproved MemberStatus = "approved"
	StatusDeclined MemberStatus = "declined"
)

type Member struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Phone     string
	Status    MemberStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	return re.MatchString(email)
}

func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^7\d{10}$`)
	return re.MatchString(phone)
}
