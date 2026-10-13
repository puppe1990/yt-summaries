-- 002_videos.sql — resumos de vídeos do YouTube (yt-resumos-zap).
-- up

CREATE TABLE IF NOT EXISTS videos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    youtube_id TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    channel TEXT NOT NULL DEFAULT '',
    url TEXT NOT NULL DEFAULT '',
    thumbnail TEXT NOT NULL DEFAULT '',
    published_at TEXT,
    imported_at TEXT NOT NULL DEFAULT (datetime('now')),
    tags TEXT NOT NULL DEFAULT '',
    search_text TEXT NOT NULL DEFAULT '',
    payload TEXT NOT NULL
);

-- down

DROP TABLE IF EXISTS videos;
