package database

import (
	"fmt"
	"shortfyurl/internal/models"
	"time"

	"github.com/gocql/gocql"
)

type CassandraDB struct {
	session *gocql.Session
}

func NewCassandraDB(hosts []string, keyspace string) (*CassandraDB, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no Cassandra: %v", err)
	}

	return &CassandraDB{session: session}, nil
}

func (db *CassandraDB) InitSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id BIGINT PRIMARY KEY,
			short_code TEXT,
			original_url TEXT,
			created_at TIMESTAMP,
			clicks BIGINT
		)
	`
	if err := db.session.Query(query).Exec(); err != nil {
		return fmt.Errorf("erro ao criar tabela: %v", err)
	}

	indexQuery := `CREATE INDEX IF NOT EXISTS ON urls (short_code)`
	if err := db.session.Query(indexQuery).Exec(); err != nil {
		return fmt.Errorf("erro ao criar índice: %v", err)
	}

	return nil
}

func (db *CassandraDB) SaveURL(url *models.URL) error {
	query := `INSERT INTO urls (id, short_code, original_url, created_at, clicks) VALUES (?, ?, ?, ?, ?)`
	return db.session.Query(query, url.ID, url.ShortCode, url.OriginalURL, url.CreatedAt, url.Clicks).Exec()
}

func (db *CassandraDB) GetURLByShortCode(shortCode string) (*models.URL, error) {
	var url models.URL
	query := `SELECT id, short_code, original_url, created_at, clicks FROM urls WHERE short_code = ? LIMIT 1`
	
	if err := db.session.Query(query, shortCode).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.Clicks); err != nil {
		return nil, err
	}

	return &url, nil
}

func (db *CassandraDB) GetAllURLs() ([]models.URL, error) {
	var urls []models.URL
	query := `SELECT id, short_code, original_url, created_at, clicks FROM urls`
	
	iter := db.session.Query(query).Iter()
	var url models.URL
	for iter.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.Clicks) {
		urls = append(urls, url)
		url = models.URL{}
	}
	
	if err := iter.Close(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (db *CassandraDB) IncrementClicks(shortCode string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?`
	return db.session.Query(query, shortCode).Exec()
}

func (db *CassandraDB) Close() {
	db.session.Close()
}
