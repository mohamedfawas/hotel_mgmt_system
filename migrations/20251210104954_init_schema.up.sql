CREATE TABLE businesses (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE hotels (
    id UUID PRIMARY KEY,
    hotel_code BIGSERIAL UNIQUE NOT NULL,  
    business_id UUID NOT NULL REFERENCES businesses(id),
    name TEXT NOT NULL,
    address TEXT NOT NULL
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY,
    business_id UUID NOT NULL REFERENCES businesses(id),
    hotel_id UUID NOT NULL REFERENCES hotels(id),
    room_number INTEGER UNIQUE NOT NULL,
    room_type TEXT NOT NULL
);