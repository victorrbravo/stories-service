package social

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"strings"

	"stories/constants"
	"stories/database"
	"stories/scanners"
	"stories/types"
	"stories/users"
)

const (
	SocialAddTask = `
	INSERT INTO "{{.EntityName}}" ("{{.Keyfield}}", "flow")
VALUES ($1, jsonb_set('{"name": "","variable": "","needs": [],"token": "","tasks": []}'::jsonb, '{tasks}', '{"name": "","variable": "","needs": [],"token": "","tasks": []}'::jsonb -> 'tasks' || $2, false)::jsonb)
ON CONFLICT ("{{.Keyfield}}")
DO UPDATE SET "flow" = jsonb_set("{{.EntityName}}"."flow", '{tasks}',
"{{.EntityName}}"."flow" -> 'tasks' || $2, false)::jsonb
RETURNING "flow" ->> 'tasks'
`

	SocialRemoveTask = `
with tt as (SELECT
jsonb_agg(j) FILTER (WHERE j->>'name'!='$2') AS new_tasks
FROM ta_flow t, jsonb_array_elements(flow->'tasks') j WHERE pk_flow = $1
GROUP BY t.flow) update ta_flow  
set flow=jsonb_set(flow, '{tasks}', tt.new_tasks::jsonb, false)::jsonb from tt;
`

	SocialCountTask = `
select jsonb_array_length((flow->>'tasks')::jsonb) from ta_flow;
`
	SocialAddLike = `
INSERT INTO "{{.EntityName}}" ("{{.Keyfield}}")
VALUES ($1)
ON CONFLICT ("{{.Keyfield}}")
DO UPDATE SET "flow" = jsonb_set("{{.EntityName}}"."flow", '{likes}',
"{{.EntityName}}"."flow" -> 'likes' || $2, false)::jsonb
RETURNING "flow" ->> 'likes'
`
	SocialAddComment = `
	INSERT INTO "{{.EntityName}}" ("{{.Keyfield}}", "flow")
VALUES ($1, jsonb_set('{"likes": [],"dislikes": [],"shares": [],"tags": [],"comments": []}'::jsonb, '{comments}', '{"likes": [],"dislikes": [],"shares": [],"tags": [],"comments": []}'::jsonb -> 'comments' || $2, false)::jsonb)
ON CONFLICT ("{{.Keyfield}}")
DO UPDATE SET "flow" = jsonb_set("{{.EntityName}}"."flow", '{comments}',
"{{.EntityName}}"."flow" -> 'comments' || $2, false)::jsonb
RETURNING "flow" ->> 'comments'
`
	SocialGetSocialLikes = `
SELECT "flow" ->> 'likes'
FROM "{{.EntityName}}"
WHERE "{{.Keyfield}}" = $1
`

	SocialGetSocial = `
SELECT "flow"::varchar
FROM "{{.EntityName}}"
WHERE "{{.Keyfield}}" = $1
`
)

const SocialSchema string = `
{
	"likes": [],
	"dislikes": [],
	"shares": [],
	"tags": [],
	"comments": []	
}
`

type Link struct {
	To    string `json:"to"`
	Query string `json:"queries"`
}

type Task struct {
	Name      string   `json:"name"`
	TokenLink string   `json:"token_link"`
	Pattern   string   `json:"pattern"`
	Links     []Link   `json:"links"`
	Needs     []string `json:"needs"`
}

type Workflow struct {
	Name     string   `json:"name"`
	Variable string   `json:"variable"`
	Token    string   `json:"token"`
	Needs    []string `json:"needs"`
	Tasks    []Task   `json:"tasks"`
}

type Social struct {
	EntityName string `json:"entity_name"`
	Keyfield   string `json:"keyfield"`
}

type SocialData struct {
	Likes    []Like    `json:"likes"`
	Dislikes []DisLike `json:"dislikes"`
	Shares   []Share   `json:"shares"`
	Tags     []Tag     `json:"tags"`
	Comment  []Comment `json:"comments"`
}

type Comment struct {
	UserEmail string         `json:"user_email"`
	FullName  string         `json:"fullname"`
	Text      string         `json:"text"`
	Date      types.DateTime `json:"date"`
}
type Tag struct {
	UserID  types.Int `json:"user_id"`
	Keyword string    `json:"keyword"`
}

type Like struct {
	UserID types.Int      `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

type DisLike struct {
	UserID types.UUID     `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

type Share struct {
	UserID types.Int      `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

func (obj *Social) AddTask(strId string, name string, tokenLink string, pattern string, user users.User, query database.IQuery) (Task, error) {
	var err error
	var result scanners.String
	var data []byte
	task := Task{name, tokenLink, pattern, []Link{}, []string{}}
	data, err = json.Marshal(task)
	if err != nil {
		return task, err
	}
	strQuery, _ := obj.useTemplate(SocialAddTask)

	_, err = query.QueryOne(strQuery, &result, strId, data)
	return task, err
}

func (obj *Social) RemoveTask(strId string, name string, tokenLink string, pattern string, user users.User, query database.IQuery) (Task, error) {
	var err error
	var result scanners.String
	var data []byte
	task := Task{name, tokenLink, pattern, []Link{}, []string{}}
	data, err = json.Marshal(task)
	if err != nil {
		return task, err
	}
	strQuery, _ := obj.useTemplate(SocialAddTask)

	_, err = query.QueryOne(strQuery, &result, strId, data)
	return task, err
}

func (obj *Social) AddComment(value types.String, text string, user users.User, query database.IQuery) (Comment, error) {
	var err error
	var result scanners.String
	var data []byte
	comment := Comment{types.AsString(user.GetEmail()), types.AsString(user.GetFullName()), text, types.Now()}
	data, err = json.Marshal(comment)
	if err != nil {
		return comment, err
	}
	strquery, _ := obj.useTemplate(SocialAddComment)
	_, err = query.QueryOne(strquery, &result, value, data)
	log.Println("comment result", result)
	return comment, err
}

func (obj Social) getComments(uid types.UUID) []string {
	var comments []string
	return comments
}

func (obj Social) useTemplate(strquery string) (string, error) {
	var buffer bytes.Buffer
	var err error
	var tpl *template.Template
	tpl = template.New("sql_query")
	tpl, _ = tpl.Parse(strquery)
	err = tpl.Execute(&buffer, obj)
	return buffer.String(), err
}

func (obj Social) ExistsLikeForThisUser(value types.Int, user users.User, query database.IQuery) error {
	var likes types.String
	var sameUser bool
	var err error
	likes, err = obj.GetLikes(value, query)
	if err != nil {
		return err
	}
	sameUser = strings.Contains(types.AsString(likes), user.GetID().String())
	// check if user has already liked
	if sameUser == true {
		// Yes, he is
		return errors.New(constants.UserHasAlreadyLiked)
	}
	return nil
}

func MarshalLikes(value types.String) (likes []Like, err error) {
	likes = make([]Like, 0)
	str := types.AsString(value)
	err = json.Unmarshal([]byte(str), &likes)
	if err != nil {
		return likes, err
	}
	return likes, nil
}

func (obj Social) AddLike(value types.Int, user users.User, query database.IQuery) (types.String, error) {
	var err error
	var result scanners.String
	var data []byte
	err = obj.ExistsLikeForThisUser(value, user, query)
	if err != nil {
		return types.NewString("[]"), err
	}
	like := Like{user.GetID(), types.Now()}
	data, _ = json.Marshal(like)
	strquery, _ := obj.useTemplate(SocialAddLike)
	_, err = query.QueryOne(strquery, &result, value, data)
	return result.Value, err
}

func (obj Social) GetLikes(value types.Int, query database.IQuery) (types.String, error) {
	var scanner scanners.String
	var err error
	strquery, err := obj.useTemplate(SocialGetSocialLikes)
	_, err = query.QueryOne(strquery, &scanner, value)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.NewString("[]"), nil
		}
		return scanner.Value, err
	}
	return scanner.Value, err
}

func (obj Social) GetSocial(value types.UUID) (SocialData, error) {
	var socialData SocialData
	var scanner scanners.String
	var query database.IQuery
	var err error
	strquery, _ := obj.useTemplate(SocialGetSocial)
	_, err = query.QueryOne(strquery, &scanner, value)
	if err != nil {
		return socialData, err
	}
	str := types.AsString(scanner.Value)
	err = json.Unmarshal([]byte(str), &socialData)
	return socialData, err
}
