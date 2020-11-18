package models

import "fmt"

type Gunpla struct {
	BoxArtImageLink string
	Title           string
	Subtitle        string
	Classification  string
	LineupNo        string
	Scale           string
	Franchise       string
	ReleaseDate     string
	JanISBN         string
	Run             string
	Price           string
	Includes        string
	Features        string
}

func (g *Gunpla) Print() {
	fmt.Printf(
		`
Title: %s
Subtitle: %s
Classification: %s
LineupNo: %s
Scale: %s
Franchise: %s
ReleaseDate: %s
JanISBN: %s
Run: %s
Price: %s
BoxArtLink: %s
		 
`,
		g.Title, g.Subtitle, g.Classification, g.LineupNo, g.Scale, g.Franchise, g.ReleaseDate, g.JanISBN, g.Run, g.Price, g.BoxArtImageLink,
	)
}
