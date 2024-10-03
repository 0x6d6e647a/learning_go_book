package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
)

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Name  string
	Teams map[string]Team
	Wins  map[string]int
}

func (l *League) MatchResult(teamA string, scoreA int, teamB string, scoreB int) error {
	if _, ok := l.Teams[teamA]; !ok {
		return fmt.Errorf("teamA '%s' is not in league '%s'", teamA, l.Name)
	}

	if _, ok := l.Teams[teamB]; !ok {
		return fmt.Errorf("teamB '%s' is not in league '%s'", teamB, l.Name)
	}

	if scoreA > scoreB {
		l.Wins[teamA] += 1
	} else if scoreB > scoreA {
		l.Wins[teamB] += 1
	}

	return nil
}

func (l League) Ranking() []string {
	names := make([]string, 0, len(l.Teams))
	for name := range l.Teams {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		return l.Wins[names[i]] > l.Wins[names[j]]
	})

	return names
}

type Ranker interface {
	Ranking() []string
}

func RankPrinter(ranker Ranker, writer io.Writer) {
	for index, name := range ranker.Ranking() {
		io.WriteString(writer, fmt.Sprintf("%d. %s\n", index+1, name))
	}
}

func main() {
	// MLB 2024
	teams := [...]Team{
		{
			"Los Angeles Dodgers",
			[]string{
				"Anthony Banda",
				"Ryan Brasier",
				"Walker Buehler",
				"Jack Flaherty",
				"Nick Frasso",
				"Edgardo Henriquez",
				"Daniel Hudson",
				"Joe Kelly",
				"Landon Knack",
				"Michael Kopech",
				"Evan Phillips",
				"Blake Treinen",
				"Alex Vesia",
				"Yoshinobu Yamamoto",
				"Austin Barnes",
				"Will Smith",
				"Freddie Freeman",
				"Enrique Hernandez",
				"Gavin Lux",
				"Max Muncy",
				"Miguel Rojas",
				"Tommy Edman",
				"Teoscar Hernandez",
				"Kevin Kiermaier",
				"James Outman",
				"Andy Pages",
				"Chris Taylor",
				"Shohei Ohtani",
			},
		},
		{
			"Philadelphia Phillies",
			[]string{
				"Kolby Allard",
				"Jose Alvarado",
				"Tanner Banks",
				"Jose Cuas",
				"Carlos Estevez",
				"Jeff Hoffman",
				"Orion Kerkering",
				"Max Lazar",
				"Aaron Nola",
				"Jose Ruiz",
				"Cristopher Sanchez",
				"Matt Strahm",
				"Ranger Suarez",
				"Taijuan Walker",
				"Zack Wheeler",
				"J.T. Realmuto",
				"Garrett Stubbs",
				"Alec Bohm",
				"Kody Clemens",
				"Bryce Harper",
				"Buddy Kennedy",
				"Edmundo Sosa",
				"Bryson Stott",
				"Trea Turner",
				"Nick Castellanos",
				"Austin Hays",
				"Brandon Marsh",
				"Johan Rojas",
				"Weston Wilson",
				"Kyle Schwarber",
			},
		},
		{
			"New York Yankees",
			[]string{
				"Clayton Beeter",
				"Gerrit Cole",
				"Luis Gil",
				"Ian Hamilton",
				"Tim Hill",
				"Clay Holmes",
				"Tommy Kahnle",
				"Mark Leiter Jr.",
				"Tim Mayza",
				"Carlos Rodon",
				"Clarke Schmidt",
				"Marcus Stroman",
				"Luke Weaver",
				"Jose Trevino",
				"Austin Wells",
				"Jon Berti",
				"Oswaldo Cabrera",
				"Jazz Chisholm Jr.",
				"Ben Rice",
				"Anthony Rizzo",
				"Gleyber Torres",
				"Anthony Volpe",
				"Jasson Dominguez",
				"Trent Grisham",
				"Aaron Judge",
				"Juan Soto",
				"Alex Verdugo",
				"Giancarlo Stanton",
			},
		},
		{
			"Milwaukee Brewers",
			[]string{
				"Aaron Ashby",
				"Aaron Civale",
				"DL Hall",
				"Jared Koenig",
				"Nick Mears",
				"Trevor Megill",
				"Frankie Montas",
				"Tobias Myers",
				"Joel Payamps",
				"Freddy Peralta",
				"Joe Ross",
				"Devin Williams",
				"William Contreras",
				"Eric Haase",
				"Gary Sanchez",
				"Willy Adames",
				"Jake Bauers",
				"Rhys Hoskins",
				"Andruw Monasterio",
				"Joey Ortiz",
				"Brice Turang",
				"Jackson Chourio",
				"Isaac Collins",
				"Sal Frelick",
				"Garrett Mitchell",
				"Blake Perkins",
			},
		},
		{
			"San Diego Padres",
			[]string{
				"Jason Adam",
				"Dylan Cease",
				"Yu Darvish",
				"Jeremiah Estrada",
				"Bryan Hoeing",
				"Michael King",
				"Yuki Matsui",
				"Adrian Morejon",
				"Joe Musgrove",
				"Wandy Peralta",
				"Tanner Scott",
				"Robert Suarez",
				"Elias Diaz",
				"Kyle Higashioka",
				"Nick Ahmed",
				"Luis Arraez",
				"Xander Bogaerts",
				"Jake Cronenworth",
				"Manny Machado",
				"Donovan Solano",
				"Tyler Wade",
				"Brandon Lockridge",
				"Jackson Merrill",
				"David Peralta",
				"Jurickson Profar",
				"Fernando Tatis Jr.",
			},
		},
	}
	league := League{
		"Major League Baseball 2024",
		make(map[string]Team, len(teams)),
		make(map[string]int),
	}
	for _, team := range teams {
		league.Teams[team.Name] = team
	}

	// Simulate some games.
	const numGames = 30
	teamNames := make([]string, 0, len(league.Teams))
	for teamName := range league.Teams {
		teamNames = append(teamNames, teamName)
	}
	numTeams := len(teamNames)

	for i := 0; i < numGames; i++ {
		// Randomly pick two teams.
		var indexA, indexB uint

		for {
			indexA = uint(rand.Intn(numTeams))
			indexB = uint(rand.Intn(numTeams))

			if indexA != indexB {
				break
			}
		}

		teamA := teamNames[indexA]
		teamB := teamNames[indexB]

		// Randomly generate score.
		scoreA := rand.Int()
		scoreB := rand.Int()
		err := league.MatchResult(teamA, scoreA, teamB, scoreB)
		if err != nil {
			panic(err)
		}
	}

	// Print ranking.
	RankPrinter(league, os.Stdout)
}
