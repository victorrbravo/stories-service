package pg


import (
    "gopkg.in/yaml.v3"
    "stories/database"
    "stories/generated/models"
)


var ValuesTaFlow  []models.TaFlow

var ValuesTaStory  []models.TaStory

var ValuesTaTracking  []models.TaTracking

var ValuesTaUser  []models.TaUser



func MainProcessor(action string, order string, query *database.IQuery) (string, error) {
    var err error
    var value string
    
    if ValuesTaFlow  == nil {
        ValuesTaFlow = make([]models.TaFlow, 0)
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
    
	switch action {
    
    			case "add_flow":
    			    err = InsertTaFlow(order, query)
    			case "delete_flow":
                    err = DeleteTaFlow(order, query)
    			case "update_flow":
                    err = UpdateTaFlow(order, query)
                case "all_flow":
                    ValuesTaFlow, _ = GetAllTaFlow(query)
                case "get_flow_by_PkFlow":
                    ValueTaFlow, _ := GetTaFlowByPkFlow(order, query)
                    ValuesTaFlow = append(ValuesTaFlow, *ValueTaFlow)
                
    
    			case "add_story":
    			    err = InsertTaStory(order, query)
    			case "delete_story":
                    err = DeleteTaStory(order, query)
    			case "update_story":
                    err = UpdateTaStory(order, query)
                case "all_story":
                    ValuesTaStory, _ = GetAllTaStory(query)
                case "get_story_by_PkStory":
                    ValueTaStory, _ := GetTaStoryByPkStory(order, query)
                    ValuesTaStory = append(ValuesTaStory, *ValueTaStory)
                
                case "get_story_by_user_id":
                    ValueTaStory, _ := GetTaStoryByUserID(order, query)
                    ValuesTaStory = append(ValuesTaStory, *ValueTaStory)
                
    
    			case "add_tracking":
    			    err = InsertTaTracking(order, query)
    			case "delete_tracking":
                    err = DeleteTaTracking(order, query)
    			case "update_tracking":
                    err = UpdateTaTracking(order, query)
                case "all_tracking":
                    ValuesTaTracking, _ = GetAllTaTracking(query)
                case "get_tracking_by_PkTracking":
                    ValueTaTracking, _ := GetTaTrackingByPkTracking(order, query)
                    ValuesTaTracking = append(ValuesTaTracking, *ValueTaTracking)
                
    
    			case "add_user":
    			    err = InsertTaUser(order, query)
    			case "delete_user":
                    err = DeleteTaUser(order, query)
    			case "update_user":
                    err = UpdateTaUser(order, query)
                case "all_user":
                    ValuesTaUser, _ = GetAllTaUser(query)
                case "get_user_by_PkUser":
                    ValueTaUser, _ := GetTaUserByPkUser(order, query)
                    ValuesTaUser = append(ValuesTaUser, *ValueTaUser)
                
    
	}
    if err != nil  {
        return value, err
    }
	return value, nil
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



