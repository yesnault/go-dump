package test

import (
	"bytes"
	"testing"

	"github.com/ovh/cds/engine/api/test"
	"github.com/ovh/cds/sdk"
	"github.com/stretchr/testify/assert"

	"encoding/json"

	"github.com/fsamin/go-dump"
)

func TestDumpStruct(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
__Type__: T
`
	assert.Equal(t, expected, out.String())

}

type T struct {
	A int
	B string
	C Tbis
}

type Tbis struct {
	Cbis string
	Cter string
}

func TestDumpStruct_Nested(t *testing.T) {

	a := T{23, "foo bar", Tbis{"lol", "lol"}}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
T.C.Cbis: lol
T.C.Cter: lol
T.C.__Type__: Tbis
__Type__: T
`
	assert.Equal(t, expected, out.String())

}

type TP struct {
	A *int
	B string
	C *Tbis
}

func TestDumpStruct_NestedWithPointer(t *testing.T) {
	i := 23
	a := TP{&i, "foo bar", &Tbis{"lol", "lol"}}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `TP.A: 23
TP.B: foo bar
TP.C.Cbis: lol
TP.C.Cter: lol
TP.C.__Type__: Tbis
__Type__: TP
`
	assert.Equal(t, expected, out.String())

}

type TM struct {
	A int
	B string
	C map[string]Tbis
}

func TestDumpStruct_Map(t *testing.T) {

	a := TM{A: 23, B: "foo bar"}
	a.C = map[string]Tbis{}
	a.C["bar"] = Tbis{"lel", "lel"}
	a.C["foo"] = Tbis{"lol", "lol"}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `TM.A: 23
TM.B: foo bar
TM.C.__Len__: 2
TM.C.__Type__: Map
TM.C.bar.Tbis.Cbis: lel
TM.C.bar.Tbis.Cter: lel
TM.C.bar.Tbis.__Type__: Tbis
TM.C.foo.Tbis.Cbis: lol
TM.C.foo.Tbis.Cter: lol
TM.C.foo.Tbis.__Type__: Tbis
__Type__: TM
`
	assert.Equal(t, expected, out.String())

}

func TestDumpArray(t *testing.T) {
	a := []T{
		{23, "foo bar", Tbis{"lol", "lol"}},
		{24, "fee bor", Tbis{"lel", "lel"}},
	}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `0.A: 23
0.B: foo bar
0.C.Cbis: lol
0.C.Cter: lol
0.C.__Type__: Tbis
0.__Type__: T
1.A: 24
1.B: fee bor
1.C.Cbis: lel
1.C.Cter: lel
1.C.__Type__: Tbis
1.__Type__: T
__Len__: 2
__Type__: Array
`
	assert.Equal(t, expected, out.String())
}

type TS struct {
	A int
	B string
	C []T
	D []bool
}

func TestDumpStruct_Array(t *testing.T) {
	a := TS{
		A: 0,
		B: "here",
		C: []T{
			{23, "foo bar", Tbis{"lol", "lol"}},
			{24, "fee bor", Tbis{"lel", "lel"}},
		},
		D: []bool{true, false},
	}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)
	expected := `TS.A: 0
TS.B: here
TS.C.C0.A: 23
TS.C.C0.B: foo bar
TS.C.C0.C.Cbis: lol
TS.C.C0.C.Cter: lol
TS.C.C0.C.__Type__: Tbis
TS.C.C0.__Type__: T
TS.C.C1.A: 24
TS.C.C1.B: fee bor
TS.C.C1.C.Cbis: lel
TS.C.C1.C.Cter: lel
TS.C.C1.C.__Type__: Tbis
TS.C.C1.__Type__: T
TS.C.__Len__: 2
TS.C.__Type__: Array
TS.D.D0: true
TS.D.D1: false
TS.D.__Len__: 2
TS.D.__Type__: Array
__Type__: TS
`
	assert.Equal(t, expected, out.String())
}

func TestToMap(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	m, err := dump.ToMap(a)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(m))
	var m1Found, m2Found bool
	for k, v := range m {
		t.Logf("%s: %v (%T)", k, v, v)
		if k == "T.A" {
			m1Found = true
			assert.Equal(t, 23, v)
		}
		if k == "T.B" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "T.A not found in map")
	assert.True(t, m2Found, "T.B not found in map")
}

func TestToMapWithFormatter(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	m, err := dump.ToMap(a, dump.WithDefaultLowerCaseFormatter())
	t.Log(m)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(m))
	var m1Found, m2Found bool
	for k, v := range m {
		if k == "t.a" {
			m1Found = true
			assert.Equal(t, 23, v)
		}
		if k == "t.b" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "t.a not found in map")
	assert.True(t, m2Found, "t.b not found in map")
}

func TestComplex(t *testing.T) {
	p := sdk.Pipeline{
		Name: "MyPipeline",
		Type: sdk.BuildPipeline,
		Stages: []sdk.Stage{
			{
				BuildOrder: 1,
				Name:       "stage 1",
				Enabled:    true,
				Jobs: []sdk.Job{
					{
						Action: sdk.Action{
							Name:        "Job 1",
							Description: "This is job 1",
							Actions: []sdk.Action{
								{

									Type: sdk.BuiltinAction,
									Name: sdk.ScriptAction,
									Parameters: []sdk.Parameter{
										{
											Name:  "script",
											Type:  sdk.TextParameter,
											Value: "echo lol",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Log(dump.MustSdump(p))

	out := &bytes.Buffer{}
	err := dump.Fdump(out, p)
	assert.NoError(t, err)
	expected := `Pipeline.GroupPermission.__Len__: 0
Pipeline.GroupPermission.__Type__: Array
Pipeline.ID: 0
Pipeline.LastModified: 0
Pipeline.LastPipelineBuild:
Pipeline.Name: MyPipeline
Pipeline.Parameter.__Len__: 0
Pipeline.Parameter.__Type__: Array
Pipeline.Permission: 0
Pipeline.ProjectID: 0
Pipeline.ProjectKey:
Pipeline.Stages.Stages0.BuildOrder: 1
Pipeline.Stages.Stages0.Enabled: true
Pipeline.Stages.Stages0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Actions.__Len__: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Actions.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.AlwaysExecuted: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Deprecated: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Description:
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Name: Script
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Optional: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Description:
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Name: script
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Type: text
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Value: echo lol
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.__Type__: Parameter
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.__Len__: 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Requirements.__Len__: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Requirements.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Type: Builtin
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.__Type__: Action
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.__Len__: 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.AlwaysExecuted: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Deprecated: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Description: This is job 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Name: Job 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Optional: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Parameters.__Len__: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Parameters.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Requirements.__Len__: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Requirements.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Type:
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.__Type__: Action
Pipeline.Stages.Stages0.Jobs.Jobs0.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.PipelineActionID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.PipelineStageID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Warnings.__Len__: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Warnings.__Type__: Array
Pipeline.Stages.Stages0.Jobs.Jobs0.__Type__: Job
Pipeline.Stages.Stages0.Jobs.__Len__: 1
Pipeline.Stages.Stages0.Jobs.__Type__: Array
Pipeline.Stages.Stages0.LastModified: 0
Pipeline.Stages.Stages0.Name: stage 1
Pipeline.Stages.Stages0.PipelineBuildJobs.__Len__: 0
Pipeline.Stages.Stages0.PipelineBuildJobs.__Type__: Array
Pipeline.Stages.Stages0.PipelineID: 0
Pipeline.Stages.Stages0.Prerequisites.__Len__: 0
Pipeline.Stages.Stages0.Prerequisites.__Type__: Array
Pipeline.Stages.Stages0.RunJobs.__Len__: 0
Pipeline.Stages.Stages0.RunJobs.__Type__: Array
Pipeline.Stages.Stages0.Status:
Pipeline.Stages.Stages0.Warnings.__Len__: 0
Pipeline.Stages.Stages0.Warnings.__Type__: Array
Pipeline.Stages.Stages0.__Type__: Stage
Pipeline.Stages.__Len__: 1
Pipeline.Stages.__Type__: Array
Pipeline.Type: build
Pipeline.Usage:
__Type__: Pipeline
`

	assert.Equal(t, expected, out.String())
	assert.NoError(t, err)
}

func TestMapStringInterface(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["id"] = "ID"
	myMap["name"] = "foo"
	myMap["value"] = "bar"

	result, err := dump.ToStringMap(myMap)
	t.Log(dump.Sdump(myMap))
	assert.NoError(t, err)
	assert.Equal(t, 5, len(result))

	expected := `__len__: 3
__type__: Map
id: ID
name: foo
value: bar
`
	out := &bytes.Buffer{}
	err = dump.Fdump(out, myMap, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestFromJSON(t *testing.T) {
	js := []byte(`{
    "blabla": "lol log", 
    "boubou": {
        "yo": 1
    } 
}`)

	var i interface{}
	test.NoError(t, json.Unmarshal(js, &i))

	result, err := dump.ToStringMap(i)
	t.Log(dump.Sdump(i))
	t.Log(result)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(result))
	assert.Equal(t, "lol log", result["blabla"])
	assert.Equal(t, "1", result["boubou.yo"])
}
func TestComplexWithFormatter(t *testing.T) {
	p := sdk.Pipeline{
		Name: "MyPipeline",
		Type: sdk.BuildPipeline,
		Stages: []sdk.Stage{
			{
				BuildOrder: 1,
				Name:       "stage 1",
				Enabled:    true,
				Jobs: []sdk.Job{
					{
						Action: sdk.Action{
							Name:        "Job 1",
							Description: "This is job 1",
							Actions: []sdk.Action{
								{

									Type: sdk.BuiltinAction,
									Name: sdk.ScriptAction,
									Parameters: []sdk.Parameter{
										{
											Name:  "script",
											Type:  sdk.TextParameter,
											Value: "echo lol",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Log(dump.MustSdump(p))

	out := &bytes.Buffer{}
	err := dump.Fdump(out, p, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	expected := `__type__: Pipeline
pipeline.grouppermission.__len__: 0
pipeline.grouppermission.__type__: Array
pipeline.id: 0
pipeline.lastmodified: 0
pipeline.lastpipelinebuild:
pipeline.name: MyPipeline
pipeline.parameter.__len__: 0
pipeline.parameter.__type__: Array
pipeline.permission: 0
pipeline.projectid: 0
pipeline.projectkey:
pipeline.stages.__len__: 1
pipeline.stages.__type__: Array
pipeline.stages.stages0.__type__: Stage
pipeline.stages.stages0.buildorder: 1
pipeline.stages.stages0.enabled: true
pipeline.stages.stages0.id: 0
pipeline.stages.stages0.jobs.__len__: 1
pipeline.stages.stages0.jobs.__type__: Array
pipeline.stages.stages0.jobs.jobs0.__type__: Job
pipeline.stages.stages0.jobs.jobs0.action.__type__: Action
pipeline.stages.stages0.jobs.jobs0.action.actions.__len__: 1
pipeline.stages.stages0.jobs.jobs0.action.actions.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.__type__: Action
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.actions.__len__: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.actions.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.alwaysexecuted: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.deprecated: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.description:
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.enabled: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.id: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.name: Script
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.optional: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.__len__: 1
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.__type__: Parameter
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.description:
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.id: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.name: script
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.type: text
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.value: echo lol
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.requirements.__len__: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.requirements.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.type: Builtin
pipeline.stages.stages0.jobs.jobs0.action.alwaysexecuted: false
pipeline.stages.stages0.jobs.jobs0.action.deprecated: false
pipeline.stages.stages0.jobs.jobs0.action.description: This is job 1
pipeline.stages.stages0.jobs.jobs0.action.enabled: false
pipeline.stages.stages0.jobs.jobs0.action.id: 0
pipeline.stages.stages0.jobs.jobs0.action.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.action.name: Job 1
pipeline.stages.stages0.jobs.jobs0.action.optional: false
pipeline.stages.stages0.jobs.jobs0.action.parameters.__len__: 0
pipeline.stages.stages0.jobs.jobs0.action.parameters.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.requirements.__len__: 0
pipeline.stages.stages0.jobs.jobs0.action.requirements.__type__: Array
pipeline.stages.stages0.jobs.jobs0.action.type:
pipeline.stages.stages0.jobs.jobs0.enabled: false
pipeline.stages.stages0.jobs.jobs0.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.pipelineactionid: 0
pipeline.stages.stages0.jobs.jobs0.pipelinestageid: 0
pipeline.stages.stages0.jobs.jobs0.warnings.__len__: 0
pipeline.stages.stages0.jobs.jobs0.warnings.__type__: Array
pipeline.stages.stages0.lastmodified: 0
pipeline.stages.stages0.name: stage 1
pipeline.stages.stages0.pipelinebuildjobs.__len__: 0
pipeline.stages.stages0.pipelinebuildjobs.__type__: Array
pipeline.stages.stages0.pipelineid: 0
pipeline.stages.stages0.prerequisites.__len__: 0
pipeline.stages.stages0.prerequisites.__type__: Array
pipeline.stages.stages0.runjobs.__len__: 0
pipeline.stages.stages0.runjobs.__type__: Array
pipeline.stages.stages0.status:
pipeline.stages.stages0.warnings.__len__: 0
pipeline.stages.stages0.warnings.__type__: Array
pipeline.type: build
pipeline.usage:
`
	assert.Equal(t, expected, out.String())
	assert.NoError(t, err)
}

type Result struct {
	Body     string      `json:"body,omitempty" yaml:"body,omitempty"`
	BodyJSON interface{} `json:"bodyjson,omitempty" yaml:"bodyjson,omitempty"`
}

func TestMapStringInterfaceInStruct(t *testing.T) {

	r := Result{}
	r.Body = "foo"
	r.BodyJSON = map[string]interface{}{
		"cardID": "1234",
		"items":  []string{"foo", "beez"},
		"test": Result{
			Body: "12",
			BodyJSON: map[string]interface{}{
				"card": "@",
				"yolo": 3,
				"beez": true,
			},
		},
		"description": "yolo",
	}

	expected := `__type__: Result
result.body: foo
result.bodyjson.__len__: 4
result.bodyjson.__type__: Map
result.bodyjson.cardid: 1234
result.bodyjson.description: yolo
result.bodyjson.items.__len__: 2
result.bodyjson.items.__type__: Array
result.bodyjson.items.items0: foo
result.bodyjson.items.items1: beez
result.bodyjson.test.result.__type__: Result
result.bodyjson.test.result.body: 12
result.bodyjson.test.result.bodyjson.__len__: 3
result.bodyjson.test.result.bodyjson.__type__: Map
result.bodyjson.test.result.bodyjson.beez: true
result.bodyjson.test.result.bodyjson.card: @
result.bodyjson.test.result.bodyjson.yolo: 3
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, r, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestWeird(t *testing.T) {
	testJSON := `{
	"beez": null,
	"foo" : "bar",
	"bou" : [null, "hello"]
  }`

	var test interface{}
	json.Unmarshal([]byte(testJSON), &test)
	expected := `__len__: 3
__type__: Map
beez:
bou.__len__: 2
bou.__type__: Array
bou.bou0:
bou.bou1: hello
foo: bar
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, test, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())

}

type ResultUnexported struct {
	body *string
	Foo  string
}

func TestUnexportedField(t *testing.T) {

	test := ResultUnexported{
		body: nil,
		Foo:  "bar",
	}

	expected := `__type__: ResultUnexported
resultunexported.foo: bar
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, test, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestWithDetailedStruct(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	enc := dump.NewDefaultEncoder(new(bytes.Buffer))
	enc.ExtraFields.DetailedStruct = true
	enc.ExtraFields.Type = false
	res, _ := enc.Sdump(a)
	t.Log(res)
	assert.Equal(t, `T: {23 foo bar}
T.A: 23
T.B: foo bar
T.__Len__: 2
`, res)
}
