package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	applemusic "github.com/minchao/go-apple-music"
)

type config struct {
	DeveloperToken string `json:"developer_token"`
	UserToken      string `json:"user_token"`
}

type playlist struct {
	Playlist applemusic.LibraryPlaylist
	Songs    []applemusic.Song
}

func main() {
	configFile := flag.String("config", "", "Config file location (default ~/.amp.json)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || (args[0] != "backup" && args[0] != "restore") {
		fmt.Fprintln(os.Stderr, "Usage: amp [-config file] <backup|restore> [options]")
		return
	}

	action := args[0]

	if *configFile == "" {
		*configFile = os.ExpandEnv("$HOME/.amp.json")
	}

	var cfg config
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if cfg.DeveloperToken == "" {
		cfg.DeveloperToken = os.Getenv("DEVELOPER_TOKEN")
	}
	if cfg.UserToken == "" {
		cfg.UserToken = os.Getenv("USER_TOKEN")
	}

	if cfg.DeveloperToken == "" || cfg.UserToken == "" {
		log.Fatal("Developer token and user token must be provided either in the config file or as environment variables.")
	}

	ctx := context.Background()
	tp := applemusic.Transport{
		Token:          cfg.DeveloperToken,
		MusicUserToken: cfg.UserToken,
	}
	client := applemusic.NewClient(tp.Client())

	switch action {
	case "backup":
		playlists, _, err := client.Me.GetAllLibraryPlaylists(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to get playlists: %v", err)
		}
		log.Println("found", len(playlists.Data), "playlists")
		for _, pl := range playlists.Data {
			songs, _, err := client.Me.GetLibraryPlaylistTracks(ctx, pl.Id, &applemusic.PageOptions{
				Limit: 0,
			})
			if err != nil {
				log.Fatalf("Failed to get songs for playlist %s: %v", pl.Id, err)
			}
			file := fmt.Sprintf("playlist.%s.json", pl.Id)
			log.Println("writing playlist", pl.Attributes.Name, "to", file)
			f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("Failed to open file: %v", err)
			}
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			err = enc.Encode(playlist{
				Playlist: pl,
				Songs:    songs,
			})
			if err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
			f.Close()
		}
	case "restore":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: amp [-config file] restore <file> [optional new playlist name]")
			return
		}
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		var pl playlist
		if err := json.Unmarshal(data, &pl); err != nil {
			log.Fatalf("Failed to parse JSON: %v", err)
		}
		var name string = pl.Playlist.Attributes.Name
		if len(args) > 2 && args[2] != "" {
			name = args[2]
		}
		var description string
		if pl.Playlist.Attributes.Description != nil {
			description = pl.Playlist.Attributes.Description.Standard
		}
		var tracks []applemusic.CreateLibraryPlaylistTrack
		for _, song := range pl.Songs {
			tracks = append(tracks, applemusic.CreateLibraryPlaylistTrack{
				Id:   song.Id,
				Type: song.Type,
			})
		}
		playlists, _, err := client.Me.CreateLibraryPlaylist(ctx, applemusic.CreateLibraryPlaylist{
			Attributes: applemusic.CreateLibraryPlaylistAttributes{
				Name:        name,
				Description: description,
			},
			Relationships: &applemusic.CreateLibraryPlaylistRelationships{
				Tracks: applemusic.CreateLibraryPlaylistTrackData{
					Data: tracks,
				},
			},
		}, nil)
		if err != nil {
			log.Fatalf("Failed to create playlist: %v", err)
		}
		log.Println("created playlist [ name:", playlists.Data[0].Attributes.Name, ", id:", playlists.Data[0].Id, "] with", len(tracks), "tracks")
	}
}
