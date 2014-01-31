
package main

import (
  "os"
  "net/http"
  "strconv"
  "io/ioutil"
  "encoding/json"
  "net/url"
  "fmt"
)

type Artist struct {
  Name string
}

type Track struct {
  Name string `json:"name"`
  Artist Artist
}

type Meta struct {
  TotalPages string `json:"totalPages"`
}


type Tracks struct {
  Tracks []Track `json:"track"`
  Meta Meta `json:"@attr"`
}

type UserLovedTracks struct {
  Tracks Tracks `json:"lovedtracks"`
}

type DeezerTrack struct{
  Id int
  Title string
  Artist Artist
}

type DeezerSearchTrack struct {
  Data []DeezerTrack `json:"data"`
}

func getTracks(page int, user_name string) ([]Track, Meta) {
  lovedTracksURL := "http://ws.audioscrobbler.com/2.0/?limit=100&method=user.getlovedtracks&user="+user_name+"&api_key=a6deb27fe252484a367a729f2c85a18b&format=json&page="+strconv.Itoa(page)

  response, err := http.Get(lovedTracksURL)
  if err != nil {
    panic(err)
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    panic(err)
  }
  var loved UserLovedTracks
  json.Unmarshal(contents, &loved)
  return loved.Tracks.Tracks, loved.Tracks.Meta
}

func searchTrack(artist string, name string) {
  deezer := "http://api.deezer.com/search"
  deezerUrl, _ := url.Parse(deezer)

	q := deezerUrl.Query()
	q.Set("q", name + " " + artist)
	deezerUrl.RawQuery = q.Encode()

  response, err := http.Get(deezerUrl.String())
  if err != nil {
    panic(err)
  }
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
    panic(err)
  }

  var tracks DeezerSearchTrack
  json.Unmarshal(contents, &tracks)

  fmt.Println("Found to", name, " - ", artist)
  fmt.Println(deezerUrl)
  for _, track := range tracks.Data {
    fmt.Println("     ->", track.Title, " - ", track.Artist.Name)
  }
}

func processLastTracks(tracks []Track){
  for _, track := range tracks {
    searchTrack(track.Artist.Name, track.Name)
  }
}

func main(){
  last_user_name :=  os.Getenv("LAST_USER_NAME")
  tracks, meta := getTracks(1, last_user_name)
  processLastTracks(tracks)

  pages, _ := strconv.Atoi(meta.TotalPages)

  for page := 2; page <= pages; page++ {
    fmt.Println("Search Page:", page)
    tracks, meta = getTracks(page, last_user_name)
    processLastTracks(tracks)
  }
}
