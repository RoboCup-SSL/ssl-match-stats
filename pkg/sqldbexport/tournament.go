package sqldbexport

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
)

func (p *SqlDbExporter) AddTournamentIfNotPresent(tournamentName string) (*uuid.UUID, error) {
	tournamentId, err := p.FindTournamentId(tournamentName)
	if err != nil {
		return nil, err
	}

	if tournamentId == nil {
		tournamentId = new(uuid.UUID)
		*tournamentId = uuid.New()
		if _, err := p.db.Exec(
			"INSERT INTO tournaments (id, name) VALUES ($1, $2)",
			tournamentId,
			tournamentName,
		); err != nil {
			return nil, err
		}
		log.Printf("New tournament %v inserted with id %v", tournamentName, tournamentId)
	}
	return tournamentId, nil
}

func (p *SqlDbExporter) FindTournamentId(tournamentName string) (*uuid.UUID, error) {
	id := new(uuid.UUID)
	err := p.db.QueryRow(
		"SELECT id FROM tournaments WHERE name=$1",
		tournamentName).Scan(id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "Could not query tournaments:")
	}
	return id, nil
}
