-- running CREATE TYPE multiple times throws an error
-- todo: when migrations are properly implemented call create type directly
create or replace function create_book_status_type()
    returns boolean
    language plpgsql
as
$$
begin
  if not exists(select 1 from pg_type where typname = 'book_status') then
      create type book_status as enum ('checked_out', 'checked_in');
  end if;
  return true;
end;
$$;

select create_book_status_type();

create table if not exists books (
    id bigserial primary key,
    title text not null,
    author text not null,
    publisher text not null,
    publish_date time not null,
    rating smallint not null CHECK (rating >= 1 AND rating <= 3),
    _status book_status not null,
    created_at timestamp(0) with time zone not null default now(),
    deleted_at timestamp(0)
);

create index if not exists books_deleted on books (deleted_at);
