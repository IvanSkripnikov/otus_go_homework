CREATE TABLE events (
                        id serial primary key,
                        title text,
                        date_start date not null,
                        date_end date not null,
                        description text,
                        user_id integer,
                        remember_time integer
);
create index user_id_idx on events (user_id);
create index date_start_idx on events (date_start);
