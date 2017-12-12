package control

// App represents an app on the TV
type App struct {
	Name string
	ID   string
	tv   *LgTv
}

// Launch launches the app on the TV. It returns the ID of the new session
func (app *App) Launch() (string, error) {
	return app.tv.LaunchApp(app.ID)
}
