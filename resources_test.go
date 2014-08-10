package hat

import (
	"net/http"
	"testing"
)

func TestHat(t *testing.T) {
	if s, err := NewServer(Root{}); err != nil {
		t.Error(err)
	} else {
		println("Listening on :8080")
		t.Fatal(http.ListenAndServe(":8080", s))
	}
}

type Root struct {
	Hello  string
	Apps   Apps   `hat:"embed()"`
	Health Health `hat:"link()"`
}

type Health struct {
	Hello string
}

type Apps map[string]App

func (entity *Apps) Manifest(_ *Root, _ string) error {
	(*entity) = the_apps
	return nil
}

// func (*Apps) Page(_ *Root, _ string, number int, maxItems int) ([]string, error) {
// 	ids := []string{}
// 	for k, _ := range the_apps {
// 		ids := append(ids, k)
// 	}
// 	return ids, nil
// }

type App struct {
	Name     string
	Versions Versions `hat:"embed()"`
}

func (entity *App) Manifest(parent *Apps, id string) error {
	if app, ok := (*parent)[id]; ok {
		(*entity) = app
	}
	return nil
}

type Versions map[string]Version

type Version struct {
	ID      string
	Version string
	Date    string
}

func (entity *Root) Manifest(_ interface{}, _ string) error {
	(*entity) = Root{
		Hello: "Wecome to the test API.",
	}
	return nil
}

func (entity *Health) Manifest(_ *Root, _ string) error {
	(*entity) = Health{
		Hello: "Feelin' good!",
	}
	return nil
}

var the_apps = Apps{
	"test-app": App{
		"Test App",
		Versions{
			"0.0.1": Version{"test-app-v0-0-1", "0.0.1", "May 2013"},
			"0.0.2": Version{"test-app-v0-0-2", "0.0.2", "August 2014"},
		},
	},
	"other-app": App{
		"Other App",
		Versions{
			"0.1.0": Version{"other-app-v0-0-1", "0.1.0", "June 2014"},
			"0.4.0": Version{"other-app-v0-0-2", "0.4.0", "July 2014"},
		},
	},
}

func (entity *Versions) Manifest(parent *App, _ string) error {
	(*entity) = parent.Versions
	return nil
}

func (entity *Version) Manifest(parent *Versions, id string) error {
	if version, ok := (*parent)[id]; ok {
		(*entity) = version
	}
	return nil
}
