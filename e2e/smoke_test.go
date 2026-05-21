//go:build e2e

package e2e

import (
	"context"
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("account_status", func(t *testing.T) {
		acct, err := client.GetAccountStatus(ctx)
		requireNoError(t, err, "get account status")
		requireNotEmpty(t, acct.ID, "account id is empty")
		requireNotEmpty(t, acct.Email, "account email is empty")
		requireNotEmpty(t, acct.DefaultWorkspace, "default workspace is empty")
		t.Logf("authenticated as %s (%s)", acct.Email, acct.DefaultWorkspace)
	})
}
