CREATE TABLE files(
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    filename varchar not null,
    original_path varchar not null,
    file_size bigint not null,
    mime_type varchar not null,
    status varchar not null default 'uploaded',
    created_at timestamp with time zone default now()
);

CREATE INDEX idx_files_user_id ON files(user_id);