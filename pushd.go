package main

import "os"

type popdir func() error

func pushdir(dir string) (popdir, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	return func() error { return os.Chdir(pwd) }, nil
}
