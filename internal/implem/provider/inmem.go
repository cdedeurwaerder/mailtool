package provider

import (
	"time"

	"github.com/cdedeurwaerder/mailtool/internal/business"
	"github.com/google/uuid"
)

type InMemApi struct {
	users []business.User
	mails map[string][]business.Email
}

func NewInMemApi() InMemApi {
	users := []business.User{
		{ID: uuid.MustParse("48ff42b8-f93e-4b8b-92ed-5e3c3b5b710f"), Firstname: "Michel", Lastname: "Dudu", Email: "michel.dudu@aol.com"},
		{ID: uuid.MustParse("884611ea-6e1f-457d-88a8-a0d9ca2e2e31"), Firstname: "Bernard", Lastname: "De La Haute", Email: "bernard.delahaute@wanadoo.fr"},
	}

	mails := map[string][]business.Email{
		users[0].ID.String(): {
			business.Email{
				ID:           "abc",
				From:         "Paulo",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@samplefake.com",
				Title:        "My emergency request",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "def",
				From:         "Paulo",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@samplefake.com",
				Title:        "My second emergency request",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "fgh",
				From:         "Tata",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "tata@samplefake.com",
				Title:        "My email",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "ijk",
				From:         "Luigi",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@mario.com",
				Title:        "My normal request",
				Body:         "this is a dummy text",
			},
		},
		users[1].ID.String(): {
			business.Email{
				ID:           "abc",
				From:         "Paulo",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@samplefake.com",
				Title:        "My emergency request",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "def",
				From:         "Paulo",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@samplefake.com",
				Title:        "My second emergency request",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "fgh",
				From:         "Tata",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "tata@samplefake.com",
				Title:        "My email",
				Body:         "this is a dummy text",
			},
			business.Email{
				ID:           "ijk",
				From:         "Luigi",
				To:           users[0].Firstname,
				RcptTo:       users[0].Email,
				EnvelopeFrom: "toto@mario.com",
				Title:        "My normal request",
				Body:         "this is a dummy text",
			},
		},
	}

	return InMemApi{
		users: users,
		mails: mails,
	}
}

func (ima InMemApi) GetUsers(tenantID uuid.UUID) ([]business.User, error) {
	return ima.users, nil
}

func (ima InMemApi) GetEmails(userID uuid.UUID, receivedAfter time.Time) ([]business.Email, error) {
	return ima.mails[userID.String()], nil
}
