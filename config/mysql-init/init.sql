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
    PRIMARY KEY pk_user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='班级表' ;

CREATE TABLE IF NOT EXISTS user_with_class (
    uid  BIGINT UNSIGNED COMMENT 'uid',
    class_id  BIGINT UNSIGNED COMMENT 'class_id',
    standing int8 NOT NULL DEFAULT 0 COMMENT '0 学生，1 教师'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户班级关联表' ;


