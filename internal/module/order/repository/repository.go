package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/module/order/ports"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.OrderRepository = &orderRepository{}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) InsertNewOrder(ctx context.Context, tx *sqlx.Tx, data *entity.Order) (*entity.Order, error) {
	var res = new(entity.Order)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewOrder),
		data.ID,
		data.UserID,
		data.Status,
		data.ShippingName,
		data.ShippingAddress,
		data.ShippingPhone,
		data.ShippingNumber,
		data.ShippingType,
		data.TotalQuantity,
		data.TotalWeight,
		data.VoucherDiscount,
		data.AddressID,
		data.CostName,
		data.CostService,
		data.VoucherID,
		data.TotalProductAmount,
		data.TotalShippingCost,
		data.TotalShippingAmount,
		data.TotalAmount,
		data.Notes,
		data.OrderDate,
		data.CreatedAt,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.Status,
		&res.ShippingName,
		&res.ShippingAddress,
		&res.ShippingPhone,
		&res.ShippingNumber,
		&res.ShippingType,
		&res.TotalQuantity,
		&res.TotalWeight,
		&res.VoucherDiscount,
		&res.AddressID,
		&res.CostName,
		&res.CostService,
		&res.VoucherID,
		&res.TotalProductAmount,
		&res.TotalShippingCost,
		&res.TotalShippingAmount,
		&res.TotalAmount,
		&res.Notes,
		&res.OrderDate,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewOrder - Failed to insert new order")
		return nil, err
	}

	return res, nil
}

func (r *orderRepository) CountByFilter(ctx context.Context, userID uuid.UUID, search, status string) (int, error) {
	var (
		clauses []string
		args    []interface{}
	)

	clauses = append(clauses, "deleted_at IS NULL")
	clauses = append(clauses, "user_id = ?")
	args = append(args, userID)

	if search != "" {
		clauses = append(clauses, `(id::text ILIKE ? OR status ILIKE ? OR shipping_number ILIKE ?)`)
		term := "%" + search + "%"
		args = append(args, term, term, term)
	}
	if status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, status)
	}

	where := strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`SELECT COUNT(*) FROM orders WHERE %s`, where)

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Error().Err(err).
			Msgf("repository::CountByFilter - failed to count orders by filter: %v", args)
		return 0, err
	}

	return count, nil
}

func (r *orderRepository) FindByFilter(ctx context.Context, userID uuid.UUID, offset, limit int, search, status string) ([]entity.Order, error) {
	var (
		clauses []string
		args    []interface{}
	)

	clauses = append(clauses, "deleted_at IS NULL")
	clauses = append(clauses, "o.user_id = ?")
	args = append(args, userID)

	if search != "" {
		clauses = append(clauses, `(o.id::text ILIKE ? OR o.status ILIKE ? OR o.shipping_number ILIKE ?)`)
		term := "%" + search + "%"
		args = append(args, term, term, term)
	}
	if status != "" {
		clauses = append(clauses, "o.status = ?")
		args = append(args, status)
	}

	where := strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`
        SELECT
            o.id,
            o.order_date,
            o.status,
            o.total_quantity,
            o.total_amount,
            o.shipping_number,
            o.total_shipping_amount,
            o.voucher_id,
            o.voucher_discount,
            o.address_id,
            o.user_id,
            o.notes,
            o.total_product_amount,
            o.cost_name,
            o.cost_service
        FROM orders o
        WHERE %s
        ORDER BY o.order_date DESC
        LIMIT ? OFFSET ?
    `, where)
	args = append(args, limit, offset)

	var rows []entity.Order
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Error().Err(err).Msgf("repository::FindByFilter - failed to get orders by filter: %v", args)
		return nil, err
	}

	return rows, nil
}

func (r *orderRepository) FindOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var res = new(entity.Order)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindOrderByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msgf("repository::FindOrderByID - order not found: %s", id)
			return nil, errors.New(constants.ErrOrderNotFound)
		}

		log.Error().Err(err).Msgf("repository::FindOrderByID - failed to get order by id: %s", id)
		return nil, err
	}

	return res, nil
}

func (r *orderRepository) CountListOrderTransactionByFilter(ctx context.Context, search, status string) (int, error) {
	var (
		clauses []string
		args    []interface{}
	)

	clauses = append(clauses, "deleted_at IS NULL")

	if search != "" {
		clauses = append(clauses, `(id::text ILIKE ? OR status ILIKE ? OR shipping_number ILIKE ?)`)
		term := "%" + search + "%"
		args = append(args, term, term, term)
	}
	if status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, status)
	}

	where := strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`SELECT COUNT(*) FROM orders WHERE %s`, where)

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Error().Err(err).Msgf("repository::CountListOrderTransactionByFilter - failed to count orders by filter: %v", args)
		return 0, err
	}

	return count, nil
}

func (r *orderRepository) FindListTransactionOrderByFilter(ctx context.Context, offset, limit int, search, status string) ([]entity.Order, error) {
	var (
		clauses []string
		args    []interface{}
	)

	clauses = append(clauses, "deleted_at IS NULL")

	if search != "" {
		clauses = append(clauses, `(o.id::text ILIKE ? OR o.status ILIKE ? OR o.shipping_number ILIKE ?)`)
		term := "%" + search + "%"
		args = append(args, term, term, term)
	}
	if status != "" {
		clauses = append(clauses, "o.status = ?")
		args = append(args, status)
	}

	where := strings.Join(clauses, " AND ")
	query := fmt.Sprintf(`
        SELECT
            o.id,
            o.order_date,
            o.status,
            o.total_quantity,
            o.total_amount,
            o.shipping_number,
            o.total_shipping_amount,
            o.voucher_id,
            o.voucher_discount,
            o.address_id,
            o.user_id,
            o.notes,
            o.total_product_amount,
            o.cost_name,
            o.cost_service
        FROM orders o
        WHERE %s
        ORDER BY o.order_date DESC
        LIMIT ? OFFSET ?
    `, where)
	args = append(args, limit, offset)

	var rows []entity.Order
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Error().Err(err).Msgf("repository::FindListTransactionOrderByFilter - failed to get orders by filter: %v", args)
		return nil, err
	}

	return rows, nil
}

func (r *orderRepository) UpdateShippingNumber(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, shippingNumber string) (*entity.Order, error) {
	var res = new(entity.Order)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateShippingNumber),
		shippingNumber,
		id,
	).Scan(
		&res.ID,
		&res.UserID,
		&res.Status,
		&res.ShippingName,
		&res.ShippingAddress,
		&res.ShippingPhone,
		&res.ShippingNumber,
		&res.ShippingType,
		&res.TotalQuantity,
		&res.TotalWeight,
		&res.VoucherDiscount,
		&res.AddressID,
		&res.CostName,
		&res.CostService,
		&res.VoucherID,
		&res.TotalProductAmount,
		&res.TotalShippingCost,
		&res.TotalShippingAmount,
		&res.TotalAmount,
		&res.Notes,
		&res.OrderDate,
	)
	if err != nil {
		log.Error().Err(err).Msgf("repository::UpdateShippingNumber - failed to update shipping number: %s", id)
		return nil, err
	}

	return res, nil
}

func (r *orderRepository) UpdateOrderTransactionStatus(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, status string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateOrderTransactionStatus),
		status,
		id,
	)
	if err != nil {
		log.Error().Err(err).Msgf("repository::UpdateOrderTransactionStatus - failed to update order transaction status: %s", id)
		return err
	}

	return nil
}
