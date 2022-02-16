package iamgo

import (
    "github.com/stretchr/testify/assert"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestFull(t *testing.T) {

    tests := []struct {
        name    string
        json    string
        assert  func(t *testing.T, d *Document)
        wantErr bool
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
  "Id": "testing",
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "s3:ListBucket",
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::mybucket",
      "Condition": {"StringLike": {"s3:prefix": ["David/*"]}},
      "Principal": "*"
    },
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:s3:::mybucket/David/*", 
        "*"
      ],
      "Principal": { 
        "AWS": "blah",
        "Federated": ["a", "b"]
      }
    }
  ]
}
            `,
            assert: func(t *testing.T, d *Document) {
                id, r := d.ID()
                assert.Equal(t, "testing", id)
                assert.Equal(t, 3, r.StartLine)
                assert.Equal(t, 3, r.EndLine)

                v, r := d.Version()
                assert.Equal(t, "2012-10-17", v)
                assert.Equal(t, 4, r.StartLine)
                assert.Equal(t, 4, r.EndLine)

                statements, r := d.Statements()
                assert.Equal(t, 5, r.StartLine)
                assert.Equal(t, 28, r.EndLine)
                require.Len(t, statements, 2)

                // first statement
                {
                    assert.Equal(t, 6, statements[0].Range().StartLine)
                    assert.Equal(t, 12, statements[0].Range().EndLine)
                    actions, r := statements[0].Actions()
                    assert.Equal(t, 7, r.StartLine)
                    assert.Equal(t, 7, r.EndLine)
                    require.Len(t, actions, 1)
                    assert.Equal(t, "s3:ListBucket", actions[0])
                    effect, r := statements[0].Effect()
                    assert.Equal(t, 8, r.StartLine)
                    assert.Equal(t, 8, r.EndLine)
                    assert.Equal(t, "Allow", effect)
                    resources, r := statements[0].Resources()
                    assert.Equal(t, 9, r.StartLine)
                    assert.Equal(t, 9, r.EndLine)
                    require.Len(t, resources, 1)
                    assert.Equal(t, "arn:aws:s3:::mybucket", resources[0])
                    conditions, r := statements[0].Conditions()
                    assert.Equal(t, 10, r.StartLine)
                    assert.Equal(t, 10, r.EndLine)
                    require.Len(t, conditions, 1)
                    key, r := conditions[0].Key()
                    assert.Equal(t, 10, r.StartLine)
                    assert.Equal(t, 10, r.EndLine)
                    assert.Equal(t, "s3:prefix", key)
                    operator, r := conditions[0].Operator()
                    assert.Equal(t, 10, r.StartLine)
                    assert.Equal(t, 10, r.EndLine)
                    assert.Equal(t, "StringLike", operator)
                    val, r := conditions[0].Value()
                    assert.Equal(t, 10, r.StartLine)
                    assert.Equal(t, 10, r.EndLine)
                    require.Len(t, val, 1)
                    assert.Equal(t, "David/*", val[0])
                    principals, r := statements[0].Principals()
                    assert.Equal(t, 11, r.StartLine)
                    assert.Equal(t, 11, r.EndLine)
                    all, r := principals.All()
                    assert.Equal(t, 11, r.StartLine)
                    assert.Equal(t, 11, r.EndLine)
                    assert.True(t, all)
                }

                // second statement
                {
                    assert.Equal(t, 13, statements[1].Range().StartLine)
                    assert.Equal(t, 27, statements[1].Range().EndLine)
                    actions, r := statements[1].Actions()
                    assert.Equal(t, 14, r.StartLine)
                    assert.Equal(t, 17, r.EndLine)
                    require.Len(t, actions, 2)
                    assert.Equal(t, "s3:GetObject", actions[0])
                    assert.Equal(t, "s3:PutObject", actions[1])
                    effect, r := statements[1].Effect()
                    assert.Equal(t, 18, r.StartLine)
                    assert.Equal(t, 18, r.EndLine)
                    assert.Equal(t, "Allow", effect)
                    resources, r := statements[1].Resources()
                    assert.Equal(t, 19, r.StartLine)
                    assert.Equal(t, 22, r.EndLine)
                    require.Len(t, resources, 2)
                    assert.Equal(t, "arn:aws:s3:::mybucket/David/*", resources[0])
                    assert.Equal(t, "*", resources[1])
                    conditions, _ := statements[1].Conditions()
                    require.Len(t, conditions, 0)
                    principals, r := statements[1].Principals()
                    assert.Equal(t, 23, r.StartLine)
                    assert.Equal(t, 26, r.EndLine)
                    aws, r := principals.AWS()
                    assert.Equal(t, 24, r.StartLine)
                    assert.Equal(t, 24, r.EndLine)
                    require.Len(t, aws, 1)
                    assert.Equal(t, "blah", aws[0])
                    federated, r := principals.Federated()
                    assert.Equal(t, 25, r.StartLine)
                    assert.Equal(t, 25, r.EndLine)
                    require.Len(t, federated, 2)
                    assert.Equal(t, "a", federated[0])
                    assert.Equal(t, "b", federated[1])
                    all, _ := principals.All()
                    assert.False(t, all)
                }
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            doc, err := Parse([]byte(test.json))
            if test.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
            if test.assert != nil {
                test.assert(t, doc)
            }
        })
    }

}
