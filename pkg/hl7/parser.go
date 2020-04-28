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

// Package hl7 provides functions for manipulating and handling HL7 messages, native
// HL7 types, and reading and writing messages over MLLP.
//
// The HL7 2.3 specification is defined here:
// http://www.hl7.org/implement/standards/product_brief.cfm?product_id=140.
//
// The MLLP specification is defined here:
// http://www.hl7.org/documentcenter/public/wg/inm/mllp_transport_specification.PDF
//
package hl7

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding"
	"github.com/google/simhospital/pkg/constants"
)

const (
	// SegmentTerminator is the character used to terminate a HL7 segment,
	// defined in section 2.7 of the HL7 2.3 specification.
	SegmentTerminator = constants.SegmentTerminator
	// SegmentTerminatorStr is the string representation of SegmentTerminator.
	SegmentTerminatorStr = constants.SegmentTerminatorStr

	asciiNewLine byte = 0xa
)

var (
	// Timezone is a timezone for dates in the generated HL7 messages.
	// Its value is set by TimezoneAndLocation.
	Timezone string
	// Location is a location loaded for the Timezone used for dates in the
	// generated HL7 messages.
	// Its value is set by TimezoneAndLocation.
	Location *time.Location
	// encodings is a mapping between string character set names on page 2-96 of the
	// HL7 2.3.1 specification to the corresponding Go encodings.
	encodings = map[string]encoding.Encoding{
		// The HL7 spec says messages that are ASCII but use characters outside the
		// printable 7-bit should be rejected. Right now, we pass them through
		// blindly so encoding issues can be seen upstream, and on the assumption
		// that ASCII means UTF-8 in practice.
		"ASCII":   encoding.Nop,
		"8859/1":  charmap.ISO8859_1,
		"8859/2":  charmap.ISO8859_2,
		"8859/3":  charmap.ISO8859_3,
		"8859/4":  charmap.ISO8859_4,
		"8859/5":  charmap.ISO8859_5,
		"8859/6":  charmap.ISO8859_6,
		"8859/7":  charmap.ISO8859_7,
		"8859/8":  charmap.ISO8859_8,
		"8859/9":  charmap.ISO8859_9,
		"8859/15": charmap.ISO8859_15,
		// We don't handle the following character sets, which are explicitly
		// mentioned in the spec:
		// ISO IR6
		// GB 18030-2000
		// KS X 1001
		// CNS 11643-1992
		// BIG-5
		// JAS2020: need to determine differences with x/text/encoding/japanese
		// JIS X 0202: likewise
		// ISO IR14/JIS X 0201-1976: likewise
		// ISO IR87/JIS X 0208-1990: likewise
		// ISO IR159/JIS X 0212-1990: likewise
		"UNICODE":       encoding.Nop, // Legacy charmap retained for v2.5 compatibility
		"UNICODE UTF-8": encoding.Nop, // Backport from v2.8
	}
	// null is a "null" value as defined in HL7 spec, ie: two double quotes without content.
	null = []byte(`""`)

	// segmentTypeRegex is a regex to parse the main type for a segment.
	segmentTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3}$`)
	// messageTypeRegex is a regex to parse the main type for a message.
	messageTypeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3}_[a-zA-Z0-9]{3}$`)

	defaultDelimiters = &Delimiters{
		Field:        '|',
		Component:    '^',
		Subcomponent: '&',
		Repetition:   '~',
		Escape:       '\\',
	}
)

// TimezoneAndLocation sets the Timezone and Location based on tz provided.
// Returns an error if the location for the given timezone cannot be loaded.
func TimezoneAndLocation(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return errors.Wrapf(err, "invalid timezone: %q", tz)
	}
	Timezone = tz
	Location = loc
	return nil
}

// Segment is an interface implemented by the generated structs that represent
// concrete HL7 segment types, eg, MSH, PID, etc.
type Segment interface {
	SegmentName() string
}

// Context represents state about the HL7 message as a whole (eg derived from
// the header) which is necessary to parse individual values.
type Context struct {
	Decoder    *encoding.Decoder
	Delimiters *Delimiters
	// Nesting represents how deep in the HL7 parsing we are. HL7 only allows two levels of
	// nesting, so we keep track of this as an int in order to be able to detect if the code is
	// behaving incorrectly (e.g. nesting contexts too many times).
	// `0` represents the initial nesting level.
	Nesting int
	// Timezone that TS values will be parsed in.
	TimezoneLoc *time.Location
}

// Nested returns a Context identical to the original one, but where the nesting level is
// incremented by one.
func (c *Context) Nested() *Context {
	newContext := *c
	newContext.Nesting++
	return &newContext
}

// ParseMessageOptions contains optional parameters to ParseMessage.
type ParseMessageOptions struct {
	TimezoneLoc *time.Location
	// SegmentTerminator contains characters used as an end of segment terminator.
	SegmentTerminator []byte
}

// NewParseMessageOptions returns a ParseMessageOptions, which can be used to
// configure the parser's behaviour.
func NewParseMessageOptions() *ParseMessageOptions {
	return &ParseMessageOptions{
		TimezoneLoc:       Location,
		SegmentTerminator: []byte{SegmentTerminator},
	}
}

// Delimiters are the delimiter characters used within a message, defined in
// section 2.7 of the HL7 2.3 specification.
type Delimiters struct {
	Field        byte
	Component    byte
	Subcomponent byte
	Repetition   byte
	Escape       byte
}

var _ Primitive = (*Delimiters)(nil)

// Marshal marshals Delimiters.
func (d *Delimiters) Marshal(_ *Context) ([]byte, error) {
	return []byte{d.Component, d.Repetition, d.Escape, d.Subcomponent}, nil
}

// Unmarshal Delimiters, replacing the values currently used in message unmarshaling.
func (d *Delimiters) Unmarshal(field []byte, c *Context) error {
	if len(field) < 4 {
		return ErrBadValue
	}
	*d = Delimiters{
		// Obtain the field delimiter from the context.
		Field:        c.Delimiters.Field,
		Component:    field[0],
		Repetition:   field[1],
		Escape:       field[2],
		Subcomponent: field[3],
	}
	// Also replace the delimiters in the current context.
	c.Delimiters = d
	return nil
}

func (d Delimiters) splitFields(segment Token) []Token {
	return split(segment, d.Field)
}

func (d Delimiters) splitComponents(field Token, nesting int) []Token {
	switch nesting {
	case 0:
		return split(field, d.Component)
	case 1:
		return split(field, d.Subcomponent)
	default:
		// Out of nesting levels: we can't split anymore. This occurs in a small
		// number of poorly defined HL7 types. See TestNestingDepthNeverExceedsTwo.
		return []Token{field}
	}
}

func (d Delimiters) joinComponents(components [][]byte, nesting int) []byte {
	if len(components) == 1 {
		// This masks the case of attempting to join components when we're out of
		// nesting levels - which is useful, since some types are broken. See
		// TestNestingDepthNeverExceedsTwo.
		return components[0]
	}
	switch nesting {
	case 0:
		return bytes.Join(components, []byte{d.Component})
	case 1:
		return bytes.Join(components, []byte{d.Subcomponent})
	default:
		// Out of nesting levels: we can't join anymore.
		panic("Too many nesting levels")
	}
}

func (d Delimiters) splitRepeated(field Token) []Token {
	return split(field, d.Repetition)
}

func (d Delimiters) joinRepeated(repetitions [][]byte) []byte {
	return bytes.Join(repetitions, []byte{d.Repetition})
}

func split(input Token, delimiter byte) []Token {
	return splitMultiCharDelimiter(input, []byte{delimiter})
}

func splitMultiCharDelimiter(input Token, delimiter []byte) []Token {
	r := make([]Token, 0, 16)
	start := 0
	tks := bytes.Split(input.Value, delimiter)
	for _, tk := range tks {
		r = append(r, Token{tk, input.Offset + start, input.Location})
		start += len(tk) + len(delimiter)
	}
	return r
}

// BadSegmentError occurs when we find a segment name we're not aware of.
type BadSegmentError struct {
	Name string
}

func (e *BadSegmentError) Error() string {
	return fmt.Sprintf("bad segment %q", e.Name)
}

// BadMessageTypeError occurs when we find a message type name we're not aware of.
type BadMessageTypeError struct {
	Name string
}

func (e *BadMessageTypeError) Error() string {
	errMsg := "bad message type"
	if segmentTypeRegex.MatchString(e.Name) || messageTypeRegex.MatchString(e.Name) {
		return fmt.Sprintf("%s: %s", errMsg, e.Name)
	}
	return errMsg
}

// ErrBadValue occurs when we can't parse the value for a primitive HL7 type.
// We typically can't pass through the underlying reason (eg something like:
// strconv.ParseFloat: parsing ""E83039"": invalid syntax) as that reason
// may contain patient identifiable data.
var ErrBadValue = errors.New("bad value for primitive HL7 type")

// ParseError occurs when we can't parse parse of a HL7 message.
type ParseError struct {
	Offset   int
	Location string
	Cause    error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("error in %s: %v", e.Location, e.Cause.Error())
}

// ParseErrors occurs when we encounter multiple errors while attempting to parse
// a HL7 message.
type ParseErrors []ParseError

func (e ParseErrors) Error() string {
	s := make([]string, len(e))
	for i, err := range e {
		s[i] = err.Error()
	}
	return fmt.Sprintf("errors (%d): %s", len(e), strings.Join(s, ", "))
}

// Token represents a substring within a HL7 message, together with the offset
// in bytes at which that token starts within the message.
type Token struct {
	Value  []byte
	Offset int
	// Location is a string describing the HL7 field that corresponds the
	// location of this token, built from segment/type name, field number, and
	// description, separated path style, eg: PID-2-Patient ID/CX-2-Check Digit
	Location string
}

// Error returns a ParseError at the location of this token, due to err.
func (t Token) Error(err error) *ParseError {
	return &ParseError{t.Offset, t.Location, err}
}

// Errors returns a ParseErrors with a single ParseError at the location
// of this token, due to err.
func (t Token) Errors(err error) ParseErrors {
	return ParseErrors{*t.Error(err)}
}

func parseTag(f reflect.StructField) (string, bool) {
	parts := strings.SplitN(f.Tag.Get("hl7"), ",", 2)
	if len(parts) == 2 {
		return parts[1], parts[0] == "true"
	}
	return parts[0], false
}

// fieldLocation returns a string describing field i within the segment type t,
// eg: PD1-2-Patient Primary Facility.
func fieldLocation(t reflect.Type, i int) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name, _ := parseTag(t.Field(i))
	if name != "" {
		name = "-" + name
	}
	return fmt.Sprintf("%s-%d%s", t.Name(), i+1, name)
}

func appendLocation(location string, next string) string {
	if location != "" {
		return location + "/" + next
	}
	return next
}

// Message is an HL7 message.
type Message struct {
	*Context
	Segments []Token
	msh      MSH
}

// ParseMessage returns an object representing the HL7 message in input,
// ensuring it has a correct header, and returning an error if not.
func ParseMessage(input []byte) (*Message, error) {
	mo := NewParseMessageOptions()
	return ParseMessageWithOptions(input, mo)
}

// ParseMessageWithOptions returns an object representing the HL7 message in
// input, ensuring it has a correct header, and returning an error if not.
// If options.TimezoneLoc is populated, the given timezone is used to interpret
// dates from the message. If not specified, Timezone is used.
// options.segmentTerminator is used as the segment terminator or delimiter. The default value is
// \r. The spec doesn't allow custom values for this delimiter, but it might be necessary to change
// it to deal with some messages that use a non-standard terminator.
func ParseMessageWithOptions(input []byte, options *ParseMessageOptions) (*Message, error) {
	delimiters := defaultDelimiters
	// Messages start with "MSH" header and 5 delimiter characters.
	if len(input) < 8 || !bytes.HasPrefix(input, []byte("MSH")) {
		return nil, errors.New("bad HL7 MSH header")
	}
	delimiters = &Delimiters{
		Field: input[3],
		// The remaining delimiters are filled in when the MSH segment is parsed.
	}

	m := &Message{
		Context: &Context{
			Decoder:     encoding.Nop.NewDecoder(),
			Delimiters:  delimiters,
			Nesting:     0,
			TimezoneLoc: options.TimezoneLoc,
		},
		Segments: splitMultiCharDelimiter(Token{input, 0, ""}, options.SegmentTerminator),
	}
	err := parseSegment(m.Segments[0].Value, m.Context, &m.msh)
	if err != nil {
		return nil, err
	}
	if len(m.msh.CharacterSet) > 0 && m.msh.CharacterSet[0] != "" {
		enc, ok := encodings[strings.TrimSpace(string(m.msh.CharacterSet[0]))]
		if !ok {
			return nil, fmt.Errorf("bad character set: %q", string(m.msh.CharacterSet[0]))
		}
		m.Context.Decoder = enc.NewDecoder()
	}
	return m, nil
}

func parseCompositeValue(input Token, c *Context, v reflect.Value) error {
	components := c.Delimiters.splitComponents(input, c.Nesting)
	errs := ParseErrors{}
	for i := 0; i < v.NumField(); i++ {
		if i >= len(components) {
			// We're out of fields. Some of the missing fields could be required and we could
			// potentially check them and throw an error because they're absent, but we choose to
			// ignore this type of error for simplicity.
			break
		}
		component := components[i]
		component.Location = appendLocation(component.Location, fieldLocation(v.Type(), i))
		err := parseValue(component, c.Nested(), v.Field(i))
		if err != nil {
			if perr, ok := err.(ParseErrors); ok {
				errs = append(errs, perr...)
			} else {
				return err
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// parseSegment parses a HL7 segment, ef PID. v should be a pointer to an
// instance of the corresponding struct, for example:
//   var pid PID
//   err := parseSegment(segment, c, &pid)
func parseSegment(input []byte, c *Context, v interface{}) error {
	return parseSegmentValue(Token{input, 0, ""}, c, reflect.ValueOf(v).Elem())
}

func parseSegmentValue(input Token, c *Context, v reflect.Value) error {
	input.Location = v.Type().Name()
	vType := v.Type()
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	input.Location = vType.Name()
	fields := c.Delimiters.splitFields(input)
	errs := ParseErrors{}
	for i := 0; i < v.NumField(); i++ {
		if i+1 >= len(fields) {
			// We're out of fields. Some of the missing fields could be required and we could
			// potentially check them and throw an error because they're absent, but we choose to
			// ignore this type of error for simplicity.
			break
		}
		// i+1 to skip the segment type that's in field[0] (eg MSH, PID), and not
		// in the struct.
		field := fields[i+1]
		field.Location = fieldLocation(v.Type(), i)
		err := parseValue(field, c, v.Field(i))
		if err != nil {
			if perr, ok := err.(ParseErrors); ok {
				errs = append(errs, perr...)
			} else {
				return err
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func parseRepeatedValue(input Token, c *Context, v reflect.Value) error {
	elements := c.Delimiters.splitRepeated(input)
	errs := ParseErrors{}
	slice := reflect.MakeSlice(v.Type(), len(elements), len(elements))
	for i := 0; i < len(elements); i++ {
		err := parseValue(elements[i], c, slice.Index(i))
		if err != nil {
			if perr, ok := err.(ParseErrors); ok {
				errs = append(errs, perr...)
			} else {
				return err
			}
		}
	}
	v.Set(slice)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func isHL7Null(field []byte) bool {
	return bytes.Equal(field, null)
}

func parseValue(input Token, c *Context, v reflect.Value) error {
	if !v.CanSet() {
		panic("Can't set value") // Implies a bug in the parser.
	}
	if len(input.Value) == 0 {
		return nil
	}
	var primitive Primitive
	primitiveType := reflect.TypeOf((*Primitive)(nil)).Elem()
	if v.Type().Implements(primitiveType) {
		n := reflect.New(v.Type().Elem())
		v.Set(n)
		primitive = n.Interface().(Primitive)
	} else if reflect.PtrTo(v.Type()).Implements(primitiveType) {
		primitive = v.Addr().Interface().(Primitive)
	}
	if primitive != nil {
		err := primitive.Unmarshal(input.Value, c)
		if err != nil {
			return input.Errors(err)
		}
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		// Pointers to primitives are handled earlier, so anything here must
		// be a composite.
		n := reflect.New(v.Type().Elem())
		v.Set(n)
		return parseCompositeValue(input, c, n.Elem())
	case reflect.Slice:
		return parseRepeatedValue(input, c, v)
	case reflect.Struct:
		return parseCompositeValue(input, c, v)
	default:
		panic("Unexpected kind: " + string(v.Kind()) + " type: " + v.Type().Name()) // Implies a bug in the parser.
	}
}

// segmentName returns the name of the given segment.
// Segment names are either 3 characters long, and are followed by the field delimiter.
// If the value is exactly 3 characters long without a trailing field delimiter, we treat that as correct.
func segmentName(segment Token, d *Delimiters) (string, *ParseError) {
	if len(segment.Value) < 4 || segment.Value[3] != d.Field {
		l := len(segment.Value)
		if l == 3 {
			return string(segment.Value[:3]), nil
		}
		if l > 3 {
			l = 3
		}
		return "", segment.Error(&BadSegmentError{string(segment.Value[0:l])})
	}
	return string(segment.Value[:3]), nil
}

func isSegment(s Token, expected string, d *Delimiters) bool {
	name, err := segmentName(s, d)
	return err == nil && name == expected
}

// parse returns a pointer to the parsed representation of the first segment
// of the specified segmentType.
func (m *Message) parse(name string) (interface{}, error) {
	t, ok := Types[name]
	if !ok {
		return nil, &BadSegmentError{name}
	}
	for _, s := range m.Segments {
		if !isSegment(s, name, m.Context.Delimiters) {
			continue
		}
		ps := reflect.New(t)
		err := parseSegmentValue(s, m.Context, ps.Elem())
		return ps.Interface(), err
	}
	return nil, nil
}

// ParseAll returns a slice of pointers to the parsed representations of every
// segment of the specified segmentType.
func (m *Message) ParseAll(name string) (interface{}, error) {
	t, ok := Types[name]
	if !ok {
		return nil, &BadSegmentError{name}
	}
	v := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(t)), 0, 1)
	errs := ParseErrors{}
	for _, s := range m.Segments {
		if !isSegment(s, name, m.Context.Delimiters) {
			continue
		}
		ps := reflect.New(t)
		err := parseSegmentValue(s, m.Context, ps.Elem())
		if err != nil {
			if perr, ok := err.(ParseErrors); ok {
				errs = append(errs, perr...)
			} else {
				return nil, err
			}
		}
		v = reflect.Append(v, ps)
	}
	if len(errs) > 0 {
		return v.Interface(), errs
	}
	return v.Interface(), nil
}

// All returns a slice of pointers to the parsed representations of
// every segment within the message.
func (m *Message) All() ([]interface{}, error) {
	v := make([]interface{}, 0, len(m.Segments))
	errs := ParseErrors{}
	for _, s := range m.Segments {
		if len(s.Value) == 0 {
			continue
		}
		name, pe := segmentName(s, m.Context.Delimiters)
		if pe != nil {
			errs = append(errs, *pe)
			continue
		}
		if strings.HasPrefix(name, "Z") {
			v = append(v, &GenericHL7Segment{s.Value})
			continue
		}
		t, ok := Types[name]
		if !ok {
			errs = append(errs, *s.Error(&BadSegmentError{name}))
			continue
		}
		ps := reflect.New(t)
		err := parseSegmentValue(s, m.Context, ps.Elem())
		if err != nil {
			if perr, ok := err.(ParseErrors); ok {
				errs = append(errs, perr...)
			} else {
				return nil, err
			}
		}
		v = append(v, ps.Interface())
	}
	if len(errs) > 0 {
		return v, errs
	}
	return v, nil
}

// endOfFieldsWithValues returns the index of the last field within v for which
// that and all subsequent fields have a nil value.
func endOfFieldsWithValues(v reflect.Value) int {
	var last int
	for last = v.NumField() - 1; last >= 0; last-- {
		if !v.Field(last).IsNil() {
			break
		}
	}
	return last + 1
}

// Primitive represents a primitive HL7 type, eg ST or ID.
type Primitive interface {
	Marshal(*Context) ([]byte, error)
	Unmarshal([]byte, *Context) error
}

func marshalValue(v reflect.Value, c *Context) ([]byte, error) {
	var primitive Primitive
	primitiveType := reflect.TypeOf((*Primitive)(nil)).Elem()
	if v.Type().Kind() == reflect.Ptr && v.IsNil() {
		return []byte{}, nil
	} else if v.Type().Implements(primitiveType) {
		primitive = v.Interface().(Primitive)
	} else if reflect.PtrTo(v.Type()).Implements(primitiveType) {
		primitive = v.Addr().Interface().(Primitive)
	}
	if primitive != nil {
		return primitive.Marshal(c)
	}

	switch v.Kind() {
	case reflect.Ptr:
		// Primitives are handled earlier, so anything here must by a composite.
		return marshalCompositeValue(v.Elem(), c)
	case reflect.Struct:
		return marshalCompositeValue(v, c)
	case reflect.Slice:
		return marshalRepeatedValue(v, c)
	default:
		// Implies a bug in the marshaller.
		panic("Unexpected kind: " + string(v.Kind()) + " type: " + v.Type().Name())
	}
}

func marshalCompositeValue(v reflect.Value, c *Context) ([]byte, error) {
	var err error
	end := endOfFieldsWithValues(v)
	fields := make([][]byte, end)
	for i := 0; i < end; i++ {
		fields[i], err = marshalValue(v.Field(i), c.Nested())
		if err != nil {
			return nil, err
		}
	}
	return c.Delimiters.joinComponents(fields, c.Nesting), nil
}

func marshalRepeatedValue(v reflect.Value, c *Context) ([]byte, error) {
	repetitions := make([][]byte, v.Len())
	var err error
	for i := 0; i < v.Len(); i++ {
		repetitions[i], err = marshalValue(v.Index(i), c)
		if err != nil {
			panic(err)
		}
	}
	return c.Delimiters.joinRepeated(repetitions), nil
}

// parseText parses a text field.
func parseText(field []byte, c *Context) (string, error) {
	if c == nil {
		panic("nil context")
	}
	if c.Decoder == nil {
		panic("nil decoder")
	}
	decoded, err := c.Decoder.String(string(field))
	if err != nil {
		return "", err
	}
	return decoded, nil
}

// marshalText marshals a text field.
// The returned slice will contain escaped characters.
func marshalText(field []byte, c *Context) []byte {
	dst := make([]byte, 0, len(field))
	for _, b := range field {
		switch b {
		case c.Delimiters.Field:
			dst = append(dst, []byte("\\F\\")...)
		case c.Delimiters.Component:
			dst = append(dst, []byte("\\S\\")...)
		case c.Delimiters.Subcomponent:
			dst = append(dst, []byte("\\T\\")...)
		case c.Delimiters.Repetition:
			dst = append(dst, []byte("\\R\\")...)
		case c.Delimiters.Escape:
			dst = append(dst, []byte("\\E\\")...)
		case asciiNewLine:
			dst = append(dst, []byte("\\.br\\")...)
		default:
			dst = append(dst, b)
		}
	}
	return dst
}
