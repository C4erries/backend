CREATE TABLE IF NOT EXISTS forms
(
username CHARACTER VARYING(40) PRIMARY KEY,
fio CHARACTER VARYING(150),
tel CHARACTER VARYING(30),
email CHARACTER VARYING(65),
birth_date DATE,
gender CHARACTER VARYING(6),
bio TEXT
);

CREATE TABLE IF NOT EXISTS langs
(
lang_id SERIAL PRIMARY KEY,
lang_name CHARACTER VARYING(30) UNIQUE
);

CREATE TABLE IF NOT EXISTS favlangs(
username CHARACTER VARYING(40),
lang_id INTEGER,
PRIMARY KEY (username, lang_id),
FOREIGN KEY (username) REFERENCES forms (username),
FOREIGN KEY (lang_id) REFERENCES langs (lang_id)
);

CREATE TABLE IF NOT EXISTS userinfo(
username CHARACTER VARYING(40) PRIMARY KEY,
enc_password CHARACTER VARYING(300),
FOREIGN KEY (username) REFERENCES forms (username)
);

INSERT INTO langs VALUES
(1, 'Prolog'),
(2, 'JavaScript'),
(3, 'PHP'),
(4, 'C++'),
(5, 'Java'),
(6, 'C#'),
(7, 'Haskell'),
(8, 'Clojure'),
(9, 'Scala'),
(10, 'Pascal'),
(11, 'Python');