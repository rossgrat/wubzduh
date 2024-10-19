--Create artist table
CREATE TABLE IF NOT EXISTS artists (
    id SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    UNIQUE(spotify_id)
);
--Create album table
CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    artist_id INT REFERENCES artists(id) NOT NULL,
    cover_art_url TEXT NOT NULL,
    release_date DATE NOT NULL,
    "type" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    UNIQUE(spotify_id)
);
--Create track table
CREATE TABLE IF NOT EXISTS tracks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    spotify_id TEXT NOT NULL,
    album_id INT REFERENCES albums(id) NOT NULL,
    length_ms INT NOT NULL,
    "number" INT NOT NULL,
    UNIQUE(spotify_id)
);

INSERT INTO artists 
("name", spotify_id) 
VALUES
('Tipper', '1soJ22UMyjIw6SYFtoFJwe'),
('Bob Moses', '6LHsnRBUYhFyt01PdKXAF5'),
('Tycho', '5oOhM2DFWab8XhSdQiITry'),
('Supertask', '47qa2xx9Xlw1oGkKbMq8Zt'),
('Lab Group ', '4VSPQ1ufWQpHYbIIbRguSV'),
('CharlestheFirst', '2FTj5ijy8lP59d2V9dHR6I'),
('ZHU', '28j8lBWDdDSHSSt5oPlsX2'),
('GRiZ', '25oLRSUjJk4YHNUsQXk7Ut'),
('CloZee', '1496XxkytEk26FUJLfpVZr'),
('TOKiMONSTA', '3VwKSHAfgzV1DOHV0aANCI'),
('LP Giobbi', '3oKnyRhYWzNsTiss5n4Z1J'),
('Sofi Tukker', '586uxXMyD5ObPuzjtrzO1Q'),
('Flying Lotus', '29XOeO6KIWxGthejQqn793'),
('ODESZA', '21mKp7DqtSNHhCAU2ugvUw');

CREATE ROLE "server";
CREATE DATABASE "server";

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "server";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "server";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO "server";
GRANT ALL PRIVILEGES ON SCHEMA public TO "server";

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO "server";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO "server";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON FUNCTIONS TO "server";