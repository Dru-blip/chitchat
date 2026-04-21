package mailer

import "github.com/stretchr/testify/mock"

type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendMagicLink(recipient, link string) error {
	return m.Called(recipient, link).Error(0)
}
