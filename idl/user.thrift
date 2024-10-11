namespace go api.user

include "base.thrift"

struct User {
    1: optional i64 ID (api.body="id");
    2: optional string Name (api.body="name");
    3: optional string StudentNumber (api.body="student_number")
    4: optional string Avatar (api.body="avatar");
    5: optional string PhoneNumber (api.body="phone_number");
    6: optional string PasswordDigest (api.body="password_digest");
}

struct BaseUser {
    1: optional string StudentNumber (api.body="student_number")
    2: optional string Name (api.body="name");
    3: optional string Avatar (api.body="avatar");
    4: optional i64 UID (api.body="uid")
    5: optional i32 Weight (api.body="weight")
}

struct RegisterReq {
    1: optional string Name (api.body="name",api.form="name");
    2: optional string Password (api.body="password",api.form="password");
    3: optional string StudentNumber (api.body="student_number",api.form="student_number")
    4: optional string PhoneNumber (api.body="phone_number",api.form="phone_number");
    5: optional string VerifyCode (api.body="verify_code",api.form="verify_code")
}

struct RegisterResp {
    1: optional base.Base Base (api.body="base");
}

struct GetVerifyCodeReq {
    1: optional string PhoneNumber (api.body="phone_number",api.form="phone_number");
}

struct GetVerifyCodeResp {
    1: optional base.Base Base (api.body="base");
    2: optional string VerifyCode (api.body="verify_code");
}

struct LoginReq {
    1: optional string PhoneNumber (api.body="phone_number",api.form="phone_number");
    2: optional string Password (api.body="password",api.form="password");
}

struct LoginResp {
    1: optional base.Base Base (api.body="base");
    2: optional string Token (api.body="token");
}

struct ChangePasswordReq {
    1: optional i64 UID;
    2: optional string OldPassword (api.body="old_password",api.form="old_password")
    3: optional string NewPassword (api.body="new_passowrd",api.form="new_password")
}

struct ChangePasswordResp {
    1: optional base.Base Base (api.body="base");
}

struct UserInfoReq {
    1: optional i64 UID (api.body="uid");
}

struct UserInfoResp {
    1: optional base.Base Base (api.body="base");
    2: optional i64 UID (api.body="uid");
    3: optional string Name (api.body="name");
    4: optional string StudentNumber (api.body="student_number");
    5: optional string PhoneNumber (api.body="phone_number");
    6: optional i32 Point (api.body="point");
}

service UserService {
     RegisterResp Register(1: RegisterReq req) (api.post="/user/register")
     GetVerifyCodeResp GetVerifyCode(1: GetVerifyCodeReq req) (api.post="/user/get-verifycode")
     LoginResp Login(1: LoginReq req) (api.post="/user/login")
     ChangePasswordResp ChangePassword(1: ChangePasswordReq req) (api.post="/user/auth/change-password")
     UserInfoResp UserInfo(1: UserInfoReq req) (api.get="/user/auth/info")
}
