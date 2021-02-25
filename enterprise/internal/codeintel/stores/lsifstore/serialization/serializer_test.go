package serialization

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/lsifstore"
)

func TestDocumentData(t *testing.T) {
	expected := lsifstore.DocumentData{
		Ranges: map[lsifstore.ID]lsifstore.RangeData{
			lsifstore.ID("7864"): {
				StartLine:          541,
				StartCharacter:     10,
				EndLine:            541,
				EndCharacter:       12,
				DefinitionResultID: ID("1266"),
				ReferenceResultID:  ID("15871"),
				HoverResultID:      ID("1269"),
				MonikerIDs:         nil,
			},
			lsifstore.ID("8265"): {
				StartLine:          266,
				StartCharacter:     10,
				EndLine:            266,
				EndCharacter:       16,
				DefinitionResultID: lsifstore.ID("311"),
				ReferenceResultID:  lsifstore.ID("15500"),
				HoverResultID:      lsifstore.ID("317"),
				MonikerIDs:         []lsifstore.ID{lsifstore.ID("314")},
			},
		},
		HoverResults: map[lsifstore.ID]string{
			lsifstore.ID("1269"): "```go\nvar id string\n```",
			lsifstore.ID("317"):  "```go\ntype Vertex struct\n```\n\n---\n\nVertex contains information of a vertex in the graph.\n\n---\n\n```go\nstruct {\n    Element\n    Label VertexLabel \"json:\\\"label\\\"\"\n}\n```",
		},
		Monikers: map[lsifstore.ID]lsifstore.MonikerData{
			lsifstore.ID("314"): {
				Kind:                 "export",
				Scheme:               "gomod",
				Identifier:           "github.com/sourcegraph/lsif-go/protocol:Vertex",
				PackageInformationID: lsifstore.ID("213"),
			},
			lsifstore.ID("2494"): {
				Kind:                 "export",
				Scheme:               "gomod",
				Identifier:           "github.com/sourcegraph/lsif-go/protocol:VertexLabel",
				PackageInformationID: lsifstore.ID("213"),
			},
		},
		PackageInformation: map[lsifstore.ID]lsifstore.PackageInformationData{
			lsifstore.ID("213"): {
				Name:    "github.com/sourcegraph/lsif-go",
				Version: "v0.0.0-ad3507cbeb18",
			},
		},
	}

	serializer := NewSerializer()

	recompressed, err := serializer.MarshalDocumentData(expected)
	if err != nil {
		t.Fatalf("unexpected error marshalling document data: %s", err)
	}

	roundtripActual, err := serializer.UnmarshalDocumentData(recompressed)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling document data: %s", err)
	}

	if diff := cmp.Diff(expected, roundtripActual); diff != "" {
		t.Errorf("unexpected document data (-want +got):\n%s", diff)
	}
}

func TestResultChunkData(t *testing.T) {
	expected := lsifstore.ResultChunkData{
		DocumentPaths: map[lsifstore.ID]string{
			lsifstore.ID("4"):   "internal/gomod/module.go",
			lsifstore.ID("302"): "protocol/protocol.go",
			lsifstore.ID("305"): "protocol/writer.go",
		},
		DocumentIDRangeIDs: map[lsifstore.ID][]lsifstore.DocumentIDRangeID{
			lsifstore.ID("34"): {
				{DocumentID: lsifstore.ID("4"), RangeID: lsifstore.ID("31")},
			},
			lsifstore.ID("14040"): {
				{DocumentID: lsifstore.ID("3978"), RangeID: lsifstore.ID("4544")},
			},
			lsifstore.ID("14051"): {
				{DocumentID: lsifstore.ID("3978"), RangeID: lsifstore.ID("4568")},
				{DocumentID: lsifstore.ID("3978"), RangeID: lsifstore.ID("9224")},
				{DocumentID: lsifstore.ID("3978"), RangeID: lsifstore.ID("9935")},
				{DocumentID: lsifstore.ID("3978"), RangeID: lsifstore.ID("9996")},
			},
		},
	}

	serializer := NewSerializer()

	recompressed, err := serializer.MarshalResultChunkData(expected)
	if err != nil {
		t.Fatalf("unexpected error marshalling result chunk data: %s", err)
	}

	roundtripActual, err := serializer.UnmarshalResultChunkData(recompressed)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling result chunk data: %s", err)
	}

	if diff := cmp.Diff(expected, roundtripActual); diff != "" {
		t.Errorf("unexpected document data (-want +got):\n%s", diff)
	}
}

func TestLocations(t *testing.T) {
	expected := []lsifstore.LocationData{
		{
			URI:            "internal/index/indexer.go",
			StartLine:      36,
			StartCharacter: 26,
			EndLine:        36,
			EndCharacter:   32,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      100,
			StartCharacter: 9,
			EndLine:        100,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      115,
			StartCharacter: 9,
			EndLine:        115,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      95,
			StartCharacter: 9,
			EndLine:        95,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      130,
			StartCharacter: 9,
			EndLine:        130,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      155,
			StartCharacter: 9,
			EndLine:        155,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      80,
			StartCharacter: 9,
			EndLine:        80,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      36,
			StartCharacter: 9,
			EndLine:        36,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      135,
			StartCharacter: 9,
			EndLine:        135,
			EndCharacter:   15,
		},
		{
			URI:            "protocol/writer.go",
			StartLine:      12,
			StartCharacter: 5,
			EndLine:        12,
			EndCharacter:   11,
		},
	}

	serializer := NewSerializer()

	recompressed, err := serializer.MarshalLocations(expected)
	if err != nil {
		t.Fatalf("unexpected error marshalling locations: %s", err)
	}

	roundtripActual, err := serializer.UnmarshalLocations(recompressed)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling locations: %s", err)
	}

	if diff := cmp.Diff(expected, roundtripActual); diff != "" {
		t.Errorf("unexpected locations (-want +got):\n%s", diff)
	}
}
