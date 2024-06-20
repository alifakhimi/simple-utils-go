package enums

// MilitaryStatus ...
type MilitaryStatus uint

// Military status kinds ...
const (
	ConscriptionExempt MilitaryStatus = iota + 1
	ConscriptionFinished
	Soldier
	Inductee
)

// Gender ...
type Gender uint

// Gender kinds ...
const (
	GENDER_UNKNOWN        = iota
	GENDER_MALE    Gender = iota
	GENDER_FEMALE
)

// EduDegree ...
type EduDegree int

// Education Degree kinds
const (
	UnknownEducationType EduDegree = iota
	BelowDiploma
	Diploma
	AssociateDegree
	Bachelor
	MastersDegree
	Doctor
)

// DatabaseDriver ...
type DatabaseDriver int

const (
	// Postgres SQL ...
	PostgresSQL DatabaseDriver = iota + 1
	// SQL Server ...
	SQLServer
	// SQLite ...
	SQLite
)
