drop table IF exists "ta_tracking";
drop table IF exists "ta_flow";
drop table IF exists "ta_story";
drop table IF exists "ta_user";
drop table IF exists "ta_signed_in";
drop table IF exists "ta_session";


create table "ta_user"
(
    "pk_user"                     serial         PRIMARY KEY,
    "first_name"             VARCHAR        NOT NULL,
    "last_name"              VARCHAR        NOT NULL,
    "email"                  VARCHAR        NOT NULL UNIQUE,
    "created_at"             TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
    "account_type"           VARCHAR NOT NULL,
    "password"               VARCHAR        NOT NULL,
    "activated"              BOOLEAN        DEFAULT FALSE,
    "status"                VARCHAR DEFAULT ''::varchar
);

create table "ta_signed_in"
(
    "token"   character(64)                                      NOT NULL,
    "started" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "ended"   timestamp with time zone,
    PRIMARY KEY ("token")
);

create table "ta_session"
(
    "token"      character(64),
    "fk_user"    int REFERENCES "ta_user" ("pk_user"),
    "started_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "invalid"    boolean                  DEFAULT false
);


create table "ta_story"
(
    "pk_story" serial primary key,
    "name"     varchar not null,
    "status"   varchar not null

);

create table "ta_flow"
(
    "pk_flow"  serial primary key,
    "name"     varchar not null unique,
    "flow"     jsonb DEFAULT '{}'::jsonb,
    "fk_story" int REFERENCES "ta_story" ("pk_story")
);

create table "ta_tracking"
(
    "pk_tracking" serial primary key,
    "images"      varchar[],
    "videos"      varchar[],
    "fk_story"    int REFERENCES "ta_story" ("pk_story")
);


CREATE TABLE IF NOT EXISTS "ta_session"
(
    "token"      character(64),
    "user_id"    uuid,
    "started_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "invalid"    boolean                  DEFAULT false
);