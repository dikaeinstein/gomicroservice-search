package data

import (
	"database/sql"

	// Register mysql driver for db/sql interface
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// MySQLStore is a MongoDB data store which implements the Store interface
type MySQLStore struct {
	session *sql.DB
}

// NewMySQLStore creates an instance of MySQLStore with the given connection string
func NewMySQLStore(connection string) (*MySQLStore, error) {
	log.Println("Opening connection to:", connection)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Check connection is up
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &MySQLStore{session: db}, nil
}

// Search returns Kittens from the MySQL instance which have the name name
func (m *MySQLStore) Search(name string) ([]Kitten, error) {
	log.Println("Search for:", name)
	var results []Kitten

	rows, err := m.session.Query("SELECT id, name, weight FROM kittens WHERE name=?", name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		kitten := Kitten{}
		rows.Scan(&kitten.ID, &kitten.Name, &kitten.Weight)
		results = append(results, kitten)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteAllKittens deletes all the kittens from the datastore
func (m *MySQLStore) DeleteAllKittens() {
	m.session.Exec("DELETE FROM kittens")
}

// InsertKittens inserts a slice of kittens into the datastore
func (m *MySQLStore) InsertKittens(kittens []Kitten) error {
	for _, kitten := range kittens {
		_, err := m.session.Exec("INSERT INTO kittens (id, name, weight) VALUES (?, ?, ?)",
			kitten.ID,
			kitten.Name,
			kitten.Weight,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// CreateSchema creates the initial datbase schema
func (m *MySQLStore) CreateSchema() {
	m.session.Exec("DROP TABLE kittens")
	m.session.Exec("CREATE TABLE kittens (id varchar(50), name varchar(200), weight int)")
}
