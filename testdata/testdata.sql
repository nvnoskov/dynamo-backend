INSERT INTO flight (
    id,
    name,
    number,
    departure,
    departure_time,
    destination,
    arrival_time,
    fare,
    duration,
    created_at,
    updated_at
)
VALUES (
    '967d5bb5-3a7a-4d5e-8a6c-febc8c5b3f13', 
    'BOEING 737-400 ', 
    'UR-CSV', 
    'MALMÖ, SWEDEN', 
    '2019-10-01 15:36:38'::timestamp, 
    'MERZIFON, TURKEY', 
    '2019-10-01 15:36:38'::timestamp, 
    '100EUR', 
    '3h20m', 
    '2019-10-02 11:16:12'::timestamp, 
    '2019-10-02 11:16:12'::timestamp
);

INSERT INTO "user" (
    id,
    name,
    password,
    email
) VALUES (
    'd67d5bb5-3a7a-4d5e-8a6c-febc8c5b3f13', 
    'nvnoskov',
    '$2a$10$eDUmXWENcjQGnsPy87xfw.QjSkltZUr4nvIxOUWJutEdkNvmMikQS',
    'me@noskov.dev'
)