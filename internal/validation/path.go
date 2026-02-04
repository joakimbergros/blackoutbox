// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package validation

import (
	"errors"
	"fmt"
	"os"
)

func ValidatePrintableFile(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("file is a symlink")
	}

	if !info.Mode().IsRegular() {
		return errors.New("not a regular file")
	}

	if info.Size() == 0 {
		return errors.New("file is empty")
	}

	const maxSize = 20 << 20 // 20 MB (adjust as needed)
	if info.Size() > maxSize {
		return errors.New("file too large to print")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	f.Close()

	return nil
}
