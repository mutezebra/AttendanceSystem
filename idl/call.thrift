namespace go api.call

include "base.thrift"
include "user.thrift"
include "class.thrift"

struct CallEvent {
    1: optional i64 ID (api.body="id");
    2: optional string CallEventName (api.body="call_event_name");
    3: optional i64 ClassID (api.body="class_id");
    4: optional string ClassName (api.body="class_name");
    5: optional i64 CallerID (api.body="caller")
    6: optional string CallerName (api.body="caller_name");
    7: optional i64 StartTime (api.body="start_time");
    8:optional i64 EndTime (api.body="end_time");
}

struct CallEventWithUser {
    1: optional i64 CallEventID (api.body="call_event_id");
    2: optional i64 UID (api.body="uid");
    3: optional string ClassName (api.body="class_name");
    4: optional string CallEventName (api.body="call_event_name");
    5: optional bool Done (api.body="done");
}

struct CallAllStudentReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
    3: optional i16 Deadline (api.body="deadline",api.form="deadline")
    4: optional string CallEventName (api.body="call_event_name",api.form="call_event_name")
}

struct CallAllStudentResp {
    1: optional base.Base Base (api.body="base");
    2: optional i64 EventID (api.body="event_id")
}

struct DoCallEventReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 EventID (api.body="event_id",api.form="event_id")
    3: optional i64 ClassID (api.body="class_id",api.form="class_id")
}

struct DoCallEventResp {
    1: optional base.Base Base (api.body="base");
}

struct UndoCallEventsReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
}

struct UndoCallEventsResp {
    1: optional base.Base Base (api.body="base");
    2: optional CallEvent Event (api.body="event");
    3: optional bool Exist (api.body="exist");
}

struct RandomCallReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
    3: optional i64 CallNumber (api.body="call_number",api.form="call_number");
    4: optional i16 Deadline (api.body="deadline",api.form="deadline");
    5: optional string CallEventName (api.body="call_event_name",api.form="call_event_name");
}

struct RandomCallResp {
    1: optional base.Base Base (api.body="base");
    2: optional i64 Count (api.body="count");
    3: optional list<user.BaseUser> Users (api.body="users");
    4: optional i64 EventID (api.body="event_id",api.form="event_id")
}

struct HistoryCallEventReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
}

struct HistoryCallEventResp {
    1: optional base.Base Base (api.body="base");
    2: optional list<CallEvent> Events (api.body="events");
}

service CallService {
    CallAllStudentResp CallAllStudent(1: CallAllStudentReq req) (api.post="/call/auth/call-all-student")
    DoCallEventResp DoCallEvent(1: DoCallEventReq req) (api.post="/call/auth/do-call-event")
    UndoCallEventsResp UndoCallEvents(1: UndoCallEventsReq req) (api.get="/call/auth/undo-call-events")
    RandomCallResp RandomCall(1: RandomCallReq req) (api.post="/call/auth/random-call")
    HistoryCallEventResp HistoryCallEvent(1: HistoryCallEventReq req) (api.get="/call/auth/history-call-event")
}
