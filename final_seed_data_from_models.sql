-- USERS


-- HOTELS
INSERT INTO hotels (
  name, address, inn, phone, region,
  photo_urls, amenities,
  view_count, total_rating, review_count, revenue
)
VALUES
  (
    'Отель Север', 'Москва, Тверская 10', '7736050003', '89267222182', 'Москва',
    '["https://avatars.mds.yandex.net/get-altay/374295/2a0000015b1791c294a596d0698883284a49/XXXL"]', '["wifi", "parking"]',
    0, 0, 0, 0
  ),
  (
    'Отель Юг', 'Сочи, Приморская 3', '2233445566', '89112223344', 'Краснодарский край',
    '["https://cdn.worldota.net/t/x500/content/dc/8c/dc8ca038fcd1385d03d04a4bf7121816cedecf0e.jpeg"]', '["pool"]',
    0, 0, 0, 0
  );

-- ADMIN HOTEL
INSERT INTO admin_hotels (
  user_id, hotel_id, photo_urls, amenities
)
VALUES
  (2, 1, '[]', '[]'),
  (2, 2, '[]', '[]');

-- ROOMS
INSERT INTO rooms (
  hotel_id, type, description, price, capacity,
  photo_urls, amenities
)
VALUES
  (
    1, 'Стандарт', 'Номер с видом на парк', 3500, 2,
    '["https://avatars.mds.yandex.net/get-altay/216588/2a0000015b1717eba4aa8df0d55d8b4b0ddf/XXXL"]', '["tv", "minibar"]'
  ),
  (
    1, 'Люкс', 'Просторный номер с джакузи', 7500, 3,
    '["https://avatars.mds.yandex.net/i?id=c31ea19d4febbbb1cd5f860367b160cf70556746-13219587-images-thumbs&n=13"]', '["jacuzzi", "balcony"]'
  ),
  (
    2, 'Эконом', 'Дешевый номер', 2000, 2,
    '["https://sovcominvest.ru/uploads/photo/12394/image/c66f6053ce97df86b5fb25f7412b31ca.jpg"]', '[]'
  );

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
  hotel_id, user_id, rating, text,
  photo_urls, amenities
)
VALUES
  (
    1, 1, 5, 'Отличный отель! Всё понравилось.',
    '["/static/review1.jpg"]', '[]'
  ),
  (
    2, 1, 4, 'Хороший сервис, но пляж далеко.',
    '[]', '["wifi"]'
  );