USE isys2099_group9_app;


-- Warehouse admin
INSERT INTO wh_admin (username, refresh_token, password_hash)
VALUES ('admin', '', 'admin');


-- Lazada user (buyer and seller)
INSERT INTO lazada_user (username, refresh_token, password_hash)
VALUES ('tony', '', 'tony'),
       ('loi', '', 'loi'),
       ('mike', '', 'mike');

INSERT INTO buyer (username)
VALUES ('tony');

INSERT INTO seller (username, shop_name, city)
VALUES ('loi', 'loi', 'HCM'),
       ('mike', 'vo', 'HCM');


-- Product category and attribute

/*
 product_category
 ├── Electronic Devices
 │   ├── Mobiles
 │   ├── Computer
 │   ├── Computer Peripherals
 │   │   ├── Computer Mice
 │   │   ├── Keyboards
 │   │   │   └── Mechanical Keyboards
 │   │   ├── Mechanical Keyboard Key-caps
 │   │   ├── Mechanical Keyboard Switches
 │   │   └── Monitors
 │   └── Audio
 │       ├── Speaker
 │       ├── Earphones
 │       └── Headphones
 └── TV & Home Appliances
     ├── Televisions
     ├── Small Appliances
     └── Large Appliances
 */
INSERT INTO product_category (category_name, parent)
VALUES ('Electronic Devices', null),
       ('Mobiles', 'Electronic Devices'),
       ('Computer', 'Electronic Devices'),
       ('Computer Peripherals', 'Electronic Devices'),
       ('Computer Mice', 'Computer Peripherals'),
       ('Keyboards', 'Computer Peripherals'),
       ('Mechanical Keyboards', 'Keyboards'),
       ('Mechanical Keyboard Key-caps', 'Computer Peripherals'),
       ('Mechanical Keyboard Switches', 'Computer Peripherals'),
       ('Monitors', 'Computer Peripherals'),
       ('Audio', 'Electronic Devices'),
       ('Speaker', 'Audio'),
       ('Earphones', 'Audio'),
       ('Headphones', 'Audio'),
       ('TV & Home Appliances', null),
       ('Televisions', 'TV & Home Appliances'),
       ('Small Appliances', 'TV & Home Appliances'),
       ('Large Appliances', 'TV & Home Appliances');

INSERT INTO product_attribute (attribute_name, attribute_type, required)
VALUES ('Width', 'Number', FALSE),
       ('Length', 'Number', FALSE),
       ('Height', 'Number', FALSE),
       ('Thickness', 'Number', FALSE),
       ('Color', 'String', FALSE),
       ('RAM', 'String', TRUE),
       ('CPU', 'String', TRUE),
       ('GPU', 'String', FALSE),
       ('Storage capacity', 'String', FALSE),
       ('Brand', 'String', FALSE),
       ('Screen size', 'String', FALSE),
       ('Resolution', 'String', FALSE),
       ('Refresh rate', 'String', FALSE),
       ('Wireless', 'Boolean', FALSE),
       ('Battery powered', 'Boolean', FALSE),
       ('Material', 'String', FALSE),
       ('Key-cap profile', 'String', TRUE),
       ('Switch type', 'String', TRUE),
       ('Latency', 'String', FALSE),
       ('Noise cancelling', 'Boolean', FALSE),
       ('Frequency range', 'String', FALSE),
       ('Smart TV', 'Boolean', FALSE),
       ('Voltage', 'String', TRUE),
       ('Wattage', 'String', TRUE),
       ('Warranty', 'Boolean', TRUE),
       ('Warranty period', 'String', FALSE);

INSERT INTO product_category_attribute_association (category, attribute)
VALUES ('Electronic Devices', 'Warranty'),
       ('Electronic Devices', 'Warranty period'),
       ('Electronic Devices', 'Brand'),
       ('Mobiles', 'Width'),
       ('Mobiles', 'Length'),
       ('Mobiles', 'Thickness'),
       ('Mobiles', 'Color'),
       ('Mobiles', 'RAM'),
       ('Mobiles', 'CPU'),
       ('Mobiles', 'Storage capacity'),
       ('Mobiles', 'Resolution'),
       ('Computer', 'Width'),
       ('Computer', 'Length'),
       ('Computer', 'Thickness'),
       ('Computer', 'Color'),
       ('Computer', 'RAM'),
       ('Computer', 'CPU'),
       ('Computer', 'GPU'),
       ('Computer', 'Storage capacity'),
       ('Computer', 'Screen size'),
       ('Computer', 'Resolution'),
       ('Computer Peripherals', 'Wireless'),
       ('Computer Peripherals', 'Battery powered'),
       ('Computer Mice', 'Color'),
       ('Computer Mice', 'Latency'),
       ('Keyboards', 'Color'),
       ('Keyboards', 'Latency'),
       ('Mechanical Keyboards', 'Key-cap profile'),
       ('Mechanical Keyboard Switches', 'Switch type'),
       ('Mechanical Keyboard Key-caps', 'Key-cap profile'),
       ('Mechanical Keyboard Key-caps', 'Material'),
       ('Monitors', 'Screen size'),
       ('Monitors', 'Resolution'),
       ('Monitors', 'Refresh rate'),
       ('Audio', 'Wireless'),
       ('Audio', 'Battery powered'),
       ('Audio', 'Latency'),
       ('Audio', 'Noise cancelling'),
       ('Audio', 'Frequency range'),
       ('TV & Home Appliances', 'Warranty'),
       ('TV & Home Appliances', 'Warranty period'),
       ('TV & Home Appliances', 'Brand'),
       ('TV & Home Appliances', 'Color'),
       ('TV & Home Appliances', 'Voltage'),
       ('TV & Home Appliances', 'Wattage'),
       ('Televisions', 'Smart TV'),
       ('Televisions', 'Screen size'),
       ('Televisions', 'Resolution'),
       ('Televisions', 'Width'),
       ('Televisions', 'Height'),
       ('Televisions', 'Thickness'),
       ('Small Appliances', 'Width'),
       ('Small Appliances', 'Length'),
       ('Small Appliances', 'Height'),
       ('Large Appliances', 'Width'),
       ('Large Appliances', 'Length'),
       ('Large Appliances', 'Height');


-- Warehouse and products
INSERT INTO warehouse (warehouse_name, volume, province, city, district, street, street_number)
VALUES ('Toronto LAZ', 5000000, 'Ontario', 'Toronto', NULL, NULL, NULL),
       ('Montreal LAZ', 4600000, 'Quebec', 'Montreal', NULL, NULL, NULL),
       ('Vancouver LAZ', 4000000, 'British Columbia', 'Vancouver', NULL, NULL, NULL),
       ('Victoria LAZ', 1800000, 'British Columbia', 'Victoria', NULL, NULL, NULL),
       ('Winnipeg LAZ', 3000000, 'Manitoba', 'Winnipeg', NULL, NULL, NULL),
       ('Edmonton LAZ', 2000000, 'Alberta', 'Edmonton', NULL, NULL, NULL),
       ('St. John\'s LAZ', 1200000, 'Newfoundland and Labrador', 'St. John\'s', NULL, NULL, NULL),
       ('Regina LAZ', 2200000, 'Saskatchewan', 'Regina', NULL, NULL, NULL),
       ('Charlottetown LAZ', 800000, 'Prince Edward Island', 'Charlottetown', NULL, NULL, NULL);

INSERT INTO product (id, image, title, product_description, category, price, width, length, height, seller)
VALUES (1, 'http://54.145.255.207/uploads/mock-data/1.jpg', 'Smartphone Model X', 'High-end smartphone with top-notch features.', 'Mobiles', 899.99, 6, 12, 1, 'loi'),
       (2, 'http://54.145.255.207/uploads/mock-data/2.jpg', 'Gaming Laptop Pro', 'Powerful gaming laptop for enthusiasts.', 'Computer', 1499.99, 35, 25, 2, 'loi'),
       (3, 'http://54.145.255.207/uploads/mock-data/3.jpg', 'RGB Mechanical Keyboard', 'Mechanical keyboard with customizable RGB lighting.', 'Mechanical Keyboards', 89.99, 44, 14, 4, 'loi'),
       (4, 'http://54.145.255.207/uploads/mock-data/4.jpg', 'Membrane Keyboard', 'Standard membrane keyboard for daily use.', 'Keyboards', 39.99, 42, 15, 3, 'loi'),
       (5, 'http://54.145.255.207/uploads/mock-data/5.jpg', '27-inch Gaming Monitor', 'High-refresh rate gaming monitor for smooth gameplay.', 'Monitors', 399.99, 60, 33, 5, 'loi'),
       (6, 'http://54.145.255.207/uploads/mock-data/6.jpg', 'Wireless Bluetooth Speaker', 'Portable speaker with long battery life.', 'Speaker', 49.99, 12, 8, 8, 'loi'),
       (7, 'http://54.145.255.207/uploads/mock-data/7.jpg', 'Noise-Cancelling Headphones', 'Over-ear headphones with active noise cancellation.', 'Headphones', 199.99, 18, 15, 25, 'loi'),
       (8, 'http://54.145.255.207/uploads/mock-data/8.jpg', 'Smart LED TV', 'High-definition smart TV with various features.', 'Televisions', 799.99, 48, 28, 3, 'loi'),
       (9, 'http://54.145.255.207/uploads/mock-data/9.jpg', 'Compact Refrigerator', 'Small refrigerator for dorm rooms and offices.', 'Small Appliances', 179.99, 20, 18, 33, 'loi'),
       (10, 'http://54.145.255.207/uploads/mock-data/10.jpg', 'Front-Load Washing Machine', 'Energy-efficient washing machine with multiple programs.', 'Large Appliances', 599.99, 30, 27, 39, 'loi'),
       (11, 'http://54.145.255.207/uploads/mock-data/11.jpg', 'High-Performance Desktop PC', 'Custom-built desktop PC for gaming and heavy tasks.', 'Computer', 2499.99, 45, 50, 25, 'mike'),
       (12, 'http://54.145.255.207/uploads/mock-data/12.jpg', 'Ultra-Narrow Curved Monitor', 'Curved narrow monitor for disruptive viewing experience.', 'Monitors', 599.99, 80, 35, 10, 'mike'),
       (13, 'http://54.145.255.207/uploads/mock-data/13.jpg', 'Wireless Gaming Mouse', 'Precision gaming mouse with customizable buttons.', 'Computer Mice', 69.99, 12, 7, 4, 'mike'),
       (14, 'http://54.145.255.207/uploads/mock-data/14.jpg', 'Bluetooth Earbuds', 'True wireless earbuds with advanced audio quality.', 'Earphones', 119.99, 5, 3, 2, 'mike'),
       (15, 'http://54.145.255.207/uploads/mock-data/15.jpg', 'Portable Air Conditioner', 'Cool down your space with this portable AC unit.', 'Small Appliances', 299.99, 18, 15, 30, 'mike'),
       (16, 'http://54.145.255.207/uploads/mock-data/16.jpg', 'Low-Performance Desktop PC', 'Custom-built desktop potato.', 'Computer', 2499.99, 45, 50, 25, 'mike'),
       (17, 'http://54.145.255.207/uploads/mock-data/17.jpg', 'Ultra-Wide Curved Monitor', 'Curved monitor for immersive viewing experience.', 'Monitors', 599.99, 80, 35, 10, 'mike'),
       (18, 'http://54.145.255.207/uploads/mock-data/18.jpg', 'Wired Gaming Mouse', 'Low-precision gaming mouse with customizable buttons.', 'Computer Mice', 69.99, 12, 7, 4, 'mike'),
       (19, 'http://54.145.255.207/uploads/mock-data/19.jpg', 'Green-toe Earbuds', 'True wireless earbuds with advanced audio quality.', 'Earphones', 119.99, 5, 3, 2, 'mike'),
       (20, 'http://54.145.255.207/uploads/mock-data/20.jpg', 'Importable Air Conditioner', 'Cool down your space with this immovable AC unit.', 'Small Appliances', 299.99, 18, 15, 30, 'mike'),
       (21, 'http://54.145.255.207/uploads/mock-data/21.jpg', 'Office Keyboard', 'Mechanical office keyboard with customizable lighting.', 'Keyboards', 79.99, 42, 15, 4, 'mike'),
       (22, 'http://54.145.255.207/uploads/mock-data/22.jpg', 'Smartphone Model Y', 'Mid-range smartphone with sleek design.', 'Mobiles', 499.99, 5, 11, 1, 'mike'),
       (23, 'http://54.145.255.207/uploads/mock-data/23.jpg', 'Noise-Boosting Over-Ear Headphones', 'Premium over-ear headphones with cutting-edge noise boosting.', 'Headphones', 249.99, 18, 16, 25, 'mike'),
       (24, 'http://54.145.255.207/uploads/mock-data/24.jpg', 'Low Resolution Smart TV', '320p TV with mediocre visuals.', 'Televisions', 999.99, 55, 32, 5, 'mike'),
       (25, 'http://54.145.255.207/uploads/mock-data/25.jpg', 'Robotic Vacuum Dust Depositor', 'Automated dust depositor for a dirtier than ever house.', 'Small Appliances', 149.99, 14, 14, 4, 'mike'),
       (26, 'http://54.145.255.207/uploads/mock-data/26.jpeg', 'A Cube', 'A mysterious dense 1x1x1 centimeter cube of exotic material. Where is it from? Why is it here? What is its purpose? No one knows...', 'Large Appliances', 9999.99, 1, 1, 1, 'mike');

INSERT INTO stockpile (product_id, warehouse_id, quantity)
VALUES (1, 1, 100),
       (2, 2, 50),
       (3, 3, 30),
       (4, 5, 40),
       (5, 1, 20),
       (6, 4, 80),
       (7, 2, 40),
       (8, 6, 15),
       (9, 8, 25),
       (10, 3, 10),
       (11, 3, 15),
       (12, 5, 10),
       (13, 7, 30),
       (14, 1, 50),
       (15, 6, 5),
       (16, 3, 15),
       (17, 5, 10),
       (18, 7, 30),
       (19, 1, 50),
       (20, 6, 5),
       (21, 3, 20),
       (22, 1, 75),
       (23, 2, 25),
       (24, 6, 10),
       (25, 8, 15),
       (26, 9, 1);



-- Orders
INSERT INTO inbound_order (quantity, product_id, created_date, created_time, fulfilled_date, fulfilled_time, seller)
VALUES (50, 1, '2022-08-01', '10:00:00', '2022-08-02', '14:30:00', 'loi'),
       (30, 2, '2022-08-02', '14:30:00', '2022-08-03', '09:45:00', 'loi'),
       (20, 3, '2022-08-03', '09:45:00', null, null, 'loi'),
       (15, 4, '2022-08-04', '13:20:00', null, null, 'loi'),
       (60, 5, '2022-08-05', '16:15:00', null, null, 'loi'),
       (50, 11, '2022-08-01', '10:00:00', '2022-08-02', '14:30:00', 'mike'),
       (30, 12, '2022-08-02', '14:30:00', '2022-08-03', '09:45:00', 'mike'),
       (20, 13, '2022-08-03', '09:45:00', null, null, 'mike'),
       (15, 14, '2022-08-04', '13:20:00', null, null, 'mike'),
       (60, 15, '2022-08-05', '16:15:00', null, null, 'mike');

INSERT INTO buyer_order (quantity, product_id, created_date, created_time, order_status, fulfilled_date, fulfilled_time, buyer)
VALUES (3, 1, '2023-08-06', '11:30:00', 'A', '2023-08-07', '15:45:00', 'tony'),
       (1, 4, '2023-08-07', '15:45:00', 'A', '2023-08-08', '10:20:00', 'tony'),
       (2, 2, '2023-08-08', '10:20:00', 'A', '2023-08-09', '12:00:00', 'tony'),
       (5, 3, '2023-08-09', '12:00:00', 'R', null, null, 'tony'),
       (4, 5, '2023-08-10', '14:10:00', 'P', null, null, 'tony');
