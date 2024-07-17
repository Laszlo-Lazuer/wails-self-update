package main

import (
	"embed"
	"io/ioutil"

	"encoding/json"
	"fmt"
	"io"

	"github.com/inconshreveable/go-update"

	// "io/ioutil"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

const repoURL = "https://api.github.com/repos/yourusername/yourrepo/releases/latest"

type Release struct {
	TagName string `json:"tag_name"`
}

// type App struct {
// 	ctx context.Context
// }

// func NewApp() *App {
// 	return &App{}
// }

func (a *App) getLatestRelease() (string, error) {
	resp, err := http.Get(repoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func (a *App) downloadAndUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tempFile, err := ioutil.TempFile("", "update")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	err = tempFile.Close()
	if err != nil {
		return err
	}

	if err := update.Apply(tempFile, update.Options{}); err != nil {
		if rollbackErr := update.RollbackError(err); rollbackErr != nil {
			return fmt.Errorf("update failed: %v, rollback failed: %v", err, rollbackErr)
		}
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
}

func (a *App) CheckForUpdates() {
	latestRelease, err := a.getLatestRelease()
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	currentVersion := "v1.0.0" // Replace with your app's current version
	if latestRelease != currentVersion {
		fmt.Println("New version available:", latestRelease)
		// URL to the binary file of the latest release
		downloadURL := fmt.Sprintf("https://github.com/yourusername/yourrepo/releases/download/%s/yourapp", latestRelease)
		if err := a.downloadAndUpdate(downloadURL); err != nil {
			fmt.Println("Error updating application:", err)
			return
		}
		fmt.Println("Application updated successfully. Restarting...")
		os.Exit(0)
	} else {
		fmt.Println("You are already using the latest version.")
	}
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "crossword",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
