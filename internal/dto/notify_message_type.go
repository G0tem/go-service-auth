package dto

type NotifyMessageType uint8

const (
	NotifyUnknown NotifyMessageType = iota
	NotifyOnLiveBroadcast
	NotifyOnVOD
	NotifyOnMention
	NotifyOnNewSubscriber
	NotifyOnRecommendations
	NotifyOnNews
)

func ParseNotifyMessageType(smtpMessageType string) NotifyMessageType {
	switch smtpMessageType {
	case "@notify.OnLiveBroadcast":
		return NotifyOnLiveBroadcast
	case "@notify.OnVOD":
		return NotifyOnVOD
	case "@notify.OnMention":
		return NotifyOnMention
	case "@notify.OnNewSubscriber":
		return NotifyOnNewSubscriber
	case "@notify.OnRecommendations":
		return NotifyOnRecommendations
	case "@notify.OnNews":
		return NotifyOnNews
	}
	return NotifyUnknown
}

func (smtpMessageType NotifyMessageType) String() string {
	switch smtpMessageType {
	case NotifyOnLiveBroadcast:
		return "@notify.OnLiveBroadcast"
	case NotifyOnVOD:
		return "@notify.OnVOD"
	case NotifyOnMention:
		return "@notify.OnMention"
	case NotifyOnNewSubscriber:
		return "@notify.OnNewSubscriber"
	case NotifyOnRecommendations:
		return "@notify.OnRecommendations"
	case NotifyOnNews:
		return "@notify.OnNews"
	}
	return ""
}
