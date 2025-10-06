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
    TaPolicyCreateQuery = "INSERT INTO \"ta_policy\" (\"resource\",\"name\",\"condition\") VALUES ($1, $2, $3) RETURNING (\"pk_policy\")"
    AllTaPolicyQuery = "SELECT \"pk_policy\", \"resource\", \"name\", \"condition\" FROM \"ta_policy\""
    TaPolicyByPkPolicyQuery = "SELECT \"pk_policy\", \"resource\", \"name\", \"condition\" FROM \"ta_policy\" WHERE \"pk_policy\" = $1"
    TaPolicyUpdateQuery = "UPDATE \"ta_policy\" SET \"resource\" = COALESCE($1, \"resource\"), \"name\" = COALESCE($2, \"name\"), \"condition\" = COALESCE($3, \"condition\") WHERE \"pk_policy\" = $4"
    TaPolicyDeleteQuery = "DELETE FROM \"ta_policy\" WHERE \"pk_policy\" = $1"
    TaPolicyByNameQuery = "SELECT \"pk_policy\", \"resource\", \"name\", \"condition\" FROM \"ta_policy\" WHERE \"name\" = $1 "
    TaPolicyListOfFields = " ta_policy.\"pk_policy\", ta_policy.\"resource\", ta_policy.\"name\", ta_policy.\"condition\""
)

type TaPolicy struct {
	PkPolicy types.Int `json:"pk_policy"`
	Resource types.String `json:"resource"`
	Name types.String `json:"name"`
	Condition types.Bytes `json:"condition"`
}

func (taPolicy TaPolicy) ListOfFields() string {
	return TaPolicyListOfFields
}


func (taPolicy TaPolicy) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taPolicy.PkPolicy,
		&taPolicy.Resource,
		&taPolicy.Name,
		&taPolicy.Condition,
	)
	if err != nil {
		return nil, err
	}
	return taPolicy, nil
}

func CreateTaPolicy(
    query database.IQuery,
	resource types.String,
	name types.String,
	condition types.Bytes,
) (*TaPolicy, error) {
	var err error
    var scanner scanners.Int
	var taPolicy *TaPolicy
	_, err = query.QueryOne(
	    TaPolicyCreateQuery,
	    &scanner,
		resource,
		name,
		condition,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taPolicy = &TaPolicy{
	    PkPolicy: scanner.Value,
	    Resource: resource,
	    Name: name,
	    Condition: condition,
	}
	return taPolicy, nil
}

func (taPolicy *TaPolicy) UpdateByOne(
    query database.IQuery,
    options ...func() (string, any),
) error {
var err error
    resource := types.NilString()
    name := types.NilString()
    condition := types.NilBytes
    for _, option := range options {
    		actionName, value := option()
    		switch actionName {
            case "resource":
            resource = value.(types.String)
            case "name":
            name = value.(types.String)
            case "condition":
            condition = value.(types.Bytes)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
	    TaPolicyUpdateQuery,
        resource,
        name,
        condition,
        taPolicy.PkPolicy,
    )
	if err != nil {
    	return err
	}
    return nil
}


func (taPolicy *TaPolicy) Update(
    query database.IQuery,
    resource types.String,
    name types.String,
    condition types.Bytes,
) error {
    var err error
	err = query.Exec(
	    TaPolicyUpdateQuery,
        resource,
        name,
        condition,
        taPolicy.PkPolicy,
    )
	if err != nil {
    	return err
	}
    taPolicy.Resource = resource
    taPolicy.Name = name
    taPolicy.Condition = condition
    return nil
}

func GetTaPolicyByPkPolicy(query database.IQuery, value types.Int) (*TaPolicy, error) {
	var err error
	var taPolicy TaPolicy
	var valueTaPolicy interface{}
	valueTaPolicy, err = query.QueryOne(TaPolicyByPkPolicyQuery, &taPolicy, value)
	if err != nil {
		return  nil, err
	}
    taPolicy = valueTaPolicy.(TaPolicy)
	return &taPolicy, nil
}



func GetTaPolicyByName(query database.IQuery,
name types.String,
) ([]TaPolicy, error) {
	    var err error
        var rows []interface{}
        var list []TaPolicy
	rows, err = query.QueryMany(TaPolicyByNameQuery, &TaPolicy{},
    name,
)
	if err != nil {
        if err == sql.ErrNoRows {
        	return make([]TaPolicy, 0), nil
        }
        return nil, err
       }
       for _, row := range rows {
        	list = append(list, row.(TaPolicy))
        }
       return list, nil
}

func GetTaPoliciesByNamePaged(query database.IQuery,
name types.String,
pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaPolicyByNameQuery,  &TaPolicy{},
	pageBase.PageNumber,
	pageBase.PageSize,
        name,
	)
	if err != nil {
    		return  page, err
    }
	return page, nil
}






func (taPolicy TaPolicy) GetID() types.Int {
	return taPolicy.PkPolicy}


func (taPolicy TaPolicy) Getresource() types.String {
	return taPolicy.Resource}


func (taPolicy TaPolicy) Getaction() types.String {
	return taPolicy.Name}


func (taPolicy *TaPolicy) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaPolicyDeleteQuery,
        taPolicy.PkPolicy,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaPolicies(query database.IQuery) ([]TaPolicy, error) {
	var err error
	var rows []interface{}
	var list []TaPolicy
	rows, err = query.QueryMany(AllTaPolicyQuery, &TaPolicy{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaPolicy, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaPolicy))
	}
	return list, nil
}
