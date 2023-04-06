package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Field struct {
	Name      string
	Type      string
	Necessity string
	Info      string
	Default   interface{}
}

var mandatory = []Field{
	{
		Name:      "rpc.algodUrl",
		Type:      "string",
		Necessity: "always needed",
		Info:      "URL of the algorand RPC node to connect to",
	},
	{
		Name:      "rpc.indexerUrl",
		Type:      "string",
		Necessity: "always needed",
		Info:      "URL of the algorand RPC node to connect to",
	},
	{
		Name:      "rpc.token",
		Type:      "string",
		Necessity: "always needed",
		Info:      "Token to use for RPC authentication",
	},
}

func validateConfig() error {
	for _, mf := range mandatory {
		err := validateField(mf)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateField(f Field) error {
	if !viper.IsSet(f.Name) {
		return errors.New("config error: unset mandatory field: " + f.Name + " (" + f.Type + "); Info: " + f.Info)
	}

	ok := true
	item := viper.Get(f.Name)

	switch f.Type {
	case "string":
		_, ok = item.(string)
	case "uint64":
		_, ok = item.(int)
	case "bool":
		_, ok = item.(bool)
	case "[]string":
		elems, ok := item.([]interface{})
		if !ok {
			break
		}
		for _, i := range elems {
			_, ok := i.(string)
			if !ok {
				break
			}
		}
	}

	if !ok {
		return errors.New("mandatory field type invalid: " + f.Name + " (" + f.Type + "); description: " + f.Info)
	}
	return nil
}
