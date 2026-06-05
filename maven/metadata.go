package maven

import (
	"bytes"
	"encoding/xml"
	"io"
)

type Metadata struct {
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Versions   []string `xml:"-"`
	Latest     string   `xml:"-"`
	Release    string   `xml:"-"`
}

type xmlMetadata struct {
	GroupID    string        `xml:"groupId"`
	ArtifactID string        `xml:"artifactId"`
	Versioning xmlVersioning `xml:"versioning"`
}

type xmlVersioning struct {
	Latest   string      `xml:"latest"`
	Release  string      `xml:"release"`
	Versions xmlVersions `xml:"versions"`
}

type xmlVersions struct {
	Version []string `xml:"version"`
}

func ParseMetadata(data []byte, groupID, artifactID string) (Metadata, error) {
	var raw xmlMetadata
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.CharsetReader = func(_ string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	if err := dec.Decode(&raw); err != nil {
		return Metadata{}, err
	}
	return Metadata{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Versions:   raw.Versioning.Versions.Version,
		Latest:     raw.Versioning.Latest,
		Release:    raw.Versioning.Release,
	}, nil
}
