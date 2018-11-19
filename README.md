## telegrambot

### 技术选型
- go
- [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [mysql](https://github.com/go-sql-driver/mysql)


Feature
==============
> 关键词回复，关键词一般显示为蓝色，且点击会自动复制并发送；(命令本身是这种格式)

> 快捷踢人，在群内发送被禁止的消息会把消息自动删除，并把该用户踢出，并屏蔽；

> 入群问答，不回答10分钟之内就会被提踢，被踢之后可以重新问答并加群；

> 定时发送提示消息。

> 每次入群都需重新回答问题


Bot准备工作
=============
> 申请机器人 （搜索BotFather,发送/newbot，跟着指示完成bot创建，并记录token）

> 设置机器人start显示界面 （/setdescription）

> 设置非私有，用来接收群里发的所有消息 （/setprivacy -> Disable）

> 设置为群管理员，以执行删除消息、踢人、屏蔽(屏蔽必须是超级群组)等功能 (群主手机客户端， 群聊右上角点击头像 -> 点击Edit -> 
    Add Admins -> 取消All Members Are Admins -> 设置机器人为管理员)

Mysql表创建
=============
```mysql
CREATE SCHEMA `telegram` DEFAULT CHARACTER SET utf8mb4 ;

DROP TABLE IF EXISTS  `telegram`.`telegram_group_user`;

CREATE TABLE `telegram`.`telegram_group_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `chat_id` bigint(20) NOT NULL COMMENT '群组ID',
  `is_bot` tinyint(1) NOT NULL DEFAULT '0' COMMENT '用户是否为机器人 0：否 1：是',
  `user_name` varchar(200) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '用户全称',
  `first_name` varchar(100) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `last_name` varchar(100) CHARACTER SET utf8mb4 DEFAULT NULL,
  `language_code` varchar(50) DEFAULT NULL,
  `is_in_group` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:不在 1:在',
  `is_complete_test` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:未完成测试 1:已完成测试',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `group_user_key` (`user_id`,`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

Mysql字符集问题
==============
```
# 由于需要支持火星文，需要修改my.cnf配置
# 增加如下配置
[client]
default-character-set = utf8mb4

[mysql]
default-character-set = utf8mb4

[mysqld]
character-set-client-handshake = FALSE
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init_connect='SET NAMES utf8mb4'
```

Build and Run
============
```go
    go run main.go
```