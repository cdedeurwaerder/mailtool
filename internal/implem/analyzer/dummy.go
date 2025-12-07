package analyzer

import "github.com/cdedeurwaerder/mailtool/internal/business"

//TODO: To be implemented.

// SP / Fraud detection
// - look for sender spoofed email
// - look for pressure / emergency in the email title / body
// - payment request
// - check links
// - check attachments
// - timing notion (night or week end instead of business hours)
// - map users interactions.

type DummyAnalyzer struct {
	isPhishing bool
}

func (da *DummyAnalyzer) IsPhishing(user business.User, email business.Email) (bool, error) {
	da.isPhishing = !da.isPhishing
	return da.isPhishing, nil
}
