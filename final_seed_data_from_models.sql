-- USERS


-- HOTELS
INSERT INTO hotels (name, address, inn, phone, region, view_count, total_rating, review_count, revenue)
VALUES 
  (
    'Отель Север', 'Москва, Тверская 10', '7736050003', '89267222182', 'Москва',
    0, 0, 0, 0
  ),
  (
    'Отель Юг', 'Сочи, Приморская 3', '2233445566', '89112223344', 'Краснодарский край',
    0, 0, 0, 0
  );

-- ADMIN HOTEL
INSERT INTO admin_hotels (user_id, hotel_id)
VALUES 
  (2, 1),
  (2, 2);

-- ROOMS
INSERT INTO rooms (
  hotel_id, type, description, price, capacity)
VALUES 
  (
    1, 'Стандарт', 'Номер с видом на парк', 3500, 2),
  (
    1, 'Люкс', 'Просторный номер с джакузи', 7500, 3),
  (
    2, 'Эконом', 'Дешевый номер', 2000, 2);

-- BOOKINGS
INSERT INTO bookings (
  user_id, room_id, start_date, end_date,
  guest_count, status, comment, paid
)
VALUES 
  (
    1, 1, '2025-10-01', '2025-10-02', 2,
    'new', 'Хочу с видом на парк', false
  ),
  (
    1, 2, '2025-11-10', '2025-11-12', 2,
    'paid', 'Подарок для жены', true
  );

-- PAYMENTS
INSERT INTO payments (
  booking_id, amount, status, payment_method,
  transaction_id
)
VALUES (
  2, 7500, 'success', 'mock', 'MOCK-1720000000'
);

-- REVIEWS
INSERT INTO reviews (
  hotel_id, user_id, rating, text)
VALUES 
  (
    1, 1, 5, 'Отличный отель! Всё понравилось.'),
  (
    2, 1, 4, 'Хороший сервис, но пляж далеко.');