package core

import (
	"encoding/json"
	"fmt"
)

/*
CustomFields holds top-level payload members that are not part of the OCPP 1.6
message definition. Some vendors extend a message in place rather than nesting
their additions under a "data" object, and dropping those members silently makes
the extension invisible to the application. Messages that opt in keep them here,
unparsed, and write them back out unchanged when marshalling.
*/
type CustomFields map[string]json.RawMessage

// Get unmarshals a single custom field into out. It reports false if the field
// was absent, so a caller can distinguish "not sent" from "sent as null".
func (f CustomFields) Get(key string, out interface{}) (bool, error) {
	raw, ok := f[key]
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return true, fmt.Errorf("custom field %q: %w", key, err)
	}
	return true, nil
}

// Set stores a custom field, marshalling value to JSON.
func (f *CustomFields) Set(key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("custom field %q: %w", key, err)
	}
	if *f == nil {
		*f = CustomFields{}
	}
	(*f)[key] = raw
	return nil
}

// extractCustomFields collects every top-level member of payload that is not one
// of the message's own known fields.
func extractCustomFields(payload []byte, known ...string) (CustomFields, error) {
	var all map[string]json.RawMessage
	if err := json.Unmarshal(payload, &all); err != nil {
		return nil, err
	}
	for _, k := range known {
		delete(all, k)
	}
	if len(all) == 0 {
		return nil, nil
	}
	return all, nil
}

/*
marshalWithCustomFields serialises msg and merges the custom fields back in at
the top level. msg must be an alias type without the custom MarshalJSON method,
otherwise this recurses. Known fields win on a name collision, so a stale custom
entry can never shadow the value the struct actually carries.
*/
func marshalWithCustomFields(msg interface{}, custom CustomFields) ([]byte, error) {
	raw, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	if len(custom) == 0 {
		return raw, nil
	}
	var merged map[string]json.RawMessage
	if err := json.Unmarshal(raw, &merged); err != nil {
		return nil, err
	}
	for k, v := range custom {
		if _, exists := merged[k]; exists {
			continue
		}
		merged[k] = v
	}
	return json.Marshal(merged)
}
