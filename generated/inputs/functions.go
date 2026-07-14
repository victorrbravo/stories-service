package pg


import (
    "context"
    "encoding/json"
    "gopkg.in/yaml.v3"
    "stories/database"
    "stories/generated/models"
)


var ValuesTaCategory  []models.TaCategory

var ValuesTaFlow  []models.TaFlow

var ValuesTaFolder  []models.TaFolder

var ValuesTaPolicy  []models.TaPolicy

var ValuesTaSession  []models.TaSession

var ValuesTaStory  []models.TaStory

var ValuesTaTracking  []models.TaTracking

var ValuesTaUser  []models.TaUser

var ValuesTaUserPolicy  []models.TaUserPolicy

var ValuesTaUserRelationship  []models.TaUserRelationship



func MainProcessor(ctx context.Context, action string, order string, query *database.IQuery) (string, error) {
    var err error
    var value string
    var data []byte
    
    if ValuesTaCategory  == nil {
        ValuesTaCategory = make([]models.TaCategory, 0)
    }
    
    if ValuesTaFlow  == nil {
        ValuesTaFlow = make([]models.TaFlow, 0)
    }
    
    if ValuesTaFolder  == nil {
        ValuesTaFolder = make([]models.TaFolder, 0)
    }
    
    if ValuesTaPolicy  == nil {
        ValuesTaPolicy = make([]models.TaPolicy, 0)
    }
    
    if ValuesTaSession  == nil {
        ValuesTaSession = make([]models.TaSession, 0)
    }
    
    if ValuesTaStory  == nil {
        ValuesTaStory = make([]models.TaStory, 0)
    }
    
    if ValuesTaTracking  == nil {
        ValuesTaTracking = make([]models.TaTracking, 0)
    }
    
    if ValuesTaUser  == nil {
        ValuesTaUser = make([]models.TaUser, 0)
    }
    
    if ValuesTaUserPolicy  == nil {
        ValuesTaUserPolicy = make([]models.TaUserPolicy, 0)
    }
    
    if ValuesTaUserRelationship  == nil {
        ValuesTaUserRelationship = make([]models.TaUserRelationship, 0)
    }
    
	switch action {
    
    			case "add_category":
    			    err = InsertTaCategory(order, query)
    			case "delete_category":
                    err = DeleteTaCategory(order, query)
    			case "update_category":
                    err = UpdateTaCategory(order, query)
                case "all_category":
                    ValuesTaCategory, _ = GetAllTaCategory(query)
                    data, err = json.Marshal(ValuesTaCategory)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_category_by_PkCategory":
                    ValueTaCategory, _ := GetTaCategoryByPkCategory(order, query)
                    ValuesTaCategory = append(ValuesTaCategory, *ValueTaCategory)
                
    			case "add_flow":
    			    err = InsertTaFlow(order, query)
    			case "delete_flow":
                    err = DeleteTaFlow(order, query)
    			case "update_flow":
                    err = UpdateTaFlow(order, query)
                case "all_flow":
                    ValuesTaFlow, _ = GetAllTaFlow(query)
                    data, err = json.Marshal(ValuesTaFlow)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_flow_by_PkFlow":
                    ValueTaFlow, _ := GetTaFlowByPkFlow(order, query)
                    ValuesTaFlow = append(ValuesTaFlow, *ValueTaFlow)
                
    			case "add_folder":
    			    err = InsertTaFolder(order, query)
    			case "delete_folder":
                    err = DeleteTaFolder(order, query)
    			case "update_folder":
                    err = UpdateTaFolder(order, query)
                case "all_folder":
                    ValuesTaFolder, _ = GetAllTaFolder(query)
                    data, err = json.Marshal(ValuesTaFolder)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_folder_by_PkFolder":
                    ValueTaFolder, _ := GetTaFolderByPkFolder(order, query)
                    ValuesTaFolder = append(ValuesTaFolder, *ValueTaFolder)
                
                case "get_folder_by_user_id":
                    ValueTaFolders, errSearch := GetTaFolderByUserID(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( ValueTaFolders )
                    value = string(data)
    			case "add_policy":
    			    err = InsertTaPolicy(order, query)
    			case "delete_policy":
                    err = DeleteTaPolicy(order, query)
    			case "update_policy":
                    err = UpdateTaPolicy(order, query)
                case "all_policy":
                    ValuesTaPolicy, _ = GetAllTaPolicy(query)
                    data, err = json.Marshal(ValuesTaPolicy)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_policy_by_PkPolicy":
                    ValueTaPolicy, _ := GetTaPolicyByPkPolicy(order, query)
                    ValuesTaPolicy = append(ValuesTaPolicy, *ValueTaPolicy)
                
    			case "add_session":
    			    err = InsertTaSession(order, query)
    			case "delete_session":
                    err = DeleteTaSession(order, query)
    			case "update_session":
                    err = UpdateTaSession(order, query)
                case "all_session":
                    ValuesTaSession, _ = GetAllTaSession(query)
                    data, err = json.Marshal(ValuesTaSession)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_session_by_PkSession":
                    ValueTaSession, _ := GetTaSessionByPkSession(order, query)
                    ValuesTaSession = append(ValuesTaSession, *ValueTaSession)
                
                case "get_session_by_token":
                    ValueTaSessions, errSearch := GetTaSessionByToken(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( ValueTaSessions )
                    value = string(data)
    			case "add_story":
    			    err = InsertTaStory(order, query)
    			case "delete_story":
                    err = DeleteTaStory(order, query)
    			case "update_story":
                    err = UpdateTaStory(order, query)
                case "all_story":
                    ValuesTaStory, _ = GetAllTaStory(query)
                    data, err = json.Marshal(ValuesTaStory)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_story_by_PkStory":
                    ValueTaStory, _ := GetTaStoryByPkStory(order, query)
                    ValuesTaStory = append(ValuesTaStory, *ValueTaStory)
                
                case "get_story_by_user_id":
                    ValueTaStories, errSearch := GetTaStoryByUserID(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( ValueTaStories )
                    value = string(data)
    			case "add_tracking":
    			    err = InsertTaTracking(order, query)
    			case "delete_tracking":
                    err = DeleteTaTracking(order, query)
    			case "update_tracking":
                    err = UpdateTaTracking(order, query)
                case "all_tracking":
                    ValuesTaTracking, _ = GetAllTaTracking(query)
                    data, err = json.Marshal(ValuesTaTracking)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_tracking_by_PkTracking":
                    ValueTaTracking, _ := GetTaTrackingByPkTracking(order, query)
                    ValuesTaTracking = append(ValuesTaTracking, *ValueTaTracking)
                
                case "get_tracking_by_fk_story_tips":
                    ValueTaTrackings, errSearch := GetTaTrackingByFkStoryTips(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( ValueTaTrackings )
                    value = string(data)
    			case "add_user":
    			    err = InsertTaUser(order, query)
    			case "delete_user":
                    err = DeleteTaUser(order, query)
    			case "update_user":
                    err = UpdateTaUser(order, query)
                case "all_user":
                    ValuesTaUser, _ = GetAllTaUser(query)
                    data, err = json.Marshal(ValuesTaUser)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_user_by_PkUser":
                    ValueTaUser, _ := GetTaUserByPkUser(order, query)
                    ValuesTaUser = append(ValuesTaUser, *ValueTaUser)
                
                case "get_user_by_email":
                    ValueTaUsers, errSearch := GetTaUserByEmail(order, query)
                    if errSearch != nil  {
                            return value, errSearch
                    }
                    data, err = json.Marshal( ValueTaUsers )
                    value = string(data)
    			case "add_user_policy":
    			    err = InsertTaUserPolicy(order, query)
    			case "delete_user_policy":
                    err = DeleteTaUserPolicy(order, query)
    			case "update_user_policy":
                    err = UpdateTaUserPolicy(order, query)
                case "all_user_policy":
                    ValuesTaUserPolicy, _ = GetAllTaUserPolicy(query)
                    data, err = json.Marshal(ValuesTaUserPolicy)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_user_policy_by_PkUserPolicy":
                    ValueTaUserPolicy, _ := GetTaUserPolicyByPkUserPolicy(order, query)
                    ValuesTaUserPolicy = append(ValuesTaUserPolicy, *ValueTaUserPolicy)
                
    			case "add_user_relationship":
    			    err = InsertTaUserRelationship(order, query)
    			case "delete_user_relationship":
                    err = DeleteTaUserRelationship(order, query)
    			case "update_user_relationship":
                    err = UpdateTaUserRelationship(order, query)
                case "all_user_relationship":
                    ValuesTaUserRelationship, _ = GetAllTaUserRelationship(query)
                    data, err = json.Marshal(ValuesTaUserRelationship)
                    if err != nil {
                        return "", err
                    }
                    value = string(data)
                case "get_user_relationship_by_PkUserRelationship":
                    ValueTaUserRelationship, _ := GetTaUserRelationshipByPkUserRelationship(order, query)
                    ValuesTaUserRelationship = append(ValuesTaUserRelationship, *ValueTaUserRelationship)
                
	}
    if err != nil  {
        return value, err
    }
	return value, nil
}


func UnmarshallTaCategory(order []byte) (TaCategory, error) {
	var taCategory TaCategory
	var err error
	err = yaml.Unmarshal(order, &taCategory)
	if err != nil {
		return taCategory, err
	}
	return taCategory, nil
}



func UnmarshallTaFlow(order []byte) (TaFlow, error) {
	var taFlow TaFlow
	var err error
	err = yaml.Unmarshal(order, &taFlow)
	if err != nil {
		return taFlow, err
	}
	return taFlow, nil
}



func UnmarshallTaFolder(order []byte) (TaFolder, error) {
	var taFolder TaFolder
	var err error
	err = yaml.Unmarshal(order, &taFolder)
	if err != nil {
		return taFolder, err
	}
	return taFolder, nil
}



func UnmarshallTaPolicy(order []byte) (TaPolicy, error) {
	var taPolicy TaPolicy
	var err error
	err = yaml.Unmarshal(order, &taPolicy)
	if err != nil {
		return taPolicy, err
	}
	return taPolicy, nil
}



func UnmarshallTaSession(order []byte) (TaSession, error) {
	var taSession TaSession
	var err error
	err = yaml.Unmarshal(order, &taSession)
	if err != nil {
		return taSession, err
	}
	return taSession, nil
}



func UnmarshallTaStory(order []byte) (TaStory, error) {
	var taStory TaStory
	var err error
	err = yaml.Unmarshal(order, &taStory)
	if err != nil {
		return taStory, err
	}
	return taStory, nil
}



func UnmarshallTaTracking(order []byte) (TaTracking, error) {
	var taTracking TaTracking
	var err error
	err = yaml.Unmarshal(order, &taTracking)
	if err != nil {
		return taTracking, err
	}
	return taTracking, nil
}



func UnmarshallTaUser(order []byte) (TaUser, error) {
	var taUser TaUser
	var err error
	err = yaml.Unmarshal(order, &taUser)
	if err != nil {
		return taUser, err
	}
	return taUser, nil
}



func UnmarshallTaUserPolicy(order []byte) (TaUserPolicy, error) {
	var taUserPolicy TaUserPolicy
	var err error
	err = yaml.Unmarshal(order, &taUserPolicy)
	if err != nil {
		return taUserPolicy, err
	}
	return taUserPolicy, nil
}



func UnmarshallTaUserRelationship(order []byte) (TaUserRelationship, error) {
	var taUserRelationship TaUserRelationship
	var err error
	err = yaml.Unmarshal(order, &taUserRelationship)
	if err != nil {
		return taUserRelationship, err
	}
	return taUserRelationship, nil
}



