package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type URLInfo struct {
	ID          int
	ShortenURL  string
	OriginalURL string
	Created     time.Time
	Updated     time.Time
}

type Repository interface {
	Save(longURL string) (int64, error)
	Update(id int64, shortURL string) (int64, error)
	SearchByShortURL(shortURL string) (*URLInfo, error)
}

type Datastore struct {
	store *sql.DB
}

// Save inserts new longURL
func (d *Datastore) Save(longURL string) (int64, error) {
	ctx := context.Background()
	var err error

	if d.store == nil {
		err = errors.New("db is not initialized")
		return -1, err
	}

	// Check if database is alive.
	err = d.store.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := `
      INSERT INTO URLInfo (LongURL) VALUES (@longURL);
      select isNull(SCOPE_IDENTITY(), -1);
    `
	stmt, err := d.store.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sql.Named("longURL", longURL))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}

// Update updates a record with shortURL
func (d *Datastore) Update(id int64, shortURL string) (int64, error) {
	ctx := context.Background()

	if d.store == nil {
		err := errors.New("db is not initialized")
		return -1, err
	}

	// Check if database is alive.
	err := d.store.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := `UPDATE URLInfo SET ShortURL = @shortURL WHERE Id = @id`

	// Execute non-query with named parameters
	result, err := d.store.ExecContext(
		ctx,
		tsql,
		sql.Named("shortURL", shortURL),
		sql.Named("id", id))
	if err != nil {
		return -1, err
	}
	log.Println("record updated")
	return result.RowsAffected()
}

func (d *Datastore) SearchByShortURL(shortURL string) (*URLInfo, error) {
	var urlInfo URLInfo
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := `SELECT Id, LongURL, ShortURL FROM URLInfo WHERE ShortURL=@shortURL;`

	// Execute query
	row := d.store.QueryRowContext(ctx, tsql, sql.Named("shortURL", shortURL))
	if err != nil {
		return nil, err
	}

	err = row.Scan(&urlInfo.ID, &urlInfo.OriginalURL, &urlInfo.ShortenURL)
	if err != nil {
		return nil, err
	}

	return &urlInfo, nil
}
