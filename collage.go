package main

import (
	"flag"
	_ "image/jpeg"
	"os"

	"github.com/muffix/flickr-collage-cli/flickrcollage"
)

func main() {
	apiKey := flag.String("api-key", os.Getenv("FLICKR_API_KEY"), "Flickr API key")
	outputPath := flag.String("output", "collage.jpg", "Output file path")
	collageWidth := flag.Int("width", 800, "Desired width of the collage")

	flag.Parse()

	if *apiKey == "" {
		panic("You need to specify a Flickr API key. Use --api-key or the FLICKR_API_KEY env var")
	}

	flickr := flickrcollage.FlickrAPI(*apiKey)
	flickrCollage := &flickrcollage.FlickrCollage{}
	flickrCollage.Create(flickr, flag.Args(), *collageWidth, *outputPath, *apiKey)
}
