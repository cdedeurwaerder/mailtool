package repository

import (
	"path/filepath"

	"github.com/cdedeurwaerder/mailtool/internal/business"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type suspiciousEmailDao struct {
	gorm.Model
	UserID       string
	EmailID      string
	From         string
	To           string
	EnvelopeFrom string
	RcptTo       string
	Title        string
	Body         string
	Flag         int
	Reason       string
}

func (s suspiciousEmailDao) TableName() string {
	return "suspicious_emails"
}

type Sqlite struct {
	db *gorm.DB
}

func NewSqliteRepository(path string) (*Sqlite, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(filepath.Join(path, "gorm.db")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&suspiciousEmailDao{})

	return &Sqlite{
		db: db,
	}, nil
}

func (s Sqlite) StoreSuspiciousEmail(se business.SuspiciousEmail) error {
	dao := suspiciousEmailToDao(se)
	tx := s.db.Save(&dao)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// func suspiciousEmailDaoToBusiness(dao suspiciousEmailDao) business.SuspiciousEmail {
// 	return business.SuspiciousEmail{
// 		UserID:       dao.UserID,
// 		EmailID:      dao.EmailID,
// 		From:         dao.From,
// 		To:           dao.To,
// 		EnvelopeFrom: dao.EnvelopeFrom,
// 		RcptTo:       dao.RcptTo,
// 		Title:        dao.Title,
// 		Body:         dao.Body,
// 		Flag:         business.EmailFlag(dao.Flag),
// 		Reason:       dao.Reason,
// 	}
// }

func suspiciousEmailToDao(se business.SuspiciousEmail) suspiciousEmailDao {
	return suspiciousEmailDao{
		UserID:       se.UserID,
		EmailID:      se.EmailID,
		From:         se.From,
		To:           se.To,
		EnvelopeFrom: se.EnvelopeFrom,
		RcptTo:       se.RcptTo,
		Title:        se.Title,
		Body:         se.Body,
		Flag:         int(se.Flag),
		Reason:       se.Reason,
	}
}
