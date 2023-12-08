package models


import (
	"database/sql"
    "errors"
	"stories/database"
	"stories/scanners"
	"stories/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaFlowCreateQuery = "INSERT INTO \"ta_flow\" (\"name\",\"flow\") VALUES ($1, $2) RETURNING (\"pk_flow\")"
    AllTaFlowQuery = "SELECT \"pk_flow\", \"name\", \"flow\" FROM \"ta_flow\""
    TaFlowByPkFlowQuery = "SELECT \"pk_flow\", \"name\", \"flow\" FROM \"ta_flow\" WHERE \"pk_flow\" = $1"
    TaFlowUpdateQuery = "UPDATE \"ta_flow\" SET \"name\" = COALESCE($1, \"name\"), \"flow\" = COALESCE($2, \"flow\") WHERE \"pk_flow\" = $3"
    TaFlowDeleteQuery = "DELETE FROM \"ta_flow\" WHERE \"pk_flow\" = $1"
)

type TaFlow struct {
	PkFlow types.Int `json:"pk_flow"`
	Name types.String `json:"name"`
	Flow []byte `json:"flow"`
}

func (taFlow TaFlow) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taFlow.PkFlow,
		&taFlow.Name,
		&taFlow.Flow,
	)
	if err != nil {
		return nil, err
	}
	return taFlow, nil
}

func CreateTaFlow(
    query database.IQuery,
	name types.String,
	flow []byte,
) (*TaFlow, error) {
	var err error
    var scanner scanners.Int
	var taFlow *TaFlow
	_, err = query.QueryOne(
	    TaFlowCreateQuery,
	    &scanner,
		name,
		flow,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taFlow = &TaFlow{
	    PkFlow: scanner.Value,
	    Name: name,
	    Flow: flow,
	}
	return taFlow, nil
}

func (taFlow *TaFlow) Update(
    query database.IQuery,
    name types.String,
    flow []byte,
) error {
    var err error
	err = query.Exec(
	    TaFlowUpdateQuery,
        name,
        flow,
        taFlow.PkFlow,
    )
	if err != nil {
    	return err
	}
    taFlow.Name = name
    taFlow.Flow = flow
    return nil
}

func GetTaFlowByPkFlow(query database.IQuery, value types.Int) (*TaFlow, error) {
	var err error
	var taFlow TaFlow
	var valueTaFlow interface{}
	valueTaFlow, err = query.QueryOne(TaFlowByPkFlowQuery, &taFlow, value)
	if err != nil {
		return  nil, err
	}
    taFlow = valueTaFlow.(TaFlow)
	return &taFlow, nil
}








func (taFlow TaFlow) GetID() types.Int {
	return taFlow.PkFlow}


func (taFlow TaFlow) GetName() types.String {
	return taFlow.Name}


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
