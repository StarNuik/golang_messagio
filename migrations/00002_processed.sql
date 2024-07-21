create table processed_workloads (
    load_id uuid primary key not null default gen_random_uuid(),
    load_msg_id uuid not null,
    load_created timestamp not null default now(),
    load_hash char(64) not null
);