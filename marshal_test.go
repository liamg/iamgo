package iamgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {

	tests := []struct {
		name     string
		document Document
		expected string
	}{
		{
			name: "simple policy",
			expected: `
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": ["s3:ListBucket"],
      "Effect": "Allow",
      "Resource": ["arn:aws:s3:::mybucket"],
      "Condition": {"StringLike": {"s3:prefix": ["David/*"]}}
    },
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": ["arn:aws:s3:::mybucket/David/*"]
    }
  ]
}
            `,
			document: Document{
				Version: Version20121017,
				Statement: Statements{
					{
						Action:   StringSet{"s3:ListBucket"},
						Effect:   "Allow",
						Resource: StringSet{"arn:aws:s3:::mybucket"},
						Condition: Conditions{
							{
								Operator: "StringLike",
								Key:      "s3:prefix",
								Value:    StringSet{"David/*"},
							},
						},
					},
					{
						Action:   StringSet{"s3:GetObject", "s3:PutObject"},
						Effect:   "Allow",
						Resource: StringSet{"arn:aws:s3:::mybucket/David/*"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := json.Marshal(test.document)
			require.NoError(t, err)

			var actualMap map[string]interface{}
			var expectedMap map[string]interface{}

			require.NoError(t, json.Unmarshal(actual, &actualMap))
			require.NoError(t, json.Unmarshal([]byte(test.expected), &expectedMap))

			assert.Equal(t, expectedMap, actualMap)
		})
	}

}
