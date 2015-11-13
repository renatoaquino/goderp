package goderp

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Info struct {
	descr string
	group string
	def   interface{}
}

type Config struct {
	Records       map[string]interface{}
	Descriptions  map[string]Info
	EnableEnvVars bool
}

func New() *Config {
	c := Config{}
	c.Records = make(map[string]interface{})
	c.Descriptions = make(map[string]Info)
	c.EnableEnvVars = false
	return &c
}

func (c *Config) EnableEnv() {
	c.EnableEnvVars = true
}

func (c *Config) Define(key string, value interface{}, descr string, group string) {
	c.Records[key] = value
	//fmt.Printf("%s:(%s) %s\n", key, reflect.TypeOf(c.Records[key]), c.Records[key])
	c.Descriptions[key] = Info{descr: descr, group: group, def: value}
}

func (c *Config) Get(key string) interface{} {
	return c.Records[key]
}

func (c *Config) GetInt(key string) int {
	return c.Records[key].(int)
}

func (c *Config) GetFloat(key string) float64 {
	return c.Records[key].(float64)
}

func (c *Config) GetString(key string) string {
	return c.Records[key].(string)
}

func (c *Config) GetBool(key string) bool {
	return c.Records[key].(bool)
}

func (c *Config) GetDescription(key string) string {
	return c.Descriptions[key].descr
}

func (c *Config) GetGroup(key string) string {
	return c.Descriptions[key].group
}

func (c *Config) GetDefault(key string) interface{} {
	return c.Descriptions[key].def
}

func (c *Config) Parse(filename string) (err error) {

	if _, err := toml.DecodeFile(filename, &c.Records); err != nil {
		return err
	}

	if c.EnableEnvVars {
		var tmp string
		for k := range c.Records {
			tmp = os.Getenv(k)
			if tmp != "" {
				c.Records[k], err = coerce(c.Records[k], tmp)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}
	return nil
}

func (c *Config) Dump() {
	groupkeys := make(map[string][]string)
	var group string
	for k := range c.Records {
		group = c.GetGroup(k)
		if _, ok := groupkeys[group]; !ok {
			groupkeys[group] = make([]string, 0)
		}
		groupkeys[group] = append(groupkeys[group], k)
	}

	const SEPARATOR = "#"
	const MAX_LEN = 80

	for k := range groupkeys {
		sep_size := ((MAX_LEN - utf8.RuneCountInString(k)) / 2) - 1
		fmt.Print(strings.Repeat(SEPARATOR, sep_size) + " " + k + " " + strings.Repeat(SEPARATOR, sep_size))

		if sep_size*2+utf8.RuneCountInString(k)+2 < MAX_LEN {
			fmt.Print(SEPARATOR)
		}
		fmt.Print("\n")

		for _, i := range groupkeys[k] {
			fmt.Printf("%s %s\n", strings.Repeat(SEPARATOR, 2), c.GetDescription(i))
			fmt.Printf("%s Defaults to: ", strings.Repeat(SEPARATOR, 2))
			fmt.Print(c.GetDefault(i))
			fmt.Printf("\n%s=", i)
			v := reflect.ValueOf(c.Get(i))
			switch v.Kind() {
			case reflect.String:
				fmt.Printf("\"%s\"\n", c.Get(i))
				break
			default:
				fmt.Println(c.Get(i))
			}
		}
		fmt.Print("\n")
	}
}

func coerce(current interface{}, replacement string) (interface{}, error) {
	v := reflect.ValueOf(current)
	switch v.Kind() {
	case reflect.String:
		return replacement, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(replacement, 0, 64)
		if err != nil {
			return current, err
		}
		return intValue, nil

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(replacement)
		if err != nil {
			return current, err
		}
		return boolValue, nil

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(replacement, 64)
		if err != nil {
			return current, err
		}
		return floatValue, nil
	}
	return current, nil
}
