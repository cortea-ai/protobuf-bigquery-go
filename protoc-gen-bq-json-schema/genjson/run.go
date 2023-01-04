package genjson

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func Run(gen *protogen.Plugin, config Config) error {
	for _, f := range gen.Files {
		if config.SchemaOptions.IncludePath != "" && !strings.HasPrefix(f.Desc.Path(), config.SchemaOptions.IncludePath) {
			continue
		}
		if !f.Generate {
			continue
		}
		for _, msg := range f.Messages {
			if err := GenerateSchemaFile(gen, msg, config); err != nil {
				return err
			}
		}
	}
	return nil
}
