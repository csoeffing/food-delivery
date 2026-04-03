-- menus

TRUNCATE TABLE menus RESTART IDENTITY;

INSERT INTO menus (created_at, updated_at, name, menu_image) VALUES (NOW(), NOW(), 'Dinner Menu', '');

-- profiles

TRUNCATE TABLE profiles RESTART IDENTITY;

INSERT INTO profiles (created_at, updated_at, user_id, first_name, last_name, profile_image, user_type, pro_type, user_name) VALUES
(NOW(), NOW(), 1, 'George', 'Washington', '', 'CLIENT', 'RIDER', 'gwashington');

-- restaurants

TRUNCATE TABLE restaurants RESTART IDENTITY;

INSERT INTO restaurants (created_at, updated_at, restaurant_image, restaurant_name, phone_number, address, location_id, profile_id, registration_status)
VALUES (NOW(), NOW(), '', 'Tony''s Pizza', '800-555-1234', '101 Main St', 1, 1, 'ACCEPTED');

-- foods

TRUNCATE TABLE foods RESTART IDENTITY;

INSERT INTO foods (created_at, updated_at, name, description, price, food_image, menu_id, restaurant_id, status) VALUES
(NOW(), NOW(), 'Pizza', 'Pepperoni Pizza', 15.75, '', 1, 1, 't');

-- users

TRUNCATE TABLE users RESTART IDENTITY;

-- inserts occur from shell script with POSTs