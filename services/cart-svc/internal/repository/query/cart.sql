-- name: CreateCartIfNotExists :one
INSERT INTO cart(id, user_id, status)
VALUES ($1, $2, 'OPEN')
ON CONFLICT (id) DO NOTHING
RETURNING *;

-- name: GetCart :one
SELECT
  id,
  user_id,
  status,
  created_at,
  updated_at
FROM cart
WHERE id = $1
  AND status IN ('OPEN', 'CHECKOUT');

-- name: ListItems :many
SELECT
  cart_id,
  product_id,
  price,
  quantity
FROM cart_item
WHERE cart_id = $1
ORDER BY product_id;


-- name: UpdateCartStatus :exec
UPDATE cart
SET status = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteExpiredCarts :execrows
DELETE FROM cart
WHERE status = 'ABANDONED'
    AND updated_at < NOW() - INTERVAL '30 days';

-- name: UpdateQuantity :execrows
UPDATE cart_item SET quantity = $3
WHERE cart_id = $1
    AND product_id = $2;

-- name: UpsertCartItem :exec
INSERT INTO cart_item (cart_id, product_id, price, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, product_id) DO UPDATE
SET price = EXCLUDED.price,
    quantity = EXCLUDED.quantity;

-- name: DeleteCartItem :execrows
DELETE FROM cart_item
WHERE cart_id = $1
    AND product_id = $2;

-- name: DeleteCart :execrows
DELETE FROM cart
WHERE id = $1;



