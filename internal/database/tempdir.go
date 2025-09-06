package database

import "os"

func (d *Database) createTempDir() error {
	temp, err := os.MkdirTemp("", "dict-down-temp-*")
	if err != nil {
		return err
	}

	d.tempDir = temp

	return nil
}

func (d *Database) cleanTempDir() {
	_ = os.RemoveAll(d.tempDir)
}
