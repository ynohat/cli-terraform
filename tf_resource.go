// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	//"io/ioutil"
	//"path/filepath"
	//"encoding/json"
	//"os"
	"reflect"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"fmt"

)

// resource config 
var gtmResourceConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_resource" `)

// Process resource resources
func processResources(resources []*gtm.Resource, rImportList map[string][]int, dcIL map[int]string, resourceDomainName string) (string) {

	resourcesString := ""
	for _, resource := range resources {
		if _, ok := rImportList[resource.Name]; !ok {
			continue
		}
        	resourceBody := ""
        	name := ""
        	rString := gtmResourceConfigP1
        	rElems := reflect.ValueOf(resource).Elem()
        	for i := 0; i < rElems.NumField(); i++ {
                	varName := rElems.Type().Field(i).Name
                	varType := rElems.Type().Field(i).Type
                	varValue := rElems.Field(i).Interface()
                	key := convertKey(varName)
                	if key == "" {
				continue
			}
                	keyVal := fmt.Sprint(varValue)
                	if key == "name" { name = keyVal }
                	if varName == "ResourceInstances" {
				keyVal = processResourceInstances(varValue.([]*gtm.ResourceInstance))
               	 	}
			if keyVal == "" && varType.Kind() == reflect.String {
				continue
			}
                	resourceBody += tab4 + key + " = "
                	if varType.Kind() == reflect.String {
                        	resourceBody += "\"" + keyVal + "\"\n"
                	} else {
                        	resourceBody += keyVal + "\n"
                	}
        	}
        	rString += "\"" + name + "\" {"
        	rString += gtmRConfigP2 + resourceDomainName + ".name}\"\n"
        	rString += resourceBody
		rString += dependsClauseP1 + resourceDomainName + "\""
                // process dc dependencies (only one type in 1.4 schema)
                for _, dcDep := range rImportList[name] {
                        rString += ",\n"
                        rString += tab8 + datacenterResource + "." + dcIL[dcDep]
                }
		rString += "\n"
		rString += tab4 + "]\n"
        	rString += "}\n"
		resourcesString += rString
	}

        return resourcesString

}

func processResourceInstances(instances []*gtm.ResourceInstance) string {

        instanceString := "[]\n"                  // assume MT
        for ii, instance := range instances {
                instanceString = "[{\n"           // at least one
                instElems := reflect.ValueOf(instance).Elem()
                for i := 0; i < instElems.NumField(); i++ {
                        varName := instElems.Type().Field(i).Name
                        varType := instElems.Type().Field(i).Type
                        varValue := instElems.Field(i).Interface()
                        key := convertKey(varName)
                        keyVal := fmt.Sprint(varValue)
                        if varName == "LoadServers" {
                                keyVal = processStringList(instance.LoadServers)
                        }
			if varType.Kind() == reflect.Struct {
				fmt.Println(varName + " is a STRUCT!")
			}
                        if varType.Kind() == reflect.String {
                                instanceString += tab8 + "\"" + key + "\" = \"" + keyVal + "\"\n"
                        } else {
                                instanceString += tab8 + "\"" + key + "\" = " + keyVal + "\n"
                        }
                }
                if ii < len(instances) {
                        instanceString += tab8 + "},\n" + tab8 + "{\n"
                } else {
                        instanceString += tab8 + "}\n"
                        instanceString += tab4 + "]"
                }
        }
        return instanceString

}

