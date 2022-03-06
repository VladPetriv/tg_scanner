package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

type noSignUp struct{}

type TermAuth struct {
	noSignUp

	UserPhone string
}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func (a TermAuth) Phone(_ context.Context) (string, error) {
	return os.Getenv("PHONE"), nil
}

func (a TermAuth) Password(_ context.Context) (string, error) {
	return os.Getenv("PASSWORD"), nil
}

func (a TermAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func Login(ctx context.Context, client *telegram.Client, cfg config.Config) (*tg.AuthAuthorization, error) {
	//Create new flow
	flow := auth.NewFlow(
		TermAuth{UserPhone: cfg.Phone},
		auth.SendCodeOptions{},
	)
	//Authorization
	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		return nil, err
	}

	//Authorization with password
	password, _ := flow.Auth.Password(ctx)

	user, err := client.Auth().Password(ctx, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
