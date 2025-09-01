SELECT st.pk_story, st.name, st.status, st.summary, st.user_id, st.fk_flow,
       st.story_type, st.updated_at, st.content_type FROM
    ta_story st WHERE st.user_id
        IN (SELECT tr.fk_followed_id from ta_user_relationship tr
    WHERE
    tr.fk_follower_id = $1)

