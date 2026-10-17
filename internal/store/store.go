// Package store persists history in SQLite with the §3.4 recipe: the file is
// pre-created 0600 inside a 0700 dir before SQLite opens it (MyGo would
// create it 0644), journal mode is pinned to DELETE, and secure_delete is on
// so deleted clips do not linger in free pages.
package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/egoist/mygo/plugins/sqlite"

	"github.com/L92-Labs/myclip/internal/history"
)

const (
	DirMode  fs.FileMode = 0o700
	FileMode fs.FileMode = 0o600
)

// Sidecars are the files SQLite may create beside a database.
var Sidecars = []string{"-journal", "-wal", "-shm"}

// Store is the history database. It is safe from any goroutine; the
// connection serializes operations.
type Store struct {
	db   *sqlite.DB
	path string
}

const schema = `
CREATE TABLE IF NOT EXISTS items (
	id               INTEGER PRIMARY KEY,
	application      TEXT    NOT NULL DEFAULT '',
	first_copied_at  INTEGER NOT NULL,
	last_copied_at   INTEGER NOT NULL,
	number_of_copies INTEGER NOT NULL DEFAULT 1,
	pin              TEXT    NOT NULL DEFAULT '',
	title            TEXT    NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS contents (
	item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
	ord     INTEGER NOT NULL,
	type    TEXT    NOT NULL,
	value   BLOB,
	PRIMARY KEY (item_id, ord)
);
CREATE INDEX IF NOT EXISTS items_last_copied ON items(last_copied_at DESC);
`

// PrepareFile creates path O_EXCL 0600, or tightens an existing regular file
// owned by this user to 0600. Existing sidecars get the same check. It
// reports whether the file was created now. It never follows a symlink.
func PrepareFile(path string) (created bool, err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, FileMode)
	if err == nil {
		// umask cannot widen 0600, but a pre-set ACL-free fchmod pins it.
		err = f.Chmod(FileMode)
		return true, errors.Join(err, f.Close())
	}
	if !errors.Is(err, fs.ErrExist) {
		return false, fmt.Errorf("store: create %s: %w", path, err)
	}
	if err := tighten(path); err != nil {
		return false, err
	}
	for _, s := range Sidecars {
		if err := tighten(path + s); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return false, err
		}
	}
	return false, nil
}

// tighten opens path without following links, requires a regular file owned
// by this user, and sets it to 0600.
func tighten(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("store: refuse %s: %w", path, err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("store: stat %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("store: refuse %s: not a regular file", path)
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok && int(st.Uid) != os.Getuid() {
		return fmt.Errorf("store: refuse %s: owned by uid %d", path, st.Uid)
	}
	if info.Mode().Perm() != FileMode {
		if err := f.Chmod(FileMode); err != nil {
			return fmt.Errorf("store: chmod %s: %w", path, err)
		}
	}
	return nil
}

// Open prepares path (see PrepareFile), opens it and pins the connection's
// security pragmas. The parent dir must exist; MyClip creates it 0700.
func Open(ctx context.Context, path string) (*Store, error) {
	if _, err := PrepareFile(path); err != nil {
		return nil, err
	}
	db, err := sqlite.Open(ctx, path, sqlite.OpenOptions{})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.pin(ctx); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Transaction(ctx, statements(schema)); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: schema: %w", err)
	}
	return s, nil
}

// pin sets journal_mode=DELETE and secure_delete=ON, and asserts both. A
// stray -wal (a crash of a build that used WAL) is checkpointed first so its
// pages are not lost; only a failed re-pin refuses startup.
func (s *Store) pin(ctx context.Context) error {
	if _, err := os.Lstat(s.path + "-wal"); err == nil {
		if _, err := s.db.Query(ctx, "PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
			return fmt.Errorf("store: checkpoint stray wal: %w", err)
		}
	}
	mode, err := s.text(ctx, "PRAGMA journal_mode=DELETE")
	if err != nil {
		return err
	}
	if !strings.EqualFold(mode, "delete") {
		return fmt.Errorf("store: journal_mode is %q, want delete", mode)
	}
	if _, err := s.db.Execute(ctx, "PRAGMA secure_delete=ON"); err != nil {
		return err
	}
	if v, err := s.text(ctx, "PRAGMA secure_delete"); err != nil || v != "1" {
		return fmt.Errorf("store: secure_delete is %q (%v), want 1", v, err)
	}
	return nil
}

// JournalMode reports the connection's journal mode (tests assert it).
func (s *Store) JournalMode(ctx context.Context) (string, error) {
	return s.text(ctx, "PRAGMA journal_mode")
}

func (s *Store) text(ctx context.Context, q string) (string, error) {
	r, err := s.db.Query(ctx, q)
	if err != nil {
		return "", err
	}
	if len(r.Values) == 0 || len(r.Values[0]) == 0 {
		return "", fmt.Errorf("store: %s returned no rows", q)
	}
	return fmt.Sprint(r.Values[0][0]), nil
}

// Path is the database file.
func (s *Store) Path() string { return s.path }

// Close closes the connection.
func (s *Store) Close() error { return s.db.Close() }

// Load returns all items, newest first.
func (s *Store) Load(ctx context.Context) ([]*history.Item, error) {
	r, err := s.db.Query(ctx, `SELECT id, application, first_copied_at, last_copied_at,
		number_of_copies, pin, title FROM items ORDER BY last_copied_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	items := make([]*history.Item, 0, len(r.Values))
	byID := make(map[int64]*history.Item, len(r.Values))
	for _, row := range r.Values {
		it := &history.Item{
			ID:             asInt(row[0]),
			Application:    asString(row[1]),
			FirstCopiedAt:  time.Unix(0, asInt(row[2])),
			LastCopiedAt:   time.Unix(0, asInt(row[3])),
			NumberOfCopies: int(asInt(row[4])),
			Pin:            asString(row[5]),
			Title:          asString(row[6]),
		}
		items = append(items, it)
		byID[it.ID] = it
	}
	c, err := s.db.Query(ctx, "SELECT item_id, type, value FROM contents ORDER BY item_id, ord")
	if err != nil {
		return nil, err
	}
	for _, row := range c.Values {
		if it := byID[asInt(row[0])]; it != nil {
			v, _ := row[2].([]byte)
			it.Contents = append(it.Contents, history.Content{Type: asString(row[1]), Value: v})
		}
	}
	return items, nil
}

// Save inserts it (ID 0, then sets ID) or replaces its stored row and
// contents, atomically.
func (s *Store) Save(ctx context.Context, it *history.Item) error {
	fields := []any{it.Application, it.FirstCopiedAt.UnixNano(), it.LastCopiedAt.UnixNano(),
		int64(it.NumberOfCopies), it.Pin, it.Title}
	var stmts []sqlite.Statement
	// New rows take max(id)+1 (INTEGER PRIMARY KEY without AUTOINCREMENT),
	// so contents of a new item reference max(id) inside the transaction.
	ref := "(SELECT max(id) FROM items)"
	if it.ID == 0 {
		stmts = append(stmts, sqlite.Statement{SQL: `INSERT INTO items(application, first_copied_at,
			last_copied_at, number_of_copies, pin, title) VALUES (?, ?, ?, ?, ?, ?)`, Args: fields})
	} else {
		ref = "?"
		stmts = append(stmts,
			sqlite.Statement{SQL: `INSERT INTO items(id, application, first_copied_at, last_copied_at,
				number_of_copies, pin, title) VALUES (?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT(id) DO UPDATE SET application=excluded.application,
				first_copied_at=excluded.first_copied_at, last_copied_at=excluded.last_copied_at,
				number_of_copies=excluded.number_of_copies, pin=excluded.pin, title=excluded.title`,
				Args: append([]any{it.ID}, fields...)},
			sqlite.Statement{SQL: "DELETE FROM contents WHERE item_id = ?", Args: []any{it.ID}})
	}
	for i, c := range it.Contents {
		args := []any{int64(i), c.Type, blob(c.Value)}
		if it.ID != 0 {
			args = append([]any{it.ID}, args...)
		}
		stmts = append(stmts, sqlite.Statement{
			SQL:  "INSERT INTO contents(item_id, ord, type, value) VALUES (" + ref + ", ?, ?, ?)",
			Args: args,
		})
	}
	res, err := s.db.Transaction(ctx, stmts)
	if err != nil {
		return err
	}
	if it.ID == 0 {
		it.ID = res[0].LastInsertID
	}
	return nil
}

// Delete removes items by ID; their contents cascade.
func (s *Store) Delete(ctx context.Context, ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}
	stmts := make([]sqlite.Statement, len(ids))
	for i, id := range ids {
		stmts[i] = sqlite.Statement{SQL: "DELETE FROM items WHERE id = ?", Args: []any{id}}
	}
	_, err := s.db.Transaction(ctx, stmts)
	return err
}

// statements splits a schema script on semicolons (no literals contain one).
func statements(script string) []sqlite.Statement {
	var out []sqlite.Statement
	for _, q := range strings.Split(script, ";") {
		if q = strings.TrimSpace(q); q != "" {
			out = append(out, sqlite.Statement{SQL: q})
		}
	}
	return out
}

// blob keeps an empty value a zero-length BLOB rather than NULL.
func blob(b []byte) []byte {
	if b == nil {
		return []byte{}
	}
	return b
}

func asInt(v any) int64 {
	n, _ := v.(int64)
	return n
}

func asString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	}
	return ""
}
