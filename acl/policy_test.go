package acl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func errStartsWith(t *testing.T, actual error, expected string) {
	t.Helper()
	require.Error(t, actual)
	require.Truef(t, strings.HasPrefix(actual.Error(), expected), "Received unexpected error: %#v\nExpecting an error with the prefix: %q", actual, expected)
}

func TestPolicySourceParse(t *testing.T) {
	ljoin := func(lines ...string) string {
		return strings.Join(lines, "\n")
	}
	cases := []struct {
		Name     string
		Syntax   SyntaxVersion
		Rules    string
		Expected *Policy
		Err      string
	}{
		{
			"Legacy Basic",
			SyntaxLegacy,
			ljoin(
				`agent "foo" {      `,
				`	policy = "read"  `,
				`}                  `,
				`agent "bar" {      `,
				`	policy = "write" `,
				`}                  `,
				`event "" {         `,
				`	policy = "read"  `,
				`}                  `,
				`event "foo" {      `,
				`	policy = "write" `,
				`}                  `,
				`event "bar" {      `,
				`	policy = "deny"  `,
				`}                  `,
				`key "" {           `,
				`	policy = "read"  `,
				`}                  `,
				`key "foo/" {       `,
				`	policy = "write" `,
				`}                  `,
				`key "foo/bar/" {   `,
				`	policy = "read"  `,
				`}                  `,
				`key "foo/bar/baz" {`,
				`	policy = "deny"  `,
				`}                  `,
				`keyring = "deny"   `,
				`node "" {          `,
				`	policy = "read"  `,
				`}                  `,
				`node "foo" {       `,
				`	policy = "write" `,
				`}                  `,
				`node "bar" {       `,
				`	policy = "deny"  `,
				`}                  `,
				`operator = "deny"  `,
				`service "" {       `,
				`	policy = "write" `,
				`}                  `,
				`service "foo" {    `,
				`	policy = "read"  `,
				`}                  `,
				`session "foo" {    `,
				`	policy = "write" `,
				`}                  `,
				`session "bar" {    `,
				`	policy = "deny"  `,
				`}                  `,
				`query "" {         `,
				`	policy = "read"  `,
				`}                  `,
				`query "foo" {      `,
				`	policy = "write" `,
				`}                  `,
				`query "bar" {      `,
				`	policy = "deny"  `,
				`}                  `),
			&Policy{PolicyRules: PolicyRules{
				AgentPrefixes: []*AgentRule{
					&AgentRule{
						Node:   "foo",
						Policy: PolicyRead,
					},
					&AgentRule{
						Node:   "bar",
						Policy: PolicyWrite,
					},
				},
				EventPrefixes: []*EventRule{
					&EventRule{
						Event:  "",
						Policy: PolicyRead,
					},
					&EventRule{
						Event:  "foo",
						Policy: PolicyWrite,
					},
					&EventRule{
						Event:  "bar",
						Policy: PolicyDeny,
					},
				},
				Keyring: PolicyDeny,
				KeyPrefixes: []*KeyRule{
					&KeyRule{
						Prefix: "",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "foo/",
						Policy: PolicyWrite,
					},
					&KeyRule{
						Prefix: "foo/bar/",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "foo/bar/baz",
						Policy: PolicyDeny,
					},
				},
				NodePrefixes: []*NodeRule{
					&NodeRule{
						Name:   "",
						Policy: PolicyRead,
					},
					&NodeRule{
						Name:   "foo",
						Policy: PolicyWrite,
					},
					&NodeRule{
						Name:   "bar",
						Policy: PolicyDeny,
					},
				},
				Operator: PolicyDeny,
				PreparedQueryPrefixes: []*PreparedQueryRule{
					&PreparedQueryRule{
						Prefix: "",
						Policy: PolicyRead,
					},
					&PreparedQueryRule{
						Prefix: "foo",
						Policy: PolicyWrite,
					},
					&PreparedQueryRule{
						Prefix: "bar",
						Policy: PolicyDeny,
					},
				},
				ServicePrefixes: []*ServiceRule{
					&ServiceRule{
						Name:   "",
						Policy: PolicyWrite,
					},
					&ServiceRule{
						Name:   "foo",
						Policy: PolicyRead,
					},
				},
				SessionPrefixes: []*SessionRule{
					&SessionRule{
						Node:   "foo",
						Policy: PolicyWrite,
					},
					&SessionRule{
						Node:   "bar",
						Policy: PolicyDeny,
					},
				},
			}},
			"",
		},
		{
			"Legacy (JSON)",
			SyntaxLegacy,
			ljoin(
				`{                         `,
				`	"agent": {              `,
				`		"foo": {             `,
				`			"policy": "write" `,
				`		},                   `,
				`		"bar": {             `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	},                      `,
				`	"event": {              `,
				`		"": {                `,
				`			"policy": "read"  `,
				`		},                   `,
				`		"foo": {             `,
				`			"policy": "write" `,
				`		},                   `,
				`		"bar": {             `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	},                      `,
				`	"key": {                `,
				`		"": {                `,
				`			"policy": "read"  `,
				`		},                   `,
				`		"foo/": {            `,
				`			"policy": "write" `,
				`		},                   `,
				`		"foo/bar/": {        `,
				`			"policy": "read"  `,
				`		},                   `,
				`		"foo/bar/baz": {     `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	},                      `,
				`	"keyring": "deny",      `,
				`	"node": {               `,
				`		"": {                `,
				`			"policy": "read"  `,
				`		},                   `,
				`		"foo": {             `,
				`			"policy": "write" `,
				`		},                   `,
				`		"bar": {             `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	},                      `,
				`	"operator": "deny",     `,
				`	"query": {              `,
				`		"": {                `,
				`			"policy": "read"  `,
				`		},                   `,
				`		"foo": {             `,
				`			"policy": "write" `,
				`		},                   `,
				`		"bar": {             `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	},                      `,
				`	"service": {            `,
				`		"": {                `,
				`			"policy": "write" `,
				`		},                   `,
				`		"foo": {             `,
				`			"policy": "read"  `,
				`		}                    `,
				`	},                      `,
				`	"session": {            `,
				`		"foo": {             `,
				`			"policy": "write" `,
				`		},                   `,
				`		"bar": {             `,
				`			"policy": "deny"  `,
				`		}                    `,
				`	}                       `,
				`}                         `),
			&Policy{PolicyRules: PolicyRules{
				AgentPrefixes: []*AgentRule{
					&AgentRule{
						Node:   "foo",
						Policy: PolicyWrite,
					},
					&AgentRule{
						Node:   "bar",
						Policy: PolicyDeny,
					},
				},
				EventPrefixes: []*EventRule{
					&EventRule{
						Event:  "",
						Policy: PolicyRead,
					},
					&EventRule{
						Event:  "foo",
						Policy: PolicyWrite,
					},
					&EventRule{
						Event:  "bar",
						Policy: PolicyDeny,
					},
				},
				Keyring: PolicyDeny,
				KeyPrefixes: []*KeyRule{
					&KeyRule{
						Prefix: "",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "foo/",
						Policy: PolicyWrite,
					},
					&KeyRule{
						Prefix: "foo/bar/",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "foo/bar/baz",
						Policy: PolicyDeny,
					},
				},
				NodePrefixes: []*NodeRule{
					&NodeRule{
						Name:   "",
						Policy: PolicyRead,
					},
					&NodeRule{
						Name:   "foo",
						Policy: PolicyWrite,
					},
					&NodeRule{
						Name:   "bar",
						Policy: PolicyDeny,
					},
				},
				Operator: PolicyDeny,
				PreparedQueryPrefixes: []*PreparedQueryRule{
					&PreparedQueryRule{
						Prefix: "",
						Policy: PolicyRead,
					},
					&PreparedQueryRule{
						Prefix: "foo",
						Policy: PolicyWrite,
					},
					&PreparedQueryRule{
						Prefix: "bar",
						Policy: PolicyDeny,
					},
				},
				ServicePrefixes: []*ServiceRule{
					&ServiceRule{
						Name:   "",
						Policy: PolicyWrite,
					},
					&ServiceRule{
						Name:   "foo",
						Policy: PolicyRead,
					},
				},
				SessionPrefixes: []*SessionRule{
					&SessionRule{
						Node:   "foo",
						Policy: PolicyWrite,
					},
					&SessionRule{
						Node:   "bar",
						Policy: PolicyDeny,
					},
				},
			}},
			"",
		},
		{
			"Service No Intentions (Legacy)",
			SyntaxLegacy,
			ljoin(
				`service "foo" {    `,
				`   policy = "write"`,
				`}                  `),
			&Policy{PolicyRules: PolicyRules{
				ServicePrefixes: []*ServiceRule{
					{
						Name:   "foo",
						Policy: "write",
					},
				},
			}},
			"",
		},
		{
			"Service Intentions (Legacy)",
			SyntaxLegacy,
			ljoin(
				`service "foo" {       `,
				`   policy = "write"   `,
				`   intentions = "read"`,
				`}                     `),
			&Policy{PolicyRules: PolicyRules{
				ServicePrefixes: []*ServiceRule{
					{
						Name:       "foo",
						Policy:     "write",
						Intentions: "read",
					},
				},
			}},
			"",
		},
		{
			"Service Intention: invalid value (Legacy)",
			SyntaxLegacy,
			ljoin(
				`service "foo" {      `,
				`   policy = "write"  `,
				`   intentions = "foo"`,
				`}                    `),
			nil,
			"Invalid service intentions policy",
		},
		{
			"Bad Policy - ACL",

			SyntaxCurrent,
			`acl = "list"`, // there is no list policy but this helps to exercise another check in isPolicyValid
			nil,
			"Invalid acl policy",
		},
		{
			"Bad Policy - Agent",
			SyntaxCurrent,
			`agent "foo" { policy = "nope" }`,
			nil,
			"Invalid agent policy",
		},
		{
			"Bad Policy - Agent Prefix",
			SyntaxCurrent,
			`agent_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid agent_prefix policy",
		},
		{
			"Bad Policy - Key",
			SyntaxCurrent,
			`key "foo" { policy = "nope" }`,
			nil,
			"Invalid key policy",
		},
		{
			"Bad Policy - Key Prefix",
			SyntaxCurrent,
			`key_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid key_prefix policy",
		},
		{
			"Bad Policy - Node",
			SyntaxCurrent,
			`node "foo" { policy = "nope" }`,
			nil,
			"Invalid node policy",
		},
		{
			"Bad Policy - Node Prefix",
			SyntaxCurrent,
			`node_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid node_prefix policy",
		},
		{
			"Bad Policy - Service",
			SyntaxCurrent,
			`service "foo" { policy = "nope" }`,
			nil,
			"Invalid service policy",
		},
		{
			"Bad Policy - Service Prefix",
			SyntaxCurrent,
			`service_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid service_prefix policy",
		},
		{
			"Bad Policy - Session",
			SyntaxCurrent,
			`session "foo" { policy = "nope" }`,
			nil,
			"Invalid session policy",
		},
		{
			"Bad Policy - Session Prefix",
			SyntaxCurrent,
			`session_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid session_prefix policy",
		},
		{
			"Bad Policy - Event",
			SyntaxCurrent,
			`event "foo" { policy = "nope" }`,
			nil,
			"Invalid event policy",
		},
		{
			"Bad Policy - Event Prefix",
			SyntaxCurrent,
			`event_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid event_prefix policy",
		},
		{
			"Bad Policy - Prepared Query",
			SyntaxCurrent,
			`query "foo" { policy = "nope" }`,
			nil,
			"Invalid query policy",
		},
		{
			"Bad Policy - Prepared Query Prefix",
			SyntaxCurrent,
			`query_prefix "foo" { policy = "nope" }`,
			nil,
			"Invalid query_prefix policy",
		},
		{
			"Bad Policy - Keyring",
			SyntaxCurrent,
			`keyring = "nope"`,
			nil,
			"Invalid keyring policy",
		},
		{
			"Bad Policy - Operator",
			SyntaxCurrent,
			`operator = "nope"`,
			nil,
			"Invalid operator policy",
		},
		{
			"Keyring Empty",
			SyntaxCurrent,
			`keyring = ""`,
			&Policy{PolicyRules: PolicyRules{Keyring: ""}},
			"",
		},
		{
			"Operator Empty",
			SyntaxCurrent,
			`operator = ""`,
			&Policy{PolicyRules: PolicyRules{Operator: ""}},
			"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			req := require.New(t)
			actual, err := NewPolicyFromSource("", 0, tc.Rules, tc.Syntax, nil, nil)
			if tc.Err != "" {
				errStartsWith(t, err, tc.Err)
			} else {
				req.Equal(tc.Expected, actual)
			}
		})
	}
}

func TestMergePolicies(t *testing.T) {
	type mergeTest struct {
		name     string
		input    []*Policy
		expected *Policy
	}

	tests := []mergeTest{
		{
			name: "Agents",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Agents: []*AgentRule{
						&AgentRule{
							Node:   "foo",
							Policy: PolicyWrite,
						},
						&AgentRule{
							Node:   "bar",
							Policy: PolicyRead,
						},
						&AgentRule{
							Node:   "baz",
							Policy: PolicyWrite,
						},
					},
					AgentPrefixes: []*AgentRule{
						&AgentRule{
							Node:   "000",
							Policy: PolicyWrite,
						},
						&AgentRule{
							Node:   "111",
							Policy: PolicyRead,
						},
						&AgentRule{
							Node:   "222",
							Policy: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Agents: []*AgentRule{
						&AgentRule{
							Node:   "foo",
							Policy: PolicyRead,
						},
						&AgentRule{
							Node:   "baz",
							Policy: PolicyDeny,
						},
					},
					AgentPrefixes: []*AgentRule{
						&AgentRule{
							Node:   "000",
							Policy: PolicyRead,
						},
						&AgentRule{
							Node:   "222",
							Policy: PolicyDeny,
						},
					},
				},
				}},
			expected: &Policy{PolicyRules: PolicyRules{
				Agents: []*AgentRule{
					&AgentRule{
						Node:   "foo",
						Policy: PolicyWrite,
					},
					&AgentRule{
						Node:   "bar",
						Policy: PolicyRead,
					},
					&AgentRule{
						Node:   "baz",
						Policy: PolicyDeny,
					},
				},
				AgentPrefixes: []*AgentRule{
					&AgentRule{
						Node:   "000",
						Policy: PolicyWrite,
					},
					&AgentRule{
						Node:   "111",
						Policy: PolicyRead,
					},
					&AgentRule{
						Node:   "222",
						Policy: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Events",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Events: []*EventRule{
						&EventRule{
							Event:  "foo",
							Policy: PolicyWrite,
						},
						&EventRule{
							Event:  "bar",
							Policy: PolicyRead,
						},
						&EventRule{
							Event:  "baz",
							Policy: PolicyWrite,
						},
					},
					EventPrefixes: []*EventRule{
						&EventRule{
							Event:  "000",
							Policy: PolicyWrite,
						},
						&EventRule{
							Event:  "111",
							Policy: PolicyRead,
						},
						&EventRule{
							Event:  "222",
							Policy: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Events: []*EventRule{
						&EventRule{
							Event:  "foo",
							Policy: PolicyRead,
						},
						&EventRule{
							Event:  "baz",
							Policy: PolicyDeny,
						},
					},
					EventPrefixes: []*EventRule{
						&EventRule{
							Event:  "000",
							Policy: PolicyRead,
						},
						&EventRule{
							Event:  "222",
							Policy: PolicyDeny,
						},
					},
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				Events: []*EventRule{
					&EventRule{
						Event:  "foo",
						Policy: PolicyWrite,
					},
					&EventRule{
						Event:  "bar",
						Policy: PolicyRead,
					},
					&EventRule{
						Event:  "baz",
						Policy: PolicyDeny,
					},
				},
				EventPrefixes: []*EventRule{
					&EventRule{
						Event:  "000",
						Policy: PolicyWrite,
					},
					&EventRule{
						Event:  "111",
						Policy: PolicyRead,
					},
					&EventRule{
						Event:  "222",
						Policy: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Node",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Nodes: []*NodeRule{
						&NodeRule{
							Name:   "foo",
							Policy: PolicyWrite,
						},
						&NodeRule{
							Name:   "bar",
							Policy: PolicyRead,
						},
						&NodeRule{
							Name:   "baz",
							Policy: PolicyWrite,
						},
					},
					NodePrefixes: []*NodeRule{
						&NodeRule{
							Name:   "000",
							Policy: PolicyWrite,
						},
						&NodeRule{
							Name:   "111",
							Policy: PolicyRead,
						},
						&NodeRule{
							Name:   "222",
							Policy: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Nodes: []*NodeRule{
						&NodeRule{
							Name:   "foo",
							Policy: PolicyRead,
						},
						&NodeRule{
							Name:   "baz",
							Policy: PolicyDeny,
						},
					},
					NodePrefixes: []*NodeRule{
						&NodeRule{
							Name:   "000",
							Policy: PolicyRead,
						},
						&NodeRule{
							Name:   "222",
							Policy: PolicyDeny,
						},
					},
				},
				}},
			expected: &Policy{PolicyRules: PolicyRules{
				Nodes: []*NodeRule{
					&NodeRule{
						Name:   "foo",
						Policy: PolicyWrite,
					},
					&NodeRule{
						Name:   "bar",
						Policy: PolicyRead,
					},
					&NodeRule{
						Name:   "baz",
						Policy: PolicyDeny,
					},
				},
				NodePrefixes: []*NodeRule{
					&NodeRule{
						Name:   "000",
						Policy: PolicyWrite,
					},
					&NodeRule{
						Name:   "111",
						Policy: PolicyRead,
					},
					&NodeRule{
						Name:   "222",
						Policy: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Keys",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Keys: []*KeyRule{
						&KeyRule{
							Prefix: "foo",
							Policy: PolicyWrite,
						},
						&KeyRule{
							Prefix: "bar",
							Policy: PolicyRead,
						},
						&KeyRule{
							Prefix: "baz",
							Policy: PolicyWrite,
						},
						&KeyRule{
							Prefix: "zoo",
							Policy: PolicyList,
						},
					},
					KeyPrefixes: []*KeyRule{
						&KeyRule{
							Prefix: "000",
							Policy: PolicyWrite,
						},
						&KeyRule{
							Prefix: "111",
							Policy: PolicyRead,
						},
						&KeyRule{
							Prefix: "222",
							Policy: PolicyWrite,
						},
						&KeyRule{
							Prefix: "333",
							Policy: PolicyList,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Keys: []*KeyRule{
						&KeyRule{
							Prefix: "foo",
							Policy: PolicyRead,
						},
						&KeyRule{
							Prefix: "baz",
							Policy: PolicyDeny,
						},
						&KeyRule{
							Prefix: "zoo",
							Policy: PolicyRead,
						},
					},
					KeyPrefixes: []*KeyRule{
						&KeyRule{
							Prefix: "000",
							Policy: PolicyRead,
						},
						&KeyRule{
							Prefix: "222",
							Policy: PolicyDeny,
						},
						&KeyRule{
							Prefix: "333",
							Policy: PolicyRead,
						},
					},
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				Keys: []*KeyRule{
					&KeyRule{
						Prefix: "foo",
						Policy: PolicyWrite,
					},
					&KeyRule{
						Prefix: "bar",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "baz",
						Policy: PolicyDeny,
					},
					&KeyRule{
						Prefix: "zoo",
						Policy: PolicyList,
					},
				},
				KeyPrefixes: []*KeyRule{
					&KeyRule{
						Prefix: "000",
						Policy: PolicyWrite,
					},
					&KeyRule{
						Prefix: "111",
						Policy: PolicyRead,
					},
					&KeyRule{
						Prefix: "222",
						Policy: PolicyDeny,
					},
					&KeyRule{
						Prefix: "333",
						Policy: PolicyList,
					},
				},
			}},
		},
		{
			name: "Services",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Services: []*ServiceRule{
						&ServiceRule{
							Name:       "foo",
							Policy:     PolicyWrite,
							Intentions: PolicyWrite,
						},
						&ServiceRule{
							Name:       "bar",
							Policy:     PolicyRead,
							Intentions: PolicyRead,
						},
						&ServiceRule{
							Name:       "baz",
							Policy:     PolicyWrite,
							Intentions: PolicyWrite,
						},
					},
					ServicePrefixes: []*ServiceRule{
						&ServiceRule{
							Name:       "000",
							Policy:     PolicyWrite,
							Intentions: PolicyWrite,
						},
						&ServiceRule{
							Name:       "111",
							Policy:     PolicyRead,
							Intentions: PolicyRead,
						},
						&ServiceRule{
							Name:       "222",
							Policy:     PolicyWrite,
							Intentions: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Services: []*ServiceRule{
						&ServiceRule{
							Name:       "foo",
							Policy:     PolicyRead,
							Intentions: PolicyRead,
						},
						&ServiceRule{
							Name:       "baz",
							Policy:     PolicyDeny,
							Intentions: PolicyDeny,
						},
					},
					ServicePrefixes: []*ServiceRule{
						&ServiceRule{
							Name:       "000",
							Policy:     PolicyRead,
							Intentions: PolicyRead,
						},
						&ServiceRule{
							Name:       "222",
							Policy:     PolicyDeny,
							Intentions: PolicyDeny,
						},
					},
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				Services: []*ServiceRule{
					&ServiceRule{
						Name:       "foo",
						Policy:     PolicyWrite,
						Intentions: PolicyWrite,
					},
					&ServiceRule{
						Name:       "bar",
						Policy:     PolicyRead,
						Intentions: PolicyRead,
					},
					&ServiceRule{
						Name:       "baz",
						Policy:     PolicyDeny,
						Intentions: PolicyDeny,
					},
				},
				ServicePrefixes: []*ServiceRule{
					&ServiceRule{
						Name:       "000",
						Policy:     PolicyWrite,
						Intentions: PolicyWrite,
					},
					&ServiceRule{
						Name:       "111",
						Policy:     PolicyRead,
						Intentions: PolicyRead,
					},
					&ServiceRule{
						Name:       "222",
						Policy:     PolicyDeny,
						Intentions: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Sessions",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					Sessions: []*SessionRule{
						&SessionRule{
							Node:   "foo",
							Policy: PolicyWrite,
						},
						&SessionRule{
							Node:   "bar",
							Policy: PolicyRead,
						},
						&SessionRule{
							Node:   "baz",
							Policy: PolicyWrite,
						},
					},
					SessionPrefixes: []*SessionRule{
						&SessionRule{
							Node:   "000",
							Policy: PolicyWrite,
						},
						&SessionRule{
							Node:   "111",
							Policy: PolicyRead,
						},
						&SessionRule{
							Node:   "222",
							Policy: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					Sessions: []*SessionRule{
						&SessionRule{
							Node:   "foo",
							Policy: PolicyRead,
						},
						&SessionRule{
							Node:   "baz",
							Policy: PolicyDeny,
						},
					},
					SessionPrefixes: []*SessionRule{
						&SessionRule{
							Node:   "000",
							Policy: PolicyRead,
						},
						&SessionRule{
							Node:   "222",
							Policy: PolicyDeny,
						},
					},
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				Sessions: []*SessionRule{
					&SessionRule{
						Node:   "foo",
						Policy: PolicyWrite,
					},
					&SessionRule{
						Node:   "bar",
						Policy: PolicyRead,
					},
					&SessionRule{
						Node:   "baz",
						Policy: PolicyDeny,
					},
				},
				SessionPrefixes: []*SessionRule{
					&SessionRule{
						Node:   "000",
						Policy: PolicyWrite,
					},
					&SessionRule{
						Node:   "111",
						Policy: PolicyRead,
					},
					&SessionRule{
						Node:   "222",
						Policy: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Prepared Queries",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					PreparedQueries: []*PreparedQueryRule{
						&PreparedQueryRule{
							Prefix: "foo",
							Policy: PolicyWrite,
						},
						&PreparedQueryRule{
							Prefix: "bar",
							Policy: PolicyRead,
						},
						&PreparedQueryRule{
							Prefix: "baz",
							Policy: PolicyWrite,
						},
					},
					PreparedQueryPrefixes: []*PreparedQueryRule{
						&PreparedQueryRule{
							Prefix: "000",
							Policy: PolicyWrite,
						},
						&PreparedQueryRule{
							Prefix: "111",
							Policy: PolicyRead,
						},
						&PreparedQueryRule{
							Prefix: "222",
							Policy: PolicyWrite,
						},
					},
				}},
				&Policy{PolicyRules: PolicyRules{
					PreparedQueries: []*PreparedQueryRule{
						&PreparedQueryRule{
							Prefix: "foo",
							Policy: PolicyRead,
						},
						&PreparedQueryRule{
							Prefix: "baz",
							Policy: PolicyDeny,
						},
					},
					PreparedQueryPrefixes: []*PreparedQueryRule{
						&PreparedQueryRule{
							Prefix: "000",
							Policy: PolicyRead,
						},
						&PreparedQueryRule{
							Prefix: "222",
							Policy: PolicyDeny,
						},
					},
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				PreparedQueries: []*PreparedQueryRule{
					&PreparedQueryRule{
						Prefix: "foo",
						Policy: PolicyWrite,
					},
					&PreparedQueryRule{
						Prefix: "bar",
						Policy: PolicyRead,
					},
					&PreparedQueryRule{
						Prefix: "baz",
						Policy: PolicyDeny,
					},
				},
				PreparedQueryPrefixes: []*PreparedQueryRule{
					&PreparedQueryRule{
						Prefix: "000",
						Policy: PolicyWrite,
					},
					&PreparedQueryRule{
						Prefix: "111",
						Policy: PolicyRead,
					},
					&PreparedQueryRule{
						Prefix: "222",
						Policy: PolicyDeny,
					},
				},
			}},
		},
		{
			name: "Write Precedence",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					ACL:      PolicyRead,
					Keyring:  PolicyRead,
					Operator: PolicyRead,
				}},
				&Policy{PolicyRules: PolicyRules{
					ACL:      PolicyWrite,
					Keyring:  PolicyWrite,
					Operator: PolicyWrite,
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				ACL:      PolicyWrite,
				Keyring:  PolicyWrite,
				Operator: PolicyWrite,
			}},
		},
		{
			name: "Deny Precedence",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					ACL:      PolicyWrite,
					Keyring:  PolicyWrite,
					Operator: PolicyWrite,
				}},
				&Policy{PolicyRules: PolicyRules{
					ACL:      PolicyDeny,
					Keyring:  PolicyDeny,
					Operator: PolicyDeny,
				}},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				ACL:      PolicyDeny,
				Keyring:  PolicyDeny,
				Operator: PolicyDeny,
			}},
		},
		{
			name: "Read Precedence",
			input: []*Policy{
				&Policy{PolicyRules: PolicyRules{
					ACL:      PolicyRead,
					Keyring:  PolicyRead,
					Operator: PolicyRead,
				}},
				&Policy{},
			},
			expected: &Policy{PolicyRules: PolicyRules{
				ACL:      PolicyRead,
				Keyring:  PolicyRead,
				Operator: PolicyRead,
			}},
		},
	}

	req := require.New(t)

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			act := MergePolicies(tcase.input)
			exp := tcase.expected
			req.Equal(exp.ACL, act.ACL)
			req.Equal(exp.Keyring, act.Keyring)
			req.Equal(exp.Operator, act.Operator)
			req.ElementsMatch(exp.Agents, act.Agents)
			req.ElementsMatch(exp.AgentPrefixes, act.AgentPrefixes)
			req.ElementsMatch(exp.Events, act.Events)
			req.ElementsMatch(exp.EventPrefixes, act.EventPrefixes)
			req.ElementsMatch(exp.Keys, act.Keys)
			req.ElementsMatch(exp.KeyPrefixes, act.KeyPrefixes)
			req.ElementsMatch(exp.Nodes, act.Nodes)
			req.ElementsMatch(exp.NodePrefixes, act.NodePrefixes)
			req.ElementsMatch(exp.PreparedQueries, act.PreparedQueries)
			req.ElementsMatch(exp.PreparedQueryPrefixes, act.PreparedQueryPrefixes)
			req.ElementsMatch(exp.Services, act.Services)
			req.ElementsMatch(exp.ServicePrefixes, act.ServicePrefixes)
			req.ElementsMatch(exp.Sessions, act.Sessions)
			req.ElementsMatch(exp.SessionPrefixes, act.SessionPrefixes)
		})
	}

}

func TestRulesTranslate(t *testing.T) {
	input := `
# top level comment

# block comment
agent "" {
  # policy comment
  policy = "write"
}

# block comment
key "" {
  # policy comment
  policy = "write"
}

# block comment
node "" {
  # policy comment
  policy = "write"
}

# block comment
event "" {
  # policy comment
  policy = "write"
}

# block comment
service "" {
  # policy comment
  policy = "write"
}

# block comment
session "" {
  # policy comment
  policy = "write"
}

# block comment
query "" {
  # policy comment
  policy = "write"
}

# comment
keyring = "write"

# comment
operator = "write"
`

	expected := `
# top level comment

# block comment
agent_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
key_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
node_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
event_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
service_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
session_prefix "" {
  # policy comment
  policy = "write"
}

# block comment
query_prefix "" {
  # policy comment
  policy = "write"
}

# comment
keyring = "write"

# comment
operator = "write"
`

	output, err := TranslateLegacyRules([]byte(input))
	require.NoError(t, err)
	require.Equal(t, strings.Trim(expected, "\n"), string(output))
}

func TestRulesTranslate_GH5493(t *testing.T) {
	input := `
{
	"key": {
		"": {
			"policy": "read"
		},
		"key": {
			"policy": "read"
		},
		"policy": {
			"policy": "read"
		},
		"privatething1/": {
			"policy": "deny"
		},
		"anapplication/private/": {
			"policy": "deny"
		},
		"privatething2/": {
			"policy": "deny"
		}
	},
	"session": {
		"": {
			"policy": "write"
		}
	},
	"node": {
		"": {
			"policy": "read"
		}
	},
	"agent": {
		"": {
			"policy": "read"
		}
	},
	"service": {
		"": {
			"policy": "read"
		}
	},
	"event": {
		"": {
			"policy": "read"
		}
	},
	"query": {
		"": {
			"policy": "read"
		}
	}
}`
	expected := `
key_prefix "" {
  policy = "read"
}

key_prefix "key" {
  policy = "read"
}

key_prefix "policy" {
  policy = "read"
}

key_prefix "privatething1/" {
  policy = "deny"
}

key_prefix "anapplication/private/" {
  policy = "deny"
}

key_prefix "privatething2/" {
  policy = "deny"
}

session_prefix "" {
  policy = "write"
}

node_prefix "" {
  policy = "read"
}

agent_prefix "" {
  policy = "read"
}

service_prefix "" {
  policy = "read"
}

event_prefix "" {
  policy = "read"
}

query_prefix "" {
  policy = "read"
}
`
	output, err := TranslateLegacyRules([]byte(input))
	require.NoError(t, err)
	require.Equal(t, strings.Trim(expected, "\n"), string(output))
}

func TestPrecedence(t *testing.T) {
	type testCase struct {
		name     string
		a        string
		b        string
		expected bool
	}

	cases := []testCase{
		{
			name:     "Deny Over Write",
			a:        PolicyDeny,
			b:        PolicyWrite,
			expected: true,
		},
		{
			name:     "Deny Over List",
			a:        PolicyDeny,
			b:        PolicyList,
			expected: true,
		},
		{
			name:     "Deny Over Read",
			a:        PolicyDeny,
			b:        PolicyRead,
			expected: true,
		},
		{
			name:     "Deny Over Unknown",
			a:        PolicyDeny,
			b:        "not a policy",
			expected: true,
		},
		{
			name:     "Write Over List",
			a:        PolicyWrite,
			b:        PolicyList,
			expected: true,
		},
		{
			name:     "Write Over Read",
			a:        PolicyWrite,
			b:        PolicyRead,
			expected: true,
		},
		{
			name:     "Write Over Unknown",
			a:        PolicyWrite,
			b:        "not a policy",
			expected: true,
		},
		{
			name:     "List Over Read",
			a:        PolicyList,
			b:        PolicyRead,
			expected: true,
		},
		{
			name:     "List Over Unknown",
			a:        PolicyList,
			b:        "not a policy",
			expected: true,
		},
		{
			name:     "Read Over Unknown",
			a:        PolicyRead,
			b:        "not a policy",
			expected: true,
		},
		{
			name:     "Write Over Deny",
			a:        PolicyWrite,
			b:        PolicyDeny,
			expected: false,
		},
		{
			name:     "List Over Deny",
			a:        PolicyList,
			b:        PolicyDeny,
			expected: false,
		},
		{
			name:     "Read Over Deny",
			a:        PolicyRead,
			b:        PolicyDeny,
			expected: false,
		},
		{
			name:     "Deny Over Unknown",
			a:        PolicyDeny,
			b:        "not a policy",
			expected: true,
		},
		{
			name:     "List Over Write",
			a:        PolicyList,
			b:        PolicyWrite,
			expected: false,
		},
		{
			name:     "Read Over Write",
			a:        PolicyRead,
			b:        PolicyWrite,
			expected: false,
		},
		{
			name:     "Unknown Over Write",
			a:        "not a policy",
			b:        PolicyWrite,
			expected: false,
		},
		{
			name:     "Read Over List",
			a:        PolicyRead,
			b:        PolicyList,
			expected: false,
		},
		{
			name:     "Unknown Over List",
			a:        "not a policy",
			b:        PolicyList,
			expected: false,
		},
		{
			name:     "Unknown Over Read",
			a:        "not a policy",
			b:        PolicyRead,
			expected: false,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			require.Equal(t, tcase.expected, takesPrecedenceOver(tcase.a, tcase.b))
		})
	}
}
