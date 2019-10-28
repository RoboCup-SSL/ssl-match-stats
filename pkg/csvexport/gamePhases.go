package csvexport

import (
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"strings"
)

func WriteGamePhases(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {
	header := []string{
		"File",
		"Type",
		"For Team",
		"Duration",
		"Stage",
		"Stage time left entry",
		"Stage time left exit",
		"Command Entry",
		"Command Entry Team",
		"Command Exit",
		"Command Exit Team",
		"Next Command",
		"Next Command Team",
		"Primary Game Event Entry",
		"All Game Events Entry",
		"All Game Events Exit",
		"Command Prev",
		"Command Prev Team",
	}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		for _, gamePhase := range matchStats.GamePhases {
			primaryGameEvent := ""
			if len(gamePhase.GameEventsEntry) > 0 {
				primaryGameEvent = gamePhase.GameEventsEntry[0].Type.String()
			}
			nextCommandType := ""
			nextCommandForTeam := ""
			if gamePhase.NextCommandProposed != nil {
				nextCommandType = gamePhase.NextCommandProposed.Type.String()[8:]
				nextCommandForTeam = gamePhase.NextCommandProposed.ForTeam.String()[5:]
			}

			var gameEventsEntry []string
			for _, event := range gamePhase.GameEventsEntry {
				gameEventsEntry = append(gameEventsEntry, event.Type.String())
			}
			var gameEventsExit []string
			for _, event := range gamePhase.GameEventsExit {
				gameEventsExit = append(gameEventsExit, event.Type.String())
			}

			record := []string{
				matchStats.Name,
				gamePhase.Type.String()[6:],
				gamePhase.ForTeam.String()[5:],
				uintToStr(gamePhase.Duration),
				gamePhase.Stage.String()[6:],
				intToStr(gamePhase.StageTimeLeftEntry),
				intToStr(gamePhase.StageTimeLeftExit),
				gamePhase.CommandEntry.Type.String()[8:],
				gamePhase.CommandEntry.ForTeam.String()[5:],
				gamePhase.CommandExit.Type.String()[8:],
				gamePhase.CommandExit.ForTeam.String()[5:],
				nextCommandType,
				nextCommandForTeam,
				primaryGameEvent,
				strings.Join(gameEventsEntry, "|"),
				strings.Join(gameEventsExit, "|"),
				gamePhase.CommandPrev.Type.String()[8:],
				gamePhase.CommandPrev.ForTeam.String()[5:],
			}
			records = append(records, record)
		}
	}

	return writeCsv(header, records, filename)
}
