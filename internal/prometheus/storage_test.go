package prometheus_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth/internal/log"
	"github.com/slok/sloth/internal/prometheus"
)

func TestIOWriterGroupedRulesYAMLRepoStore(t *testing.T) {
	tests := map[string]struct {
		slos    []prometheus.StorageSLO
		expYAML string
		expErr  bool
	}{
		"Having 0 SLO rules should fail.": {
			slos:   []prometheus.StorageSLO{},
			expErr: true,
		},

		"Having 0 SLO rules generated should fail.": {
			slos: []prometheus.StorageSLO{
				{},
			},
			expErr: true,
		},

		"Having a single SLI recording rule should render correctly.": {
			slos: []prometheus.StorageSLO{
				{
					SLO: prometheus.SLO{ID: "test1"},
					Rules: prometheus.SLORules{
						SLIErrorRecRules: []rulefmt.Rule{
							{
								Record: "test:record",
								Expr:   "test-expr",
								Labels: map[string]string{"test-label": "one"},
							},
						},
					},
				},
			},
			expYAML: `
---
# Code generated by Sloth (dev): https://github.com/slok/sloth.
# DO NOT EDIT.

groups:
- name: sloth-slo-sli-recordings-test1
  rules:
  - record: test:record
    expr: test-expr
    labels:
      test-label: one
`,
		},
		"Having a single metadata recording rule should render correctly.": {
			slos: []prometheus.StorageSLO{
				{
					SLO: prometheus.SLO{ID: "test1"},
					Rules: prometheus.SLORules{
						MetadataRecRules: []rulefmt.Rule{
							{
								Record: "test:record",
								Expr:   "test-expr",
								Labels: map[string]string{"test-label": "one"},
							},
						},
					},
				},
			},
			expYAML: `
---
# Code generated by Sloth (dev): https://github.com/slok/sloth.
# DO NOT EDIT.

groups:
- name: sloth-slo-meta-recordings-test1
  rules:
  - record: test:record
    expr: test-expr
    labels:
      test-label: one
`,
		},
		"Having a single SLO alert rule should render correctly.": {
			slos: []prometheus.StorageSLO{
				{
					SLO: prometheus.SLO{ID: "test1"},
					Rules: prometheus.SLORules{
						AlertRules: []rulefmt.Rule{
							{
								Alert:       "testAlert",
								Expr:        "test-expr",
								Labels:      map[string]string{"test-label": "one"},
								Annotations: map[string]string{"test-annot": "one"},
							},
						},
					},
				},
			},
			expYAML: `
---
# Code generated by Sloth (dev): https://github.com/slok/sloth.
# DO NOT EDIT.

groups:
- name: sloth-slo-alerts-test1
  rules:
  - alert: testAlert
    expr: test-expr
    labels:
      test-label: one
    annotations:
      test-annot: one
`,
		},

		"Having a multiple SLO alert and recording rules should render correctly.": {
			slos: []prometheus.StorageSLO{
				{
					SLO: prometheus.SLO{ID: "testa"},
					Rules: prometheus.SLORules{
						SLIErrorRecRules: []rulefmt.Rule{
							{
								Record: "test:record-a1",
								Expr:   "test-expr-a1",
								Labels: map[string]string{"test-label": "a-1"},
							},
							{
								Record: "test:record-a2",
								Expr:   "test-expr-a2",
								Labels: map[string]string{"test-label": "a-2"},
							},
						},
						MetadataRecRules: []rulefmt.Rule{
							{
								Record: "test:record-a3",
								Expr:   "test-expr-a3",
								Labels: map[string]string{"test-label": "a-3"},
							},
							{
								Record: "test:record-a4",
								Expr:   "test-expr-a4",
								Labels: map[string]string{"test-label": "a-4"},
							},
						},
						AlertRules: []rulefmt.Rule{
							{
								Alert:       "testAlertA1",
								Expr:        "test-expr-a1",
								Labels:      map[string]string{"test-label": "a-1"},
								Annotations: map[string]string{"test-annot": "a-1"},
							},
							{
								Alert:       "testAlertA2",
								Expr:        "test-expr-a2",
								Labels:      map[string]string{"test-label": "a-2"},
								Annotations: map[string]string{"test-annot": "a-2"},
							},
						},
					},
				},
				{
					SLO: prometheus.SLO{ID: "testb"},
					Rules: prometheus.SLORules{
						SLIErrorRecRules: []rulefmt.Rule{
							{
								Record: "test:record-b1",
								Expr:   "test-expr-b1",
								Labels: map[string]string{"test-label": "b-1"},
							},
						},
						MetadataRecRules: []rulefmt.Rule{
							{
								Record: "test:record-b2",
								Expr:   "test-expr-b2",
								Labels: map[string]string{"test-label": "b-2"},
							},
						},
						AlertRules: []rulefmt.Rule{
							{
								Alert:       "testAlertB1",
								Expr:        "test-expr-b1",
								Labels:      map[string]string{"test-label": "b-1"},
								Annotations: map[string]string{"test-annot": "b-1"},
							},
						},
					},
				},
			},
			expYAML: `
---
# Code generated by Sloth (dev): https://github.com/slok/sloth.
# DO NOT EDIT.

groups:
- name: sloth-slo-sli-recordings-testa
  rules:
  - record: test:record-a1
    expr: test-expr-a1
    labels:
      test-label: a-1
  - record: test:record-a2
    expr: test-expr-a2
    labels:
      test-label: a-2
- name: sloth-slo-meta-recordings-testa
  rules:
  - record: test:record-a3
    expr: test-expr-a3
    labels:
      test-label: a-3
  - record: test:record-a4
    expr: test-expr-a4
    labels:
      test-label: a-4
- name: sloth-slo-alerts-testa
  rules:
  - alert: testAlertA1
    expr: test-expr-a1
    labels:
      test-label: a-1
    annotations:
      test-annot: a-1
  - alert: testAlertA2
    expr: test-expr-a2
    labels:
      test-label: a-2
    annotations:
      test-annot: a-2
- name: sloth-slo-sli-recordings-testb
  rules:
  - record: test:record-b1
    expr: test-expr-b1
    labels:
      test-label: b-1
- name: sloth-slo-meta-recordings-testb
  rules:
  - record: test:record-b2
    expr: test-expr-b2
    labels:
      test-label: b-2
- name: sloth-slo-alerts-testb
  rules:
  - alert: testAlertB1
    expr: test-expr-b1
    labels:
      test-label: b-1
    annotations:
      test-annot: b-1
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var gotYAML bytes.Buffer
			repo := prometheus.NewIOWriterGroupedRulesYAMLRepo(&gotYAML, log.Noop)
			err := repo.StoreSLOs(context.TODO(), test.slos)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expYAML, gotYAML.String())
			}
		})
	}
}
