package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/victorrbravo/stories-service/constants"
)

// Organizaation Extended Filter
type DataSource string

const (
	DataSourceFinvar                 DataSource = "Finvar"
	DataSourceSEC                    DataSource = "SEC"
	DataSourceCIAWorldFactbook       DataSource = "CIA World Factbook"
	DataSourceQuandlSharadar         DataSource = "Quandl - Sharadar"
	DataSourceWorldBank              DataSource = "World Bank"
	DataSourceIMF                    DataSource = "IMF"
	DataSourceBBC                    DataSource = "BBC"
	DataSourceWorldTradeOrganization DataSource = "World Trade Organization"
	DataSourceCNBC                   DataSource = "CNBC"
	DataSourceFinancialTimes         DataSource = "Financial Times"
	DataSourceWallStreetJounral      DataSource = "Wall Street Jounral"
	DataSourceCNN                    DataSource = "CNN"
	DataSourceFortune                DataSource = "Fortune"
	DataSourceWallStreetJournal      DataSource = "Wall Street Journal"
)

func (mt DataSource) IsValid() error {
	switch mt {
	case DataSourceFinvar, DataSourceSEC, DataSourceCIAWorldFactbook, DataSourceQuandlSharadar, DataSourceWorldBank, DataSourceIMF, DataSourceBBC, DataSourceWorldTradeOrganization,
		DataSourceCNBC, DataSourceFinancialTimes, DataSourceWallStreetJounral, DataSourceCNN, DataSourceFortune, DataSourceWallStreetJournal:
		return nil
	}
	return fmt.Errorf(constants.FieldMustBeType, "DataSource", "Finvar, SEC, CIAWorldFactbook, QuandlSharadar, WorldBank, IMF, BBC, WorldTradeOrganization,CNBC, FinancialTimes, WallStreetJounral, CNN, Fortune, WallStreetJournal")
}

func (lt DataSource) isNil() bool {
	if len(string(lt)) == 0 {
		return true
	}
	return false
}

func (lt DataSource) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}
