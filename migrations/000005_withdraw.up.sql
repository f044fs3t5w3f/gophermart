CREATE TABLE withdraws(
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    user_id integer,
    order_ varchar(40),
    sum BIGINT,
    processed_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE INDEX withdraws_user_id ON public.orders USING btree (user_id);
CREATE UNIQUE INDEX order__1766698640648_index ON "withdraws" USING btree ("order_");