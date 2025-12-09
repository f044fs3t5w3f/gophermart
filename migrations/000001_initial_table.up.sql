CREATE TABLE users(
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    login varchar(50),
    password_hash varchar(255),
    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX login_1762799044404_index ON public.users USING btree (login);