package models


import (
	"database/sql"
    "errors"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/scanners"
	"github.com/victorrbravo/stories-service/types"
)
/*
 * Generated queries (these can be used perfectly instead of writing them by hand
 */
const (
    TaCategoryCreateQuery = "INSERT INTO \"ta_category\" (\"name\",\"fk_main_category\") VALUES ($1, $2) RETURNING (\"pk_category\")"
    AllTaCategoryQuery = "SELECT \"pk_category\", \"name\", \"fk_main_category\" FROM \"ta_category\""
    TaCategoryByPkCategoryQuery = "SELECT \"pk_category\", \"name\", \"fk_main_category\" FROM \"ta_category\" WHERE \"pk_category\" = $1"
    TaCategoryUpdateQuery = "UPDATE \"ta_category\" SET \"name\" = COALESCE($1, \"name\"), \"fk_main_category\" = COALESCE($2, \"fk_main_category\") WHERE \"pk_category\" = $3"
    TaCategoryDeleteQuery = "DELETE FROM \"ta_category\" WHERE \"pk_category\" = $1"
    TaCategoryByUserIDQuery = "SELECT \"pk_category\", \"name\", \"fk_main_category\" FROM \"ta_category\" WHERE \"user_id\" = $1 "
    TaCategoryListOfFields = " ta_category.\"pk_category\", ta_category.\"name\", ta_category.\"fk_main_category\""
)

type TaCategory struct {
	PkCategory types.Int `json:"pk_category"`
	Name types.String `json:"name"`
	FkMainCategory types.Int `json:"fk_main_category"`
}

func (taCategory TaCategory) ListOfFields() string {
	return TaCategoryListOfFields
}


func (taCategory TaCategory) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taCategory.PkCategory,
		&taCategory.Name,
		&taCategory.FkMainCategory,
	)
	if err != nil {
		return nil, err
	}
	return taCategory, nil
}

func CreateTaCategory(
    query database.IQuery,
	name types.String,
	fkMainCategory types.Int,
) (*TaCategory, error) {
	var err error
    var scanner scanners.Int
	var taCategory *TaCategory
	_, err = query.QueryOne(
	    TaCategoryCreateQuery,
	    &scanner,
		name,
		fkMainCategory,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taCategory = &TaCategory{
	    PkCategory: scanner.Value,
	    Name: name,
	    FkMainCategory: fkMainCategory,
	}
	return taCategory, nil
}

func (taCategory *TaCategory) UpdateByOne(
    query database.IQuery,
    options ...func() (string, any),
) error {
var err error
    name := types.NilString()
    fkMainCategory := types.NilInt
    for _, option := range options {
    		actionName, value := option()
    		switch actionName {
            case "name":
            name = value.(types.String)
            case "fkMainCategory":
            fkMainCategory = value.(types.Int)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
	    TaCategoryUpdateQuery,
        name,
        fkMainCategory,
        taCategory.PkCategory,
    )
	if err != nil {
    	return err
	}
    return nil
}


func (taCategory *TaCategory) Update(
    query database.IQuery,
    name types.String,
    fkMainCategory types.Int,
) error {
    var err error
	err = query.Exec(
	    TaCategoryUpdateQuery,
        name,
        fkMainCategory,
        taCategory.PkCategory,
    )
	if err != nil {
    	return err
	}
    taCategory.Name = name
    taCategory.FkMainCategory = fkMainCategory
    return nil
}

func GetTaCategoryByPkCategory(query database.IQuery, value types.Int) (*TaCategory, error) {
	var err error
	var taCategory TaCategory
	var valueTaCategory interface{}
	valueTaCategory, err = query.QueryOne(TaCategoryByPkCategoryQuery, &taCategory, value)
	if err != nil {
		return  nil, err
	}
    taCategory = valueTaCategory.(TaCategory)
	return &taCategory, nil
}



func GetTaCategoryByUserID(query database.IQuery,
userID types.Int,
) ([]TaCategory, error) {
	    var err error
        var rows []interface{}
        var list []TaCategory
	rows, err = query.QueryMany(TaCategoryByUserIDQuery, &TaCategory{},
    userID,
)
	if err != nil {
        if err == sql.ErrNoRows {
        	return make([]TaCategory, 0), nil
        }
        return nil, err
       }
       for _, row := range rows {
        	list = append(list, row.(TaCategory))
        }
       return list, nil
}

func GetTaCategoriesByUserIDPaged(query database.IQuery,
userID types.Int,
pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaCategoryByUserIDQuery,  &TaCategory{},
	pageBase.PageNumber,
	pageBase.PageSize,
        userID,
	)
	if err != nil {
    		return  page, err
    }
	return page, nil
}






func (taCategory TaCategory) GetID() types.Int {
	return taCategory.PkCategory}


func (taCategory TaCategory) GetName() types.String {
	return taCategory.Name}


func (taCategory *TaCategory) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaCategoryDeleteQuery,
        taCategory.PkCategory,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaCategories(query database.IQuery) ([]TaCategory, error) {
	var err error
	var rows []interface{}
	var list []TaCategory
	rows, err = query.QueryMany(AllTaCategoryQuery, &TaCategory{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaCategory, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaCategory))
	}
	return list, nil
}
