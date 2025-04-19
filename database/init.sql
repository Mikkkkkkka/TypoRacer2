-- Table: public.users

-- DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(50) COLLATE pg_catalog."default" NOT NULL,
    token character varying(50) COLLATE pg_catalog."default",
    token_expiration timestamp without time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;


-- Table: public.quotes

-- DROP TABLE IF EXISTS public.quotes;

CREATE TABLE IF NOT EXISTS public.quotes
(
    id integer NOT NULL DEFAULT nextval('quotes_id_seq'::regclass),
    text text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT quotes_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.quotes
    OWNER to postgres;


-- Table: public.plays

-- DROP TABLE IF EXISTS public.plays;

CREATE TABLE IF NOT EXISTS public.plays
(
    user_id integer NOT NULL,
    quote_id integer NOT NULL,
    words_per_minute real NOT NULL,
    accuracy real NOT NULL,
    consistency real NOT NULL,
    CONSTRAINT plays_quote_id_fkey FOREIGN KEY (quote_id)
        REFERENCES public.quotes (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT plays_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT precentage_test CHECK (accuracy <= 100::double precision AND consistency <= 100::double precision)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.plays
    OWNER to postgres;