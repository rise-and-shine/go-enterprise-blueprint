//nolint:forbidigo // using fmt.Printf is allowed for CLI commands
package cli

import (
	"bufio"
	"context"
	"fmt"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/createsuperadmin"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"golang.org/x/term"
)

func (c *Controller) CreateSuperadminCmd() error {
	const (
		executionTimeout = 30 * time.Second
	)

	reader := bufio.NewReader(os.Stdin)

	username, err := askUsername(reader)
	if err != nil {
		return errx.Wrap(err)
	}

	password, err := askPassword()
	if err != nil {
		return errx.Wrap(err)
	}

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	// Set trace ID to context
	ctx = context.WithValue(ctx, meta.TraceID, tracing.GetStartingTraceID(ctx))

	err = c.usecaseContainer.CreateSuperadmin().Execute(ctx, createsuperadmin.Input{
		Username: username,
		Password: password,
	})

	return errx.Wrap(err)
}

func askUsername(reader *bufio.Reader) (string, error) {
	const (
		minUsernameLen = 3
		maxUsernameLen = 30
	)

	for {
		fmt.Print("Enter username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", errx.Wrap(err)
		}
		username := strings.TrimSpace(input)

		if len(username) < minUsernameLen {
			fmt.Printf("Username must be at least %d characters\n", minUsernameLen)
			continue
		}
		if len(username) > maxUsernameLen {
			fmt.Printf("Username must be at most %d characters\n", maxUsernameLen)
			continue
		}

		return username, nil
	}
}

func askPassword() (string, error) {
	const (
		minPasswordLen = 5
	)

	for {
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return "", errx.Wrap(err)
		}
		password := string(passwordBytes)

		if len(password) < minPasswordLen {
			fmt.Printf("Password must be at least %d characters\n", minPasswordLen)
			continue
		}

		return password, nil
	}
}
