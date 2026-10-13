package handlers

// indexCard is one item on the root index (list of resumos). The first item of
// the unfiltered list is the "próximo" — project invariant.
type indexCard struct {
	Slug     string
	Kicker   string
	TitleA   string
	TitleEm  string
	Lede     string
	Channel  string
	Date     string
	Duration string
	Threads  int
	Tags     []string
	Proximo  bool
}
