CREATE TABLE routes (
    id VARCHAR(256) PRIMARY KEY,
    flight_id VARCHAR(256),
    source_airport_code VARCHAR(256),
    source_country VARCHAR(256),
    destiny_airport_code VARCHAR(256),
    destiny_country VARCHAR(256),
    bag_cost int,
    planned_start_date DATE,
    planned_end_date DATE,
    created_at DATE,
    updated_at DATE,
    deleted_at DATE
  );