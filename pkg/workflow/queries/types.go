package queries

type Wrapper struct {
}

type Token interface {
	int32 | string
}

type DataTemplate struct {
	Key1Link      string
	Key2Link      string
	TimeField     string
	InfoField     string
	AmountField   string
	TrackingTable string
	PkField       string
	PkStatus      string
	TableName     string
	Statuses      string
	Status        string
	Token         string
	Filter        string
	Fields        string
}
