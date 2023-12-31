// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hl7

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/encoding/unicode"
)

func TestPrimitive(t *testing.T) {
	TimezoneAndLocation("UTC")
	c := &Context{
		Decoder:     unicode.UTF8.NewDecoder(),
		Delimiters:  DefaultDelimiters,
		Nesting:     0,
		TimezoneLoc: Location,
	}

	cases := []struct {
		name string
		p    Primitive
		want Primitive
		// got is an empty placeholder of the given type,
		// where the value will be unmarshalled to.
		got Primitive
	}{{
		name: "ST",
		p:    NewST("value"),
		want: NewST("value"),
		got:  NewST(""),
	}, {
		name: "ID",
		p:    NewID("value"),
		want: NewID("value"),
		got:  NewID(""),
	}, {
		name: "SI",
		p:    NewSI(44),
		want: NewSI(44),
		got:  &SI{Valid: false},
	}, {
		name: "NM",
		p:    NewNM(44),
		want: NewNM(44),
		got:  &NM{Valid: false},
	}, {
		name: "IS",
		p:    NewIS("value"),
		want: NewIS("value"),
		got:  NewIS(""),
	}, {
		name: "DT",
		p:    NewDT("value"),
		want: NewDT("value"),
		got:  NewDT(""),
	}, {
		name: "TM",
		p:    NewTM("value"),
		want: NewTM("value"),
		got:  NewTM(""),
	}, {
		name: "TS_YearPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: YearPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Precision: YearPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_MonthPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: MonthPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 01, 0, 0, 0, 0, time.UTC), Precision: MonthPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_DayPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: DayPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 0, 0, 0, 0, time.UTC), Precision: DayPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_HourPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: HourPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 0, 0, 0, time.UTC), Precision: HourPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_MinutePrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: MinutePrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 0, 0, time.UTC), Precision: MinutePrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_SecondPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: SecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: SecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_SecondPrecision_WithNanoseconds",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: SecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 0, time.UTC), Precision: SecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_TenthSecondPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: TenthSecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 100000000, time.UTC), Precision: TenthSecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_HundredthSecondPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: HundredthSecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 120000000, time.UTC), Precision: HundredthSecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_ThousandthSecondPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: ThousandthSecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 123000000, time.UTC), Precision: ThousandthSecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "TS_TenThousandthSecondPrecision",
		p:    &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: TenThousandthSecondPrecision},
		want: &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 123400000, time.UTC), Precision: TenThousandthSecondPrecision},
		got:  &TS{IsHL7Null: true},
	}, {
		name: "DTM_YearPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: YearPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC), Precision: YearPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_MonthPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: MonthPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 01, 0, 0, 0, 0, time.UTC), Precision: MonthPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_DayPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: DayPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 0, 0, 0, 0, time.UTC), Precision: DayPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_HourPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: HourPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 0, 0, 0, time.UTC), Precision: HourPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_MinutePrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: MinutePrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 0, 0, time.UTC), Precision: MinutePrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_SecondPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: SecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, time.UTC), Precision: SecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_SecondPrecision_WithNanoseconds",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: SecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 0, time.UTC), Precision: SecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_TenthSecondPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: TenthSecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 100000000, time.UTC), Precision: TenthSecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_HundredthSecondPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: HundredthSecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 120000000, time.UTC), Precision: HundredthSecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_ThousandthSecondPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: ThousandthSecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 123000000, time.UTC), Precision: ThousandthSecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "DTM_TenThousandthSecondPrecision",
		p:    &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, time.UTC), Precision: TenThousandthSecondPrecision},
		want: &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 35, 123400000, time.UTC), Precision: TenThousandthSecondPrecision},
		got:  &DTM{IsHL7Null: true},
	}, {
		name: "TN",
		p:    NewTN("value"),
		want: NewTN("value"),
		got:  NewTN(""),
	}, {
		name: "FT",
		p:    NewFT("value"),
		want: NewFT("value"),
		got:  NewFT(""),
	}, {
		name: "TX",
		p:    NewTX("value"),
		want: NewTX("value"),
		got:  NewTX(""),
	}, {
		name: "SNM",
		p:    &SNM{Value: "1234567"},
		want: &SNM{Value: "1234567", Valid: true},
		got:  &SNM{Valid: false},
	}, {
		name: "CM",
		p:    NewCM([]byte("value")),
		want: NewCM([]byte("value")),
		got:  NewCM([]byte{}),
	}, {
		name: "Any",
		p:    NewAny([]byte("value")),
		want: NewAny([]byte("value")),
		got:  NewAny([]byte{}),
	}, {
		name: "NUL",
		p:    NewNUL("value"),
		want: NewNUL("value"),
		got:  NewNUL(""),
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.p.Marshal(c)
			if err != nil {
				t.Fatalf("[%+v].Marshal(%v) failed with %v", tc.p, c, err)
			}

			if err := tc.got.Unmarshal(b, c); err != nil {
				t.Fatalf("[%+v].Unmarshal(%s, %v) failed with %v", tc.got, string(b), c, err)
			}
			if diff := cmp.Diff(tc.want, tc.got); diff != "" {
				t.Errorf("[%+v].Unmarshal(%s, %v) got diff (-want, +got):\n%s", tc.got, string(b), c, diff)
			}
		})
	}
}

func TestSanitizedString(t *testing.T) {
	type sanitizable interface {
		SanitizedString() string
	}

	cases := []struct {
		name string
		s    sanitizable
		want string
	}{
		{name: "ST", s: NewST("value"), want: "value"},
		{name: "ST null", s: NewST(`""`), want: ""},
		{name: "ID", s: NewID("value"), want: "value"},
		{name: "ID null", s: NewID(`""`), want: ""},
		{name: "IS", s: NewIS("value"), want: "value"},
		{name: "IS null", s: NewIS(`""`), want: ""},
		{
			name: "HD",
			s: &HD{
				NamespaceID:     NewIS("NamespaceID"),
				UniversalID:     NewST("UniversalID"),
				UniversalIDType: NewID("UniversalIDType"),
			},
			want: "NamespaceID^UniversalID^UniversalIDType",
		}, {
			name: "HD every field null",
			s: &HD{
				NamespaceID:     NewIS(`""`),
				UniversalID:     NewST(`""`),
				UniversalIDType: NewID(`""`),
			},
			want: `""^""^""`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.s.SanitizedString(); got != tc.want {
				t.Errorf("[%+v].SanitizedString()=%q, want %q", tc.s, got, tc.want)
			}
		})
	}
}

type empty interface {
	Empty() bool
}

func TestEmpty(t *testing.T) {
	cases := []struct {
		name string
		f    func(s string) empty
	}{
		{name: "ST", f: func(s string) empty { return NewST(ST(s)) }},
		{name: "ID", f: func(s string) empty { return NewID(ID(s)) }},
		{name: "IS", f: func(s string) empty { return NewIS(IS(s)) }},
		{name: "SNM", f: func(s string) empty { return &SNM{Value: s, Valid: true} }},
	}

	for _, tc := range cases {
		for k, want := range map[empty]bool{tc.f("value"): false, tc.f(""): true} {
			t.Run(fmt.Sprintf("%s-%s", tc.name, k), func(t *testing.T) {
				if got := k.Empty(); got != want {
					t.Errorf("[%+v].Empty()=%t, want %t", k, got, want)
				}
			})
		}
	}
}

func TestEmptyNilEmpty(t *testing.T) {
	cases := []struct {
		name string
		e    empty
	}{
		{name: "ST", e: new(ST)},
		{name: "ID", e: new(ID)},
		{name: "IS", e: new(IS)},
		{name: "SNM", e: new(SNM)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.e.Empty(); !got {
				t.Errorf("[%+v].Empty()=%t, want true", tc.e, got)
			}
		})
	}
}

func TestParseTS(t *testing.T) {
	tests := []struct {
		in           string
		outTime      string
		outPrecision TSPrecision
	}{
		{"20141128001635", "2014-11-28T00:16:35Z", SecondPrecision},
		// The precision specified in the second component ("^M" = Minute) overrides the precision
		// implicitly specified in the first component (Second).
		{"20141128001635^M", "2014-11-28T00:16:35Z", MinutePrecision},
		{"20141128001635.1", "2014-11-28T00:16:35.1Z", TenthSecondPrecision},
		{"20141128001635.12", "2014-11-28T00:16:35.12Z", HundredthSecondPrecision},
		{"20141128001635.123", "2014-11-28T00:16:35.123Z", ThousandthSecondPrecision},
		{"20141128001635.1234", "2014-11-28T00:16:35.1234Z", TenThousandthSecondPrecision},
		// The following examples are from the HL7 specification
		{"19760704010159-0600", "1976-07-04T01:01:59-06:00", SecondPrecision},
		{"19760704010159-0500", "1976-07-04T01:01:59-05:00", SecondPrecision},
		{"198807050000", "1988-07-04T23:00:00Z", MinutePrecision},
		{"19880705", "1988-07-05T00:00:00Z", DayPrecision},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var ts TS
			err := ts.Unmarshal([]byte(tt.in), testContext)
			if err != nil {
				t.Errorf("ParseTS(%q) got error %v, want err=<nil>", tt.in, err)
			}
			want, _ := time.Parse(time.RFC3339Nano, tt.outTime)
			if !want.Equal(ts.Time) {
				t.Errorf("ParseTS(%q).Time got %v, want %v", tt.in, ts.Time, want)
			}
			if diff := cmp.Diff(tt.outPrecision, ts.Precision); diff != "" {
				t.Errorf("ParseTS(%q).Precision mismatch (-want, +got)=\n%s", tt.in, diff)
			}
		})
	}
}

func TestParseDTM(t *testing.T) {
	tests := []struct {
		in           string
		outTime      string
		outPrecision TSPrecision
	}{
		{"20141128001635", "2014-11-28T00:16:35Z", SecondPrecision},
		{"20141128001635.1", "2014-11-28T00:16:35.1Z", TenthSecondPrecision},
		{"20141128001635.12", "2014-11-28T00:16:35.12Z", HundredthSecondPrecision},
		{"20141128001635.123", "2014-11-28T00:16:35.123Z", ThousandthSecondPrecision},
		{"20141128001635.1234", "2014-11-28T00:16:35.1234Z", TenThousandthSecondPrecision},
		// The following examples are from the HL7 specification
		{"19760704010159-0600", "1976-07-04T01:01:59-06:00", SecondPrecision},
		{"19760704010159-0500", "1976-07-04T01:01:59-05:00", SecondPrecision},
		{"198807050000", "1988-07-04T23:00:00Z", MinutePrecision},
		{"19880705", "1988-07-05T00:00:00Z", DayPrecision},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var dtm DTM
			err := dtm.Unmarshal([]byte(tt.in), testContext)
			if err != nil {
				t.Errorf("ParseDTM(%q) got error %v, want err=<nil>", tt.in, err)
			}
			want, _ := time.Parse(time.RFC3339Nano, tt.outTime)
			if !want.Equal(dtm.Time) {
				t.Errorf("ParseDTM(%q).Time got %v, want %v", tt.in, dtm.Time, want)
			}
			if diff := cmp.Diff(tt.outPrecision, dtm.Precision); diff != "" {
				t.Errorf("ParseDTM(%q).Precision mismatch (-want, +got)=\n%s", tt.in, diff)
			}
		})
	}
}

func TestMarshalTS(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   *TS
		want string
	}{{
		name: "TS_YearPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: YearPrecision},
		want: "2020",
	}, {
		name: "TS_MonthPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: MonthPrecision},
		want: "202002",
	}, {
		name: "TS_DayPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: DayPrecision},
		want: "20200224",
	}, {
		name: "TS_HourPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: HourPrecision},
		want: "2020022412",
	}, {
		name: "TS_MinutePrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: MinutePrecision},
		want: "202002241255",
	}, {
		name: "TS_SecondPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: SecondPrecision},
		want: "20200224125530",
	}, {
		name: "TS_SecondPrecision_WithNanoseconds",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: SecondPrecision},
		want: "20200224125535",
	}, {
		name: "TS_TenthSecondPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: TenthSecondPrecision},
		want: "20200224125535.1",
	}, {
		name: "TS_HundredthSecondPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: HundredthSecondPrecision},
		want: "20200224125535.12",
	}, {
		name: "TS_ThousandthSecondPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: ThousandthSecondPrecision},
		want: "20200224125535.123",
	}, {
		name: "TS_TenThousandthSecondPrecision",
		in:   &TS{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: TenThousandthSecondPrecision},
		want: "20200224125535.1234",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Marshal(testContext)
			if err != nil {
				t.Fatalf("[%+v].Marshal(%+v) failed with %+v", tt.in, testContext, err)
			}
			if got, want := string(got), tt.want; got != want {
				t.Errorf("[%+v].Marshal(%+v)=%q, want %q", tt.in, testContext, got, want)
			}
		})
	}
}

func TestMarshalDTM(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   *DTM
		want string
	}{{
		name: "YearPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: YearPrecision},
		want: "2020",
	}, {
		name: "MonthPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: MonthPrecision},
		want: "202002",
	}, {
		name: "DayPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: DayPrecision},
		want: "20200224",
	}, {
		name: "HourPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: HourPrecision},
		want: "2020022412",
	}, {
		name: "MinutePrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: MinutePrecision},
		want: "202002241255",
	}, {
		name: "SecondPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 0, testLocation), Precision: SecondPrecision},
		want: "20200224125530",
	}, {
		name: "SecondPrecision_WithNanoseconds",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: SecondPrecision},
		want: "20200224125535",
	}, {
		name: "TenthSecondPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: TenthSecondPrecision},
		want: "20200224125535.1",
	}, {
		name: "HundredthSecondPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: HundredthSecondPrecision},
		want: "20200224125535.12",
	}, {
		name: "ThousandthSecondPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: ThousandthSecondPrecision},
		want: "20200224125535.123",
	}, {
		name: "TenThousandthSecondPrecision",
		in:   &DTM{IsHL7Null: false, Time: time.Date(2020, 02, 24, 12, 55, 30, 5123456789, testLocation), Precision: TenThousandthSecondPrecision},
		want: "20200224125535.1234",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Marshal(testContext)
			if err != nil {
				t.Fatalf("[%+v].Marshal(%+v) failed with %+v", tt.in, testContext, err)
			}
			if got, want := string(got), tt.want; got != want {
				t.Errorf("[%+v].Marshal(%+v)=%q, want %q", tt.in, testContext, got, want)
			}
		})
	}
}

func TestParseTS_ClearField(t *testing.T) {
	var ts TS
	if err := ts.Unmarshal([]byte("\"\""), testContext); err != nil {
		t.Fatalf("Unmarshal('') failed with %v", err)
	}
	if !ts.Time.IsZero() {
		t.Error("ts.Time.IsZero() is false, want true")
	}
}

func TestParseDTM_ClearField(t *testing.T) {
	var dtm DTM
	if err := dtm.Unmarshal([]byte("\"\""), testContext); err != nil {
		t.Fatalf("Unmarshal('') failed with %v", err)
	}
	if !dtm.Time.IsZero() {
		t.Error("dtm.Time.IsZero() is false, want true")
	}
}

func TestParseTS_Error(t *testing.T) {
	tests := []string{
		// The empty string
		"",
		// A two digit year
		"20",
		// Fractions of a second with seconds
		"201411280016.12",
		// More precision than thousandths of a second
		"20141128001635.12345",
		// A unknown (legacy) precision value
		"20141128001635^T",
		// A timezone without the correct number of digits
		"201411280016+010",
	}
	for _, tt := range tests {
		var ts TS
		if err := ts.Unmarshal([]byte(tt), testContext); err == nil {
			t.Errorf("Unmarshal(%q) got err=<nil>, want error", tt)
		}
	}
}

func TestParseDTM_Error(t *testing.T) {
	tests := []string{
		// The empty string
		"",
		// A two digit year
		"20",
		// Fractions of a second with seconds
		"201411280016.12",
		// More precision than thousandths of a second
		"20141128001635.12345",
		// A unknown (legacy) precision value
		"20141128001635^T",
		// A timezone without the correct number of digits
		"201411280016+010",
		// Precision specified in a second component. This is supported in TS values, but not in DTM.
		"20141128001635^M",
	}
	for _, tt := range tests {
		var dtm DTM
		if err := dtm.Unmarshal([]byte(tt), testContext); err == nil {
			t.Errorf("Unmarshal(%s) got err=<nil>, want error", tt)
		}
	}
}

func TestParseSI(t *testing.T) {
	tests := []struct {
		in  string
		out SI
	}{
		{"0", SI{Value: 0, Valid: true}},
		{"1", SI{Value: 1, Valid: true}},
		{"2", SI{Value: 2, Valid: true}},
		{"112233445566", SI{Value: 112233445566, Valid: true}},
		{`""`, SI{Valid: false}},
	}
	for _, tt := range tests {
		var si SI
		if err := si.Unmarshal([]byte(tt.in), testContext); err != nil {
			t.Fatalf("ParseSI(%q) failed with %v", tt.in, err)
		}
		if diff := cmp.Diff(tt.out, si); diff != "" {
			t.Errorf("ParseSI(%q) mismatch (-want, +got)=\n%s", tt.in, diff)
		}
	}
}

func TestParseSI_Error(t *testing.T) {
	tests := []string{
		"",
		"-",
		" ",
		"-1",  // Only non-negative numbers allowed.
		"1.2", // Only integer numbers allowed.
		"2-1",
	}
	for _, tt := range tests {
		var si SI
		if err := si.Unmarshal([]byte(tt), testContext); err == nil {
			t.Errorf("ParseSI(%q) got err=<nil>, want error", tt)
		}
	}
}

func TestParseNM(t *testing.T) {
	tests := []struct {
		in  string
		out NM
	}{
		{"0", NM{Value: 0.0, Valid: true}},
		{"-0", NM{Value: 0.0, Valid: true}},
		{"0.0", NM{Value: 0.0, Valid: true}},
		{"-0.0", NM{Value: 0.0, Valid: true}},
		{"0011.2200", NM{Value: 11.22, Valid: true}},   // Leading/trailing zeroes.
		{"-0011.2200", NM{Value: -11.22, Valid: true}}, // Leading/trailing zeroes.
		{"112233445566", NM{Value: 112233445566.0, Valid: true}},
		{"-112233445566", NM{Value: -112233445566.0, Valid: true}},
		{"112233445566.77", NM{Value: 112233445566.77, Valid: true}},
		{"-112233445566.77", NM{Value: -112233445566.77, Valid: true}},
		{`""`, NM{Valid: false}},
	}
	for _, tt := range tests {
		var nm NM
		if err := nm.Unmarshal([]byte(tt.in), testContext); err != nil {
			t.Errorf("ParseNM(%q) failed with %v", tt.in, err)
		}
		if diff := cmp.Diff(NM(tt.out), nm); diff != "" {
			t.Errorf("ParseNM(%q) mismatch (-want, +got)=\n%s", tt.in, diff)
		}
	}
}

func TestParseNM_Error(t *testing.T) {
	tests := []string{
		"",
		"-",
		" ",
		"2-1",
	}
	for _, tt := range tests {
		var nm NM
		if err := nm.Unmarshal([]byte(tt), testContext); err == nil {
			t.Errorf("ParseNM(%q) got err=<nil>, want error", tt)
		}
	}
}

func TestParseSNM(t *testing.T) {
	tests := []struct {
		in  string
		out SNM
	}{
		{"1234", SNM{Value: "1234", Valid: true}},
		{"00012345", SNM{Value: "00012345", Valid: true}},
		{"+0012345", SNM{Value: "+0012345", Valid: true}},
		{"123 456", SNM{Value: "123 456", Valid: true}},
		{"+1 123 456", SNM{Value: "+1 123 456", Valid: true}},
		{"  +1 123 456  ", SNM{Value: "+1 123 456", Valid: true}},
		{"+00 123 456", SNM{Value: "+00 123 456", Valid: true}},
		{`""`, SNM{Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var snm SNM
			if err := snm.Unmarshal([]byte(tt.in), testContext); err != nil {
				t.Errorf("Unmarshal(%q) failed with %v", tt.in, err)
			}
			if diff := cmp.Diff(SNM(tt.out), snm); diff != "" {
				t.Errorf("Unmarshal(%q) mismatch (-want, +got)=\n%s", tt.in, diff)
			}
		})
	}
}

func TestParseSNM_Error(t *testing.T) {
	tests := []string{
		"a string",
		"-",
		"-0",
		"1.5",
		"+44+32789",
		"+44 +32 789",
		"+  123 456",
	}
	for _, tt := range tests {
		var snm SNM
		if err := snm.Unmarshal([]byte(tt), testContext); err == nil {
			t.Errorf("Unmarshal(%q) got err=<nil>, want error", tt)
		}
	}
}

func TestParseST_UnescapesText(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`One\F\Escape`, "One|Escape"},
		{`Two\F\Escapes\S\`, "Two|Escapes^"},
		{`No spaces\F\\R\between escapes`, "No spaces|~between escapes"},
		{`\F\Escape at index zero`, "|Escape at index zero"},
		// Technically the following string is invalid, as raw delimiters
		// are not allowed, but choose to be more permissive.
		{`Escaped\F\and|^&~not escaped`, "Escaped|and|^&~not escaped"},
		{"", ""},
	}
	for _, test := range tests {
		var st ST
		err := st.Unmarshal([]byte(test.in), testContext)
		if err != nil {
			t.Fatalf("Unmarshal(%v) failed with %v", test.in, err)
		}
		if got, want := st, ST(test.out); got != want {
			t.Errorf("Unmarshal(%v) got %s, want %s", test.in, got, want)
		}
	}
}

// Test escape sequences that are technically invalid, but choose to be more permissive (STR-2993).
func TestParseST_UnescapesText_InvalidPermittedSequences(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`Unterminated\escape`, "Unterminated escape"},
		{`\Unterminated escape`, " Unterminated escape"},
		{`Empty\\escape`, "Empty escape"},
		{`\\Empty escape`, " Empty escape"},
		// This is supported in fields of type FT, but not in ST.
		{`New\.br\line`, "New line"},
	}
	for _, test := range tests {
		var st ST
		err := st.Unmarshal([]byte(test.in), testContext)
		if err != nil {
			t.Fatalf("Unmarshal(%v) failed with %v", test.in, err)
		}
		if got, want := st, ST(test.out); got != want {
			t.Errorf("Unmarshal(%v) got %s, want %s", test.in, got, want)
		}
	}
}

func TestParseST_UnescapesTextWithErrors(t *testing.T) {
	tests := []string{
		`Unknown\X\escape`,
		`Unknown\XX\multi character escape`,
		// Sequences not allowed in ST fields but valid in other text fields.
		`Highlighting \H\escape`,
		`Hexadecimal \X9\value`,
	}
	for _, test := range tests {
		var st ST
		err := st.Unmarshal([]byte(test), testContext)
		if err == nil {
			t.Errorf("Unmarshal(%v) got err=<nil>, want error", test)
		}
	}
}

func TestParseFT_UnescapesText(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`One\F\Escape`, "One|Escape"},
		{"", ""},
		{`Highlighting \H\escape`, "Highlighting escape"},
		{`Normal \N\text escape`, "Normal text escape"},
		{`Custom \Zarbitrary.Chars\escape`, "Custom escape"},
		{`Hexadecimal value\X000a\with X000a`, "Hexadecimal value\nwith X000a"},
		{`Hexadecimal value\X000d\with X000d`, "Hexadecimal value\rwith X000d"},
		{`Hexadecimal \X09Af\value`, "Hexadecimal value"},
		{`New\.br\line`, "New\nline"},
		{`New\.sp\line`, "New\nline"},
		{`Two\.sp2\new lines`, "Two\n\nnew lines"},
		{`Two\.sp+2\new lines`, "Two\n\nnew lines"},
	}
	for _, test := range tests {
		var ft FT
		err := ft.Unmarshal([]byte(test.in), testContext)
		if err != nil {
			t.Fatalf("Unmarshal(%v) failed with %v", test.in, err)
		}
		if got, want := ft, FT(test.out); got != want {
			t.Errorf("Unmarshal(%v) got %v, want %v", test.in, got, want)
		}
	}
}

func TestParseFT_EscapesText(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"One|Field", `One\F\Field`},
		{"Many|Fields|a|b", `Many\F\Fields\F\a\F\b`},
		{"Component1^Component2^Component3", `Component1\S\Component2\S\Component3`},
		{"Subcomponent1&Subcomponent2&", `Subcomponent1\T\Subcomponent2\T\`},
		{"Reptition1~Repetition2~Repetition3~", `Reptition1\R\Repetition2\R\Repetition3\R\`},
		{"line break 1\nline break 2\n", `line break 1\.br\line break 2\.br\`},
		{"two new lines\n\ntwo new lines\n\n", `two new lines\.br\\.br\two new lines\.br\\.br\`},
	}
	for _, test := range tests {
		ft := FT(test.in)
		out, err := ft.Marshal(testContext)
		if err != nil {
			t.Fatalf("Marshal(%q) failed with %v", test.in, err)
		}
		if got, want := string(out), test.want; got != want {
			t.Errorf("Marshal(%q) got %v, want %v", test.in, got, want)
		}
	}
}

// Test escape sequences that are technically invalid, but choose to be more permissive (STR-2993).
func TestParseFT_UnescapesText_InvalidPermittedSequences(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`Unterminated\escape`, "Unterminated escape"},
		{`\Unterminated escape`, " Unterminated escape"},
		{`Empty\\escape`, "Empty escape"},
		{`\\Empty escape`, " Empty escape"},
	}
	for _, test := range tests {
		var ft FT
		err := ft.Unmarshal([]byte(test.in), testContext)
		if err != nil {
			t.Fatalf("Unmarshal(%v) failed with %v", test.in, err)
		}
		if got, want := ft, FT(test.out); got != want {
			t.Errorf("Unmarshal(%v) got %v, want %v", test.in, got, want)
		}
	}
}

func TestParseFT_UnescapesTextWithErrors(t *testing.T) {
	tests := []string{
		`Unknown\X\escape`,
		`Unknown\XX\multi character escape`,
		`Incomplete\X\hexadecimal`,
		`Wrong\Xg\hexadecimal`,
		`Incomplete\Z\custom`,
		`SP\.sp-4\with negative count`,
	}
	for _, test := range tests {
		var ft FT
		err := ft.Unmarshal([]byte(test), testContext)
		if err == nil {
			t.Errorf("Unmarshal(%v) got err=<nil>, want error", test)
		}
	}
}

func TestTX_UnescapesText(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want string
	}{
		{`One\F\Escape`, "One|Escape"},
		{`Two\F\Escapes\S\`, "Two|Escapes^"},
		{`No spaces\F\\R\between escapes`, "No spaces|~between escapes"},
		{`\F\Escape at index zero`, "|Escape at index zero"},
		{`Escaped\F\and|^&~not escaped`, "Escaped|and|^&~not escaped"},
		{`result\.br\result`, "result\nresult"},
		{"", ""},
	} {
		t.Run(tt.in, func(t *testing.T) {
			var tx TX
			if err := tx.Unmarshal([]byte(tt.in), testContext); err != nil {
				t.Errorf("tx.Unmarshal(%q, testContext) failed with error %+v", tt.in, err)
			}
			if got, want := string(tx), tt.want; got != want {
				t.Errorf("tx.Unmarshal(%q, testContext)=%q, want %q", tt.in, got, want)
			}
		})
	}
}

func TestHDString(t *testing.T) {
	tests := []struct {
		in      HD
		wantOut string
	}{
		{HD{NamespaceID: NewIS("namespace"),
			UniversalID:     NewST("ID"),
			UniversalIDType: NewID("IDType")}, "namespace^ID^IDType"},
		{HD{NamespaceID: NewIS("namespace")}, "namespace"},
		{HD{UniversalID: NewST("UID")}, "^UID"},
		{HD{NamespaceID: NewIS("namespace"),
			UniversalIDType: NewID("IDType")}, "namespace^^IDType"},
		{HD{}, ""},
	}
	for _, tt := range tests {
		if got, want := tt.in.String(), tt.wantOut; got != want {
			t.Errorf("[%v].String() got %v, want %v", tt.in, got, want)
		}
	}
}
