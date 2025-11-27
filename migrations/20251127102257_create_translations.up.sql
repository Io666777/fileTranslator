CREATE TABLE translations(
    id bigserial not null primary key,
    file_id bigint not null references files(id) on delete cascade,
    source_lang varchar not null,
    target_lang varchar not null,
    status varchar not null default 'pending',
    translated_path varchar,
    error text,
    created_at timestamp with time zone default now(),
    completed_at timestamp with time zone
);

CREATE INDEX idx_translations_file_id ON translations(file_id);