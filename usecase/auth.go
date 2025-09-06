package usecase

import (
	"context"
	"fmt"
	"slices"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/walnuts1018/ipxe-manager/definitions"
)

func (u *Usecase) CheckAuthorization(ctx context.Context, accessToken string) error {
	introspection, err := u.authService.IntrospectToken(ctx, accessToken)
	if err != nil {
		return err
	}

	exp := synchro.Unix[tz.AsiaTokyo](introspection.Expiration, 0)
	if u.clock.Now().After(exp) {
		return fmt.Errorf("token expired: %w", definitions.ErrInvalidToken)
	}

	if !introspection.Active {
		return fmt.Errorf("token inactive: %w", definitions.ErrInvalidToken)
	}

	for _, allowed := range u.cfg.OAuth2.AllowedAudiences {
		if slices.Contains(introspection.Audience, allowed) {
			return nil
		}
	}

	return fmt.Errorf("no allowed audience found: %w", definitions.ErrInvalidToken)
}
