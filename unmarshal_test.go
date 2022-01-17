package iamgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {

	tests := []struct {
		name     string
		json     string
		expected Document
		wantErr  bool
	}{
		{
			name:    "bad json",
			json:    `{ oh no this has been corrupted`,
			wantErr: true,
		},
		{
			name: "simple policy",
			json: `
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
			expected: Document{
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
		{
			name: "principals",
			json: `
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [
          "elasticmapreduce.amazonaws.com",
          "datapipeline.amazonaws.com"
        ]
      },
      "Action": "sts:AssumeRole"
    }
  ]
}       `,
			expected: Document{
				Version: Version20121017,
				Statement: Statements{
					{
						Action: StringSet{"sts:AssumeRole"},
						Effect: "Allow",
						Principal: &Principals{
							Service: StringSet{"elasticmapreduce.amazonaws.com", "datapipeline.amazonaws.com"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual Document
			err := json.Unmarshal([]byte(test.json), &actual)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expected, actual)
		})
	}

}
