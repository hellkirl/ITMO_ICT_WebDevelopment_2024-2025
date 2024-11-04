CREATE TABLE hotels (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        address VARCHAR(255) NOT NULL,
                        city VARCHAR(100) NOT NULL,
                        country VARCHAR(100) NOT NULL,
                        phone VARCHAR(20),
                        email VARCHAR(100),
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_hotels_city ON hotels(city);
CREATE INDEX idx_hotels_country ON hotels(country);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE rooms (
                       id SERIAL PRIMARY KEY,
                       hotel_id INTEGER REFERENCES hotels(id) ON DELETE CASCADE,
                       number VARCHAR(10) NOT NULL,
                       type VARCHAR(20) NOT NULL,
                       price DECIMAL(10, 2) NOT NULL,
                       is_available BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_rooms_number_hotel ON rooms(number, hotel_id);
CREATE INDEX idx_rooms_hotel_id ON rooms(hotel_id);
CREATE INDEX idx_rooms_type ON rooms(type);
CREATE INDEX idx_rooms_is_available ON rooms(is_available);

CREATE TABLE reservations (
                              id SERIAL PRIMARY KEY,
                              user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              room_id INTEGER REFERENCES rooms(id) ON DELETE CASCADE,
                              check_in DATE NOT NULL,
                              check_out DATE NOT NULL,
                              status VARCHAR(20) NOT NULL DEFAULT 'reserved',
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_room_id ON reservations(room_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_check_in ON reservations(check_in);
CREATE INDEX idx_reservations_check_out ON reservations(check_out);

CREATE TABLE reviews (
                         id SERIAL PRIMARY KEY,
                         reservation_id INTEGER REFERENCES reservations(id) ON DELETE CASCADE,
                         rating INTEGER CHECK (rating >= 1 AND rating <= 10),
                         comment TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reviews_reservation_id ON reviews(reservation_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);