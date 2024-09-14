namespace go api.class

include "base.thrift"
include "user.thrift"

struct UserWithClass {
    1: optional i64 UID (api.body="uid");
    2: optional string ClassID (api.body="class_id");
    3: optional i32 Weight (api.body="weight")
}

struct Class {
    1: optional i64 ID (api.body="id");
    2: optional string Name (api.body="user");
    3: optional i32 UserCount (api.body="user_count");
    4: optional string InvitationCode (api.body="invitation_code")
}

struct CreateClassReq {
    1: optional i64 UID (api.body="uid");
    2: optional string Name (api.body="name",api.form="name");
}

struct CreateClassResp {
    1: optional base.Base Base (api.body="base");
    2: optional i64 ClassID (api.body="class_id");
    3: optional string InvitationCode (api.body="invitation_code")
}

struct JoinClassReq {
    1: optional i64 ClassID (api.body="class_id",api.form="class_id");
    2: optional i64 UID (api.body="uid");
    3: optional string InvitationCode (api.body="invitation_code",api.form="invitation_code")
}

struct JoinClassResp {
    1: optional base.Base Base (api.body="base");
}

struct ClassListReq {
    1: optional i64 UID (api.body="uid");
}

struct ClassListResp {
    1: optional base.Base Base (api.body="base");
    2: optional i32 ClassCount (api.body="class_count")
    3: optional list<Class> Classes (api.body="Classes");
}

struct ClassStudentListReq { // only teacher can
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
}

struct ClassStudentListResp {
    1: optional base.Base Base (api.body="base");
    2: optional i32 UserCount (api.body="user_count")
    3: optional list<user.BaseUser> Students (api.body="students");
}

struct GetClassTeacherReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
}

struct GetClassTeacherResp {
    1: optional base.Base Base (api.body="base");
    2: optional user.BaseUser Teacher (api.body="teacher");
}

struct ViewInvitationCodeReq {
    1: optional i64 UID (api.body="uid");
    2: optional i64 ClassID (api.body="class_id",api.form="class_id");
}

struct ViewInvitationCodeResp {
    1: optional base.Base Base (api.body="base");
    2: optional string InvitationCode (api.body="invitation_code")
}

service ClassService {
    CreateClassResp CreateClass(1: CreateClassReq req) (api.post="/class/auth/create-class")
    JoinClassResp JoinClass(1: JoinClassReq req) (api.post="/class/auth/join-class")
    ClassListResp ClassList(1: ClassListReq req) (api.get="/class/auth/class-list")
    ClassStudentListResp ClassStudentList(1: ClassStudentListReq req) (api.get="/class/auth/student-list")
    GetClassTeacherResp GetClassTeacher(1: GetClassTeacherReq req) (api.get="/class/auth/get-teacher")
    ViewInvitationCodeResp ViewInvitationCode(1: ViewInvitationCodeReq req) (api.get="/class/auth/view-invitation-code")
}
