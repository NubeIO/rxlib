package topics

import "fmt"

const GlobalTopic = "r/global"
const GlobalTopicResponse = "r/global/resp"
const BiosTopic = "r/bios"
const BiosGlobalTopic = "r/bios/global"

func GetGlobalTopicResponse(globalID string) string {
	return fmt.Sprintf("%s/%s", GlobalTopic, globalID)
}

func BuildBiosTopic(globalID string) string {
	return fmt.Sprintf("%s/%s", BiosTopic, globalID)
}
