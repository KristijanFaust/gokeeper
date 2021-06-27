CREATE TABLE "user"
(
    "id"       bigserial PRIMARY KEY,
    "email"    varchar NOT NULL UNIQUE,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL
);

CREATE TABLE "password"
(
    "id"       bigserial PRIMARY KEY,
    "user_id"  bigint  NOT NULL,
    "name"     varchar NOT NULL,
    "password" varchar NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY("user_id")
            REFERENCES "user"("id")
);
