package handlers

// tabView is one video link on the home tab bar. The first tab is always the
// "próximo" (newest) — project invariant.
type tabView struct {
	Slug    string
	Title   string
	Channel string
	Active  bool
	Proximo bool
}
