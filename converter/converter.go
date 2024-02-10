package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	// "io"

	// "golang.org/x/oauth2"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	// "golang.org/x/oauth2/spotify"
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	applemusic "github.com/minchao/go-apple-music"
	// "maps"
)

type song_info struct {
	name   string // the name of the song
	artist string // the name of the first listed artist
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
func GetSpotifyData(UserId, ClientId, ClientSecret string) map[string]map[string]song_info {
	config := &clientcredentials.Config{
		ClientID:     ClientId,
		ClientSecret: ClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	accessToken, err := config.Token(context.Background())
	if err != nil {
		fmt.Println("Error getting token:", err)
	}
	client := spotify.Authenticator{}.NewClient(accessToken)
	result, err := client.GetPlaylistsForUser(UserId)
	if err != nil {
		fmt.Println("Error getting playlist:", err)
	}
	lib := make(map[string]map[string]song_info)

	for _, playlist := range result.Playlists {
		log.Println(playlist.Name)
		tracks, err := client.GetPlaylistTracks(playlist.ID)
		if err != nil {
			fmt.Println("Error getting playlist:", err)
		}
		lib[playlist.Name] = make(map[string]song_info)
		for _, track := range tracks.Tracks {
			// log.Println(track.Track.Name)
			// log.Println(track.Track.Artists)
			name := track.Track.Artists[0]
			// log.Println(name.Name)
			song := song_info{name: track.Track.Name, artist: name.Name}
			// print(song.artist)

			lib[playlist.Name][track.Track.Name] = song

		}
	}
	return lib
}
func GetAppleMusicData(musicUserToken string, devToken string) map[string]map[string]song_info {
	ctx := context.Background()
	tp := applemusic.Transport{
		Token:          devToken,
		MusicUserToken: musicUserToken,
	}
	client := applemusic.NewClient(tp.Client())

	// Fetch all the storefronts in alphabetical order
	// storefronts, _, err := client.Storefront.GetAll(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(storefronts)
	// albums, _, err := client.Me.GetAllLibraryPlaylists(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	// for _, album := range albums.Data {
	// 	fmt.Println(album.Attributes.Name)
	// }
	// // fmt.Println(resp)
	// fmt.Println(albums)
	// client.
	// fmt.Println(*client.Me)
	lib := make(map[string]map[string]song_info)
	playlists, resp, err := client.Me.GetAllLibraryPlaylists(ctx, nil)
	fmt.Println(playlists, resp, err)
	// fmt.Printf("playlists: %T\n", playlists)
	for _, playlist := range playlists.Data {
		// fmt.Printf("playlist: %T\n", playlist)
		fmt.Println(playlist.Id)
		fmt.Println(playlist.Attributes.Name)
		lib[playlist.Attributes.Name] = make(map[string]song_info)
		tracks, err := client.Me.GetLibraryPlaylistCatalogTracks(ctx, playlist.Id, 250)
		if err != nil {
			log.Println(tracks, err)
		}
		for _, track := range tracks {
			// fmt.Println(track.Attributes.Name)
			// fmt.Println(track.Attributes.ArtistName)
			song := song_info{name: track.Attributes.Name, artist: track.Attributes.ArtistName}
			fmt.Println(song)
			lib[playlist.Attributes.Name][track.Attributes.Name] = song
		}
		// break
	}
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(resp)
	// log.Println(playlists)
	return lib

}
func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic(err)
	}
	// spotifyUserId := os.Getenv("spotifyUserId")
	// spotifyClientId := os.Getenv("SpotifyClientId")
	// spotifyClientSecret := os.Getenv("SpotifyClientSecret")

	// libSpotify := GetSpotifyData(spotifyUserId, spotifyClientId, spotifyClientSecret)
	// for id, _ := range libSpotify {
	// 	fmt.Println(id)
	// }
	devToken := os.Getenv("AppleMusicDevToken")
	userToken := os.Getenv("AppleMusicUserToken")

	libApple := GetAppleMusicData(userToken, devToken)
	keys := make([]string, 0, len(libApple))
	for k, _ := range libApple {
		keys = append(keys, k)
	}
	fmt.Println(keys)

}
