package business

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// GetMicrosoftUsers(tenantID uuid.UUID) ([]MicrosoftUser, error)
// GetMicrosoftEmails(userID uuid.UUID, receivedAfter time.Time) ([]MicrosoftEmail, error)
// GetGoogleUsers(tenantID uuid.UUID) ([]GoogleUser, error)
// GetGoogleEmails(userID uuid.UUID, receivedAfter time.Time) ([]GoogleEmail, error)

type Analyzer interface {
	IsPhishing(user User, email Email) (bool, error)
}

type Repository interface {
	StoreSuspiciousEmail(se SuspiciousEmail) error
}

type ProviderApi interface {
	GetUsers(tenantID uuid.UUID) ([]User, error)
	GetEmails(userID uuid.UUID, receivedAfter time.Time) ([]Email, error)
}

type Service struct {
	pa           ProviderApi
	analyzer     Analyzer
	repo         Repository
	tenantID     uuid.UUID
	rate         time.Duration
	workersCount int
}

func NewService(
	pa ProviderApi,
	tenanID uuid.UUID,
	analyzer Analyzer,
	repo Repository,
	workersCount int,
	rate time.Duration,
) *Service {

	return &Service{
		pa:           pa,
		tenantID:     tenanID,
		rate:         rate,
		workersCount: workersCount,
		analyzer:     analyzer,
		repo:         repo,
	}
}

func (s Service) Start() {
	go s.run()
}

func (s Service) Stop() {

}

func (s Service) fetchUsers(c chan AnalyzeInfo, lastPoll time.Time) {
	users, err := s.pa.GetUsers(s.tenantID)
	if err != nil {
		// TODO: implement retry logic
		log.Printf("unable to get users list from provider api: %s\n", err.Error())
		return
	}

	log.Info().Msgf("got %d users to process", len(users))

	for _, u := range users {
		c <- AnalyzeInfo{user: u, lastPoll: lastPoll}
	}
}

func (s Service) processor(id int, ais <-chan AnalyzeInfo, results chan<- SuspiciousEmail) {
	for ai := range ais {
		fmt.Println("worker", id, "started  user", ai.user.Email)
		s.processEmails(ai.user, ai.lastPoll)
		fmt.Println("worker", id, "finished user", ai.user.Email)
	}
}

type AnalyzeInfo struct {
	user     User
	lastPoll time.Time
}

func (s Service) run() {
	var (
		lastPoll    time.Time
		usersChan   = make(chan AnalyzeInfo, s.workersCount)
		resultsChan = make(chan SuspiciousEmail, 10)
	)

	for i := range s.workersCount {
		go s.processor(i, usersChan, resultsChan)
	}

	for {
		start := time.Now()
		go s.fetchUsers(usersChan, lastPoll)
		if time.Since(start) < s.rate {
			delay := time.Until(start.Add(s.rate))
			log.Info().Msgf("waiting %s to next loop", delay.Round(time.Second))
			time.Sleep(delay)
		}
	}
}

func (s Service) processEmails(user User, limit time.Time) {
	emails, err := s.pa.GetEmails(user.ID, limit)
	if err != nil {
		log.Error().Err(err).Msgf("cannot load user %s mails", user.Email)
	}

	for _, e := range emails {
		s.processEmail(user, e)
	}
}

func (s Service) processEmail(user User, email Email) {
	is, err := s.analyzer.IsPhishing(user, email)
	if err != nil {
		log.Error().Err(err).Msg("email analysis error")
		return
	}

	if !is {
		log.Info().Str("user", user.Email).Str("email", email.ID).Msg("clear")
		return
	}

	log.Info().Str("user", user.Email).Str("email", email.ID).Msg("marking email as potential phishing")
	se := SuspiciousEmail{
		UserID:       user.ID.String(),
		EmailID:      email.ID,
		From:         email.From,
		To:           email.To,
		EnvelopeFrom: email.EnvelopeFrom,
		RcptTo:       email.RcptTo,
		Title:        email.Title,
		Body:         email.Body,
		Flag:         EmailFraudSuspicious,
	}
	if err := s.repo.StoreSuspiciousEmail(se); err != nil {
		log.Error().Err(err).Str("user_id", se.UserID).Str("email_id", email.ID).Msgf("cannot store suspicious email")
	}
	// TODO also push in M365 / GW quarantine.
}
