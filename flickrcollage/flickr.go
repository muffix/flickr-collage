package flickrcollage

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
)

// flickr is the flickr interface for this package
type flickr interface {
	fetchTopRatedPhoto(string, chan<- image.Image, chan<- error)
}

const baseAPI = "https://api.flickr.com/services/rest/?format=json&nojsoncallback=1"
const searchAPI = baseAPI + "&method=flickr.photos.search&sort=interestingness-desc&per_page=1&page=1"
const sizeAPI = baseAPI + "&method=flickr.photos.getSizes"

// FlickrAPI is the Flickr API object
type FlickrAPI string

type photo struct {
	ID string `json:"id"`
}

type photoResult struct {
	Photos []photo `json:"photo"`
}

type photoSize struct {
	ImageURL string `json:"source"`
}

type availableSizes struct {
	Available []photoSize `json:"size"`
}

type searchResponse struct {
	PhotoResult photoResult `json:"photos"`
	Status      string      `json:"stat"`
}

type sizesResponse struct {
	Sizes  availableSizes `json:"sizes"`
	Status string         `json:"stat"`
}

func (f FlickrAPI) fetchTopRatedPhoto(searchTerm string, photos chan<- image.Image, errors chan<- error) {
	fmt.Printf("Searching for %s...\n", searchTerm)
	id, err := f.fetchPhotoID(searchTerm)

	if err != nil {
		errors <- err
		return
	}
	fmt.Printf("Found photo with ID %s for %s\n", id, searchTerm)

	image, err := f.fetchPhoto(id)

	if err != nil {
		fmt.Println(err)
		errors <- err
		return
	}

	photos <- image
}

// fetchPhotoID retrieves the ID of a the top-rated photo for the given search term
func (f FlickrAPI) fetchPhotoID(searchTerm string) (string, error) {
	searchResp := &searchResponse{}
	if err := f.makeRequest(searchAPI+"&text="+searchTerm, searchResp); err != nil {
		return "", err
	}

	if len(searchResp.PhotoResult.Photos) < 1 {
		fmt.Printf("Nothing found for %s.\n", searchTerm)
		return "", fmt.Errorf("Couldn't find a photo for %s", searchTerm)
	}

	return searchResp.PhotoResult.Photos[0].ID, nil
}

// fetchPhoto downloads the photo with the given ID and returns an Image
func (f FlickrAPI) fetchPhoto(id string) (image.Image, error) {
	sizesResp := &sizesResponse{}

	if err := f.makeRequest(sizeAPI+"&photo_id="+id, sizesResp); err != nil {
		return nil, err
	}

	if sizesResp.Status != "ok" {
		return nil, fmt.Errorf("API call failed with status: %s", sizesResp.Status)
	}

	if len(sizesResp.Sizes.Available) == 0 {
		return nil, fmt.Errorf("No image found for ID " + id)
	}

	resp, err := http.Get(sizesResp.Sizes.Available[len(sizesResp.Sizes.Available)/2].ImageURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	image, _, err := image.Decode(resp.Body)

	if err != nil {
		return nil, err
	}

	return image, err
}

// makeRequest makes a GET request to the Flickr API and Unmarshals it into the given struct
func (f FlickrAPI) makeRequest(url string, responseStruct interface{}) error {
	url = fmt.Sprintf("%s&api_key=%s", url, f)

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, responseStruct); err != nil {
		return err
	}

	return nil
}
