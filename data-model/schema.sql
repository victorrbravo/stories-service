    drop table IF exists "ta_tracking" cascade;
    drop table IF exists "ta_flow" cascade;
    drop table IF exists "ta_story" cascade;
    drop table IF exists "ta_user" cascade ;
    drop table IF exists "ta_signed_in" cascade;
    drop table IF exists "ta_session" cascade;
    drop table IF exists "ta_permissions" cascade;
    drop table IF exists "ta_social_story" cascade;
    drop table IF exists "ta_stories_shared" cascade;

    create table "ta_user"
    (
        "pk_user"                     serial         PRIMARY KEY,
        "firstname"             VARCHAR        NOT NULL,
        "lastname"              VARCHAR        NOT NULL,
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




    create table "ta_flow"
    (
        "pk_flow"  serial primary key,
        "name"     varchar not null unique,
        "flow"     jsonb DEFAULT '{}'::jsonb
    );

    create table "ta_story"
    (
        "pk_story" serial primary key,
        "name"     varchar not null,
        "status"   varchar not null,
        "summary"      varchar,
        "user_id" int REFERENCES "ta_user" ("pk_user"),
        "fk_flow" int REFERENCES "ta_flow" ("pk_flow"),
        constraint "story_status_unique" unique (pk_story, status)
    );

    create table "ta_category"
    (
        "pk_category" serial primary key,
        "name"     varchar not null,
        "user_id" int REFERENCES "ta_user" ("pk_user"),
        "fk_main_category"    int
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
        amount_spent money DEFAULT  '0.0'::money,
        distance float DEFAULT 0.0,
        "fk_status" varchar,
        "fk_story"    int REFERENCES "ta_story" ("pk_story"),
        constraint "tracking_status_unique" unique (fk_story, fk_status)
    );

    ALTER TABLE "ta_tracking" ADD column  "float_value" float default 0.0;


    CREATE TABLE IF NOT EXISTS "ta_session"
    (
        "token"      character(64),
        "user_id"    uuid,
        "started_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
        "invalid"    boolean                  DEFAULT false
    );

ALTER TABLE ta_session ADD COLUMN "pk_session" SERIAL PRIMARY KEY;
ALTER TABLE ta_user ADD COLUMN username VARCHAR(255) DEFAULT '';
ALTER TABLE ta_user ALTER COLUMN username SET NOT NULL;
ALTER TABLE ta_user
        ADD CONSTRAINT unique_username UNIQUE (username);
alter table ta_story add constraint unique_name UNIQUE(user_id,name);
alter table ta_tracking add column fk_story_tips int REFERENCES ta_story(pk_story);
alter table ta_story add column story_type varchar default 'story';
alter table ta_tracking add column fk_story_sub int REFERENCES ta_story(pk_story);
alter table ta_tracking drop column distance;
alter table ta_tracking add column position  GEOMETRY(Point, 4326);
alter table ta_tracking ADD column label varchar(50);
alter table ta_tracking ADD column meta_order varchar;
alter table ta_story ADD column updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL;
alter table ta_story ADD column content_type VARCHAR(10) DEFAULT 'Senda' NOT NULL;
alter table ta_tracking add column tags varchar[] DEFAULT '{}';

    create table "ta_permissions"
    (
        "pk_permission" serial primary key,
        "name"          varchar not null unique,
        "fk_story"       int REFERENCES "ta_story" ("pk_story"),
        "description"   varchar,
        "permissions"     jsonb DEFAULT '{}'::jsonb
    );

    create table "ta_social_story"
    (
        "pk_social" serial primary key,
        "name"      varchar not null,
        "fk_story" int REFERENCES "ta_story" ("pk_story"),
        "social"   jsonb DEFAULT '{}'::jsonb
    );

    create table "ta_stories_shared"
    (
        "pk_shared" serial primary key,
        "fk_story"  int REFERENCES "ta_story" ("pk_story"),
        "fk_from_user"   int REFERENCES "ta_user" ("pk_user"),
        "fk_to_user"   int REFERENCES "ta_user" ("pk_user"),
        "permissions" jsonb DEFAULT '{}'::jsonb,
        "shared_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
    );


alter table ta_user add column name varchar;

-- Data for table `ta_user`
insert into ta_story (name, status) values ('story1', 'Start');
insert into ta_story (name, status) values ('story1', 'ToDo');
insert into ta_story (name, status) values ('story1', 'Progress');
insert into ta_story (name, status) values ('story1', 'ToDo');
insert into ta_story (name, status) values ('story1', 'Done');

    insert into ta_flow (pk_flow,name, flow) values (1,'flow1','
    {
      "id": "3f12f5a5-57ac-4861-8771-f85de203f862",
      "name": "Stories",
      "variable": {
        "name": "basic",
        "queries": "pk_story,name,status"
      },
      "tracking_table_name": "ta_tracking",
      "tracking_keys": "fk_story,fk_status",
      "tracking_id": "pk_tracking",
      "field_token": "pk_story",
      "field_status": "status",
      "table_name": "ta_story",
      "initial": "fk_status=''Start''",
      "tasks": {
        "Start": {
          "info_field": "images[1]",
          "time_field": "updated_at", "amount_field": "amount_spent",
          "transition": {
            "pattern": "none",
            "links": [
              {
                "task": "END",
                "queries": "END"
              }
            ]
          }
        }
      }
    }
    ');



insert into ta_flow (pk_flow,name, flow) values (2,'flow2','
{
  "id": "3f12f5a5-57ac-4861-8771-f85de203f863",
  "name": "Workflow1",
  "variable": {
    "name": "basic",
    "queries": "pk_story,name,status"
  },
  "tracking_table_name": "ta_tracking",
  "tracking_keys": "fk_story,fk_status",
  "tracking_id": "pk_tracking",
  "field_token": "pk_story",
  "field_status": "status",
  "table_name": "ta_story",
  "initial": "fk_status=''ToDo''",
  "tasks": {
    "ToDo": {
      "info_field": "images[1]",
      "time_field": "updated_at","amount_field": "amount_spent",
      "transition": {
        "pattern": "none",
        "links": [
          {
            "task": "Progress",
            "queries": "fk_status=''Progress''"
          }
        ]
      }
    },
    "Progress": {
      "info_field": "images[1]",
      "time_field": "updated_at","amount_spent": "amount_spent",
      "transition": {
        "pattern": "or",
        "links": [
          {
            "task": "Done",
            "queries": "fk_status=''Done''"
          }
        ]
      }
    },
    "Done": {
      "info_field": "images[1]",
      "time_field": "updated_at","amount_field": "amount_spent",
      "transition": {
        "pattern": "or",
        "links": [
          {
            "task": "END",
            "queries": "END"
          }
        ]
      }
    }
  }
}
');
insert into ta_user (firstname, lastname, email, account_type, password, activated)  values ('vbravo','vbravo', 'victorrbravo@gmail.com', 'admin', 'vbravo', 'true');

insert into ta_tracking (images, videos, roles, fk_status, fk_story) values ('{"image1", "image2"}', '{"video1"}', '{"role1"}', 'Done', 1);
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'Progress', 1, '2023-10-30 12:05:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'ToDo', 1, '2023-10-30 10:05:12.119 -0500');

insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image1", "image2"}', '{"video1"}', '{"role1"}', 'Done', 2, '2023-10-31 18:25:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'Progress', 2, '2023-10-31 17:18:12.119 -0500');
insert into ta_tracking (images, videos, roles, fk_status, fk_story, updated_at) values ('{"image3", "image6"}', '{"video3"}', '{"role3"}', 'ToDo', 2, '2023-10-31 16:08:12.119 -0500');
