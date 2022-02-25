package data4all

import (
	"os"
)

func Finished() {
	<-threads
}

func WriteFile(path string, content []byte) error {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(string(content) + "\n"); err != nil {
		return err
	}

	return nil

}
