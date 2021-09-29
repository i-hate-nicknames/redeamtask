create type book_status as enum ('checked_out', 'checked_in');

create table if not exists books (
    id bigserial primary key,
    title text not null,
    author text not null,
    publisher text not null,
    publish_date time not null,
    rating int not null,
    _status book_status not null,
    created_at timestamp(0) with time zone not null default now(),
    deleted_at timestamp(0)
);

create index books_deleted on books (deleted_at);
