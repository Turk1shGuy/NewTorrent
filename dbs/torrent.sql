CREATE TABLE IF NOT EXISTS torrentdb (
    tid INTEGER,
    uid INTEGER,
    name VARCHAR(255), 
    cat INTEGER,
    link VARCHAR(3000), 
    desc TEXT,
    uploadtime DATETIME,

    PRIMARY KEY (tid, uid)
);
