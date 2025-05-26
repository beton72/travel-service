select * from hotels

select * from users

SELECT * FROM admin_hotels;

SELECT * FROM rooms;

SELECT * FROM Bookings;

SELECT * FROM reviews;

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

SELECT 
  h.id, h.name, 
  AVG(r.price) as avg_price
FROM hotels h
JOIN rooms r ON h.id = r.hotel_id
GROUP BY h.id
HAVING AVG(r.price) BETWEEN 3000 AND 10000;