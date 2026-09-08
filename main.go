package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"language-refresher/internal/execution"
	"language-refresher/internal/exercises"
	"language-refresher/internal/languages"
	"language-refresher/internal/progress"
	"language-refresher/internal/sessions"
	"language-refresher/internal/settings"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

//go:embed content/python/*.json
var contentAssets embed.FS

func main() {
	repository, err := exercises.Load(contentAssets, "content/*/*-refresher.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := repository.LoadTranslations(contentAssets, "content/*/*.i18n.json"); err != nil {
		log.Fatal(err)
	}
	track, ok := repository.Track("python-core-refresher")
	if !ok {
		log.Fatal("trilha python-core-refresher não encontrada")
	}
	progressPath, err := progress.DefaultPath()
	if err != nil {
		log.Fatal(err)
	}
	session, err := sessions.New(track, progress.NewStore(progressPath))
	if err != nil {
		log.Fatal(err)
	}
	registry := languages.NewRegistry("python")
	settingsPath, err := settings.DefaultPath()
	if err != nil {
		log.Fatal(err)
	}
	app := NewApp(session, execution.NewRunner(), registry, repository, settings.NewStore(settingsPath))

	err = wails.Run(&options.App{
		Title:            "Language Refresher",
		Width:            1440,
		Height:           900,
		MinWidth:         640,
		MinHeight:        460,
		BackgroundColour: &options.RGBA{R: 12, G: 17, B: 26, A: 1},
		AssetServer:      &assetserver.Options{Assets: frontendAssets},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []interface{}{app},
		Windows:          &windows.Options{Theme: windows.Dark},
	})
	if err != nil {
		log.Fatal(err)
	}
}
