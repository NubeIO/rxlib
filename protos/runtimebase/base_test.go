package runtimebase

import (
	"testing"
)

// JSON input
var jsonString = `

{
    "objectDeploy": {
   "deleted": [
            "abc",
            "123"
        ],
        "new": [
            {
                "id": "trigger",
                "inputs": [],
                "outputs": [
                    {
                        "id": "output",
                        "name": "output",
                        "direction": "output",
                        "dataType": "float"
                    }
                ],
                "connections": [
                    {
                        "source": "triggerABC",
                        "sourceHandle": "output",
                        "target": "mathABC",
                        "targetHandle": "in-1",
                        "flowDirection": "publisher"
                    }
                ],
                "meta": {
                    "uuid": "triggerABC",
                    "name": "triggerABC",
                    "position": {
                        "positionY": -38,
                        "positionX": 155
                    }
                }
            }
        ]
    },
    "timeout": 10
    }`

func TestConvertStructConnectionToProto(t *testing.T) {
	// Unmarshal JSON to proto message
	//objDeploy := &runtime.ObjectDeployRequest{}
	//if err := jsonpb.UnmarshalString(jsonString, objDeploy); err != nil {
	//	fmt.Println("Error unmarshalling JSON:", err)
	//	return
	//}
	//
	//fmt.Println(proto.MarshalTextString(objDeploy))

}