package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal(msg)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func strToInt(key string) (int, error) {
	s := os.Getenv(key)
	v, err := strconv.Atoi(s)

	if err != nil {
		return 0, err
	}

	return v, nil
}