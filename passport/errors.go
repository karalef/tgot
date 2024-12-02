package passport

import (
	"context"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tgpassport"
)

// SetPassportDataErrors informs a user that some of the Telegram Passport elements they provided contains errors.
// The user will not be able to re-submit their Passport to you until the errors are fixed.
func SetPassportDataErrors(ctx context.Context, a *api.API, userID int, errs []tgpassport.PassportElementError) error {
	data := api.NewData()
	data.SetInt("user_id", userID)
	data.SetJSON("errors", errs)
	return a.Request(ctx, "setPassportDataErrors", data)
}
