CREATE TABLE songs (
    id SERIAL PRIMARY KEY,            -- Уникальный идентификатор песни
    group_name VARCHAR(255) NOT NULL, -- Название группы или исполнителя
    song_name VARCHAR(255) NOT NULL   -- Название песни
);

CREATE TABLE song_details (
    song_id INT REFERENCES songs(id) ON DELETE CASCADE, -- Ссылка на таблицу songs
    release_date DATE,                                  -- Дата выпуска песни
    text TEXT,                                          -- Текст песни
    link TEXT,                                          -- Ссылка на песню в Genius API
    artist VARCHAR(255),                                -- Исполнитель
    album VARCHAR(255),                                 -- Альбом
    genre VARCHAR(100),                                 -- Жанр
    duration VARCHAR(50),                               -- Продолжительность
    key VARCHAR(10),                                   -- Ключ
    tempo VARCHAR(10),                                  -- Темп
    PRIMARY KEY (song_id)                               -- song_id также уникален в этой таблице
);
