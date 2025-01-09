package testhelper

import (
	"fmt"

	"github.com/flohansen/auther/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type userMatcher struct {
	username *string
	password []byte
}

func UserMatches(opts ...UserMatcherOpt) *userMatcher {
	m := &userMatcher{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m userMatcher) Matches(x any) bool {
	user, ok := x.(service.User)
	if !ok {
		return false
	}

	if m.password != nil {
		if bcrypt.CompareHashAndPassword(user.PasswordHash, m.password) != nil {
			return false
		}
	}

	if m.username != nil {
		if user.Username != *m.username {
			return false
		}
	}

	return true
}

func (m userMatcher) String() string {
	return fmt.Sprintf("username: %s, password: %s", *m.username, m.password)
}

type UserMatcherOpt func(*userMatcher)

func PasswordHashOf(password string) UserMatcherOpt {
	return func(m *userMatcher) {
		m.password = []byte(password)
	}
}

func Username(username string) UserMatcherOpt {
	return func(m *userMatcher) {
		m.username = &username
	}
}
