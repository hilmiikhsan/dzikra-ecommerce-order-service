package repository

const (
	querySumProductCapital = `
		SELECT
			COALESCE(SUM(i.product_capital_price), 0) AS sum_capital
		FROM order_item_histories AS i
		JOIN order_histories AS h
		ON i.order_id = h.id
		WHERE h.order_date BETWEEN ? AND ?
	`
)
