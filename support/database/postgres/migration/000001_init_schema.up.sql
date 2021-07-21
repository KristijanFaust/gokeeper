CREATE TABLE "user"
(
    "id"       bigserial PRIMARY KEY,
    "email"    varchar(254) NOT NULL UNIQUE, -- RFC standard prohibits e-mails longer than 254 characters
    "username" varchar(32) NOT NULL,
    "password" bytea NOT NULL
);

CREATE TABLE "password"
(
    "id"       bigserial PRIMARY KEY,
    "user_id"  bigint  NOT NULL,
    "name"     varchar(64) NOT NULL,
    "password" varchar NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY("user_id")
            REFERENCES "user"("id")
);
