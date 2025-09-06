package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (u *Usecase) GetIpxeScript(ctx context.Context) (*os.File, error) {
	osstr, err := u.OSRepository.GetOS(ctx)
	if err != nil {
		return nil, err
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
		f, err := os.Open(filepath.Join(u.cfg.Ipxe.ScriptDir, fmt.Sprintf("%s.ipxe", u.cfg.Ipxe.DefaultOS)))
		if err != nil {
			return nil, fmt.Errorf("failed to open default OS script: %w", err)
		}
		return f, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open script: %w", err)
	}

	return f, nil
}
