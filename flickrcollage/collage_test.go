package flickrcollage

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"testing"

	_ "image/png"
)

type fakeFlickr struct{}
type fakeErrorFlickr struct{}

var dummyImage image.Image

func downloadDummyImage() {
	resp, err := http.Get("https://dummyimage.com/600x400/000/fff")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	dummyImage, _, err = image.Decode(resp.Body)

	if err != nil {
		panic(err)
	}
}

func (f fakeFlickr) fetchTopRatedPhoto(term string, fakePhotos chan<- image.Image, fakeErrors chan<- error) {
	if dummyImage == nil {
		downloadDummyImage()
	}

	for i := 0; i < CollageSize; i++ {
		fakePhotos <- dummyImage
	}
}

func (f fakeErrorFlickr) fetchTopRatedPhoto(term string, fakePhotos chan<- image.Image, fakeErrors chan<- error) {
	if dummyImage == nil {
		downloadDummyImage()
	}

	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			fakePhotos <- dummyImage
		} else {
			fakeErrors <- fmt.Errorf("Simulated error when fetching image")
		}
	}
}

func reset() {
	dictionary = []string{}
	searchTermsInUse = map[string]bool{}
}
func TestRandomWord(t *testing.T) {
	defer reset()

	actual := randomWord()

	if actual == "" {
		t.Errorf("Expected a random word from the dictionary.")
	}
}

func TestCreateGeneratesCollage(t *testing.T) {
	defer reset()

	testFilePath := "underTest.jpg"
	underTest := &FlickrCollage{}
	underTest.Create(fakeFlickr{}, []string{}, 1920, testFilePath, "fake")

	defer os.Remove(testFilePath)

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("Expected Create to generate an output file")
	}
}

func TestCreateFillsSearchTerms(t *testing.T) {
	defer reset()

	testFilePath := "underTest.jpg"
	underTest := &FlickrCollage{}
	underTest.Create(fakeFlickr{}, []string{"badger", "badger", "badger", "badger", "badger", "badger", "mushroom", "mushroom"}, 1920, testFilePath, "fake")

	defer os.Remove(testFilePath)

	if len(searchTermsInUse) < CollageSize {
		t.Errorf("Expected Create to search for at least %d unique terms. Got: %d", CollageSize, len(searchTermsInUse))
	}

	if found := searchTermsInUse["badger"]; !found {
		t.Errorf("Expected badger in search terms")
	}

	if found := searchTermsInUse["mushroom"]; !found {
		t.Errorf("Expected mushroom in search terms")
	}
}

func TestCollageRetriesOnFailure(t *testing.T) {
	defer reset()

	testFilePath := "underTest.jpg"
	underTest := &FlickrCollage{}
	underTest.Create(fakeErrorFlickr{}, []string{}, 1920, testFilePath, "fake")

	defer os.Remove(testFilePath)

	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Errorf("Expected Create to generate an output file")
	}
}
