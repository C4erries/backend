CREATE TABLE IF NOT EXISTS forms
(
	form_id SERIAL PRIMARY KEY,
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
	lang_name CHARACTER VARYING(30)
);

CREATE TABLE IF NOT EXISTS favlangs(
	form_id INTEGER,
	lang_id INTEGER,
	PRIMARY KEY (form_id, lang_id),
	FOREIGN KEY (form_id) REFERENCES forms (form_id),
	FOREIGN KEY (lang_id) REFERENCES langs (lang_id)
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
