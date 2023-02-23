SELECT pk_story, name, status, images, videos, roles from ta_story ts join ta_tracking tt  on
    tt.fk_story = ts.pk_story  where name = $1

