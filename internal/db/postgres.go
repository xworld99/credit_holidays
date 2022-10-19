package db

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/knadh/koanf"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg *koanf.Koanf) (*PostgresDB, error) {
	db := &PostgresDB{}

	if err := db.init(cfg); err != nil {
		return nil, err
	}

	return db, nil
}

func (p *PostgresDB) init(cfg *koanf.Koanf) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s, dbname=%s",
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

func (p *PostgresDB) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (p *PostgresDB) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func (p *PostgresDB) GetUserById(ctx context.Context, u models.User) (models.User, error) {
	query := sq.Select("id, balance, frozen_balance").From("users").Where("id = ?", u.Id)
	err := query.RunWith(p.db).QueryRowContext(ctx).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, fmt.Errorf("user with id %d isnt exists", u.Id)
	}
	return u, err
}

func (p *PostgresDB) GetCreateUser(ctx context.Context, tx *sql.Tx, u models.User) (models.User, error) {
	query := "INSERT INTO user(id, balance, frozen_balance) VALUES (?, 0, 0) ON CONFLICT DO NOTHING RETURNING *;"
	err := tx.QueryRowContext(ctx, query, u.Id).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, fmt.Errorf("cant create/get user")
	}

	return u, nil
}

func (p *PostgresDB) UpdateUser(ctx context.Context, tx *sql.Tx, u models.User) (models.User, error) {
	query := "UPDATE users SET balance = ?, frozen_balance = ? WHERE id = ? RETURNING id, balance, frozen_balance;"
	err := tx.QueryRowContext(ctx, query, u.Balance, u.FrozenBalance, u.Id).Scan(&u.Id, &u.Balance, &u.FrozenBalance)
	if err != nil {
		return models.User{}, fmt.Errorf("order with id %d isnt exists", u.Id)
	}

	return u, nil
}

func (p *PostgresDB) GetServicesList(ctx context.Context) ([]models.Service, error) {
	query := sq.Select("id, name, description, confirmation_needed")
	query = query.From("services").OrderBy("id")
	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant read data from services table")
	}

	defer rows.Close()
	var res []models.Service

	for rows.Next() {
		tmp := models.Service{}
		err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Description, &tmp.ConfNeeded)
		if err != nil {
			return nil, fmt.Errorf("cant read data from services table")
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (p *PostgresDB) GetServiceById(ctx context.Context, s models.Service) (models.Service, error) {
	query := sq.Select("id, name, description, confirmation_needed, service_type")
	query = query.From("services").Where("id = ?", s.Id)
	err := query.RunWith(p.db).QueryRowContext(ctx).Scan(&s.Id, &s.Name, &s.Description, &s.ConfNeeded, &s.ServiceType)
	if err != nil {
		return models.Service{}, fmt.Errorf("service with id %d isnt exists", s.Id)
	}

	return s, nil
}

func (p *PostgresDB) CreateOrder(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	query := `INSERT INTO orders(user_id, service_id, amount)
              VALUES (?, ?, ?) RETURNING id, user_id, service_id, amount, status;`
	err := tx.QueryRowContext(ctx, query, order.UserId, order.ServiceId, order.Amount).Scan(&order.Id, &order.UserId, &order.ServiceId, &order.Amount, &order.Status)
	if err != nil {
		return models.Order{}, fmt.Errorf("cant create order")
	}

	order.ProofedAtStr = "null"
	return order, nil
}

func (p *PostgresDB) UpdateOrder(ctx context.Context, tx *sql.Tx, order models.Order) (models.Order, error) {
	query := `UPDATE orders SET proofed_at = ?, status = ? WHERE id = ?
              RETURNING id, created_at, proofed_at, user_id, service_id, amount, status;`
	builder := tx.QueryRowContext(ctx, query, order.ProofedAtStr, order.Status, order.Id)
	err := builder.Scan(&order.Id, &order.CreatedAt, &order.ProofedAt, &order.UserId,
		&order.ServiceId, &order.Amount, &order.Status)
	if err != nil {
		return models.Order{}, fmt.Errorf("order with id %d isnt exists", order.Id)
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
              FROM orders o JOIN users u on o.user_id = u.id JOIN services s on o.service_id = s.id WHERE o.id = ?;`
	builder := tx.QueryRowContext(ctx, query, order.Id)
	err := builder.Scan(&order.Amount, &order.Status, &user.Balance, &user.FrozenBalance, &service.ServiceType)
	if err != nil {
		return order, user, service, fmt.Errorf("order with id %d isnt exists", order.Id)
	}

	return order, user, service, nil
}
