package data

import "os"

func (d *DataLoader) createTempDir() error {
	temp, err := os.MkdirTemp("", "dict-down-temp-*")
	if err != nil {
		return err
	}

	d.tempDir = temp

	return nil
}

func (d *DataLoader) cleanTempDir() {
	_ = os.RemoveAll(d.tempDir)
}
