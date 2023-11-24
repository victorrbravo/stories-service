package models


import (
	"database/sql"
	"stories/database"
	"stories/scanners"
	"stories/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaFlowCreateQuery = "INSERT INTO \"ta_flow\" (\"name\",\"flow\",\"fk_story\") VALUES ($1, $2, $3) RETURNING (\"pk_flow\")"
    AllTaFlowQuery = "SELECT \"pk_flow\", \"name\", \"flow\", \"fk_story\" FROM \"ta_flow\""
    TaFlowByPkFlowQuery = "SELECT \"pk_flow\", \"name\", \"flow\", \"fk_story\" FROM \"ta_flow\" WHERE \"pk_flow\" = $1"
    TaFlowByFkStoryQuery = "SELECT \"pk_flow\", \"name\", \"flow\", \"fk_story\" FROM \"ta_flow\" WHERE \"fk_story\" = $1"
    TaFlowUpdateQuery = "UPDATE \"ta_flow\" SET \"name\" = COALESCE($1, \"name\"), \"flow\" = COALESCE($2, \"flow\"), \"fk_story\" = COALESCE($3, \"fk_story\") WHERE \"pk_flow\" = $4"
    TaFlowDeleteQuery = "DELETE FROM \"ta_flow\" WHERE \"pk_flow\" = $1"
)

type TaFlow struct {
	PkFlow types.Int `json:"pk_flow"`
	Name types.String `json:"name"`
	Flow []byte `json:"flow"`
	FkStory types.String `json:"fk_story"`
}

func (taFlow *TaFlow) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taFlow.PkFlow,
		&taFlow.Name,
		&taFlow.Flow,
		&taFlow.FkStory,
	)
	if err != nil {
		return nil, err
	}
	return *taFlow, nil
}

func CreateTaFlow(
    query database.IQuery,
	name types.String,
	flow []byte,
	fkStory types.String,
) (*TaFlow, error) {
	var err error
    var scanner scanners.Int
	var taFlow *TaFlow
	err = query.QueryOne(
	    TaFlowCreateQuery,
	    &scanner,
		name,
		flow,
		fkStory,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	taFlow = &TaFlow{
	    PkFlow: scanner.Value,
	    Name: name,
	    Flow: flow,
	    FkStory: fkStory,
	}
	return taFlow, nil
}

func (taFlow *TaFlow) Update(
    query database.IQuery,
    name types.String,
    flow []byte,
    fkStory types.String,
) error {
    var err error
	err = query.Exec(
	    TaFlowUpdateQuery,
        name,
        flow,
        fkStory,
        taFlow.PkFlow,
    )
	if err != nil {
    	return err
	}
    taFlow.Name = name
    taFlow.Flow = flow
    taFlow.FkStory = fkStory
    return nil
}

func GetTaFlowByPkFlow(query database.IQuery, value types.Int) (*TaFlow, error) {
	var err error
	var taFlow TaFlow
	err = query.QueryOne(TaFlowByPkFlowQuery, &taFlow, value)
	if err != nil {
		return  nil, err
	}
	return &taFlow, nil
}




func (taFlow TaFlow) GetTaFlowByFkStory(query database.IQuery, value types.String) (interface{}, error) {
	var err error
	err = query.QueryOne(TaFlowByFkStoryQuery, &taFlow, value)
	if err != nil {
		return  nil, err
	}
	return &taFlow, nil
}




func (taFlow TaFlow) GetID() types.Int {
	return taFlow.PkFlow}


func (taFlow TaFlow) GetName() types.String {
	return taFlow.Name}

func (taFlow TaFlow) GetFkStory() types.String {
	return taFlow.FkStory}


func (taFlow *TaFlow) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaFlowDeleteQuery,
        taFlow.PkFlow,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaFlows(query database.IQuery) ([]TaFlow, error) {
	var err error
	var rows []interface{}
	var list []TaFlow
	rows, err = query.QueryMany(AllTaFlowQuery, &TaFlow{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaFlow, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaFlow))
	}
	return list, nil
}
