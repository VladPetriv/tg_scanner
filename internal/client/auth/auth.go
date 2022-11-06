package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/pkg/config"
)

type noSignUp struct{}

type TermAuth struct {
	noSignUp

	UserPhone string
}

var ErrNotImplemented = errors.New("not implemented")

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, ErrNotImplemented
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
	fmt.Println("Enter code: ")

	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("create new reader error: %w", err)
	}

	return strings.TrimSpace(code), nil
}

func Login(ctx context.Context, client *telegram.Client, cfg *config.Config) (*tg.AuthAuthorization, error) {
	// Create new flow
	flow := auth.NewFlow(
		TermAuth{noSignUp: noSignUp{}, UserPhone: cfg.Phone},
		auth.SendCodeOptions{},
	)
	// Authorization
	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	// Authorization with password
	password, _ := flow.Auth.Password(ctx)

	user, err := client.Auth().Password(ctx, password)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	return user, nil
}
