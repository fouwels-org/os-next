package config

import (
	"bytes"
	"encoding/json"
)

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

//UnmarshalJSON to error on missing fields
func (f *PrimaryFile) UnmarshalJSON(data []byte) error {

	f2 := primaryFile{}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force errors

	if err := dec.Decode(&f2); err != nil {
		return err
	}

	f.primaryFile = f2
	return nil
}

//UnmarshalJSON to error on missing fields
func (f *SecondaryFile) UnmarshalJSON(data []byte) error {

	f2 := secondaryFile{}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force errors

	if err := dec.Decode(&f2); err != nil {
		return err
	}

	f.secondaryFile = f2
	return nil
}
