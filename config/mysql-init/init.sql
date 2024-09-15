CREATE TABLE IF NOT EXISTS user (
    id  BIGINT UNSIGNED COMMENT 'PK',
    name varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称',
    student_number varchar(255) NOT NULL DEFAULT '' COMMENT '用户学号',
    avatar varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像的相对路径',
    phone_number varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像的手机号',
    password_digest varchar(255) NOT NULL DEFAULT '' COMMENT '加密后的密码',
    PRIMARY KEY pk_user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户基础表' ;

CREATE TABLE IF NOT EXISTS class (
    id  BIGINT UNSIGNED COMMENT 'PK',
    name varchar(255) NOT NULL DEFAULT '' COMMENT '班级名称',
    user_count int UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户数量',
    invitation_code varchar(6) NOT NULL DEFAULT ''COMMENT '邀请码',
    PRIMARY KEY pk_user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级表' ;

CREATE TABLE IF NOT EXISTS user_with_class (
    uid  BIGINT UNSIGNED COMMENT 'uid',
    class_id  BIGINT UNSIGNED COMMENT 'class_id',
    weight INT UNSIGNED DEFAULT 1 COMMENT '权重',
    INDEX idx_class_id (class_id),
    INDEX idx_user_id (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户班级关联表' ;


CREATE TABLE IF NOT EXISTS class_owner (
    uid  BIGINT UNSIGNED COMMENT 'uid',
    class_id  BIGINT UNSIGNED COMMENT 'class_id',
    INDEX idx_class_id (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师班级关联表' ;

CREATE TABLE IF NOT EXISTS call_event (
    id  BIGINT UNSIGNED COMMENT 'PK',
    call_event_name varchar(255) NOT NULL DEFAULT '' COMMENT '点名名称',
    class_id  BIGINT UNSIGNED COMMENT 'class_id',
    class_name varchar(255) NOT NULL DEFAULT '' COMMENT 'class_name',
    caller_id  BIGINT UNSIGNED COMMENT 'caller_id',
    caller_name varchar(255) NOT NULL DEFAULT '' COMMENT 'caller_name',
    start_time BIGINT NOT NULL default 0 COMMENT '任务开始时间',
    end_time BIGINT NOT NULL default 0 COMMENT '任务结束时间',
    PRIMARY KEY pk_user(id),
    INDEX idx_class_id (class_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点名表' ;

CREATE TABLE IF NOT EXISTS call_event_with_user (
    call_event_id  BIGINT UNSIGNED COMMENT 'call_event_id',
    uid  BIGINT UNSIGNED COMMENT 'uid',
    done BOOLEAN DEFAULT false COMMENT '0: 未签到, 1: 已签到',
    INDEX idx_call_event_id (call_event_id),
    INDEX idx_user_id (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点名用户关联表' ;

