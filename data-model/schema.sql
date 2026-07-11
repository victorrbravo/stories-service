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


    create table "ta_student"
    (
        "pk_student"                     serial         PRIMARY KEY,
        "firstname"             VARCHAR        NOT NULL,
        "lastname"              VARCHAR        NOT NULL,
        "email"                  VARCHAR        NOT NULL UNIQUE,
        "created_at"             TIMESTAMP               DEFAULT CURRENT_TIMESTAMP,
        "account_type"           VARCHAR NOT NULL,
        "password"               VARCHAR        NOT NULL,
        "activated"              BOOLEAN        DEFAULT FALSE,
        "student_type"           VARCHAR        NOT NULL,
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

    alter table ta_flow add column "fk_user" integer;
ALTER TABLE ta_flow
        ADD CONSTRAINT fk_user
            FOREIGN KEY (fk_user)
                REFERENCES ta_user(pk_user);

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

    create table "ta_permission"
    (
        "pk_permission" serial primary key,
        "fk_from_user"   int REFERENCES "ta_user" ("pk_user"),
        "category"   varchar not null,
        "permission" jsonb DEFAULT '{}'::jsonb,
        "updated_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

    create table ta_user_relationship (
                                        pk_user_relationship SERIAL PRIMARY KEY,
                                        fk_follower_id INTEGER NOT NULL REFERENCES ta_user(pk_user),
                                        fk_followed_id INTEGER NOT NULL REFERENCES ta_user(pk_user),
                                        status varchar(30) default 'Requested',
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT unique_relationship UNIQUE (fk_follower_id, fk_followed_id),
                                        CONSTRAINT no_self_follow CHECK (fk_follower_id <> fk_followed_id)
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

-- AUTHORIZATION SCHEMA (Tables)
-- Role table
    CREATE TABLE ta_role (
                             "pk_role" SERIAL PRIMARY KEY,
                             "name" TEXT UNIQUE NOT NULL
    );

-- Policy table
    CREATE TABLE ta_policy (
                               "pk_policy" SERIAL PRIMARY KEY,
                               "resource" TEXT NOT NULL,      -- e.g. "document"
                               "name" TEXT NOT NULL,        -- e.g. "edit"
                               "condition" JSONB NOT NULL     -- e.g. {"field": "department", "equals": "Finance"}
    );
    ALTER TABLE ta_policy
        ADD CONSTRAINT unique_name_condition UNIQUE (name, condition);

-- User ↔ Role (many-to-many)
    CREATE TABLE ta_user_role (
                                  "fk_user" INT REFERENCES ta_user("pk_user") ON DELETE CASCADE,
                                  "fk_role" INT REFERENCES ta_role("pk_role") ON DELETE CASCADE,
                                  PRIMARY KEY ("fk_user", "fk_role")
    );

-- Role ↔ Policy (many-to-many)
    CREATE TABLE ta_role_policy (
                                    "fk_role" INT REFERENCES ta_role("pk_role") ON DELETE CASCADE,
                                    "fk_policy" INT REFERENCES ta_policy("pk_policy") ON DELETE CASCADE,
                                    PRIMARY KEY ("fk_role", "fk_policy")
    );

-- (Optional) User ↔ Policy (direct assignment)
    CREATE TABLE ta_user_policy (
                                    "pk_user_policy" serial   PRIMARY KEY,
                                    "fk_user" INT REFERENCES ta_user("pk_user"),
                                    "fk_policy" INT REFERENCES ta_policy("pk_policy"),
                                    CONSTRAINT uq_user_policy UNIQUE ("fk_user", "fk_policy")
    );

    -- Folders: a folder can contain many stories, a story belongs to one folder
    create table "ta_folder"
    (
        "pk_folder"  serial primary key,
        "name"       varchar not null,
        "user_id"    int REFERENCES "ta_user" ("pk_user"),
        "created_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
        "updated_at" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

    alter table ta_story add column fk_folder int REFERENCES ta_folder(pk_folder);

    -- Insert roles
    INSERT INTO ta_role ("name") VALUES
                                     ('admin'),
                                     ('user');


-- Insert policies
    INSERT INTO ta_policy ("resource", "action", "condition") VALUES
-- Admin policies
('user', 'create', '{"field": "role", "operator": "=", "value": "admin"}'),
('user', 'delete', '{"field": "role", "operator": "=", "value": "admin"}'),
('report', 'view', '{"field": "role", "operator": "=", "value": "admin"}'),

-- User policy
('report', 'view', '{"field": "role", "operator": "=", "value": "user"}');

    -- Link roles to policies
-- (assumes pk_role and pk_policy are auto-generated sequentially)
-- Let's fetch them explicitly for clarity:

-- Admin role → 3 policies
    INSERT INTO ta_role_policy ("fk_role", "fk_policy")
    SELECT r."pk_role", p."pk_policy"
    FROM ta_role r, ta_policy p
    WHERE r."name" = 'admin' AND p."action" IN ('create', 'delete', 'view') AND p."resource" IN ('user','report');

-- User role → 1 policy
    INSERT INTO ta_role_policy ("fk_role", "fk_policy")
    SELECT r."pk_role", p."pk_policy"
    FROM ta_role r, ta_policy p
    WHERE r."name" = 'user' AND p."action" = 'view' AND p."resource" = 'report';

    INSERT INTO __ta_policy__
    (resource, "action", condition)
VALUES (
    'flows',
    'basic',
    '{"action": "add_flow_step", "permises": {"execute": "is_owner"}}'::jsonb
);