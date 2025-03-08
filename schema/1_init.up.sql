CREATE TABLE IF NOT EXISTS adress
(
    id SERIAL PRIMARY KEY,
    country VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    street VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS "user" 
(
    id SERIAL PRIMARY KEY,
    email VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(30) NOT NULL UNIQUE,
    password bytea NOT NULL,
    first_name VARCHAR(50),
    second_name VARCHAR(50),
    adress_id INT,
    FOREIGN KEY (adress_id) REFERENCES adress (id)
);

CREATE TABLE IF NOT EXISTS product
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(300),
    price DECIMAL(8, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS category
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS product_category
(   
    product_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (category_id) REFERENCES category (id)
);

CREATE TABLE IF NOT EXISTS "order"
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amout DECIMAL(10, 2),
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE TABLE IF NOT EXISTS order_details
(
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(8, 2) NOT NULL,
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES "order" (id),
    FOREIGN KEY (product_id) REFERENCES product(id)
);

INSERT INTO product (name, description, price) VALUES
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