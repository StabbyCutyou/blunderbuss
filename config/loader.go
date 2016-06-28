package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//
// The concepts and code here were heavily influenced by https://github.com/caarlos0/env
// Due to the significant level of changes we needed to make, and in an effort to reduce
// the number of forked libraries and external dependencies our default service implementation
// requires, we've chosen to roll it directly into the core service scaffold
//

// ErrNotAStructPtr is an error indicating that the user did not pass a pointer to a struct to the Initialize function
var ErrNotAStructPtr = errors.New("Expected a pointer to a Struct")

// Initialize takes a struct decorated with field tags to help dynamically load the environment
// into any of the Config structs inside the service.
// You must restart to pickup environment changes, as this will only happen when the service boots

func Initialize(configDesc interface{}) error {
	// get the value behind the interface
	ptrRef := reflect.ValueOf(configDesc)
	// make sure it's a pointer
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	// defref pointer to get the value of what it points to
	ref := ptrRef.Elem()
	// make sure it points to a struct
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}

	// now parse the fields of the struct
	return doParse(ref)
}

// ReadFromEnv is a light wrapper around Initialize which creates and returns a
// Config pointer, read from the environment
func ReadFromEnv() (*Config, error) {
  c := &Config{}
  if err := Initialize(c); err != nil {
    return nil, err
  }
  return c, nil
}

func doParse(ref reflect.Value) error {
	var unsetEnv []string

	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		envKey := field.Tag.Get("env")
		if envKey == "" {
			// the env field tag doesn't exist, ignore this field
			continue
		}
		envValue := os.Getenv(envKey)

		defaultValue := field.Tag.Get("default")

		// This tag is necessary in the case where you want an env var to default to the empty
		// string (like passwords). Because field.Tag.Get returns a string which can't be nil
		// it can't tell the difference between a default of "" and no default.
		// So in the case where default is explicitly "" and you don't want the program to panic
		// you can add the option:"true" to indicate that this env var can be missing
		optionalValue := field.Tag.Get("optional")

		if envValue == "" {
			envValue = defaultValue
		}

		if envValue == "" && optionalValue != "true" {
			// add env variable to list of unset variables for reporting an error.
			unsetEnv = append(unsetEnv, envKey)
			continue
		}

		err := set(ref.Field(i), envValue)
		if err != nil {
			return err
		}
	}

	// TODO: check to make sure there is at least one field tag `env` and panic otherwise.

	if unsetEnv != nil {
		errString := "the following environment variables must be set: " + strings.Join(unsetEnv, " ")
		return errors.New(errString)
	}
	return nil // no errors
}

func set(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Float32:
		floatVal, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	default:
		return fmt.Errorf("There was a type in the config struct that is not supported. Field: %s, Value: %s", field.Kind(), value)
	}
	return nil
}
