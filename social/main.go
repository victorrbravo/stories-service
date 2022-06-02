package social

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"stories/constants"
	"stories/database"
	"stories/scanners"
	"stories/types"
	"stories/users"
	"strings"
)

const (
	SocialAddLike = `
INSERT INTO "{{.EntityName}}" ("{{.Keyfield}}")
VALUES ($1)
ON CONFLICT ("{{.Keyfield}}")
DO UPDATE SET "social" = jsonb_set("{{.EntityName}}"."social", '{likes}',
"{{.EntityName}}"."social" -> 'likes' || $2, false)::jsonb
RETURNING "social" ->> 'likes'
`
	SocialAddComment = `
	INSERT INTO "{{.EntityName}}" ("{{.Keyfield}}", "social")
VALUES ($1, jsonb_set('{"likes": [],"dislikes": [],"shares": [],"tags": [],"comments": []}'::jsonb, '{comments}', '{"likes": [],"dislikes": [],"shares": [],"tags": [],"comments": []}'::jsonb -> 'comments' || $2, false)::jsonb)
ON CONFLICT ("{{.Keyfield}}")
DO UPDATE SET "social" = jsonb_set("{{.EntityName}}"."social", '{comments}',
"{{.EntityName}}"."social" -> 'comments' || $2, false)::jsonb
RETURNING "social" ->> 'comments'
`
	SocialGetSocialLikes = `
SELECT "social" ->> 'likes'
FROM "{{.EntityName}}"
WHERE "{{.Keyfield}}" = $1
`

	SocialGetSocial = `
SELECT "social"::varchar
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
	UserID  types.UUID `json:"user_id"`
	Keyword string     `json:"keyword"`
}

type Like struct {
	UserID types.UUID     `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

type DisLike struct {
	UserID types.UUID     `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

type Share struct {
	UserID types.UUID     `json:"user_id"`
	Date   types.DateTime `json:"date"`
}

func (obj *Social) AddComment(value types.UUID, text string, user users.User, query database.IQuery) (Comment, error) {
	var err error
	var result scanners.String
	var data []byte
	comment := Comment{user.GetEmail(), user.GetFullName(), text, types.Now()}
	data, err = json.Marshal(comment)
	if err != nil {
		return comment, err
	}
	strquery, _ := obj.useTemplate(SocialAddComment)
	err = query.QueryOne(strquery, &result, value, data)
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
	tpl, err = tpl.Parse(strquery)
	err = tpl.Execute(&buffer, obj)
	return buffer.String(), err

}
func (obj Social) ExistsLikeForThisUser(value types.UUID, user users.User, query database.IQuery) error {
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

func (obj Social) AddLike(value types.UUID, user users.User, query database.IQuery) (types.String, error) {
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
	err = query.QueryOne(strquery, &result, value, data)
	return result.Value, err
}

func (obj Social) GetLikes(value types.UUID, query database.IQuery) (types.String, error) {
	var scanner scanners.String
	var err error
	strquery, err := obj.useTemplate(SocialGetSocialLikes)
	err = query.QueryOne(strquery, &scanner, value)
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
	log.Println(constants.Version, "strquery ", strquery)
	err = query.QueryOne(strquery, &scanner, value)
	if err != nil {
		return socialData, err
	}
	str := types.AsString(scanner.Value)
	err = json.Unmarshal([]byte(str), &socialData)
	return socialData, err
}
