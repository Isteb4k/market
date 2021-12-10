package auth

import (
	"errors"
	"log"
	"market/models"
	"market/repositories"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var (
	InvalidCodeLengthErr  = errors.New("invalid code length")
	InvalidPhoneNumberErr = errors.New("invalid phone number")
	IncorrectCodeErr      = errors.New("incorrect code")
	FailedToGetTokenErr   = errors.New("failed to get token")
)

type SignInRes struct {
	AuthToken string
	User      models.User
}

type Client interface {
	SignIn(phone string, code string) (*SignInRes, error)
	SendCode(phone string) (string, error)
}

type auth struct {
	usersRepo repositories.Users
	codes     map[string]string
}

func New(usersRepo repositories.Users) Client {
	return &auth{
		usersRepo: usersRepo,
		codes:     make(map[string]string),
	}
}

var phoneRegex = regexp.MustCompile(`^(\+|\d)\d{1,19}$`)

// SignIn - sign in with an existing user or register if the user does not exist
func (a *auth) SignIn(phone string, code string) (*SignInRes, error) {
	if len(code) != 4 {
		return nil, InvalidCodeLengthErr
	}

	if !phoneRegex.MatchString(strings.Replace(phone, " ", "", -1)) {
		return nil, InvalidPhoneNumberErr
	}

	phoneByCode, has := a.codes[code]

	if !has || phoneByCode != phone {
		return nil, IncorrectCodeErr
	}

	user, err := a.usersRepo.GetUserByPhone(phone)
	if err == repositories.NotFoundErr {
		user, err = a.usersRepo.CreateUser(phone)
	}

	if err != nil {
		return nil, err
	}

	token, err := user.GenToken()
	if err != nil {
		return nil, FailedToGetTokenErr
	}

	return &SignInRes{
		AuthToken: token,
		User:      user,
	}, nil
}

// SendCode - generate new code for phone
func (a *auth) SendCode(phone string) (string, error) {
	if !phoneRegex.MatchString(strings.Replace(phone, " ", "", -1)) {
		return "", InvalidPhoneNumberErr
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("0123456789")
	length := 4
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	code := b.String()

	a.codes[code] = phone

	log.Println("Code:", code)

	return code, nil
}
