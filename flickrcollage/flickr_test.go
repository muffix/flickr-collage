package flickrcollage

import (
	"image"
	"os"
	"testing"
)

func TestFlickrIntegration(t *testing.T) {
	underTest := FlickrAPI(os.Getenv("FLICKR_API_KEY"))

	photosChan := make(chan image.Image)
	errorChan := make(chan error)

	go underTest.fetchTopRatedPhoto("car", photosChan, errorChan)

	actual := <-photosChan

	if actual == nil {
		t.Errorf("Extected the FlickrAPI to return an image")
	}
}
