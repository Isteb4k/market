package auth

import (
	"github.com/golang/mock/gomock"
	"market/models"
	"market/repositories"
	mock_repositories "market/repositories/mocks"
	"testing"
)

func TestSignIn(t *testing.T) {
	t.Parallel()

	t.Run("Incorrect phone or code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		invalidPhoneNumbers := []string{
			"1",
			"+a89998887766",
			"+89998887766d",
			"+8999-8887766",
			"6899942412421432546546546557658887766",
		}
		const validCode = "0000"
		const validPhone = "+89998887766"

		usersRepo := mock_repositories.NewMockUsers(ctrl)
		auth := New(usersRepo)

		// try with invalid phone numbers
		for _, phone := range invalidPhoneNumbers {
			_, err := auth.SignIn(phone, validCode)
			if err != InvalidPhoneNumberErr {
				t.Fatal("Failed to validate phone number", phone, err)
			}
		}

		invalidCodes := []string{
			"1",
			"00000",
		}

		// try with invalid code
		for _, code := range invalidCodes {
			_, err := auth.SignIn(validPhone, code)
			if err != InvalidCodeLengthErr {
				t.Fatal("Failed to validate code")
			}
		}

		// try with non-existent code
		_, err := auth.SignIn(validPhone, validCode)
		if err != IncorrectCodeErr {
			t.Fatal("Failed to validate non-existent code")
		}
	})

	t.Run("New user", func(t *testing.T) {
		const phone = "+89998887766"

		testUser := models.User{
			ID:    1,
			Phone: phone,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_repositories.NewMockUsers(ctrl)
		usersRepo.EXPECT().GetUserByPhone(phone).Return(models.User{}, repositories.NotFoundErr)
		usersRepo.EXPECT().CreateUser(phone).Return(testUser, nil)

		auth := New(usersRepo)

		code, err := auth.SendCode(phone)

		res, err := auth.SignIn(phone, code)
		if err != nil {
			t.Fatal("Failed to SignIn", err)
		}

		if res == nil {
			t.Fatal("Failed to SignIn: empty response")
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		const phone = "+89998887766"

		testUser := models.User{
			ID:    1,
			Phone: phone,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_repositories.NewMockUsers(ctrl)
		usersRepo.EXPECT().GetUserByPhone(phone).Return(testUser, nil)

		auth := New(usersRepo)

		code, err := auth.SendCode(phone)

		res, err := auth.SignIn(phone, code)
		if err != nil {
			t.Fatal("Failed to SignIn", err)
		}

		if res == nil {
			t.Fatal("Failed to SignIn: empty response")
		}
	})
}
