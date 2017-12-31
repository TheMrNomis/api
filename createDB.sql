CREATE TABLE IF NOT EXISTS script_data (userID TEXT, dataName TEXT, data TEXT, PRIMARY KEY (userID, dataName));

INSERT INTO script_data VALUES ("n0m1s", "data1", '{foo:"data1"}');
INSERT INTO script_data VALUES ("n0m1s", "data2", '{foo:"data2"}');
