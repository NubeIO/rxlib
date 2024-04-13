package topics

import "fmt"

const GlobalTopic = "r/global"
const BiosTopic = "r/bios"
const BiosGlobalTopic = "r/bios/global"

func GlobalTopicResponse(globalID string) string {
	return fmt.Sprintf("%s/%s", GlobalTopic, globalID)
}

func BuildBiosTopic(globalID string) string {
	return fmt.Sprintf("%s/%s", BiosTopic, globalID)
}
