# Flickr Collage CLI

Simple command line interface for Flickr. It builds a collage of top-rated images for 10 search terms. 

## Getting Started

Flickr requires an API key which can be obtained in the [Account Settings](https://www.flickr.com/services/api/keys/). For tests, the key from the [API test page](https://www.flickr.com/services/api/explore/flickr.photos.search) will be good enough.

The CLI expects 10 words as the trailing positional arguments. If fewer than 10 are provided or any single one fails to return results or download, it will add a random replacement search term.

### Prerequisites

This CLI depends on [Picasso](https://github.com/deiwin/picasso). To install, run 

```
go get github.com/deiwin/picasso
```

It makes use of the word list that can be found in `/usr/share/dict/words` on Unix systems.

### Installing

After cloning the repository and installing the [prerequisites](#prerequisites), just compile the binary by running `go build collage.go` in the root directory.


## Example

```$ ./collage coffee car football watch water git computer hockey imnotagoodsearchterm
Searching for git...
Searching for imnotagoodsearchterm...
Searching for computer...
Searching for hockey...
Searching for watch...
Searching for football...
Searching for water...
Searching for car...
Searching for coffee...
Searching for rental...
Nothing found for imnotagoodsearchterm.
Found photo with ID 3417057229 for git
Searching for Buddhic...
Found photo with ID 36441297532 for coffee
Found photo with ID 5347580266 for computer
Found photo with ID 35026760141 for rental
Found photo with ID 15785619127 for watch
Found photo with ID 16104874480 for hockey
Found photo with ID 6769482215 for football
Found photo with ID 9182829574 for Buddhic
Found photo with ID 23958307713 for car
Found photo with ID 26223084353 for water
Found 10 photos.
```

#### Command line arguments

| Argument      | Description                        | Required | Default                  |
| ------------- |----------------------------------- | :------: | ------------------------ |
| `--api-key`   | Flickr API key                     | Y        | Env var `FLICKR_API_KEY` |
| `--output`    | Output path                        | N        | `collage.jpg`            |
| `--width`     | Desired width of the collage in px | N        | 800                      |

## Running the tests

After cloning the repository and installing the [prerequisites](#prerequisites), just run `go test` in the directory of the respective package.

## Built With

* [Go 1.10.3 ](https://golang.org/)

