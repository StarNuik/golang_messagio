create table messages (
    msg_id uuid primary key not null default gen_random_uuid(),
    msg_created timestamp not null default now(),
    msg_content text not null
);

