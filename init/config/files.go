package config

//PrimaryFile superset of primaryFile to fail on unknown keys
type PrimaryFile struct {
	primaryFile
}

type primaryFile struct {
	Primary Primary
	Header  Header
}

//SecondaryFile superset of secondaryFile to fail on unknown keys
type SecondaryFile struct {
	secondaryFile
}

type secondaryFile struct {
	Secondary Secondary
	Header    Header
}
