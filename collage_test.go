package main

import (
	"os"
	"testing"
)

func TestMainPanicsWithoutAPIKey(t *testing.T) {
	oldAPIKey := os.Getenv("FLICKR_API_KEY")

	defer func() {
		os.Setenv("FLICKR_API_KEY", oldAPIKey)

		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	os.Setenv("FLICKR_API_KEY", "")

	main()
}
