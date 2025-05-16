CREATE TABLE IF NOT EXISTS userdb(
    uid INTEGER PRIMARY KEY UNIQUE,
    email VARCHAR(70) UNIQUE,
    passhash VARCHAR(70),  -- sha256
    uname VARCHAR(20) UNIQUE,
    lang VARCHAR(5), 
    rdate DATETIME,
    ldate DATETIME 
);


INSERT INTO userdb (uid, email, passhash, uname, lang, rdate, ldate) VALUES
(0, 'admin@localhost', 'passhash1', 'Admin', "TR", '2010-10-10 00:00:00', '2010-10-10 00:00:00'), -- root
(214567893021, 'emily.chen@gmail.com', 'passhash123', 'echen', "EN", '2018-07-12 14:30:00', '2022-01-01 10:00:00'),
(463729138574, 'liam.jackson@outlook.com', 'passhash456', 'ljackson', "EN", '2020-03-25 18:45:00', '2022-06-01 12:00:00'),
(982734659821, 'soo.jin.lee@tutamail.com', 'passhash789', 'sjlee', "KR", '2019-02-20 09:15:00', '2022-03-01 15:00:00'),
(118237464935, 'yui.nakamura@protonmail.com', 'passhash012', 'ynakamura', "JP", '2017-11-15 11:50:00', '2022-09-01 13:00:00'),
(657382194753, 'maria.rodriguez@gmail.com', 'passhash345', 'mrodriguez', "ES", '2021-05-10 16:20:00', '2022-11-01 10:30:00'),
(982135749823, 'david.kim@outlook.com', 'passhash678', 'dkim', "EN", '2019-08-22 10:10:00', '2022-04-01 14:00:00'),
(135792468021, 'sophia.park@tutamail.com', 'passhash234', 'spark', "KR", '2020-01-05 12:40:00', '2022-07-01 11:00:00'),
(654328917463, 'kenichi.tanaka@protonmail.com', 'passhash567', 'ktanaka', "JP", '2018-04-18 15:30:00', '2022-10-01 12:30:00'),
(819273654219, 'olivia.brown@gmail.com', 'passhash890', 'obrown', "EN", '2021-09-15 17:00:00', '2022-05-01 13:00:00'),
(467832195043, 'noah.davis@outlook.com', 'passhash111', 'ndavis', "EN", '2019-06-01 10:50:00', '2022-08-01 10:00:00'),
(123456789012, 'ava.garcia@tutamail.com', 'passhash222', 'agarcia', "ES", '2020-09-20 11:30:00', '2022-02-01 14:30:00'),
(987654321098, 'ethan.hall@protonmail.com', 'passhash333', 'ehall', "EN", '2018-01-15 12:10:00', '2022-06-01 11:00:00'),
(111111111111, 'isabella.martin@gmail.com', 'passhash444', 'imartin', "EN", '2021-03-10 15:40:00', '2022-09-01 10:00:00'),
(222222222222, 'logan.brooks@outlook.com', 'passhash555', 'lbrooks', "EN", '2019-11-25 10:20:00', '2022-04-01 12:00:00'),
(333333333333, 'sophie.patel@tutamail.com', 'passhash666', 'spatel', "EN", '2020-06-15 14:50:00', '2022-07-01 13:00:00'),
(444444444444, 'mason.kim@protonmail.com', 'passhash777', 'mkim', "KR", '2018-08-20 11:40:00', '2022-10-01 11:30:00'),
(555555555555, 'charlotte.lee@gmail.com', 'passhash888', 'clee', "EN", '2021-01-05 12:30:00', '2022-05-01 14:00:00');
