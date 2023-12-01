    drop table IF exists "ta_tracking" cascade;
    drop table IF exists "ta_flow" cascade;
    drop table IF exists "ta_story" cascade;
    drop table IF exists "ta_user" cascade ;
    drop table IF exists "ta_signed_in" cascade;
    drop table IF exists "ta_session" cascade;


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
        "status"   varchar not null,
        "summary"      varchar,
        constraint "story_status_unique" unique (pk_story, status)
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
        "roles"       varchar[],
        "summary"      varchar,
        "description" varchar,
        updated_at    timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
        "fk_status" varchar,
        "fk_story"    int REFERENCES "ta_story" ("pk_story"),
        constraint "tracking_status_unique" unique (fk_story, fk_status)
    );


    CREATE TABLE IF NOT EXISTS "ta_session"
    (
        "token"      character(64),
        "user_id"    uuid,
        "started_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
        "invalid"    boolean                  DEFAULT false
    );




insert into ta_story (name, status) values ('story1', 'ToDo');
insert into ta_story (name, status) values ('story1', 'Progress');
insert into ta_story (name, status) values ('story1', 'ToDo');
insert into ta_story (name, status) values ('story1', 'Done');

insert into ta_flow (name, fk_story) values ('flow1',1);

insert into ta_tracking (images, videos, roles, fk_status, fk_story) values ('{"image1", "image2"}', '{"video1"}', '{"role1"}', 'Done', 1);
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'Progress', 1, '2023-10-30 12:05:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'ToDo', 1, '2023-10-30 10:05:12.119 -0500');

insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image1", "image2"}', '{"video1"}', '{"role1"}', 'Done', 2, '2023-10-31 18:25:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'Progress', 2, '2023-10-31 17:18:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'ToDo', 2, '2023-10-31 16:08:12.119 -0500');
