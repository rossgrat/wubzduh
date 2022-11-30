--Create artist table
CREATE TABLE IF NOT EXISTS Artists (
    ID SERIAL PRIMARY KEY,
    ArtistName TEXT NOT NULL,
    SpotifyID TEXT NOT NULL
);
--Create album table
CREATE TABLE IF NOT EXISTS Albums (
    ID SERIAL PRIMARY KEY,
    AlbumTitle TEXT NOT NULL,
    SpotifyID TEXT NOT NULL,
    ArtistID INT REFERENCES Artists(ID) NOT NULL,
    CoverartURL TEXT NOT NULL,
    ReleaseDate TIMESTAMP NOT NULL,
    AlbumType TEXT NOT NULL,
    AlbumUrl TEXT NOT NULL
);
--Create track table
CREATE TABLE IF NOT EXISTS tracks (
    ID SERIAL PRIMARY KEY,
    TrackTitle TEXT NOT NULL,
    SpotifyID TEXT NOT NULL,
    AlbumID INT REFERENCES Albums(ID) NOT NULL,
    LengthMS INT NOT NULL,
    TrackNumber INT NOT NULL,
    AddedToPlaylist BOOLEAN NOT NULL
);

INSERT INTO Artists 
(ArtistName, SpotifyID) 
VALUES
('Tipper', '1soJ22UMyjIw6SYFtoFJwe'),
('Bob Moses', '6LHsnRBUYhFyt01PdKXAF5'),
('Tycho', '5oOhM2DFWab8XhSdQiITry'),
('Bassnectar', '1JPy5PsJtkhftfdr6saN2i'),
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