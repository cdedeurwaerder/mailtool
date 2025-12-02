package provider

import (
	"time"

	"github.com/cdedeurwaerder/mailtool/internal/business"
	"github.com/google/uuid"
)

type InMemApi struct {
}

func (ima InMemApi) GetUsers(tenantID uuid.UUID) ([]business.User, error) {
	return []business.User{
		business.User{Firstname: "Michel", Lastname: "Dudu", Email: "michel.dudu@aol.com"},
		business.User{Firstname: "Bernard", Lastname: "De La Haute", Email: "bernard.delahaute@wanadoo.fr"},
	}, nil
}

func (ima InMemApi) GetEmails(userID uuid.UUID, receivedAfter time.Time) ([]business.Email, error) {
	return []business.Email{}, nil
}
