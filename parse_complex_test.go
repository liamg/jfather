package jfather

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Complex(t *testing.T) {
	target := make(map[string]any)
	input := `{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}`
	require.NoError(t, Unmarshal([]byte(input), &target))
}

type Resource struct {
	inner resourceInner
	Line  int
}

type resourceInner struct {
	Properties map[string]*Property `json:"Properties" yaml:"Properties"`
	Type       string               `json:"Type" yaml:"Type"`
}

func (r *Resource) UnmarshalJSONWithMetadata(node Node) error {
	r.Line = node.Range().Start.Line
	return node.Decode(&r.inner)
}

type Parameter struct {
	inner parameterInner
}

type parameterInner struct {
	Default any    `yaml:"Default"`
	Type    string `json:"Type" yaml:"Type"`
}

func (p *Parameter) UnmarshalJSONWithMetadata(node Node) error {
	return node.Decode(&p.inner)
}

type Property struct {
	inner propertyInner
	Line  int
}

type CFType string

type propertyInner struct {
	Value any `json:"Value" yaml:"Value"`
	Type  CFType
}

func (p *Property) UnmarshalJSONWithMetadata(node Node) error {
	p.Line = node.Range().Start.Line
	return node.Decode(&p.inner)
}

type Temp struct {
	BucketName       *Parameter
	BucketKeyEnabled *Parameter
}

type FileContext struct {
	Parameters map[string]*Parameter `json:"Parameters" yaml:"Parameters"`
	Resources  map[string]*Resource  `json:"Resources" yaml:"Resources"`
}

func Test_CloudFormation(t *testing.T) {
	var target FileContext
	input := `
{
  "Parameters": {
   "BucketName":  {
      "Type": "String",
      "Default": "naughty"
    },
	"BucketKeyEnabled": {
      "Type": "Boolean",
      "Default": false
    }
  },
  "Resources": {
    "S3Bucket": {
      "Type": "AWS::S3::Bucket",
      "Properties": {
        "BucketName": {
          "Ref": "BucketName"
        },
        "BucketEncryption": {
          "ServerSideEncryptionConfiguration": [
            {
              "BucketKeyEnabled": {
                "Ref": "BucketKeyEnabled"
              }
            }
          ]
        }
      }
    }
  }
}
`
	require.NoError(t, Unmarshal([]byte(input), &target))
}
