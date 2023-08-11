package database

func CreateProductTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    amount integer,
    name text UNIQUE,
    description text,
    category text NOT NULL
)
`)
}

func CreateUserTable() {
	DB.Query(`CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY,
  user_name VARCHAR(50), 
  email VARCHAR(100),
  is_active BOOLEAN
  is_email_verified BOOLEAN
);
`)
}
