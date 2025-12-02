package business

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// GetMicrosoftUsers(tenantID uuid.UUID) ([]MicrosoftUser, error)
// GetMicrosoftEmails(userID uuid.UUID, receivedAfter time.Time) ([]MicrosoftEmail, error)
// GetGoogleUsers(tenantID uuid.UUID) ([]GoogleUser, error)
// GetGoogleEmails(userID uuid.UUID, receivedAfter time.Time) ([]GoogleEmail, error)

type Repository interface {
}

type ProviderApi interface {
	GetUsers(tenantID uuid.UUID) ([]User, error)
	GetEmails(userID uuid.UUID, receivedAfter time.Time) ([]Email, error)
}

type Service struct {
	pa           ProviderApi
	tenantID     uuid.UUID
	rate         time.Duration
	workersCount int
}

func NewService(pa ProviderApi, tenanID uuid.UUID, workersCount int, rate time.Duration) *Service {

	return &Service{
		pa:           pa,
		tenantID:     tenanID,
		rate:         rate,
		workersCount: workersCount,
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
		close(c) // FIXME
		return
	}

	for _, u := range users {
		c <- AnalyzeInfo{user: u, lastPoll: lastPoll}
	}
}

func (s Service) processor(id int, ais <-chan AnalyzeInfo, results chan<- SuspiciousEmail) {
	for ai := range ais {
		fmt.Println("worker", id, "started  user", ai.user.Email)
		time.Sleep(time.Second)
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
			log.Printf("waiting %s to next loop\n", delay)
			time.Sleep(delay)
		}
	}
}
