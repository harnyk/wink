package auth

import (
	"fmt"

	"github.com/harnyk/wink/internal/cryptostore"
	"github.com/harnyk/wink/internal/entities"
	api "github.com/harnyk/wink/internal/peopleapi"
	"github.com/harnyk/wink/internal/ui"
)

type AuthPrompt interface {
	Get() (api.Auth, error)
}

type authPrompt struct {
	cachedAuth     api.Auth
	configFileName string
}

func NewAuthPrompt(configFileName string) AuthPrompt {
	return &authPrompt{
		configFileName: configFileName,
	}
}

func (a *authPrompt) Get() (api.Auth, error) {
	if a.cachedAuth.APIKey != "" && a.cachedAuth.EmployeeID != "" {
		return a.cachedAuth, nil
	}

	store := cryptostore.NewCryptoStore[entities.Secrets](a.configFileName)
	u := ui.NewUI()

	password, err := u.AskPassword("Please enter the password:")
	if err != nil {
		return api.Auth{}, err
	}

	record, err := store.Load(string(password))
	if err != nil {
		return api.Auth{}, err
	}

	fmt.Println("Credentials loaded")

	a.cachedAuth = api.Auth{
		APIKey:     record.APIKey,
		EmployeeID: record.EmployeeID,
	}

	return a.cachedAuth, nil
}
