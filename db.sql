CREATE TABLE `email_confirmation_tokens`
(
    `id`         MEDIUMINT(8) references users (id) on delete cascade,
    `user_id`    mediumint(8) NOT NULL,
    `token`      VARCHAR(64)  NOT NULL,
    `expiration` DATETIME     NOT NULL,
    `status`     bool DEFAULT true
);


CREATE TABLE `bids`
(
    `id`         mediumint(8) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id`    mediumint(8) UNSIGNED NOT NULL,
    `product_id` mediumint(8)          NOT NULL,
    `bid`        float(8)              NOT NULL,
    `date`       datetime              NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `users`
(
    `id`                mediumint(8) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `first_name`        varchar(20)           NOT NULL,
    `last_name`         varchar(40)           NOT NULL,
    `email`             varchar(60)           NOT NULL,
    `phone`             varchar(12)           NOT NULL,
    `dob`               varchar(60)           NOT NULL,
    `status`            bool                           DEFAULT false,
    `role`              tinyint(1) UNSIGNED   NOT NULL DEFAULT '1',
    `pass`              char(128)             NOT NULL,
    `registration_date` datetime              NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `sessions`
(
    `id`         INT AUTO_INCREMENT PRIMARY KEY,
    `user_id`    mediumint(8) UNSIGNED NOT NULL,
    `token`      VARCHAR(64)           NOT NULL,
    `expiration` DATETIME              NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

CREATE TABLE `products`
(
    `id`          mediumint(8) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(20)           NOT NULL,
    `price`       float(8)              NOT NULL,
    `description` varchar(256)          NOT NULL
);

CREATE TABLE `inventory`
(
    `product_id`  mediumint(8) references products (id) on delete cascade,
    `quantity`    int(8)       NOT NULL,
    `description` varchar(256) NOT NULL
);

CREATE TABLE `purchase_history`
(
    `user_id`       mediumint(8) NOT NULL,
    `product_id`    mediumint(8) NOT NULL,
    `purchase_date` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO `users` (`first_name`, `last_name`, `email`, `phone`, `dob`, `status`, `role`, `pass`)
VALUES ('Artem', 'Kirichenko', 'artem.kirichenko@example.com', '1234567890', '1990-01-01', true, 1, '111'),
       ('Jane', 'Smith', 'jane.smith@example.com', '9876543210', '1992-05-15', true, 2, '111'),
       ('Alice', 'Brown', 'alice.brown@example.com', '5555555555', '1985-08-20', true, 2, '111'),
       ('Michael', 'Johnson', 'michael.johnson@example.com', '4444444444', '1980-03-10', true, 2, '111'),
       ('Emily', 'Davis', 'emily.davis@example.com', '2222222222', '1995-06-30', true, 2, '111'),
       ('Chris', 'Wilson', 'chris.wilson@example.com', '6666666666', '1988-09-15', true, 3, '111'),
       ('Sophia', 'Taylor', 'sophia.taylor@example.com', '7777777777', '1993-12-05', true, 3, '111'),
       ('Daniel', 'Moore', 'daniel.moore@example.com', '8888888888', '1987-02-25', true, 3, '111'),
       ('Olivia', 'Anderson', 'olivia.anderson@example.com', '9999999999', '1999-11-10', true, 3, '111');

INSERT INTO `products` (`name`, `price`, `description`)
VALUES ('Laptop', 850.00, 'High-performance laptop with 16GB RAM and 512GB SSD.'),
       ('Smartphone', 450.00, '5G-enabled smartphone with 128GB storage.'),
       ('Office Chair', 120.00, 'Ergonomic office chair with lumbar support.'),
       ('Printer', 200.00, 'All-in-one wireless printer with scanner and copier.'),
       ('Desk', 150.00, 'Wooden office desk with built-in cable management.'),
       ('Monitor', 300.00, '27-inch 4K monitor with adjustable stand.'),
       ('Headphones', 80.00, 'Noise-cancelling over-ear headphones.'),
       ('Keyboard', 50.00, 'Mechanical keyboard with RGB backlight.'),
       ('Coffee Machine', 100.00, 'Compact coffee machine for espresso and cappuccino.'),
       ('Webcam', 40.00, 'HD webcam with built-in microphone.');

INSERT INTO `inventory` (`product_id`, `quantity`, `description`)
VALUES (1, 10, 'Available in stock, ready for shipment.'),
       (2, 25, 'Limited stock, available for pre-order.'),
       (3, 15, 'In warehouse, shipping within 2 days.'),
       (4, 5, 'Low stock, order soon.'),
       (5, 20, 'Available in multiple colors.'),
       (6, 8, 'Ready for shipment, comes with a warranty.'),
       (7, 30, 'In stock, includes travel pouch.'),
       (8, 50, 'Bulk stock available for office supplies.'),
       (9, 12, 'Compact design, last few units remaining.'),
       (10, 40, 'High demand item, shipping within 1 day.');

INSERT INTO `purchase_history` (`user_id`, `product_id`, `purchase_date`)
VALUES (1, 3, '2024-10-01 10:15:00'),
       (2, 5, '2024-10-02 14:30:00'),
       (3, 7, '2024-10-03 09:45:00'),
       (4, 2, '2024-10-04 13:20:00'),
       (5, 1, '2024-10-05 15:50:00'),
       (6, 8, '2024-10-06 16:40:00'),
       (2, 9, '2024-10-07 11:25:00'),
       (3, 4, '2024-10-08 12:15:00'),
       (4, 6, '2024-10-09 17:35:00'),
       (1, 10, '2024-10-10 18:50:00');

INSERT INTO `bids` (`user_id`, `product_id`, `bid`)
VALUES (1, 2, 250.00),
       (2, 2, 270.50),
       (3, 2, 275.75),
       (4, 2, 280.00),
       (1, 2, 300.00);
