CREATE TABLE IF NOT EXISTS admin_contact (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    soc_med_acs TEXT[] NOT NULL,
    phone_number TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);