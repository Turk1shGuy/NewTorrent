CREATE TABLE IF NOT EXISTS commentdb (
    cid INTEGER,
    tid INTEGER,
    uid INTEGER,
    comment TEXT,
    uploadtime DATETIME,
    PRIMARY KEY (cid, tid, uid)
);

INSERT INTO commentdb (cid, tid, uid, comment, uploadtime) VALUES
(398492384982, 293849384983, 542119874321, 'This is a great torrent!', '2022-01-01 12:05:00'),
(162894839412, 293849384983, 819273654219, 'I agree, thanks for sharing!', '2022-01-01 12:10:00'),
(405873048950, 495849758282, 467832195043, 'You re welcome, happy to share!', '2022-01-01 12:15:00');
