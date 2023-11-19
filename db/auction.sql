
CREATE DATABASE IF NOT EXISTS auction_db;
-- Example SQL migration script for creating the ad table
CREATE TABLE IF NOT EXISTS ad (
    id INT AUTO_INCREMENT PRIMARY KEY,
    text VARCHAR(255),
    base_price INT,
    start_time DATETIME,
    end_time DATETIME
);


-- Example SQL migration script for creating the bid table
CREATE TABLE bid (
    bidder_id INT,
    ad_id INT,
    price INT, 
    PRIMARY KEY (bidder_id, ad_id),
    FOREIGN KEY (ad_id) REFERENCES ad(id)
);

-- Example SQL migration script for creating the auction table
CREATE TABLE auction (
    ad_id INT PRIMARY KEY,
    start_time DATETIME,
    end_time DATETIME,
    state VARCHAR(255),
    FOREIGN KEY (ad_id) REFERENCES ad(id)
);
