// Package flickrcollage contains a simple interface to get the highest-rated image for up to 10 search terms and creates a collage
package flickrcollage

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/deiwin/picasso"
)

// CollageSize contains the number of photos to use for the collage
const CollageSize = 10

var searchTermsInUse = make(map[string]bool)
var r *rand.Rand
var dictionary []string

// Collage is the interface to make a collage of photos using the top-rated photo for each search term.
// If fewer terms than the minimum collage size are specified, it uses the system dictionary to fill up.
type Collage interface {
	Create(flickrAPI flickr, searchTerms []string, width int, outputPath, apiKey string)
}

// FlickrCollage can make a collage from Flickr photos
type FlickrCollage struct{}

// randomWord returns a random word from the system dictionary
func randomWord() string {
	if len(dictionary) == 0 {
		buildWordsList()
	}

	return string(dictionary[r.Intn(len(dictionary))])
}

// newSearchTerm returns a new random search term that hasn't been used before
func newSearchTerm() string {
	var searchTermCandidate string
	for ; searchTermsInUse[searchTermCandidate]; searchTermCandidate = randomWord() {
	}
	searchTermsInUse[searchTermCandidate] = true
	return searchTermCandidate
}

// buildWordsList builds reads the system dictionary into memory
func buildWordsList() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	content, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}

	dictionary = strings.Split(string(content), "\n")
}

// Create makes a collage of photos from Flickr using the top-rated photo for each search term.
func (fc FlickrCollage) Create(flickrAPI flickr, searchTerms []string, width int, outputPath, apiKey string) {
	photos := fetchPhotos(flickrAPI, searchTerms)
	collage := makeCollage(photos, width)
	saveImageToFile(collage, outputPath)
}

// saveImageToFile saves the given mage to the given output path
func saveImageToFile(image image.Image, outputPath string) {
	file, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, image, nil); err != nil {
		panic(err)
	}

	fmt.Printf("Collage written to %s.\n", outputPath)
}

// makeCollage takes a slice of Images and combines them into a single collage Image with the desired width
func makeCollage(images []image.Image, width int) image.Image {
	grey := color.RGBA{0xaf, 0xaf, 0xaf, 0xff}
	return picasso.DrawGridLayoutWithBorder(images, width, grey, 2)
}

// fetchPhotos downloads the top-rated photo for each search term and returns a channel for the goroutines to deliver them
func fetchPhotos(f flickr, searchTerms []string) []image.Image {
	photos := make(chan image.Image, CollageSize)
	errors := make(chan error)

	for _, term := range searchTerms {
		if used := searchTermsInUse[term]; !used {
			searchTermsInUse[term] = true
			go f.fetchTopRatedPhoto(term, photos, errors)
		}
	}

	for i := len(searchTermsInUse); i < CollageSize; i++ {
		go f.fetchTopRatedPhoto(newSearchTerm(), photos, errors)
	}

	images := make([]image.Image, 0, CollageSize)

	for len(images) < CollageSize {
		select {
		case img := <-photos:
			images = append(images, img)
		case err := <-errors:
			fmt.Println(err)
			go f.fetchTopRatedPhoto(newSearchTerm(), photos, errors)
		}
	}

	fmt.Printf("Found %d photos.\n", len(images))

	return images
}
