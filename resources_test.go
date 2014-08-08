package hat

type Root struct {
	Hello string
	Apps  Apps `hat:"embed()"`
}

type Apps map[string]App

type App struct {
	Name     string
	Versions Versions `hat:"embed()"`
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

func (entity *Apps) Manifest(parent *Root, _ string) error {
	(*entity) = Apps{
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
	return nil
}

func (entity *App) Manifest(parent *Apps, id string) error {
	if app, ok := parent[id]; ok {
		(*entity) = app
	}
	return nil
}

func (entity *Versions) Manifest(parent *App, _ string) error {
	(*entity) = parent.Versions
	return nil
}

func (entity *Version) Manifest(parent *Versions, id string) error {
	if version, ok := parent[id]; ok {
		(*entity) = version
	}
	return nil
}
