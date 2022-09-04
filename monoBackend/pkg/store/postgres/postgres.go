package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

func init() {
	if err := initDefaultConnEnv(); err != nil {
		panic(err)
	}
}

//	ENV to PG params from here: https://postgrespro.ru/docs/postgresql/9.6/libpq-envars
func initDefaultConnEnv() error {
	if len(os.Getenv("PGHOST")) == 0 {
		if err := os.Setenv("PGHOST", "localhost"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGPORT")) == 0 {
		if err := os.Setenv("PGPORT", "5432"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGDATABASE")) == 0 {
		if err := os.Setenv("PGDATABASE", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGUSER")) == 0 {
		if err := os.Setenv("PGUSER", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGPASSWORD")) == 0 {
		if err := os.Setenv("PGPASSWORD", "postgres"); err != nil {
			return errors.WithStack(err)
		}
	}
	if len(os.Getenv("PGSSLMODE")) == 0 {
		if err := os.Setenv("PGSSLMODE", "disable"); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

type Config struct {
	Host     string
	Port     uint16
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg *Config) toDSN() string {
	var args []string

	if len(cfg.Host) > 0 {
		args = append(args, fmt.Sprintf("host=%s", cfg.Host))
	}
	if cfg.Port > 0 {
		args = append(args, fmt.Sprintf("port=%s", strconv.Itoa(int(cfg.Port))))
	}
	if len(cfg.Username) > 0 {
		args = append(args, fmt.Sprintf("user=%s", cfg.Username))
	}
	if len(cfg.Password) > 0 {
		args = append(args, fmt.Sprintf("password=%s", cfg.Password))
	}
	if len(cfg.DBName) > 0 {
		args = append(args, fmt.Sprintf("dbname=%s", cfg.DBName))
	}
	if len(cfg.SSLMode) > 0 {
		args = append(args, fmt.Sprintf("sslmode=%s", cfg.SSLMode))
	}

	return strings.Join(args, " ")
}

type store struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, cfg *Config) (Store, error) {
	parsedCfg, err := pgxpool.ParseConfig(cfg.toDSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn, err := pgxpool.ConnectConfig(ctx, parsedCfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &store{ctx: ctx, pool: conn}, nil
}

type Store interface {
	// Exec queries the database and returns the affected rows
	Exec(sql string, args ...interface{}) (int64, error)

	// Query queries the database and return the rows
	Query(sql string, args ...interface{}) (Rows, error)

	// QueryRow queries the database and return a single row
	QueryRow(sql string, args ...interface{}) Row

	// BeginTx Starts a database Transaction
	BeginTx() (Transaction, error)
}

// Row is the result returned from a query
type Row interface {
	// Scan reads the values from the current row into dest values positionally
	Scan(dest ...interface{}) error
}

// Rows is the result set returned from a query
type Rows interface { // TODO: Cursor
	// Scan reads the values from the current row into dest values positionally
	Scan(dest ...interface{}) error

	// Next prepares the next row for reading. It returns true if there is another row and false if no more rows are available. It automatically closes rows when all rows are read.
	Next() bool

	// Close closes the rows, making the connection ready for use again. It is safe to call Close after rows is already closed.
	Close()

	// Err returns any error that occurred while reading.
	Err() error
}

// Transaction represents an SQL database transaction
type Transaction interface {
	// Commit commits the database transaction
	Commit() error

	// Rollback rollbacks the database transaction
	Rollback() error

	// Exec queries the database and returns the affected rows
	Exec(sql string, args ...interface{}) (int64, error)

	// Query queries the database and return the rows
	Query(sql string, args ...interface{}) (Rows, error)

	// QueryRow queries the database and return a single row
	QueryRow(sql string, args ...interface{}) Row
}

func (s *store) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := s.pool.Exec(s.ctx, sql, args...)
	return result.RowsAffected(), err
}

func (s *store) Query(sql string, args ...interface{}) (Rows, error) {
	rows, err := s.pool.Query(s.ctx, sql, args...)
	return newDatabaseRows(rows), err
}

func (s *store) QueryRow(sql string, args ...interface{}) Row {
	row := s.pool.QueryRow(s.ctx, sql, args...)
	return newDatabaseRow(row)
}

func (s *store) BeginTx() (Transaction, error) {
	tx, err := s.pool.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return newDatabaseTransaction(s.ctx, tx), nil
}

type databaseRow struct {
	row pgx.Row
}

func newDatabaseRow(row pgx.Row) Row {
	return &databaseRow{
		row: row,
	}
}

func (rowPtr *databaseRow) Scan(dest ...interface{}) error {
	return rowPtr.row.Scan(dest...)
}

type databaseRows struct {
	rows pgx.Rows
}

func newDatabaseRows(rows pgx.Rows) Rows {
	return &databaseRows{
		rows: rows,
	}
}

func (p *databaseRows) Scan(dest ...interface{}) error {
	return p.rows.Scan(dest...)
}

func (p *databaseRows) Next() bool {
	return p.rows.Next()
}

func (p *databaseRows) Close() {
	p.rows.Close()
}

func (p *databaseRows) Err() error {
	return p.rows.Err()
}

type transaction struct {
	ctx context.Context
	tx  pgx.Tx
}

func newDatabaseTransaction(ctx context.Context, tx pgx.Tx) *transaction {
	return &transaction{
		ctx: ctx,
		tx:  tx,
	}
}

func (p *transaction) Commit() error {
	return p.tx.Commit(p.ctx)
}

func (p *transaction) Rollback() error {
	return p.tx.Rollback(p.ctx)
}

func (p *transaction) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := p.tx.Exec(p.ctx, sql, args...)
	return result.RowsAffected(), err
}

func (p *transaction) Query(sql string, args ...interface{}) (Rows, error) {
	rows, err := p.tx.Query(p.ctx, sql, args...)
	return newDatabaseRows(rows), err
}

func (p *transaction) QueryRow(sql string, args ...interface{}) Row {
	row := p.tx.QueryRow(p.ctx, sql, args...)
	return newDatabaseRow(row)
}
