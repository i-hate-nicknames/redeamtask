-- running CREATE TYPE multiple times throws an error
-- todo: when migrations are properly implemented call CREATE type directly
CREATE OR replace FUNCTION create_book_status_type()
    RETURNS boolean
    LANGUAGE plpgsql
AS
$$
BEGIN
  IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'book_status') THEN
      CREATE TYPE book_status AS enum ('checked_out', 'checked_in');
  END IF;
  RETURN true;
END;
$$;

SELECT create_book_status_type();

CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    publisher TEXT NOT NULL,
    publish_date TIME NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 3),
    _status book_status NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL default NOW(),
    deleted_at TIMESTAMP(0)
);

CREATE index IF NOT EXISTS books_deleted ON books (deleted_at);
