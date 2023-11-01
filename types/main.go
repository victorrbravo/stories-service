package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"stories/constants"
)

// Type of Portfolio
type PortfolioType string

// Portfolio Constants related to DB types
const (
	PortfolioPublic  PortfolioType = "Public"
	PortfolioPrivate               = "Private"
)

// Type of Portfolio
type PortfolioAccessType string

// Portfolio Constants related to DB types
const (
	PortfolioAccessOwner        PortfolioAccessType = "Owner"
	PortfolioAccessReadOnly                         = "Read Only"
	PortfolioAccessWriteAndRead                     = "Write & Read"
)

// Type of Location
type LocationType string

const (
	LocationHeadquarters LocationType = "Headquarters"
	LocationOffice                    = "Office"
)

type ProbabilityType string

const (
	ProbabilityHigh       ProbabilityType = "High"
	ProbabilityMediumHigh                 = "Medium-High"
	ProbabilityMedium                     = "Medium"
	ProbabilityMediumLow                  = "Medium-Low"
	ProbabilityLow                        = "Low"
)

func (lt ProbabilityType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt ProbabilityType) IsValid() error {
	switch lt {
	case ProbabilityHigh, ProbabilityMedium, ProbabilityLow, ProbabilityMediumLow, ProbabilityMediumHigh:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "Probability", "High, Medium,Low, Medium-Low, Medium-High"))
}

// Rating Type string
type RatingType string

const (
	RatingHigh   RatingType = "High"
	RatingMedium            = "Medium"
	RatingLow               = "Low"
)

func (lt RatingType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt RatingType) IsValid() error {
	switch lt {
	case RatingHigh, RatingMedium, RatingLow:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "RatingType", "High,  Medium, Low"))
}

// Impact Type string

type ImpactType string

const (
	ImpactHigh       ImpactType = "High"
	ImpactMediumHigh            = "Medium-High"
	ImpactMedium                = "Medium"
	ImpactMediumLow             = "Medium-Low"
	ImpactLow                   = "Low"
)

func (lt ImpactType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt ImpactType) IsValid() error {
	switch lt {
	case ImpactHigh, ImpactMedium, ImpactLow, ImpactMediumLow, ImpactMediumHigh:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "Impact", "High, Medium,Low, Medium, Medium-High, Low"))
}

//  TimePeriodType

type TimePeriodType string

const (
	ShortTime TimePeriodType = "Short Term"
	LongTIme                 = "Long Term"
)

func (lt TimePeriodType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt TimePeriodType) IsValid() error {
	switch lt {
	case ShortTime, LongTIme:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "TimePeriodType", "Short Term, Long Term"))
}

// TrendDirectionType

type TrendDirectionType string

const (
	UpIncreasing   TrendDirectionType = "Up / Increasing"
	Increasing                        = "Increasing"
	Neutral                           = "Neutral"
	Decreasing                        = "Decreasing"
	DownDecreasing                    = "Down / Decreasing"
)

func (lt TrendDirectionType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt TrendDirectionType) IsValid() error {
	switch lt {
	case UpIncreasing, Increasing, Neutral, Decreasing, DownDecreasing:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "TrendDirectionType", "UpIncreasing, Increasing, Neutral, Decreasing, DownDecreasing"))
}

// ImpactOnCompanyType

type ImpactOnCompanyType string

const (
	Negative ImpactOnCompanyType = "Negative"
	Positive                     = "Positive"
)

func (lt ImpactOnCompanyType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

func (lt ImpactOnCompanyType) IsValid() error {
	switch lt {
	case Negative, Positive:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "ImpactOnCompanyType", "Negative, Positive"))
}

func IsValidPortfolioType(value PortfolioType) bool {
	if PortfolioType(value) == PortfolioPrivate || PortfolioType(value) == PortfolioPublic {
		return true
	}
	return false
}

type CountryType string

const (
	USA    CountryType = "United States"
	France             = "France"
	Mexico             = "Mexico"
)

// Notification Type
type NotificationType string

const (
	Any                            NotificationType = "Any"
	Event                          NotificationType = "Event"
	Portfolio                      NotificationType = "Portfolio"
	Suggestion                     NotificationType = "Suggestion"
	Organization                   NotificationType = "Organization"
	KRF                            NotificationType = "KRF"
	News                           NotificationType = "News"
	Account                        NotificationType = "Account"
	Post                           NotificationType = "Posts"
	PortofolioNewsAndNotifications NotificationType = "Portfolios News & Notifications"
	CompanyPriceAlert              NotificationType = "Company Price Alert"
	Connections                    NotificationType = "Connections"
	FINVARNotifications            NotificationType = "FINVAR Notifications"
)

func (lt NotificationType) IsValid() error {
	switch lt {
	case Any, Event, Portfolio, Organization, KRF, News, Account, PortofolioNewsAndNotifications, CompanyPriceAlert, Connections, FINVARNotifications:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "NotificationType", "Any, Event, Portfolio, Organization, KRF, New, Account, Portfolios News & Notifications, Company Price Alert, Connections, FINVARNotifications"))
}

func (lt NotificationType) isNil() bool {
	if len(string(lt)) == 0 {
		return true
	}
	return false
}

func (lt NotificationType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

// Category Type

type CategoryType string

const (
	FINVARNotification CategoryType = "Notifications"
	NewsAndEvents      CategoryType = "News & Events"
	UserPost           CategoryType = "User Post"
)

func (lt CategoryType) IsValid() error {
	switch lt {
	case FINVARNotification, NewsAndEvents, UserPost:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "category", "Notifications, News & Events or User Post"))
}

func (lt CategoryType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

// Category Type

type NotificationIconType string

const (
	NotifIncreasing NotificationIconType = "Increasing"
	NotifDecreasing NotificationIconType = "Decreasing"
	NotifNeutral    NotificationIconType = "Neutral"
	NotifUndefined  NotificationIconType = "Undefined"
)

func (lt NotificationIconType) IsValid() error {
	switch lt {
	case NotifIncreasing, NotifDecreasing, NotifNeutral:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "notification_icon_type", "Increasing,Decreasing, Neutral or NotifUndefined"))
}

func (lt NotificationIconType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

// Unit of Measurement Type
type MeasurementType string

const (
	UnitCurrency   MeasurementType = "currency"
	UnitDataDay    MeasurementType = "data-day"
	UnitPercentage MeasurementType = "percentage"
	UnitBoolean    MeasurementType = "boolean"
	UnitInteger    MeasurementType = "integer"
	UnitString     MeasurementType = "string"
)

func (mt MeasurementType) IsValid() error {
	switch mt {
	case UnitCurrency, UnitDataDay, UnitPercentage, UnitBoolean, UnitInteger, UnitString:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "MeasurementType", "UnitCurrency,UnitDataDay,UnitPercentage,UnitBoolean,UnitInteger,UnitString"))
}

func (lt MeasurementType) isNil() bool {
	if len(string(lt)) == 0 {
		return true
	}
	return false
}

func (lt MeasurementType) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}

// Organizaation Extended Filter
type OrganizationExtendedFilter string

const (
	ExtendedCurrency    OrganizationExtendedFilter = "currency"
	ExtendedMarketPrice OrganizationExtendedFilter = "market_price"
	ExtendedName        OrganizationExtendedFilter = "name"
	ExtendedTicker      OrganizationExtendedFilter = "ticker"
)

func (mt OrganizationExtendedFilter) IsValid() error {
	switch mt {
	case ExtendedCurrency, ExtendedMarketPrice, ExtendedName, ExtendedTicker:
		return nil
	}
	return errors.New(fmt.Sprintf(constants.FieldMustBeType, "OrganizationExtendedFilter", "currency, market_price, name, ticker"))
}

func (lt OrganizationExtendedFilter) isNil() bool {
	if len(string(lt)) == 0 {
		return true
	}
	return false
}

func (lt OrganizationExtendedFilter) Value() (driver.Value, error) {
	if len(string(lt)) == 0 {
		return nil, nil
	}
	return string(lt), nil
}
