package openapiextension_v1

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

type documentHandler func(version string, extensionName string, document string)
type extensionHandler func(name string, yamlInput string) (bool, proto.Message, error)

func forInputYamlFromOpenapic(handler documentHandler) {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("File error:", err.Error())
		os.Exit(1)
	}
	if len(data) == 0 {
		fmt.Println("No input data.")
		os.Exit(1)
	}
	request := &ExtensionHandlerRequest{}
	err = proto.Unmarshal(data, request)
	if err != nil {
		fmt.Println("Input error:", err.Error())
		os.Exit(1)
	}
	handler(request.Wrapper.Version, request.Wrapper.ExtensionName, request.Wrapper.Yaml)
}

// ProcessExtension calles the handler for a specified extension.
func ProcessExtension(handleExtension extensionHandler) {
	response := &ExtensionHandlerResponse{}
	forInputYamlFromOpenapic(
		func(version string, extensionName string, yamlInput string) {
			var newObject proto.Message
			var err error

			handled, newObject, err := handleExtension(extensionName, yamlInput)
			if !handled {
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}

			// If we reach here, then the extension is handled
			response.Handled = true
			if err != nil {
				response.Error = append(response.Error, err.Error())
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}
			response.Value, err = ptypes.MarshalAny(newObject)
			if err != nil {
				response.Error = append(response.Error, err.Error())
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}
		})

	responseBytes, _ := proto.Marshal(response)
	os.Stdout.Write(responseBytes)
}
