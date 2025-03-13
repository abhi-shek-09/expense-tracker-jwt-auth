CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);


CREATE TABLE expenses (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  amount DECIMAL(10,2) NOT NULL,
  category TEXT CHECK (category IN ('Groceries', 'Leisure', 'Electronics', 'Utilities', 'Clothing', 'Health', 'Others')),
  description TEXT,
  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
