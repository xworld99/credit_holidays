package db

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
	"github.com/knadh/koanf"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg *koanf.Koanf) (CreditHolidaysDB, error) {
	db := PostgresDB{}

	if err := db.init(cfg); err != nil {
		return nil, err
	}

	return &db, nil
}

func (p *PostgresDB) init(cfg *koanf.Koanf) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.String("postgres.host"),
		cfg.String("postgres.port"),
		cfg.String("postgres.user"),
		cfg.String("postgres.pass"),
		cfg.String("postgres.dbname"),
	)

	var err error
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("can't connect to database %s: %w", connStr, err)
	}

	log.Info("db initialized")
	return nil
}

func (p *PostgresDB) Begin(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (p *PostgresDB) Commit(tx *sql.Tx) {
	tx.Commit()
}

func (p *PostgresDB) Rollback(tx *sql.Tx) {
	tx.Rollback()
}

func (p *PostgresDB) Close() {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
}

func (p *PostgresDB) GetUserById(ctx context.Context, u models.User) (models.User, error) {
	query := `SELECT id, balance, frozen_balance FROM users WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, u.Id).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}

func (p *PostgresDB) InsertUserIfNotExists(ctx context.Context, tx *sql.Tx, u models.User) (models.User, error) {
	query := `INSERT INTO users(id, balance, frozen_balance) VALUES ($1, 0, 0) ON CONFLICT DO NOTHING RETURNING *`
	err := tx.QueryRowContext(ctx, query, u.Id).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (p *PostgresDB) UpdateUser(ctx context.Context, tx *sql.Tx, u models.User) (models.User, error) {
	query := `UPDATE users SET balance = $1, frozen_balance = $2 WHERE id = $3 RETURNING id, balance, frozen_balance`
	err := tx.QueryRowContext(ctx, query, u.Balance, u.FrozenBalance, u.Id).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (p *PostgresDB) GetServicesList(ctx context.Context) ([]models.Service, error) {
	query := `SELECT id, name, description, confirmation_needed FROM services ORDER BY id`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var res []models.Service

	for rows.Next() {
		tmp := models.Service{}
		err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Description, &tmp.ConfNeeded)
		if err != nil {
			return nil, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (p *PostgresDB) GetServiceById(ctx context.Context, s models.Service) (models.Service, error) {
	query := `SELECT id, name, description, confirmation_needed, service_type FROM services WHERE id = $1`
	err := p.db.QueryRowContext(ctx, query, s.Id).Scan(&s.Id, &s.Name, &s.Description, &s.ConfNeeded, &s.ServiceType)
	if err != nil {
		return models.Service{}, err
	}

	return s, nil
}

func (p *PostgresDB) CreateOrder(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	query := `INSERT INTO orders(user_id, service_id, amount)
              VALUES ($1, $2, $3) RETURNING id, user_id, service_id, amount, status`
	err := tx.QueryRowContext(ctx, query, order.UserId, order.ServiceId, order.Amount).Scan(&order.Id, &order.UserId, &order.ServiceId, &order.Amount, &order.Status)
	if err != nil {
		return models.Order{}, err
	}

	order.ProofedAtStr = "null"
	return order, nil
}

func (p *PostgresDB) UpdateOrder(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	query := `UPDATE orders SET proofed_at = $1, status = $2 WHERE id = $3
              RETURNING id, created_at, proofed_at, user_id, service_id, amount, status`
	builder := tx.QueryRowContext(ctx, query, order.ProofedAtStr, order.Status, order.Id)
	err := builder.Scan(&order.Id, &order.CreatedAt, &order.ProofedAt, &order.UserId,
		&order.ServiceId, &order.Amount, &order.Status)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (p *PostgresDB) GetFullOrderInfo(
	ctx context.Context,
	tx *sql.Tx,
	order models.Order,
	user models.User,
	service models.Service,
) (models.Order, models.User, models.Service, error) {
	query := `SELECT o.amount, o.status, u.balance, u.frozen_balance, s.service_type 
              FROM orders o JOIN users u on o.user_id = u.id JOIN services s on o.service_id = s.id WHERE o.id = $1`
	builder := tx.QueryRowContext(ctx, query, order.Id)
	err := builder.Scan(&order.Amount, &order.Status, &user.Balance, &user.FrozenBalance, &service.ServiceType)
	if err != nil {
		return order, user, service, err
	}

	return order, user, service, nil
}

func (p *PostgresDB) GetLastOrderMonth(ctx context.Context, period time.Time) (models.Order, error) {
	query := `SELECT created_at FROM orders 
              WHERE status = 'success' and date_part('month', created_at) = $1 and date_part('year', created_at) = $2`
	var res models.Order

	err := p.db.QueryRowContext(ctx, query, int(period.Month()), period.Year()).Scan(&res.CreatedAt)
	if err != nil {
		return models.Order{}, err
	}

	return res, nil
}

func (p *PostgresDB) GetHistoryFrame(ctx context.Context, frame models.HistoryFrame) (models.HistoryFrame, error) {
	var err error

	queryCnt := `SELECT count(*) FROM orders WHERE user_id = $1 and created_at >= $2 and created_at <= $3`
	err = p.db.QueryRowContext(ctx, queryCnt, frame.UserId, frame.FromDate, frame.ToDate).Scan(&frame.TotalOperations)
	if err != nil {
		return models.HistoryFrame{}, err
	}

	query := `SELECT o.id, o.created_at, o.proofed_at, o.status, o.amount, s.name, s.description, s.service_type
              FROM orders o
                JOIN services s on o.service_id = s.id and o.user_id = $1
              WHERE o.created_at >= $2 AND o.created_at <= $3 ORDER BY $4 LIMIT $5 OFFSET $6`
	rows, err := p.db.QueryContext(ctx, query, frame.UserId, frame.FromDate, frame.ToDate,
		frame.OrderBy, frame.Limit, frame.Offset)
	if err != nil {
		return models.HistoryFrame{}, err
	}

	frame.Operations = []models.History{}

	for rows.Next() {
		tmp := models.History{}
		err := rows.Scan(&tmp.OrderId, &tmp.CreateAt, &tmp.ProofedAt, &tmp.Status, &tmp.Amount, &tmp.ServiceName,
			&tmp.ServiceDescription, &tmp.ServiceType)
		if err != nil {
			return models.HistoryFrame{}, err
		}
		frame.Operations = append(frame.Operations, tmp)
	}

	return frame, nil
}

func (p *PostgresDB) FormReport(ctx context.Context, data models.CSVData) (models.CSVData, error) {
	query := `SELECT s.id, s.name, s.service_type, sum(o.amount)
              FROM orders o JOIN services s ON o.service_id = s.id
              WHERE o.status = 'success' and date_part('month', o.created_at) = $1 and date_part('year', o.created_at) = $2
              GROUP BY s.id, s.name, s.service_type ORDER BY s.id`
	rows, err := p.db.QueryContext(ctx, query, int(data.Period.Month()), data.Period.Year())
	if err != nil {
		return models.CSVData{}, err
	}

	for rows.Next() {
		tmp := models.CSVRow{}
		err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Type, &tmp.CashFlow)
		if err != nil {
			return models.CSVData{}, err
		}
		data.Records = append(data.Records, tmp)
	}

	return data, nil
}
