-- Sample SQL commands for testing the RDBMS

-- Create a users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT,
    age INTEGER,
    email TEXT UNIQUE
);

-- Insert some sample data
INSERT INTO users VALUES (1, 'John Doe', 25, 'john@example.com');
INSERT INTO users VALUES (2, 'Jane Smith', 30, 'jane@example.com');
INSERT INTO users VALUES (3, 'Bob Johnson', 35, 'bob@example.com');

-- Create a posts table
CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    title TEXT,
    content TEXT
);

-- Insert sample posts
INSERT INTO posts VALUES (1, 1, 'Hello World', 'This is my first post!');
INSERT INTO posts VALUES (2, 2, 'Database Fun', 'Learning about databases is awesome.');
INSERT INTO posts VALUES (3, 1, 'Go Programming', 'Go is a great language for systems programming.');

-- Query examples

-- Select all users
SELECT * FROM users;

-- Select specific columns
SELECT name, age FROM users;

-- Select with WHERE condition
SELECT * FROM users WHERE age > 25;

-- Update a user
UPDATE users SET age = 26 WHERE id = 1;

-- Delete a user
DELETE FROM users WHERE id = 3;

-- Join example
SELECT users.name, posts.title FROM users JOIN posts ON users.id = posts.user_id;