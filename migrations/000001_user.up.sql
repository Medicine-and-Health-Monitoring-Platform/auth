CREATE TABLE IF NOT EXISTS Users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       phone_number VARCHAR(20),
                       role VARCHAR(20) NOT NULL CHECK (role IN ('patient', 'doctor,', 'admin')),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       deleted_at bigint default 0
);
INSERT INTO Users (id, email, password_hash, first_name, last_name, phone_number, role)
VALUES
    ('62f41dd1-6d38-48f8-a70d-07c91c1ef196', 'john.doe@example.com', '$2a$12$wYJZcYXYp.aPvRH6Kp2/7OgU1pUdlyEp0FqXa1b.NXE0OCcJKniXK', 'John', 'Doe', '+1234567890', 'patient'),
    ('b2e22c5c-568f-4f6f-c12f-6fc015c08de2', 'jane.smith@example.com', '$2a$12$NjAxZi1zaGFuR/5oRiWTeueBBeOU9H/5V5qZbdLj/T8WJmHT5fFpq', 'Jane', 'Smith', '+1987654321', 'doctor,'),
    ('d829eb61-9d69-4e2f-99a5-2e394650b413', 'alice.wonder@example.com', '$2a$12$5KbydfGekpNfPhjMN6cPJe63U.9s2zTiQUvTSYabRhN0uCU5gDtxG', 'Alice', 'Wonder', '+1123456789', 'admin'),
    ('6f46004b-7b66-46e3-8c73-463596c8227a', 'bob.builder@example.com', '$2a$12$t9/z9VZMx.gXYO8XvODtDuAubB3qAC9DKS1iQFsn5shw3yGiEDkqi', 'Bob', 'Builder', '+1098765432', 'patient'),
    ('0aa47b51-cea0-46eb-9c4f-a14e5f072908', 'carol.danvers@example.com', '$2a$12$gy.TX.CBO0rQHzM6pC9WkuRB/WiFWNXTt2Y2HgTOoHKN5H1jF/wLi', 'Carol', 'Danvers', '+1222333444', 'doctor,');

CREATE TABLE refresh_tokens (
                                id UUID PRIMARY KEY default gen_random_uuid(),
                                user_id UUID NOT NULL ,
                                token text not null,
                                expires_at TIMESTAMP not null,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);