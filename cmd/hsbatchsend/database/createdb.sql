

-- drop table grant_account;

CREATE TABLE `grant_account`  (
  id int NOT NULL AUTO_INCREMENT,
  address       varchar(64),
  seq           int,
  balance         DECIMAL(36,18),
  created_at    datetime,
  updated_at    datetime,
  deleted_at    datetime,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_address`(`address`) USING BTREE
);



-- 待发送表
-- drop table to_send;
CREATE TABLE `to_send`  (
  id int NOT NULL AUTO_INCREMENT,
  address     varchar(64),
  amount      DECIMAL(36,18),
  status      int COMMENT '状态； 0或者空， 未处理；  1， 处理中，即已插入send_detail，在该表的状态为未发送； 2，已发送，在send_detail 的状态为已发送' ,
  created_at  datetime,
  updated_at  datetime,
  deleted_at  datetime,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_address`(`address`) USING BTREE
);




-- 发送明细表
-- drop table send_detail;
CREATE TABLE `send_detail`  (
  id int NOT NULL AUTO_INCREMENT,
  to_send_id     int  COMMENT 'to_send 待发送表的 id',
  to_address     varchar(64),
  amount         DECIMAL(36,18),
  from_address   varchar(64),
  seq            int,
  status         int  COMMENT '状态； 0或者空， 未发送；  1， 已经发送',
  tx_hash        varchar(128),
  created_at     datetime,
  updated_at     datetime,
  deleted_at     datetime,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_to_send_id`(`to_send_id`) USING BTREE,
  INDEX `idx_to_address`(`to_address`) USING BTREE,
  INDEX `idx_from_address_seq`(`from_address`, `seq`) USING BTREE
);

