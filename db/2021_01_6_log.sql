CREATE TABLE IF NOT EXISTS log(
    id int NOT NULL AUTO_INCREMENT,
    ip VARCHAR(20) NOT NULL,
    visited_date VARCHAR(40) NOT NULL,
    status_code VARCHAR(4) NOT NULL,
    visited_url VARCHAR(255) NOT NULL,
    protocol_status VARCHAR(10) NOT NULL,
    server_response VARCHAR(255) NOT NULL,
    send_bytes VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
)