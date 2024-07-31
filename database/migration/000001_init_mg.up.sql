CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    is_admin TINYINT(1) DEFAULT 0
);

CREATE TABLE books (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(17) UNIQUE NOT NULL,
    total_copies INT NOT NULL DEFAULT 1,
    checked_out_copies INT NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED,
    book_id BIGINT UNSIGNED,
    transaction_type VARCHAR(10) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP,
    return_date TIMESTAMP,
    fine DECIMAL(10, 2) DEFAULT 0,
    status VARCHAR(255) DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id)
);

CREATE TABLE admin_requests (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED,
    status ENUM('pending', 'approved', 'denied') DEFAULT 'pending',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO books (title, author, isbn) 
VALUES
  ('Pride and Prejudice', 'Jane Austen', '978-0140439516'),
  ('To Kill a Mockingbird', 'Harper Lee', '978-0446310727'),
  ('The Lord of the Rings', 'J.R.R. Tolkien', '978-0261102880'),
  ('The Hitchhiker\'s Guide to the Galaxy', 'Douglas Adams', '978-0345391803'),
  ('One Hundred Years of Solitude', 'Gabriel Garcia Marquez', '978-0307472708'),
  ('Frankenstein', 'Mary Shelley', '978-0140621820'),
  ('The Great Gatsby', 'F. Scott Fitzgerald', '978-0743273565'),
  ('Invisible Man', 'Ralph Ellison', '978-0679740329'),
  ('Jane Eyre', 'Charlotte Brontë', '978-0140620674'),
  ('Dune', 'Frank Herbert', '978-0441014041'),
  ('Beloved', 'Toni Morrison', '978-0449906167'),
  ('The Catcher in the Rye', 'J.D. Salinger', '978-0312901084'),
  ('1984', 'George Orwell', '978-0451524935'),
  ('The God of Small Things', 'Arundhati Roy', '978-0224057991'),
  ('Things Fall Apart', 'Chinua Achebe', '978-0435900559'),
  ('The Handmaid\'s Tale', 'Margaret Atwood', '978-0345359943'),
  ('The Book Thief', 'Markus Zusak', '978-0375838907'),
  ('Americanah', 'Chimamanda Ngozi Adichie', '978-1416568797'),
  ('A Short History of Nearly Everything', 'Bill Bryson', '978-0571220728'),
  ('Sapiens: A Brief History of Humankind', 'Yuval Noah Harari', '978-0062450395');
