CREATE TABLE IF NOT EXISTS music_library (
    song_id SERIAL PRIMARY KEY,
    group_name TEXT NOT NULL,
    song_name TEXT NOT NULL,
    release_date DATE DEFAULT NULL,
    song_lyrics TEXT DEFAULT NULL,
    link TEXT DEFAULT NULL,
    CONSTRAINT unique_group_song_combination UNIQUE (group_name, song_name)
);