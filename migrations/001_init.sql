CREATE TABLE IF NOT EXISTS banner_clicks (
    timestamp TIMESTAMP NOT NULL,
    banner_id INT NOT NULL,
    count INT NOT NULL,
    PRIMARY KEY (timestamp, banner_id)
);
