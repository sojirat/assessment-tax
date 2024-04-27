CREATE TABLE IF NOT EXISTS setting (
    id INT PRIMARY KEY,
    personal FLOAT(4),
    k_receipt FLOAT(4)
);

INSERT INTO setting (id, personal, k_receipt) VALUES
    (1, 60000, 50000)
ON CONFLICT (id) DO NOTHING;
