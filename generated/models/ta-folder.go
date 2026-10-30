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
    TaFolderCreateQuery = "INSERT INTO \"ta_folder\" (\"name\",\"user_id\",\"created_at\",\"updated_at\") VALUES ($1, $2, $3, $4) RETURNING (\"pk_folder\")"
    AllTaFolderQuery = "SELECT \"pk_folder\", \"name\", \"user_id\", \"created_at\", \"updated_at\" FROM \"ta_folder\" ORDER BY \"updated_at asc\" "
    TaFolderByPkFolderQuery = "SELECT \"pk_folder\", \"name\", \"user_id\", \"created_at\", \"updated_at\" FROM \"ta_folder\" WHERE \"pk_folder\" = $1"
    TaFolderByNameQuery = "SELECT \"pk_folder\", \"name\", \"user_id\", \"created_at\", \"updated_at\" FROM \"ta_folder\" WHERE \"name\" = $1"
    TaFolderByUserIDQuery = "SELECT \"pk_folder\", \"name\", \"user_id\", \"created_at\", \"updated_at\" FROM \"ta_folder\" WHERE \"user_id\" = $1"
    TaFolderUpdateQuery = "UPDATE \"ta_folder\" SET \"name\" = COALESCE($1, \"name\"), \"user_id\" = COALESCE($2, \"user_id\"), \"created_at\" = COALESCE($3, \"created_at\"), \"updated_at\" = COALESCE($4, \"updated_at\") WHERE \"pk_folder\" = $5"
    TaFolderDeleteQuery = "DELETE FROM \"ta_folder\" WHERE \"pk_folder\" = $1"
    TaFolderByUserIDAndNameQuery = "SELECT \"pk_folder\", \"name\", \"user_id\", \"created_at\", \"updated_at\" FROM \"ta_folder\" WHERE \"user_id\" = $1 AND \"name\" = $2 ORDER BY updated_at ASC"
    TaFolderListOfFields = " ta_folder.\"pk_folder\", ta_folder.\"name\", ta_folder.\"user_id\", ta_folder.\"created_at\", ta_folder.\"updated_at\""
)

type TaFolder struct {
	PkFolder types.Int `json:"pk_folder"`
	Name types.String `json:"name"`
	UserID types.Int `json:"user_id"`
	CreatedAt types.DateTime `json:"created_at"`
	UpdatedAt types.DateTime `json:"updated_at"`
}

func (taFolder TaFolder) ListOfFields() string {
	return TaFolderListOfFields
}


func (taFolder TaFolder) ScanRow(rows database.RowInterface) (interface{}, error) {
	var err error
	err = rows.Scan(
		&taFolder.PkFolder,
		&taFolder.Name,
		&taFolder.UserID,
		&taFolder.CreatedAt,
		&taFolder.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return taFolder, nil
}

func CreateTaFolder(
    query database.IQuery,
	name types.String,
	userID types.Int,
	createdAt types.DateTime,
	updatedAt types.DateTime,
) (*TaFolder, error) {
	var err error
    var scanner scanners.Int
	var taFolder *TaFolder
	_, err = query.QueryOne(
	    TaFolderCreateQuery,
	    &scanner,
		name,
		userID,
		createdAt,
		updatedAt,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	taFolder = &TaFolder{
	    PkFolder: scanner.Value,
	    Name: name,
	    UserID: userID,
	    CreatedAt: createdAt,
	    UpdatedAt: updatedAt,
	}
	return taFolder, nil
}

func (taFolder *TaFolder) UpdateByOne(
    query database.IQuery,
    options ...func() (string, any),
) error {
var err error
    name := types.NilString()
    userID := types.NilInt
    createdAt := types.NilDateTime
    updatedAt := types.NilDateTime
    for _, option := range options {
    		actionName, value := option()
    		switch actionName {
            case "name":
            name = value.(types.String)
            case "userID":
            userID = value.(types.Int)
            case "createdAt":
            createdAt = value.(types.DateTime)
            case "updatedAt":
            updatedAt = value.(types.DateTime)
		default:
			return errors.New("unsupported field")
		}
	}
	err = query.Exec(
	    TaFolderUpdateQuery,
        name,
        userID,
        createdAt,
        updatedAt,
        taFolder.PkFolder,
    )
	if err != nil {
    	return err
	}
    return nil
}


func (taFolder *TaFolder) Update(
    query database.IQuery,
    name types.String,
    userID types.Int,
    createdAt types.DateTime,
    updatedAt types.DateTime,
) error {
    var err error
	err = query.Exec(
	    TaFolderUpdateQuery,
        name,
        userID,
        createdAt,
        updatedAt,
        taFolder.PkFolder,
    )
	if err != nil {
    	return err
	}
    taFolder.Name = name
    taFolder.UserID = userID
    taFolder.CreatedAt = createdAt
    taFolder.UpdatedAt = updatedAt
    return nil
}

func GetTaFolderByPkFolder(query database.IQuery, value types.Int) (*TaFolder, error) {
	var err error
	var taFolder TaFolder
	var valueTaFolder interface{}
	valueTaFolder, err = query.QueryOne(TaFolderByPkFolderQuery, &taFolder, value)
	if err != nil {
		return  nil, err
	}
    taFolder = valueTaFolder.(TaFolder)
	return &taFolder, nil
}



func GetTaFolderByUserIDAndName(query database.IQuery,
userID types.Int,
name types.String,
) ([]TaFolder, error) {
	    var err error
        var rows []interface{}
        var list []TaFolder
	rows, err = query.QueryMany(TaFolderByUserIDAndNameQuery, &TaFolder{},
    userID,
    name,
)
	if err != nil {
        if err == sql.ErrNoRows {
        	return make([]TaFolder, 0), nil
        }
        return nil, err
       }
       for _, row := range rows {
        	list = append(list, row.(TaFolder))
        }
       return list, nil
}

func GetTaFoldersByUserIDAndNamePaged(query database.IQuery,
userID types.Int,
name types.String,
pageBase types.PageBase) (types.Page, error) {
	var err error
	var page types.Page
	page, err = query.QueryManyPaged(TaFolderByUserIDAndNameQuery,  &TaFolder{},
	pageBase.PageNumber,
	pageBase.PageSize,
        userID,
        name,
	)
	if err != nil {
    		return  page, err
    }
	return page, nil
}


func GetTaFolderByName(query database.IQuery, value types.String) ([]TaFolder, error) {
	    var err error
    	var rows []interface{}
    	var list []TaFolder
    	rows, err = query.QueryMany(TaFolderByNameQuery, &TaFolder{}, value)
    	if err != nil {
    		if err == sql.ErrNoRows {
    			return make([]TaFolder, 0), nil
    		}
    		return nil, err
    	}
    	for _, row := range rows {
    		list = append(list, row.(TaFolder))
    	}
    	return list, nil
}
func GetTaFolderByUserID(query database.IQuery, value types.Int) ([]TaFolder, error) {
	    var err error
    	var rows []interface{}
    	var list []TaFolder
    	rows, err = query.QueryMany(TaFolderByUserIDQuery, &TaFolder{}, value)
    	if err != nil {
    		if err == sql.ErrNoRows {
    			return make([]TaFolder, 0), nil
    		}
    		return nil, err
    	}
    	for _, row := range rows {
    		list = append(list, row.(TaFolder))
    	}
    	return list, nil
}




func (taFolder TaFolder) GetPkFolder() types.Int {
	return taFolder.PkFolder}


func (taFolder TaFolder) GetName() types.String {
	return taFolder.Name}


func (taFolder *TaFolder) Delete(query database.IQuery) error {
	var err error
	err = query.Exec(
	    TaFolderDeleteQuery,
        taFolder.PkFolder,
	)
	if err != nil {
		return err
	}
	return nil
}


func GetAllTaFolders(query database.IQuery) ([]TaFolder, error) {
	var err error
	var rows []interface{}
	var list []TaFolder
	rows, err = query.QueryMany(AllTaFolderQuery, &TaFolder{})
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]TaFolder, 0), nil
		}
		return nil, err
	}
	for _, row := range rows {
		list = append(list, row.(TaFolder))
	}
	return list, nil
}
