CREATE TABLE orders(
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    user_id integer,
    number varchar(40),
    status varchar(12),
    uploaded_at TIMESTAMPT NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE UNIQUE INDEX number_1765300569235_index ON public.orders USING btree (number);