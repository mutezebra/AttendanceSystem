package service

import "encoding/json"

type SvcCallEvent struct {
	EventID   int64  `json:"event_id,omitempty"`
	Name      string `json:"Name,omitempty"`
	ClassID   int64  `json:"class_id,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
}

func (s *SvcCallEvent) GetID() int64 {
	return s.EventID
}

func (s *SvcCallEvent) GetCallEventName() string {
	return s.Name
}

func (s *SvcCallEvent) GetClassID() int64 {
	return s.ClassID
}

func (s *SvcCallEvent) GetStartTime() int64 {
	return s.StartTime
}

func (s *SvcCallEvent) GetEndTime() int64 {
	return s.EndTime
}

func (s *SvcCallEvent) GetEventID() int64 {
	return s.EventID
}

func (s *SvcCallEvent) convertToJson() []byte {
	data, _ := json.Marshal(s)
	return data
}

func JsonToSvcCallEvent(data []byte) *SvcCallEvent {
	var s SvcCallEvent
	_ = json.Unmarshal(data, &s)
	return &s
}
