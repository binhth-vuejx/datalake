-- +goose Up
-- +goose StatementBegin

-- Seed customers
INSERT INTO customers (name, email, phone, created_by, created_at, updated_at) VALUES
('John Doe', 'john@example.com', '555-0001', 'system', NOW(), NOW()),
('Jane Smith', 'jane@example.com', '555-0002', 'system', NOW(), NOW()),
('Bob Johnson', 'bob@example.com', '555-0003', 'system', NOW(), NOW()),
('Alice Williams', 'alice@example.com', '555-0004', 'system', NOW(), NOW()),
('Charlie Brown', 'charlie@example.com', '555-0005', 'system', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Seed sales data
INSERT INTO sales (amount, date, customer_id, created_by, created_at, updated_at) VALUES
(150.00, '2024-01-05 10:30:00', 1, 'system', NOW(), NOW()),
(250.50, '2024-01-10 14:15:00', 2, 'system', NOW(), NOW()),
(100.00, '2024-01-15 09:45:00', 3, 'system', NOW(), NOW()),
(320.75, '2024-02-01 11:20:00', 1, 'system', NOW(), NOW()),
(180.25, '2024-02-05 16:30:00', 4, 'system', NOW(), NOW()),
(450.00, '2024-02-10 13:00:00', 2, 'system', NOW(), NOW()),
(200.00, '2024-02-15 10:15:00', 5, 'system', NOW(), NOW()),
(275.50, '2024-03-01 15:45:00', 3, 'system', NOW(), NOW()),
(125.00, '2024-03-05 09:30:00', 1, 'system', NOW(), NOW()),
(350.00, '2024-03-10 12:00:00', 4, 'system', NOW(), NOW()),
(225.75, '2024-03-15 14:20:00', 2, 'system', NOW(), NOW()),
(400.00, '2024-03-20 10:45:00', 5, 'system', NOW(), NOW());

-- Seed orders (with sale_id references)
INSERT INTO orders (customer_id, sale_id, status, total_amount, created_by, created_at, updated_at) VALUES
(1, 1, 'completed', 150.00, 'system', NOW(), NOW()),
(2, 2, 'completed', 250.50, 'system', NOW(), NOW()),
(3, 3, 'pending', 100.00, 'system', NOW(), NOW()),
(1, 4, 'completed', 320.75, 'system', NOW(), NOW()),
(4, 5, 'completed', 180.25, 'system', NOW(), NOW()),
(2, 6, 'completed', 450.00, 'system', NOW(), NOW()),
(5, 7, 'pending', 200.00, 'system', NOW(), NOW()),
(3, 8, 'completed', 275.50, 'system', NOW(), NOW()),
(1, 9, 'completed', 125.00, 'system', NOW(), NOW()),
(4, 10, 'completed', 350.00, 'system', NOW(), NOW()),
(2, 11, 'completed', 225.75, 'system', NOW(), NOW()),
(5, 12, 'completed', 400.00, 'system', NOW(), NOW());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM orders WHERE created_by = 'system';
DELETE FROM sales WHERE created_by = 'system';
DELETE FROM customers WHERE created_by = 'system';
-- +goose StatementEnd
