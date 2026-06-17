DROP TABLE IF EXISTS payments, projects, proposals, jobs, users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('CLIENT', 'FREELANCER')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    	FOREIGN KEY (client_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    budget NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'CLOSED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE proposals (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,
    	FOREIGN KEY (job_id)
        REFERENCES jobs(id)
        ON DELETE CASCADE,
    freelancer_id INT NOT NULL,
    	FOREIGN KEY (freelancer_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    cover_letter TEXT,
    bid_amount NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'ACCEPTED', 'REJECTED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,
    	FOREIGN KEY (job_id)
        REFERENCES jobs(id),
    proposal_id INT NOT NULL UNIQUE,
    	FOREIGN KEY (proposal_id)
        REFERENCES proposals(id),
    client_id INT NOT NULL,
    	FOREIGN KEY (client_id)
        REFERENCES users(id),
    freelancer_id INT NOT NULL,
    	FOREIGN KEY (freelancer_id)
        REFERENCES users(id),
    status VARCHAR(20) NOT NULL DEFAULT 'IN_PROGRESS'
        CHECK (status IN ('IN_PROGRESS','SUBMITTED','PAID','COMPLETED')),
    submitted_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    	FOREIGN KEY (project_id)
        REFERENCES projects(id),
    amount NUMERIC(12,2) NOT NULL,
    midtrans_order_id VARCHAR(255),
    transaction_id VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'SUCCESS', 'FAILED')),
    paid_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (
    name,
    email,
    password,
    role
)
VALUES (
    'Nunin Client',
    'nunin@mail.com',
    '$2a$10$iYfjQ73kvTKPFrygFFS5WOtHP7HUhNR5SRGMhEfXBjjgV.MzcrI/2', --passwordclient
    'CLIENT'
);

INSERT INTO users (
    name,
    email,
    password,
    role
)
VALUES (
    'Ula Freelancer',
    'ula@mail.com',
    '$2a$10$GZJVN5qHPp.Hqb81cZtOlO1HRoXtjVtEUkHebFVc.EM8iEV92l7rq', --passwordfreelancer
    'FREELANCER'
);

INSERT INTO jobs (
    client_id,
    title,
    description,
    budget,
    status
)
VALUES (
    1,
    'Build REST API with Golang',
    'Need backend developer for API project',
    5000000,
    'OPEN'
);

INSERT INTO proposals (
    job_id,
    freelancer_id,
    cover_letter,
    bid_amount,
    status
)
VALUES (
    1,
    2,
    'Saya berpengalaman menggunakan Golang dan PostgreSQL.',
    4500000,
    'ACCEPTED'
);

INSERT INTO projects (
    job_id,
    proposal_id,
    client_id,
    freelancer_id,
    status
)
VALUES (
    1,
    1,
    1,
    2,
    'IN_PROGRESS'
);
