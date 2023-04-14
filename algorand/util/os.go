package util

import (
	"os"
	"os/user"
)

func GetUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	if os.Geteuid() == 0 {
		// Root, try to retrieve SUDO_USER if exists
		if u := os.Getenv("SUDO_USER"); u != "" {
			usr, err = user.Lookup(u)
			if err != nil {
				return nil, err
			}
		}
	}

	return usr, nil
}

func GetUserHomedir() string {
	home, err := GetUser()
	ENOK(err)
	return home.HomeDir
}

func VerifyFileExistence(file string) error {
	_, err := os.Stat(file)
	return err
}
