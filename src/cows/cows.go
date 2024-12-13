package cows

import (
	"embed"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/vnykmshr/gowsay/src/utils"
)

var (
	// Embed all cow files
	//go:embed cows/*.txt
	cowsFS embed.FS
	// Use a single random source for the package
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Cow struct {
	Name string
	Art  string
}

// GetRandomCow returns a random Cow struct
func GetRandomCow() (*Cow, error) {
	cowName, err := GetRandomCowName()
	if err != nil {
		return nil, err
	}
	return GetCow(cowName)
}

// GetRandomCowName returns a random cow name from available cows
func GetRandomCowName() (string, error) {
	fileNames, err := GetCowNames()
	if err != nil {
		return "", fmt.Errorf("unable to get cow names: %w", err)
	}
	if len(fileNames) == 0 {
		return "", fmt.Errorf("no cow files available")
	}
	return fileNames[rnd.Intn(len(fileNames))], nil
}

// GetCowNames returns a list of all cow names (without extensions)
func GetCowNames() ([]string, error) {
	entries, err := cowsFS.ReadDir("cows")
	if err != nil {
		return nil, fmt.Errorf("unable to read cow directory: %w", err)
	}

	filenames := make([]string, len(entries))
	for i, e := range entries {
		filenames[i] = strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
	}

	return filenames, nil
}

// GetCow returns a Cow struct by cow name
func GetCow(cowName string) (*Cow, error) {
	cowArt, err := getCowArt(cowName)
	if err != nil {
		return nil, fmt.Errorf("unable to get cow %q: %w", cowName, err)
	}

	return &Cow{
		Name: cowName,
		Art:  cowArt,
	}, nil
}

// getCowArt reads the art for a given cow name
func getCowArt(name string) (string, error) {
	contents, err := cowsFS.ReadFile(fmt.Sprintf("cows/%s.txt", name))
	if err != nil {
		return "", fmt.Errorf("unable to read cow file %q: %w", name, err)
	}

	art := utils.RemoveBackticks(strings.TrimSpace(string(contents)))
	return art, nil
}
