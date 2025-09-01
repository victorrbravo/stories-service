package stories

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"stories/actions/constants"
	"stories/database"
	st_errors "stories/errors"
	"stories/generated/models"
	"stories/generated/queries"
	"stories/pkg/utils"
	"stories/pkg/workflow/common"
	"stories/pkg/workflow/queries/pg"
	"stories/scanners"
	"stories/types"
	"strings"
	"text/template"
	"time"
)

const (
	awsRegion       = "us-east-1"
	bucketName      = "mystoriesbucket"
	durationInSec   = 3600 // Duration for which the URL will be valid (in seconds)
	JsonUrlTemplate = `
	{ "url": "%s" }
	`
	JsonDeleteFlowStepResponseTemplate = `{  "meta_order": "%s" }`
	YamlTaskFieldName                  = "task:\n"
	JsonResponseTemplate               = `{ "response": "%s" }`
	DefaultStartWorkflow               = `
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
      "initial": "fk_status='Start'",
      "tasks": {
        "Start": {
          "info_field": "images[1]",
          "time_field": "updated_at",
		  "amount_field": "amount_spent",
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
`
)

func FlowProcessor(ctx context.Context, action string, order string, query *database.IQuery) (string, error) {
	var err error
	var value string
	switch action {
	case "steps_schema":
		value, err = DoGetStepsSchema(order, query)
	case "new_story_flow":
		_, err = NewStoryFlow(order, query)
	case "new_step_sub": // add sub
		err = NewStepSub(ctx, order, query)
	case "new_step_tips": // add multiple tips
		err = NewStepTips(ctx, order, query)
	case "show_step_tips":
		value, err = DoShowStepTips(order, query)
	case "change_step_field":
		value, err = DoChangeStepField(order, query)
	case "delete_story_flow":
		err = DoDeleteStoryByID(order, query)
	case "add_step":
		_, err = DoAddStep(order, query)
	case constants.AddFlowStepAction:
		err = DoAddStepFlowAt[models.TaTracking](ctx, order, None, query)
	case "add_flow_step_ordered":
		err = DoAddStepFlowAt[models.TaTracking](ctx, order, DateTime, query)
	case "delete_flow_step":
		value, err = DoDeleteFlowStep[models.TaTracking](ctx, order, query)
	case "delete_ticket":
		err = DoDeleteStoryByID(order, query)
	case "delete_step_tips":
		err = DoDeleteTipsOfStepID(order, query)
	case "change_name_of_task":
		err = DoChangeNameTask(order, query)
	case "presign_image":
		value, err = DoPresignImage(order, query)
	case "get_presigned_image":
		value, err = DoGetPresignedImage(order, query)

	case "select_the_category":
		value, err = DoAIQuestion(order, query, "")

	case "suggest_short_name_for_image":
		value, err = DoAIQuestion(order, query, "Summarize the image in two words or less")
	case "suggest_words_for_image":
		value, err = DoAIQuestion(order, query, "Can you say almost 10 words associated with the image separated by commas")
	case "follow_user":
		log.Println("follow_user")
		value, err = DoFollowUser(order, query)
	case "unfollow_user":
		log.Println("unfollow_user")
		value, err = DoUnfollowUser(order, query)
	case "get_users":
		log.Println("get_users")
	case "accept_follower":
		value, err = DoAcceptFollower(ctx, order, query)
	case "get_following_stories":
		log.Println("get_following_stories")
		value, err = DoGetStoriesFollowing(ctx, order, query)
	case "get_stories_by_email_grouped":
		value, err = DoGetStoriesByEmailGrouped(ctx, order, query)
	case "get_stories_by_email":
		value, err = DoGetStoriesByEmail(ctx, order, query)
	case "move_task":
		value, err = DoMoveTaskInFlow[models.TaTracking](order, query)
	case "change_step_report":
		value, err = DoChangeReport[models.TaTracking](order, query)
	case "filter_by_sentence":
		value, err = DoFilterBySentence[models.TaTracking](order, query)
	case "reset_step_report":
		value, err = DoResetStepReport[models.TaTracking](order, query)
	}
	if err != nil {
		return value, err
	}

	return value, nil
}

func DoAcceptFollower(ctx context.Context, order string, query *database.IQuery) (string, error) {
	value := ""
	handler, err := Unmarshall[AcceptFollower]([]byte(order))
	if err != nil {
		return value, err
	}
	userOwner, userFollow, err := GetUsersFromEmails(query, handler.FromUsername, handler.ToUsername)
	if err != nil {
		log.Println("DofolloUser err", err)
		return value, err
	}
	relationsToUpdate, err := models.GetTaUserRelationshipByFkFollowerIdAndFkFollowedId(*query, userOwner.PkUser,
		userFollow.PkUser)
	if err != nil {
		return value, err
	}
	if len(relationsToUpdate) == 0 {
		return value, errors.New("no relation to delete")
	}
	toUpdate := relationsToUpdate[0]
	err = toUpdate.UpdateByOne(*query, func() (string, any) {
		return "status", types.NewString(string(types.FollowTypeAccepted))
	})
	if err != nil {
		log.Println("DofolloUser status err", err)
		return value, err
	}
	log.Println("DoAcceptFollower", handler)

	return value, err
}

func GetUsersFromEmails(query *database.IQuery, fromEmail types.String, toEmail types.String) (models.TaUser, models.TaUser, error) {
	var fromUser, toUser models.TaUser
	userOwners, err := models.GetTaUserByUsername(*query, fromEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows for userOwner")
		}
		return fromUser, toUser, err
	}
	fromUser = userOwners[0]
	userFollows, err := models.GetTaUserByUsername(*query, toEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows for userFollows")
		}
		return fromUser, toUser, err

	}
	toUser = userFollows[0]
	return fromUser, toUser, nil

}

func DoUnfollowUser(order string, query *database.IQuery) (string, error) {
	value := ""
	handler, err := Unmarshall[FollowUser]([]byte(order))
	if err != nil {
		return value, err
	}
	log.Println("DoUnfollowUser", handler)
	userOwner, userFollow, err := GetUsersFromEmails(query, handler.FromUsername, handler.ToUsername)
	if err != nil {
		log.Println("DofolloUser err", err)
		return value, err
	}
	relationToDeletes, err := models.GetTaUserRelationshipByFkFollowerIdAndFkFollowedId(*query, userOwner.PkUser,
		userFollow.PkUser)
	if err != nil {
		return value, err
	}
	if len(relationToDeletes) == 0 {
		return value, errors.New("no relation to delete")
	}
	relatinToDelete := relationToDeletes[0]
	err = relatinToDelete.Delete(*query)
	if err != nil {
		log.Println("error deleting", err)
		return value, err
	}
	log.Println("unfollowing user", userOwner.Email, "follow", "userFollow", userFollow.PkUser, userFollow.PkUser)
	return value, nil
}

func DoFollowUser(order string, query *database.IQuery) (string, error) {
	value := ""
	handler, err := Unmarshall[FollowUser]([]byte(order))
	if err != nil {
		return value, err
	}
	log.Println("handler", handler)
	userOwners, err := models.GetTaUserByUsername(*query, handler.FromUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows for userOwner")
		}
		return value, err
	}
	userOwner := userOwners[0]
	userFollows, err := models.GetTaUserByUsername(*query, handler.ToUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no rows for userFollows")
		}
		return value, err

	}
	userFollow := userFollows[0]
	Requested := "Requested"
	_, err = models.CreateTaUserRelationship(*query, userOwner.PkUser, userFollow.PkUser, types.NewString(Requested), types.Now())
	if err != nil {
		return value, err
	}
	log.Println("user", userOwner.Email, "follow", "userFollow", userFollow.PkUser, userFollow.PkUser)
	return value, nil
}

func DoFilterBySentence[T database.Scanner](order string, query *database.IQuery) (string, error) {
	value := ""

	return value, nil

}
func DoResetStepReport[T database.Scanner](order string, query *database.IQuery) (string, error) {
	var value string
	handler, err := Unmarshall[UpdateAttrInFlow]([]byte(order))
	if err != nil {
		return value, err
	}
	var intFkFlow int
	if !types.IsNil(handler.PkFlow) {
		intFkFlow = types.AsInt(handler.PkFlow)
	} else {
		if types.IsNil(handler.PkStory) {
			return value, errors.New("story id not found")
		}
		taStory, err := models.GetTaStoryByPkStory(*query, handler.PkStory)
		if err != nil {
			return value, err
		}
		intFkFlow = types.AsInt(taStory.FkFlow)
	}
	wf, err := GetWorkflow[T](intFkFlow, query)
	if err != nil {
		return value, err
	}
	log.Println("report number of task:", len(wf.Tasks))
	for name, _ := range wf.Tasks {
		jsonPath := fmt.Sprintf("{tasks, %s,report}", name)
		err = (*query).Exec(UpdateAttributeFlow, jsonPath, types.NewBool(true), intFkFlow)
		if err != nil {
			log.Printf("DB error task:%s sql: %s", name, UpdateAttributeFlow)
		}
	}
	log.Println("all was reset")
	return value, err
}

func DoChangeReport[T database.Scanner](order string, query *database.IQuery) (string, error) {
	var value string
	handler, err := Unmarshall[UpdateAttrInFlow]([]byte(order))
	if err != nil {
		return value, err
	}
	var intFkFlow int
	if !types.IsNil(handler.PkFlow) {
		intFkFlow = types.AsInt(handler.PkFlow)
	} else {
		if types.IsNil(handler.PkStory) {
			return value, errors.New("story id not found")
		}
		taStory, err := models.GetTaStoryByPkStory(*query, handler.PkStory)
		if err != nil {
			return value, err
		}
		intFkFlow = types.AsInt(taStory.FkFlow)
	}
	wf, err := GetWorkflow[T](intFkFlow, query)
	if err != nil {
		return value, err
	}
	if wf.ContainsAll(handler.Name.String()) != nil {
		return value, errors.New(fmt.Sprintf("'%s' step not found", handler.Name))
	}
	jsonPath := fmt.Sprintf("{tasks, %s,report}", handler.Name)
	err = (*query).Exec(UpdateAttributeFlow, jsonPath, handler.Report, intFkFlow)
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoAddExcludeWrapped[T database.Scanner](order string, query *database.IQuery) (string, error) {
	var value string
	handler, err := Unmarshall[UpdateAttrInFlow]([]byte(order))
	if err != nil {
		return value, err
	}
	var intFkFlow int
	if !types.IsNil(handler.PkFlow) {
		intFkFlow = types.AsInt(handler.PkFlow)
	} else {
		if types.IsNil(handler.PkStory) {
			return value, errors.New("story id not found")
		}
		taStory, err := models.GetTaStoryByPkStory(*query, handler.PkStory)
		if err != nil {
			return value, err
		}
		intFkFlow = types.AsInt(taStory.FkFlow)
	}
	jsonPath := fmt.Sprintf("{filter,excluded_wrapped,%s}", handler.Name)
	err = (*query).Exec(UpdateAttributeFlow, jsonPath, handler.Name, intFkFlow)
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoChangeStepField(order string, query *database.IQuery) (string, error) {
	var value string
	handler, err := Unmarshall[ChangeStepField]([]byte(order))
	if err != nil {
		return value, err
	}
	taTracking, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
	if err != nil {
		log.Println("error getting step", err)
		return value, err
	}
	if !types.IsNil(handler.Position) {
		err = taTracking.UpdateByOne(*query, func() (string, any) {
			return "position", handler.Position
		})
		if err != nil {
			log.Println("position error updating position for this step", err)
			return value, err
		}
	}
	if !types.IsNil(handler.Amount) {
		err = taTracking.UpdateByOne(*query, func() (string, any) {
			return "amountSpent", handler.Amount
		})
		if err != nil {
			log.Println("amount error updating step", err)
			return value, err
		}
	}
	if !types.IsNil(handler.Tags) {
		err = taTracking.UpdateByOne(*query, func() (string, any) {
			return "tags", handler.Tags
		})
		if err != nil {
			log.Println("tags error updating step", err)
			return value, err
		}
	}
	if !types.IsNil(handler.UpdatedAt) {
		dt, err := types.NewDateTime(handler.UpdatedAt.String())
		if err != nil {
			return value, err
		}
		err = taTracking.UpdateByOne(*query, func() (string, any) {
			return "updatedAt", dt
		})
		if err != nil {
			log.Println("updated_at error updating step", err)
			return value, err
		}
		value = fmt.Sprintf(`{"update_at": "%s"}`, dt.String())
	}
	return value, nil
}

func DoMoveTaskInFlow[T database.Scanner](order string, query *database.IQuery) (string, error) {
	var value string
	var tupleThree scanners.TupleThree

	handler, err := Unmarshall[MoveTaskInFlow]([]byte(order))
	if err != nil {
		return value, err
	}
	var intFkFlow int
	if !types.IsNil(handler.PkFlow) {
		intFkFlow = types.AsInt(handler.PkFlow)
	} else {
		if types.IsNil(handler.PkStory) {
			return value, errors.New("story id not found")
		}
		taStory, err := models.GetTaStoryByPkStory(*query, handler.PkStory)
		if err != nil {
			return value, err
		}
		intFkFlow = types.AsInt(taStory.FkFlow)
	}
	wf, err := GetWorkflow[T](intFkFlow, query)
	if err != nil {
		return value, err
	}
	target := handler.Target.String()
	if types.IsNil(handler.Target) {
		target = constants.Start
	}
	if wf.Contains(target) != nil {
		return value, errors.New(fmt.Sprintf("'%s' target step not found", target))
	}
	current := handler.Current.String()
	if wf.Contains(current) != nil {
		return value, errors.New(fmt.Sprintf("'%s' current step not found", current))
	}
	previous, err := wf.PreviousStatus(current)
	if err != nil {
		return value, err
	}
	var outputBuffer bytes.Buffer
	tmpl, err := template.New("GetPoints").Parse(GetTaskPoints)
	err = tmpl.Execute(&outputBuffer, struct {
		Target   string
		Current  string
		Previous string
		PkFlow   int
	}{
		Target:   target,
		Current:  current,
		Previous: previous,
		PkFlow:   intFkFlow,
	})
	if err != nil {
		return "", err
	}

	sqlPoints := outputBuffer.String()
	_, err = (*query).QueryOne(sqlPoints, &tupleThree)
	if err != nil {
		return value, err
	}
	jsonPath := fmt.Sprintf("{tasks, %s,transition }", previous)
	err = (*query).Exec(UpdateAttributeFlow, jsonPath, tupleThree.Second, intFkFlow)
	if err != nil {
		return value, err
	}
	jsonPath = fmt.Sprintf("{tasks, %s,transition }", target)
	err = (*query).Exec(UpdateAttributeFlow, jsonPath, tupleThree.Third, intFkFlow)
	if err != nil {
		return value, err
	}
	jsonPath = fmt.Sprintf("{tasks, %s,transition }", current)
	err = (*query).Exec(UpdateAttributeFlow, jsonPath, tupleThree.First, intFkFlow)
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoChangeNameTask(order string, query *database.IQuery) error {
	handler, err := Unmarshall[ChangeNameOfTask]([]byte(order))
	if err != nil {
		return err
	}
	taTracking, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
	if err != nil {
		log.Println("error getting step", err)
		return err
	}
	taStory, err := models.GetTaStoryByPkStory(*query, taTracking.FkStory)
	if err != nil {
		log.Println("error getting story", err)
		return err
	}
	_, err = models.GetTaFlowByPkFlow(*query, taStory.FkFlow)
	if err != nil {
		log.Println("error getting flow", err)
		return err
	}

	return nil
}

func DoDeleteTipsOfStepID(order string, query *database.IQuery) error {
	handler, err := Unmarshall[DeleteTipsByStep]([]byte(order))
	if err != nil {
		return err
	}
	taTracking, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
	if err != nil {
		log.Println("error getting story", err)
		return err
	}
	storyTipsID := taTracking.FkStoryTips

	if types.IsNil(storyTipsID) {
		log.Println("error getting storyTipsID is", storyTipsID)
		return nil
	}
	err = (*query).Exec(constants.SetTipsIDNullInTracking, handler.StepID)
	if err != nil {
		log.Println("error updating story", err)
		return err
	}
	input := fmt.Sprintf(`
action: delete_ticket
ticket: %s
`, storyTipsID)

	err = DoDeleteStoryByID(input, query)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return nil
}

func DoShowStepTips(order string, query *database.IQuery) (string, error) {
	var value string
	handler, err := Unmarshall[ShowStepTips]([]byte(order))
	if err != nil {
		return value, err
	}
	tipsId := handler.TipsID
	if types.IsNil(tipsId) {
		taStep, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return value, errors.New("step not found")
			}
			return value, err
		}
		tipsId = taStep.FkStoryTips
	}
	taStory, err := models.GetTaStoryByPkStory(*query, tipsId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.EmptyStepsArray, nil
		}
		return value, err
	}
	input := fmt.Sprintf(`
action: show_flow_data
ticket: %s
`,
		taStory.PkStory.String())
	log.Println("show_flow_data action", input)
	wfData, err := DoShowWorkflowDataInt[models.TaTracking](input, query)
	if err != nil {
		log.Println("error showing workflow data", err)
		return value, err
	}
	nTasks := len(wfData.Tasks)
	var tips []Tip
	var dataTip DataTip
	tips = make([]Tip, nTasks)
	for index, task := range wfData.Tasks {
		tips[index] = Tip{
			task.Data.ID,
			task.Data.Summary,
		}
	}
	dataTip.Tips = tips
	dataTip.Count = nTasks
	byteValue, err := json.Marshal(dataTip)
	if err != nil {
		return value, err
	}
	value = string(byteValue)
	return value, nil
}

func NewStepTips(ctx context.Context, order string, query *database.IQuery) error {
	handler, err := Unmarshall[StepTips]([]byte(order))
	if err != nil {
		return err
	}
	start := types.NewString(constants.Start)

	taStep, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
	if err != nil {
		return err
	}

	var taStoryTips *models.TaStory
	if !types.IsNil(taStep.FkStoryTips) {
		taStoryTips, err = models.GetTaStoryByPkStory(*query, taStep.FkStoryTips)
		if err != nil {
			log.Println("error getting story", err)
			return err
		}

	} else {
		nameFlow := fmt.Sprintf("%s_story_for_%s", constants.FlowTipsSuffix, taStep.PkTracking.String())
		nameStory := types.NewString(nameFlow)
		taFlow, err := models.CreateTaFlow(*query,
			types.NewString(nameFlow),
			[]byte(DefaultStartWorkflow),
		)
		if err != nil {
			return err
		}
		taStory, err := models.GetTaStoryByPkStory(*query, taStep.FkStory)
		if err != nil {
			log.Println("error getting story", err)
			return err
		}
		taStoryTips, err = models.CreateTaStory(*query,
			nameStory,
			start,
			types.NewString(constants.SummaryDefault),
			taStory.UserID,
			taFlow.PkFlow,
			types.NewString("tips"),
			types.Now(),
			types.Senda,
		)
		if err != nil {
			log.Println("error creating story", err)
			return err
		}
		err = taStep.UpdateByOne(*query,
			func() (string, any) {
				return "fkStoryTips", taStoryTips.PkStory
			})
		if err != nil {
			log.Println("error updating story", err)
			return err
		}
	}

	for index, tip := range handler.Tips.Array {
		status := fmt.Sprintf("tip_%d", index)
		action := fmt.Sprintf(`
action: add_flow_step
ticket: %s
status: %s
images: 
 - ABC897B7-71DA-4C08-B7E8-B7F2F071D445
amount: 0.00
summary: %s
media_date: %s
type: step

`, taStoryTips.PkStory.String(), status, tip, types.Now().Format(time.RFC3339))
		log.Println("add_flow_step action", action)
		err = DoAddStepFlowAt[models.TaTracking](ctx, action, None, query)
		if err != nil {
			if errors.Is(err, st_errors.DuplicateObjectsError) {
				log.Println("duplicate adding step", err)
				continue
			}
			log.Println("error adding step", err)
			return err
		}
	}
	return nil
}

func NewStepSub(ctx context.Context, order string, query *database.IQuery) error {
	log.Println("entering NewStepSub new_step_sub")
	handler, err := Unmarshall[StepSub]([]byte(order))
	if err != nil {
		log.Println("new_step_sub handler err", err)
		return err
	}
	log.Println("new_step_sub handler start")

	start := types.NewString(constants.Start)
	taStep, err := models.GetTaTrackingByPkTracking(*query, handler.StepID)
	if err != nil {
		log.Println("error getting step", err)
		return err
	}
	log.Println("new_step_sub taStep ID", taStep.PkTracking)
	var taStorySub *models.TaStory
	if !types.IsNil(taStep.FkStorySub) {
		taStorySub, err = models.GetTaStoryByPkStory(*query, taStep.FkStorySub)
		if err != nil {
			log.Println("error getting story", err)
			return err
		}

	} else {
		nameFlow := fmt.Sprintf("%s_story_for_%s", constants.SubSuffix, taStep.PkTracking.String())
		nameStory := types.NewString(nameFlow)
		taFlow, err := models.CreateTaFlow(*query,
			types.NewString(nameFlow),
			[]byte(DefaultStartWorkflow),
		)
		if err != nil {
			return err
		}
		taStory, err := models.GetTaStoryByPkStory(*query, taStep.FkStory)
		if err != nil {
			log.Println("taStory error getting story", err)
			return err
		}
		taStorySub, err = models.CreateTaStory(*query,
			nameStory,
			start,
			types.NewString(constants.SummaryDefault),
			taStory.UserID,
			taFlow.PkFlow,
			types.NewString("sub"),
			types.Now(),
			types.Senda,
		)
		if err != nil {
			log.Println("error creating story", err)
			return err
		}
		err = taStep.UpdateByOne(*query,
			func() (string, any) {
				return "fkStorySub", taStorySub.PkStory
			})
		if err != nil {
			log.Println("error updating story", err)
			return err
		}
		stringArr := []string{}
		images := types.NewStringArray(stringArr)
		roles := types.NewStringArray(stringArr)
		_, err = models.CreateTaTracking(*query,
			taStorySub.PkStory,
			start,
			types.NilString(),
			types.NilString(),
			types.NewStringArray(stringArr),
			images,
			roles,
			types.Now(),
			types.NilInt,
			types.NilInt,
			types.NilCoordinates,
			types.NewString(order),
			types.NewStringArray(stringArr),
		)
		if err != nil {
			log.Println("error creating story", err)
			return err
		}
	}
	log.Println("sub-story")
	handler.Task.Ticket = taStorySub.PkStory
	err = DoAddStepFlowAtIntl[models.TaTracking](ctx, handler.Task, None, query)
	if err != nil {
		if errors.Is(err, st_errors.DuplicateObjectsError) {
			log.Println("duplicate adding step", err)
			return err
		}
		log.Println("error adding step", err)
		return err
	}
	return nil
}

func DoGetStoriesFollowing(ctx context.Context, order string, query *database.IQuery) (string, error) {
	var value string
	var err error
	byDateTime := make(map[string][]models.TaStory)

	email := types.NewString(ctx.Value("email").(string))
	taUsers, err := models.GetTaUserByEmail(*query, email)
	if err != nil {
		log.Println("DoGetStoriesByEmail error getting user")
		return value, err
	}
	if taUsers == nil || len(taUsers) == 0 {

		return "", errors.New("user not found")
	}
	taUser := taUsers[0]

	rows, err := (*query).QueryMany(queries.GetStoriesFollowing, &models.TaStory{}, taUser.PkUser)
	if err != nil {
		log.Println("DoGetStoriesByPkStory err getting story")
		return value, err
	}
	for _, row := range rows {
		story, ok := row.(models.TaStory)
		if !ok {
			log.Println("DoGetStoriesByPkStory err getting story row", row)
		}
		onlyDate := story.UpdatedAt.ToOnlyDate()
		byDateTime[onlyDate] = append(byDateTime[onlyDate], story)

	}
	var mapResponse MapResponse
	mapResponse.Values = byDateTime
	data, err := json.Marshal(mapResponse)
	if err != nil {
		log.Println("DoGetStoriesFollowing error marshalling err", err)
		return value, err
	}
	value = string(data)
	return value, err
}

func DoGetStoriesByEmailGrouped(ctx context.Context, order string, query *database.IQuery) (string, error) {
	var value string
	var err error
	byDateTime := make(map[string][]models.TaStory)
	taStories, err := DoGetStoriesByEmailIntl(ctx, query)
	if err != nil {
		log.Println("DoGetStoriesByEmail error unmarshalling")
		return value, err
	}
	for _, story := range taStories {
		onlyDate := story.UpdatedAt.ToOnlyDate()
		byDateTime[onlyDate] = append(byDateTime[onlyDate], story)
	}
	var mapResponse MapResponse
	mapResponse.Values = byDateTime
	data, err := json.Marshal(mapResponse)
	if err != nil {
		log.Println("DoGetStoriesByEmail error marshalling")
		return value, err
	}
	value = string(data)
	return value, err
}

func DoGetStoriesByEmailIntl(ctx context.Context, query *database.IQuery) ([]models.TaStory, error) {
	var taStories []models.TaStory
	email := types.NewString(ctx.Value("email").(string))
	taUsers, err := models.GetTaUserByEmail(*query, email)
	if err != nil {
		log.Println("DoGetStoriesByEmail error getting user")
		return taStories, err
	}
	if taUsers == nil || len(taUsers) == 0 {
		return taStories, errors.New("user not found")
	}
	taUser := taUsers[0]
	storyType := types.NewString("story")
	taStories, err = models.GetTaStoryByStoryTypeAndUserID(*query, storyType, taUser.PkUser)
	if err != nil {
		log.Println("DoGetStoriesByEmail error getting stories")
		return taStories, err
	}
	if taStories == nil {
		taStories = []models.TaStory{}
	}
	return taStories, nil
}

func DoGetStoriesByEmail(ctx context.Context, order string, query *database.IQuery) (string, error) {
	value := "[]"

	taStories, err := DoGetStoriesByEmailIntl(ctx, query)
	if err != nil {
		log.Println("DoGetStoriesByEmail error getting stories")
		return value, err
	}
	data, err := json.Marshal(taStories)
	if err != nil {
		log.Println("DoGetStoriesByEmail error marshalling")
		return value, err
	}
	value = string(data)
	return value, nil
}

func DoAskAboutCategories(tags string, userID types.Int, query *database.IQuery) (string, error) {
	var value, categories string
	log.Println("DoAskAboutCategories *userID", userID)
	taCategory, err := models.GetTaCategoryByUserID(*query, userID)
	if err != nil {
		log.Println("DoAskAboutCategories error getting taCategory error", err)
		return value, err
	}
	categories = ""
	for category := range taCategory {
		name := taCategory[category].Name.String()
		if len(categories) > 0 {
			categories = categories + ","
		}
		categories = categories + name
	}
	tags = strings.ReplaceAll(tags, "\"", "")
	order_category := `
action: select_the_category 
question: Can you summarize these words -- %s --  in one two or more words of these categories -- %s -- . Can you write the response in a category list separated by commas`
	value = fmt.Sprintf(order_category, tags, categories)
	log.Println("DoAskAboutCategories value", value)
	return DoAIQuestion(value, query, "")
}

func DoAIQuestion(order string, query *database.IQuery, question string) (string, error) {
	value := fmt.Sprintf(JsonResponseTemplate, "")
	var err error
	log.Println("DoAIQuestion order ", order)
	handler, err := Unmarshall[SuggestAction]([]byte(order))
	if err != nil {
		log.Println("error unmarshalling")
		return value, err
	}
	if !types.IsNil(handler.Question) {
		question = handler.Question.String()
	}

	apiKey := os.Getenv("OPEN_API_KEY")
	if apiKey == "" {
		return value, errors.New("OPEN_API_KEY is not set")
	}
	client := openai.NewClient(apiKey)

	multiContent := []openai.ChatMessagePart{
		{
			Type: "text",
			Text: question,
		},
	}
	if !handler.ImageID.IsEmpty() {
		imagePath, err := utils.DownloadFromS3(bucketName, handler.ImageID.String())
		if err != nil {
			return value, err
		}
		log.Println("imagePath", imagePath)
		image64, err := utils.EncodeImageBase64(imagePath)
		if err != nil {
			return value, err
		}
		url := fmt.Sprintf("data:image/jpeg;base64,%s", image64)
		imageURL := openai.ChatMessageImageURL{
			URL: url,
		}
		multiContent = append(multiContent, openai.ChatMessagePart{
			Type:     "image_url",
			ImageURL: &imageURL,
		})
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:         openai.ChatMessageRoleUser,
					MultiContent: multiContent,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return value, err
	}
	if len(resp.Choices) == 0 {
		return value, fmt.Errorf("no choices")
	}
	content := resp.Choices[0].Message.Content
	trimContent := strings.Trim(content, " .")
	value = fmt.Sprintf(JsonResponseTemplate, trimContent)
	return value, err
}

func DoGetStepsSchema(order string, query *database.IQuery) (string, error) {
	var value = `
 {
	"actions" : {
		"new_story_flow": {
			"description" : "create a new flow",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"add_step": {
		"description" : "add step",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"add_flow_step": {
			"description" : "add flow step",
			"fields": [
				{ "name": "name", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"delete_ticket": {
			"description" : "delete",
			"fields": [
				{ "name": "ticket", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"presign_image": {
		"description" : "presign",
			"fields": [
				{ "name": "image", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
		},
		"get_presigned_image" : {
			"description" : "get presigned image",
			"fields": [
				{ "name": "image", "type": "string", "validator" : "" },
				{ "name": "summary", "type": "string", "validator" : "" }
			]
	}
 	}
}
`

	return value, nil
}

func DoGetPresignedImage(order string, query *database.IQuery) (string, error) {
	var presignedURL string
	log.Println("DoPresigned Imaaged ")
	handler, err := Unmarshall[PresignAction]([]byte(order))
	if err != nil {
		log.Println("error unmarshalling")
		return presignedURL, err
	}

	log.Printf("DoGetPresignedImage keyName |%s|", handler.Name)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return presignedURL, err
	}
	// Create an S3 service client
	svc := s3.New(sess)

	// Specify the duration for which the URL should be valid (e.g., 1 hour)
	duration := time.Hour * 1

	// Generate the pre-signed URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(handler.Name.String()),
	})
	url, err := req.Presign(duration)
	if err != nil {
		return presignedURL, err
	}

	fmt.Println("url", url)
	value := fmt.Sprintf(JsonUrlTemplate, url)
	return value, nil
}

func DoPresignImage(order string, query *database.IQuery) (string, error) {
	var presignedURL string
	log.Println("DoPresignImage was exectuted")
	handler, err := Unmarshall[PresignAction]([]byte(order))
	if err != nil {
		log.Println("error unmarshalling")
		return presignedURL, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return presignedURL, err

	}
	// Create an S3 client
	s3Client := s3.New(sess)
	log.Println("handler.Name", handler.Name)
	// Generate a pre-signed URL for uploading the file
	presignedURL, err = generatePresignedURL(s3Client, bucketName, handler.Name.String(), durationInSec)
	if err != nil {
		fmt.Println("Error generating pre-signed URL:", err)
		return "", err
	}
	value := fmt.Sprintf(JsonUrlTemplate, presignedURL)
	return value, nil
}

// generatePresignedURL generates a pre-signed URL for uploading a file to S3.
func generatePresignedURL(s3Client *s3.S3, bucket, key string, seconds int) (string, error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// Set the expiration time for the pre-signed URL
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return "", err
	}
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}

	return url, nil
}

func UnmarshallTicket(order []byte) (Ticket, error) {
	var ticket Ticket
	var err error
	err = yaml.Unmarshal(order, &ticket)
	if err != nil {
		return ticket, err
	}

	return ticket, nil
}

func NewStoryFlow(order string, query *database.IQuery) (*models.TaStory, error) {
	var taStory *models.TaStory
	yamlOrder := utils.FormatYaml(order)
	handler, err := UnmarshallTicket([]byte(yamlOrder))
	if err != nil {
		return taStory, err
	}
	start := types.NewString(constants.Start)
	images := types.NewStringArray([]string{})
	roles := types.NewStringArray([]string{})
	nameFlow := ""
	log.Println("handler.Email", handler.Email)
	log.Println("handler.UserID", handler.UserID)
	if types.IsNil(handler.UserID) {
		if types.IsNil(handler.Email) {
			return taStory, fmt.Errorf("email is required")
		}
		taUsers, err := models.GetTaUserByEmail(*query, handler.Email)
		if err != nil {
			return taStory, err
		}
		log.Println("taUsers", taUsers)
		if taUsers == nil || len(taUsers) == 0 {
			log.Println("user not found")
			return taStory, fmt.Errorf("user not found")
		}
		log.Println("len(taUsers)", len(taUsers))
		taUser := taUsers[0]
		handler.UserID = taUser.PkUser
		username := taUser.Username.String()
		prefix := username
		log.Println("prefix", prefix)
		nameFlow = fmt.Sprintf("%s_%s", prefix, handler.Name)
	} else {
		nameFlow = fmt.Sprintf("%s_%s", handler.UserID, handler.Name)
	}
	log.Println("creating ta_flow name_flow", nameFlow)
	taFlow, err := models.CreateTaFlow(*query,
		types.NewString(nameFlow),
		[]byte(DefaultStartWorkflow),
	)
	if err != nil {
		return taStory, err
	}
	taStory, err = models.CreateTaStory(*query,
		handler.Name,
		start,
		handler.Summary,
		handler.UserID,
		taFlow.PkFlow,
		types.NewString("story"),
		types.Now(),
		handler.ContentType,
	)
	log.Println("creating ta_story err", err)
	if err != nil {
		return taStory, err
	}
	stringArr := []string{}
	log.Println("creating ta_tracking")

	models.CreateTaTracking(*query,
		taStory.PkStory,
		start,
		types.NilString(),
		types.NilString(),
		types.NewStringArray(stringArr),
		images,
		roles,
		types.Now(),
		types.NilInt,
		types.NilInt,
		types.NilCoordinates,
		types.NewString(order),
		types.NewStringArray(stringArr),
	)
	return taStory, nil
}

func DoDeleteFlowStep[T database.Scanner](ctx context.Context, order string, query *database.IQuery) (string, error) {
	var value string
	yamlOrder := utils.FormatYaml(order)
	handler, err := Unmarshall[DeleteFlowStepAction]([]byte(yamlOrder))
	if err != nil {
		return value, err
	}

	taTracking, err := models.GetTaTrackingByPkTracking(*query, handler.ID)
	if err != nil {
		return value, err
	}

	taStory, err := models.GetTaStoryByPkStory(*query, taTracking.FkStory)
	if err != nil {
		return value, err
	}
	metaOrder := taTracking.MetaOrder
	value = fmt.Sprintf(JsonDeleteFlowStepResponseTemplate, metaOrder)
	err = taTracking.Delete(*query)
	if err != nil {
		return value, err
	}

	intPkFlow := types.AsInt(taStory.FkFlow)
	var curTask AddTask
	curTask.PkFlow = intPkFlow
	curTask.Name = taTracking.FkStatus
	err = DoDeleteTaskInt[T](curTask, query)
	if err != nil {
		return value, err
	}
	wf, err := GetWorkflow[models.TaTracking](intPkFlow, query)
	taFlow, err := models.GetTaFlowByPkFlow(*query, taStory.FkFlow)
	if err != nil {
		log.Println("error getting flow", err)
		return value, err
	}
	if handler.DeleteType == TipsStepType {
		if wf == nil || len(wf.Tasks) == 1 {
			trackings, err := models.GetTaTrackingByFkStoryTips(*query, taStory.PkStory)
			if err != nil {
				return value, err
			}
			if len(trackings) == 0 {
				return value, st_errors.NoTrackingFoundError
			}
			tipTracking := trackings[0]
			err = (*query).Exec(constants.SetTipsIDNullInTracking, tipTracking.PkTracking)
			if err != nil {
				log.Println("error updating story", err)
				return value, err
			}
			err = taStory.Delete(*query)
			if err != nil {
				log.Println("error deleting story", err)
				return value, err
			}
			err = taFlow.Delete(*query)
			if err != nil {
				log.Println("error deleting flow", err)
				return value, err
			}
		}
	}
	return value, nil
}

func DoAddStepFlowAtIntl[T database.Scanner](ctx context.Context, addStep AddStoryStatus, insert OrderInsert, query *database.IQuery) error {
	var err error
	var addTask AddTask
	fkFLow, err := DoAddStepIntl(addStep, query)
	if err != nil {
		return err
	}
	intFkFlow := types.AsInt(fkFLow)
	wf, err := GetWorkflow[T](intFkFlow, query)
	if err != nil {
		return err
	}
	var prevTask *common.Task
	isLast := true
	if insert == DateTime {
		var getFlow ShowFlowInput
		getFlow.PkFlow = intFkFlow
		getFlow.Ticket = addStep.Ticket
		wfData, err := UpdateWorkflow[T](getFlow, query, []interface{}{}, pg.RoleField[T])
		if err != nil {
			return err
		}
		statuses, err := wfData.Statuses(ctx, false, true)
		if err != nil {
			return err
		}
		for _, status := range statuses {
			task, ok := wfData.Tasks[status]
			if !ok {
				log.Println("error for task", status)
				continue
			}
			var taskTS, currTS time.Time
			var err error
			if taskTS, currTS, err = areDatesCorrect(task.Data.TimeField, addStep.MediaDate); err != nil {
				continue
			}
			if taskTS.After(currTS) {
				if prevTask != nil {
					addTask.Source = types.NewString(prevTask.Name)
				}
				addTask.Target = types.NewString(task.Name)
				isLast = false
				break
			}
			fmt.Println("task name", task.Name, "ts", task.Data.TimeField)
			prevTask = task
		}
	}
	addTask.Name = types.NewString(addStep.Status.String())
	if isLast {
		lastTask, err := wf.Last(ctx)
		if err != nil {
			return err
		}
		addTask.Source = types.NewString(lastTask.Name)
	}
	addTask.PkFlow = intFkFlow
	addTask.Sentence = types.NewString(wf.Sentence(addTask.Name.String()))
	utils.SetDefaults(&addTask)
	wf, err = DoAddTaskIntl[T](intFkFlow, addTask, query)
	if err != nil {
		return err
	}
	return nil
}

func extractResponse(jsonData string) (string, error) {
	// Define an anonymous struct inline
	var res struct {
		Response string `json:"response"`
	}

	err := json.Unmarshal([]byte(jsonData), &res)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`,\s+`)
	value := re.ReplaceAllString(res.Response, `,`)
	return value, nil
}
func DoAddStepFlowAt[T database.Scanner](ctx context.Context, order string, insert OrderInsert, query *database.IQuery) error {
	yamlOrder := utils.FormatYaml(order)
	addStep, err := Unmarshall[AddStoryStatus]([]byte(yamlOrder))
	if err != nil {
		return err
	}
	strOrder := adaptStringBeforeYaml(order)
	tags := addStep.Tags.String()
	addStep.UserID = types.NewInt(2)
	categories, err := DoAskAboutCategories(tags, addStep.UserID, query)
	if err != nil {
		log.Println("get DoAskAboutCategories error", err)
		return err
	}
	newTags := addStep.Tags.String()
	if !types.IsNil(addStep.Tags) {
		newTags = newTags + ","
	}

	categoryStr, err := extractResponse(categories)
	if err != nil {
		log.Println("Error extractResponse error extracting category", err)
	} else {
		newTags = newTags + categoryStr
	}
	newTags = strings.ReplaceAll(newTags, "\"", "")
	log.Println("**DoAddStepFlowAt newTags |", newTags, "|")
	addStep.Tags = types.NewStringArray(strings.Split(newTags, ","))
	addStep.metaOrder = types.NewString(strOrder)
	err = DoAddStepFlowAtIntl[T](ctx, addStep, insert, query)
	if err != nil {
		return err
	}
	return nil
}

func areDatesCorrect(taskDate types.String, currDate types.String) (time.Time, time.Time, error) {
	var taskTS, currTS time.Time
	var err error
	if taskDate.IsEmpty() {
		return taskTS, currTS, errors.New("task date is empty")
	}
	if currDate.IsEmpty() {
		return taskTS, currTS, errors.New("current date is empty")
	}
	taskTS, err = time.Parse(time.RFC3339, taskDate.String())
	if err != nil {
		return taskTS, currTS, errors.New("error parsing task date")
	}
	currTS, err = time.Parse(time.RFC3339, currDate.String())
	if err != nil {
		return taskTS, currTS, errors.New("error parsing current date")
	}
	return taskTS, currTS, nil
}

func DoAddStepFlow(order string, query *database.IQuery) error {
	handler, err := Unmarshall[AddStoryStatus]([]byte(order))
	if err != nil {
		return err
	}

	taStory, err := models.GetTaStoryByPkStory(*query,
		handler.Ticket)
	if err != nil {
		return err
	}

	images := handler.Images
	roles := types.NewStringArray([]string{handler.Amount.String(), handler.StepType.String()})
	status := types.NewString(handler.Status.String())
	taStory.Update(*query,
		types.NilString(),
		status,
		types.NilString(),
		types.NilInt,
		handler.UserID,
		types.NilString(),
		types.Now(),
		types.Senda,
	)

	models.CreateTaTracking(*query,
		taStory.PkStory,
		status,
		types.NilString(),
		types.NilString(),
		types.NewStringArray([]string{}),
		images,
		roles,
		types.Now(),
		types.NilInt,
		types.NilInt,
		types.NilCoordinates,
		types.NewString(order),
		types.NewStringArray([]string{}),
	)

	return nil
}

func DoAddStepIntl(status AddStoryStatus, query *database.IQuery) (types.Int, error) {
	var value types.Int
	taStory, err := models.GetTaStoryByPkStory(*query,
		status.Ticket)
	if err != nil {
		return value, err
	}
	value = taStory.FkFlow
	images := status.Images
	videos := status.Videos
	mediaDate, err := time.Parse(time.RFC3339, status.MediaDate.String())
	if err != nil {
		log.Printf("error parsing date %s", err)
		return value, err
	}
	roles := types.NewStringArray([]string{status.Amount.String(), status.StepType.String()})
	statusStr := types.NewString(status.Status.String())
	err = taStory.Update(*query,
		types.NilString(),
		statusStr,
		types.NilString(),
		types.NilInt,
		types.NilInt,
		types.NilString(),
		types.Now(),
		types.NewNilContentType(),
	)
	if err != nil {
		return value, err
	}
	_, err = models.CreateTaTracking(*query,
		taStory.PkStory,
		statusStr,
		status.Summary,
		status.Amount,
		images,
		videos,
		roles,
		types.FromTime(mediaDate),
		types.NilInt,
		types.NilInt,
		status.Position,
		status.metaOrder,
		status.Tags,
	)
	if err != nil {
		return value, err
	}
	return value, nil
}

func DoAddStep(order string, query *database.IQuery) (AddStoryStatus, error) {
	handler, err := Unmarshall[AddStoryStatus]([]byte(order))
	if err != nil {
		return handler, err
	}
	_, err = DoAddStepIntl(handler, query)
	if err != nil {
		return handler, err
	}
	return handler, nil
}

func DoDeleteStoryByID(order string, query *database.IQuery) error {
	handler, err := Unmarshall[AddStoryStatus]([]byte(order))
	if err != nil {
		return err
	}
	taStory, err := models.GetTaStoryByPkStory(*query, handler.Ticket)
	if err != nil {
		log.Println("error getting story", err)
		return err
	}
	log.Println("taStory FkFlow", taStory.FkFlow)
	taFlow, err := models.GetTaFlowByPkFlow(*query, taStory.FkFlow)
	if err != nil {
		log.Println("error getting flow", err)
		return err
	}

	err = models.DeleteTaTrackingByFkStory(*query, handler.Ticket)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("error deleting tracking", err)
			return err
		}
	}
	err = taStory.Delete(*query)
	if err != nil {
		log.Println("error deleting story", err)
		return err
	}
	err = taFlow.Delete(*query)
	log.Println("error deleting flow", err)
	if err != nil {
		log.Println("error deleting flow", err)
		return err
	}
	return nil
}
