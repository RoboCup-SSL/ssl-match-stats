package csvexport

import (
	"encoding/csv"
	"github.com/RoboCup-SSL/ssl-match-stats/pkg/matchstats"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strconv"
	"strings"
)

func writeCsv(header []string, data [][]string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Could not create CSV output file")
	}

	if _, err := f.WriteString("#" + strings.Join(header, ",") + "\n"); err != nil {
		return err
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(data); err != nil {
		return err
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return f.Close()
}

func WriteGamePhaseDurations(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Extra time", "Shootout"}
	for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
		header = append(header, matchstats.GamePhaseType_name[int32(i)][6:])
	}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		record := []string{matchStats.Name, strconv.FormatBool(matchStats.ExtraTime), strconv.FormatBool(matchStats.Shootout)}
		for i := 0; i < len(matchstats.GamePhaseType_name); i++ {
			name := matchstats.GamePhaseType_name[int32(i)]
			record = append(record, uintToStr(matchStats.TimePerGamePhase[name]))
		}
		records = append(records, record)
	}

	return writeCsv(header, records, filename)
}

func WriteTeamMetricsPerGame(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"File", "Team", "Scored Goals", "Conceded Goals", "Fouls", "Yellow Cards", "Red Cards", "Timeout Time", "Timeouts", "Penalty Shots", "Ball Placement Time", "Ball Placements", "Max active Yellow Cards"}

	var records [][]string
	for _, matchStats := range matchStatsCollection.MatchStats {
		recordYellow := []string{matchStats.Name}
		recordYellow = append(recordYellow, teamNumbers(matchStats.TeamStatsYellow)...)
		records = append(records, recordYellow)
		recordBlue := []string{matchStats.Name}
		recordBlue = append(recordBlue, teamNumbers(matchStats.TeamStatsBlue)...)
		records = append(records, recordBlue)
	}

	return writeCsv(header, records, filename)
}

func WriteTeamMetricsSum(matchStatsCollection *matchstats.MatchStatsCollection, filename string) error {

	header := []string{"Team", "Scored Goals", "Conceded Goals", "Fouls", "Yellow Cards", "Red Cards", "Timeout Time", "Timeouts", "Penalty Shots", "Ball Placement Time", "Ball Placements", "Max active Yellow Cards"}

	teams := map[string]*matchstats.TeamStats{}
	for _, matchStats := range matchStatsCollection.MatchStats {
		teams[matchStats.TeamStatsYellow.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsYellow.Name}
		teams[matchStats.TeamStatsBlue.Name] = &matchstats.TeamStats{Name: matchStats.TeamStatsBlue.Name}
	}

	for _, matchStats := range matchStatsCollection.MatchStats {
		addTeamStats(teams[matchStats.TeamStatsYellow.Name], matchStats.TeamStatsYellow)
		addTeamStats(teams[matchStats.TeamStatsBlue.Name], matchStats.TeamStatsBlue)
	}

	var teamNamesSorted []string
	for k := range teams {
		teamNamesSorted = append(teamNamesSorted, k)
	}
	sort.Strings(teamNamesSorted)

	var records [][]string
	for _, teamName := range teamNamesSorted {
		teamStats := teams[teamName]
		records = append(records, teamNumbers(teamStats))
	}

	return writeCsv(header, records, filename)
}

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

func addTeamStats(to *matchstats.TeamStats, team *matchstats.TeamStats) {
	to.Goals += team.Goals
	to.ConcededGoals += team.ConcededGoals
	to.Fouls += team.Fouls
	to.YellowCards += team.YellowCards
	to.RedCards += team.RedCards
	to.TimeoutTime += team.TimeoutTime
	to.Timeouts += team.Timeouts
	to.PenaltyShotsTotal += team.PenaltyShotsTotal
	to.BallPlacementTime += team.BallPlacementTime
	to.BallPlacements += team.BallPlacements
	if to.MaxActiveYellowCards < team.MaxActiveYellowCards {
		to.MaxActiveYellowCards = team.MaxActiveYellowCards
	}
}

func teamNumbers(stats *matchstats.TeamStats) []string {
	return []string{
		stats.Name,
		uintToStr(stats.Goals),
		uintToStr(stats.ConcededGoals),
		uintToStr(stats.Fouls),
		uintToStr(stats.YellowCards),
		uintToStr(stats.RedCards),
		uintToStr(stats.TimeoutTime),
		uintToStr(stats.Timeouts),
		uintToStr(stats.PenaltyShotsTotal),
		uintToStr(stats.BallPlacementTime),
		uintToStr(stats.BallPlacements),
		uintToStr(stats.MaxActiveYellowCards),
	}
}

func uintToStr(n uint32) string {
	return strconv.FormatUint(uint64(n), 10)
}

func intToStr(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}