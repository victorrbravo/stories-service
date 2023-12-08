package pg

import (
    "stories/generated/models"
    "stories/database"
	"stories/types"
)

type TaFlow struct {
	PkFlow types.Int `yaml:"pk_flow"`
	Name types.String `yaml:"name"`
	Flow []byte `yaml:"flow"`
}

func  InsertTaFlow(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFlow([]byte(order))
    if err != nil {
        return err
    }
    _, err = models.CreateTaFlow(*query,
            handler.Name,
            handler.Flow,
    )
    if err != nil {
        return err
    }
    return nil
}

func  DeleteTaFlow(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFlow([]byte(order))
    if err != nil {
        return err
    }
    taFlow, err := models.GetTaFlowByPkFlow(*query,
            handler.PkFlow,
    )

    if err != nil {
        return err
    }
    err = taFlow.Delete(*query)
    if err != nil {
        return err
    }
    return nil
}
func  GetTaFlowByPkFlow(order string, query *database.IQuery) (*models.TaFlow, error) {
    handler, err := UnmarshallTaFlow([]byte(order))
    if err != nil {
        return nil, err
    }
    taFlow, err := models.GetTaFlowByPkFlow(*query,
            handler.PkFlow,
    )
    if err != nil {
        return taFlow, err
    }
    return taFlow, nil
}


func  UpdateTaFlow(order string, query *database.IQuery) error {
    handler, err := UnmarshallTaFlow([]byte(order))
    if err != nil {
        return err
    }
    taFlow, err := models.GetTaFlowByPkFlow(*query,
            handler.PkFlow,
    )
    if err != nil {
        return err
    }
    err = taFlow.Update(*query,
            handler.Name,
            handler.Flow,
    )
    if err != nil {
        return err
    }
    return nil
}


func  GetAllTaFlow(query *database.IQuery) ([]models.TaFlow, error) {

    taFlow, err := models.GetAllTaFlows(*query)
    if err != nil {
        return taFlow, err
    }
    return taFlow, nil
}

