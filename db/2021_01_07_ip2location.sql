CREATE TABLE IF NOT EXISTS ip2location(
    id INT AUTO_INCREMENT NOT NULL,
    ip VARCHAR(20) NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    country_code VARCHAR(25) NOT NULL,
    country VARCHAR(50) NOT NULL,
    postalCode VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    PRIMARY KEY(id),
    UNIQUE(ip)
)