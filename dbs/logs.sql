CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY AUTOINCREMENT, -- Unique identifier for each log entry
    type INTEGER, -- 0 = error; 1 = info
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of the log entry
    ip_address VARCHAR(45), -- IP address of the client (supports IPv4 and IPv6)
    method VARCHAR(10), -- HTTP method (e.g., GET, POST)
    path TEXT, -- Requested URL
    agent TEXT, -- User-Agent string
    message TEXT -- Error messages (if applicable)
);
