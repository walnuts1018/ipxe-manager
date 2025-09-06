package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/walnuts1018/ipxe-manager/definitions"
)

func (u *Usecase) SetOS(ctx context.Context, osstr string) error {
	// check existence
	if _, err := os.Stat(filepath.Join(u.cfg.Ipxe.ScriptDir, fmt.Sprintf("%s.ipxe", osstr))); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("os script not found: %w", definitions.ErrIPXEScriptNotFound)
		}
		return err
	}

	if err := u.OSRepository.SetOS(ctx, osstr); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) GetOS(ctx context.Context) (string, error) {
	osstr, err := u.OSRepository.GetOS(ctx)
	if err != nil {
		return "", err
	}

	path := filepath.Join(u.cfg.Ipxe.ScriptDir, fmt.Sprintf("%s.ipxe", osstr))

	useDefault := false
	if osstr == "" {
		useDefault = true
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		useDefault = true
	}

	if useDefault {
		// Fallback to default OS
		return u.cfg.Ipxe.DefaultOS, nil
	}

	return osstr, nil
}
