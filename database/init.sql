-- Table: public.users
-- DROP TABLE IF EXISTS public.users;
CREATE TABLE IF NOT EXISTS public.users (
    id serial NOT NULL PRIMARY KEY,
    username varchar(50) NOT NULL,
    password varchar(50) NOT NULL,
    token varchar(50),
    token_expiration timestamp without time zone
);


-- Table: public.quotes
-- DROP TABLE IF EXISTS public.quotes;
CREATE TABLE IF NOT EXISTS public.quotes (
    id serial NOT NULL PRIMARY KEY,
    text text NOT NULL
);


-- Table: public.plays
-- DROP TABLE IF EXISTS public.plays;
CREATE TABLE IF NOT EXISTS public.plays (
    user_id integer NOT NULL,
    quote_id integer NOT NULL,
    words_per_minute real NOT NULL,
    accuracy real NOT NULL,
    consistency real NOT NULL,
    CONSTRAINT plays_quote_id_fkey FOREIGN KEY (quote_id)
        REFERENCES public.quotes (id) MATCH SIMPLE,
    CONSTRAINT plays_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE,
    CONSTRAINT percentage_test CHECK (accuracy <= 100::double precision AND consistency <= 100::double precision)
);


-- Adding initial quotes to the database
INSERT INTO public.quotes ("text")
VALUES ('The quick brown fox jumps over the lazy dog.');

INSERT INTO public.quotes ("text")
VALUES ('Eat some more of these soft french puns and drink a cup of tea.');