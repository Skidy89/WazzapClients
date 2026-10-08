package command

// runGhost turns ghost mode on. It's turned off where the composer was.
func runGhost(c *Context) error {
	c.SetGhost(true)
	c.Note(&Note{Title: c.Input, Text: c.Locale.Text("ghost.mode.on")})
	return nil
}
