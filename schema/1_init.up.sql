CREATE TABLE IF NOT EXISTS item
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DECIMAL(8, 2) NOT NULL
);

INSERT INTO item (name, description, price) VALUES
('Smartphone', 'Latest model with a 6.5-inch screen, 128GB storage, and 48MP camera.', 799.99),
('Laptop', 'High-performance laptop with an Intel i7 processor, 16GB RAM, and 512GB SSD.', 1199.00),
('Headphones', 'Noise-cancelling over-ear headphones with Bluetooth connectivity.', 199.95),
('Smartwatch', 'Fitness tracking smartwatch with heart rate monitor and sleep analysis.', 149.99),
('Camera', 'Digital camera with 24.2MP resolution, 4K video recording, and WiFi connectivity.', 549.00),
('Tablet', '10.2-inch tablet with Apple A12 Bionic chip, 64GB storage, and Retina display.', 329.99),
('Bluetooth Speaker', 'Portable wireless speaker with deep bass and up to 12 hours of battery life.', 79.99),
('Keyboard', 'Mechanical keyboard with RGB lighting and programmable keys.', 89.99),
('External Hard Drive', '2TB external hard drive with USB 3.0 for fast data transfer.', 99.99),
('Drone', 'Quadcopter drone with 4K camera, 30-minute flight time, and GPS navigation.', 799.00);