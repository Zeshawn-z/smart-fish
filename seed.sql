-- ============================================================================
-- 智钓蓝海 (Smart Fish) 数据库种子数据
-- 数据库: smart_fish (MySQL)
-- 说明: 先清空旧数据再插入，保证可重复执行
-- 使用: mysql -u root -p smart_fish < seed.sql
-- ============================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ========== 清空所有表（按外键依赖顺序） ==========
TRUNCATE TABLE fishing_suggestions;
TRUNCATE TABLE water_quality_data;
TRUNCATE TABLE environment_data;
TRUNCATE TABLE historical_data;
TRUNCATE TABLE reminders;
TRUNCATE TABLE spot_notices;
TRUNCATE TABLE user_favorites;
TRUNCATE TABLE notices;
TRUNCATE TABLE fishing_spots;
TRUNCATE TABLE devices;
TRUNCATE TABLE gateways;
TRUNCATE TABLE regions;
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. 用户 (users)
-- 密码: admin123 / staff123 / user123 / fisher666
-- 注意: 前端登录时会先 SHA256 再发送，所以这里的 bcrypt 哈希是对 SHA256(明文密码) 生成的
-- ============================================================================
INSERT INTO users (id, username, password_hash, role, phone, email, register_time, created_at, updated_at) VALUES
(1, 'admin',    '$2a$12$uGR1UmPT5lAau.ZY.CVnguqMHvoMbthz62UlJEz6kpz3ISiH8el2m', 'admin', '13800000001', 'admin@smartfish.com',    NOW() - INTERVAL 90 DAY, NOW() - INTERVAL 90 DAY, NOW()),
(2, 'staff01',  '$2a$12$LPIwPjWoQModYtxXFttltuCdNQ2y1r3l7zeNzHH6dnOEyuOtxlZne', 'staff', '13800000002', 'staff01@smartfish.com',  NOW() - INTERVAL 60 DAY, NOW() - INTERVAL 60 DAY, NOW()),
(3, 'staff02',  '$2a$12$LPIwPjWoQModYtxXFttltuCdNQ2y1r3l7zeNzHH6dnOEyuOtxlZne', 'staff', '13800000010', 'staff02@smartfish.com',  NOW() - INTERVAL 45 DAY, NOW() - INTERVAL 45 DAY, NOW()),
(4, 'fisher01', '$2a$12$Y/NYrHrcWRc1SEwQD6rbwOns7zSoe0YBiC9.0fiYZmhlOI0yT60FO', 'user',  '13800000003', 'fisher01@smartfish.com', NOW() - INTERVAL 30 DAY, NOW() - INTERVAL 30 DAY, NOW()),
(5, 'fisher02', '$2a$12$Y/NYrHrcWRc1SEwQD6rbwOns7zSoe0YBiC9.0fiYZmhlOI0yT60FO', 'user',  '13800000004', 'fisher02@qq.com',       NOW() - INTERVAL 20 DAY, NOW() - INTERVAL 20 DAY, NOW()),
(6, 'fisher03', '$2a$12$0mE0JXbdPRjHfyrXc2Yv.eDAiX3LPxQGMcxOdAOqfDRxo.gXl/zgi', 'user',  '13800000005', 'fisher03@163.com',      NOW() - INTERVAL 15 DAY, NOW() - INTERVAL 15 DAY, NOW()),
(7, 'fisher04', '$2a$12$0mE0JXbdPRjHfyrXc2Yv.eDAiX3LPxQGMcxOdAOqfDRxo.gXl/zgi', 'user',  '13800000006', 'fisher04@gmail.com',    NOW() - INTERVAL 10 DAY, NOW() - INTERVAL 10 DAY, NOW()),
(8, 'fisher05', '$2a$12$Y/NYrHrcWRc1SEwQD6rbwOns7zSoe0YBiC9.0fiYZmhlOI0yT60FO', 'user',  '13800000007', 'fisher05@smartfish.com', NOW() - INTERVAL 5 DAY,  NOW() - INTERVAL 5 DAY,  NOW()),
(9, 'fisher06', '$2a$12$Y/NYrHrcWRc1SEwQD6rbwOns7zSoe0YBiC9.0fiYZmhlOI0yT60FO', 'user',  '13800000008', 'fisher06@qq.com',       NOW() - INTERVAL 3 DAY,  NOW() - INTERVAL 3 DAY,  NOW()),
(10,'fisher07', '$2a$12$0mE0JXbdPRjHfyrXc2Yv.eDAiX3LPxQGMcxOdAOqfDRxo.gXl/zgi', 'user',  '13800000009', 'fisher07@163.com',      NOW() - INTERVAL 1 DAY,  NOW() - INTERVAL 1 DAY,  NOW());

-- ============================================================================
-- 2. 区域 (regions)
-- ============================================================================
INSERT INTO regions (id, name, province, city, description, created_at, updated_at) VALUES
(1, '哈尔滨市松花江水域',   '黑龙江', '哈尔滨', '松花江哈尔滨段，冬季可冰钓，夏季鱼种丰富',   NOW(), NOW()),
(2, '哈尔滨市太阳岛水域',   '黑龙江', '哈尔滨', '太阳岛周边水域，环境优美，适合休闲垂钓',       NOW(), NOW()),
(3, '大庆市龙凤湿地',       '黑龙江', '大庆',   '龙凤湿地自然保护区周边，生态资源丰富',         NOW(), NOW()),
(4, '镜泊湖垂钓区',         '黑龙江', '牡丹江', '镜泊湖国家级风景名胜区，大型水库鱼种多样',     NOW(), NOW()),
(5, '兴凯湖垂钓区',         '黑龙江', '鸡西',   '中俄界湖兴凯湖，盛产大白鱼、鲤鱼',            NOW(), NOW()),
(6, '齐齐哈尔扎龙湿地水域', '黑龙江', '齐齐哈尔','扎龙国家级自然保护区外围，丹顶鹤之乡',        NOW(), NOW()),
(7, '五大连池水域',         '黑龙江', '黑河',   '五大连池风景区，矿泉鱼闻名遐迩',               NOW(), NOW()),
(8, '佳木斯松花江段',       '黑龙江', '佳木斯', '松花江下游佳木斯段，水面宽阔，鱼群密集',       NOW(), NOW()),
-- 北京
(9,  '北京密云水库水域',     '北京',   '北京',   '密云水库是北京最大水库，水质优良，盛产鲤鱼、鲢鳙、翘嘴', NOW(), NOW()),
(10, '北京怀柔水库水域',     '北京',   '北京',   '怀柔水库周边环境优美，适合休闲垂钓',                       NOW(), NOW()),
-- 河北石家庄
(11, '石家庄岗南水库水域',   '河北',   '石家庄', '岗南水库位于滹沱河上游，是石家庄主要水源地之一，鱼种丰富', NOW(), NOW()),
(12, '石家庄黄壁庄水库水域', '河北',   '石家庄', '黄壁庄水库紧邻岗南水库，水面开阔，鲤鱼草鱼资源丰富',     NOW(), NOW()),
-- 新疆伊犁
(13, '伊犁赛里木湖水域',     '新疆',   '伊犁',   '赛里木湖海拔2073米，高山冷水湖，盛产高白鲑和虹鳟',       NOW(), NOW()),
(14, '伊犁伊犁河水域',       '新疆',   '伊犁',   '伊犁河是新疆水量最大河流，河鲈、鲤鱼资源丰富',           NOW(), NOW());

-- ============================================================================
-- 3. 网关 (gateways)
-- ============================================================================
INSERT INTO gateways (id, name, status, mode, cpu_usage, memory_usage, disk_usage, battery_level, last_active_at, created_at, updated_at) VALUES
(1, '松花江边缘网关-01',   'online',      'online',      35.2, 48.5, 22.1, 87,  NOW(), NOW(), NOW()),
(2, '松花江边缘网关-02',   'online',      'online',      41.8, 52.1, 25.3, 79,  NOW(), NOW(), NOW()),
(3, '太阳岛边缘网关-01',   'online',      'online',      28.7, 42.3, 18.6, 92,  NOW(), NOW(), NOW()),
(4, '龙凤湿地网关-01',     'offline',     'online',      0,    0,    15.3, 15,  NULL, NOW(), NOW()),
(5, '镜泊湖边缘网关-01',   'online',      'online',      22.5, 38.9, 19.8, 95,  NOW(), NOW(), NOW()),
(6, '兴凯湖边缘网关-01',   'online',      'online',      30.1, 45.2, 20.5, 88,  NOW(), NOW(), NOW()),
(7, '扎龙湿地网关-01',     'maintenance', 'maintenance', 0,    0,    30.2, 60,  NOW() - INTERVAL 2 DAY, NOW(), NOW()),
(8, '五大连池网关-01',     'online',      'online',      18.3, 35.6, 16.7, 96,  NOW(), NOW(), NOW()),
-- 北京
(9,  '密云水库网关-01',     'online',      'online',      25.6, 40.8, 18.2, 91,  NOW(), NOW(), NOW()),
(10, '怀柔水库网关-01',     'online',      'online',      22.1, 38.5, 16.5, 94,  NOW(), NOW(), NOW()),
-- 石家庄
(11, '岗南水库网关-01',     'online',      'online',      31.5, 46.2, 21.0, 85,  NOW(), NOW(), NOW()),
(12, '黄壁庄水库网关-01',   'online',      'online',      27.3, 43.1, 19.5, 88,  NOW(), NOW(), NOW()),
-- 伊犁
(13, '赛里木湖网关-01',     'online',      'online',      19.8, 36.7, 15.8, 97,  NOW(), NOW(), NOW()),
(14, '伊犁河网关-01',       'online',      'online',      24.5, 41.3, 17.9, 90,  NOW(), NOW(), NOW());

-- ============================================================================
-- 4. 设备 (devices)
-- ============================================================================
INSERT INTO devices (id, name, gateway_id, status, description, device_type, fishing_count, water_temp, air_temp, humidity, pressure, last_active_at, created_at, updated_at) VALUES
-- 松花江设备群 (网关1,2) — spot1 cap=80, spot2 cap=40, spot3 cap=30
-- 只有绑定水域的设备才有 fishing_count，水下/探鱼设备不计人数
(1,  '松花江环境监测-A1',     1, 'online',  '北岸主环境监测站',     'environment',  7,  8.5,  3.2,  65.3, 1013.25, NOW(), NOW(), NOW()),
(2,  '松花江水下感知-B1',     1, 'online',  '北岸水下探测节点',     'underwater',   0,  7.8,  3.2,  65.3, 1013.25, NOW(), NOW(), NOW()),
(3,  '松花江探鱼辅助-C1',     1, 'online',  '北岸声呐探鱼器',       'fishfinder',   0,  8.2,  3.1,  64.8, 1013.20, NOW(), NOW(), NOW()),
(4,  '松花江环境监测-A2',     2, 'online',  '大桥段环境监测站',     'environment',  5,  8.0,  2.8,  66.1, 1013.10, NOW(), NOW(), NOW()),
(5,  '松花江水下感知-B2',     2, 'online',  '大桥段水下探测节点',   'underwater',   0,  7.5,  2.8,  66.1, 1013.10, NOW(), NOW(), NOW()),
-- 太阳岛设备群 (网关3) — spot4 cap=60, spot5 cap=30
(6,  '太阳岛环境监测-A1',     3, 'online',  '西侧湖主监测站',       'environment',  9,  9.1,  4.5,  62.1, 1012.80, NOW(), NOW(), NOW()),
(7,  '太阳岛水下感知-B1',     3, 'online',  '西侧湖水下节点',       'underwater',   0,  8.9,  4.5,  62.1, 1012.80, NOW(), NOW(), NOW()),
(8,  '太阳岛环境监测-A2',     3, 'online',  '码头区域监测站',       'environment',  4,  9.3,  4.8,  61.5, 1012.75, NOW(), NOW(), NOW()),
-- 龙凤湿地设备群 (网关4 - 离线)
(9,  '龙凤湿地监测-A1',       4, 'offline', '湿地外围监测站',       'environment',  0,  0,    -2.1, 78.5, 1015.10, NULL, NOW(), NOW()),
-- 镜泊湖设备群 (网关5) — spot7 cap=100, spot8 cap=60
(10, '镜泊湖环境监测-A1',     5, 'online',  '大坝下游监测站',       'environment',  12, 6.8,  1.5,  70.2, 1014.50, NOW(), NOW(), NOW()),
(11, '镜泊湖水下感知-B1',     5, 'online',  '大坝深水区探测',       'underwater',   0,  6.2,  1.5,  70.2, 1014.50, NOW(), NOW(), NOW()),
(12, '镜泊湖探鱼辅助-C1',     5, 'online',  '大坝声呐系统',         'fishfinder',   0,  6.5,  1.3,  71.0, 1014.55, NOW(), NOW(), NOW()),
-- 兴凯湖设备群 (网关6) — spot9 cap=120
(13, '兴凯湖环境监测-A1',     6, 'online',  '北岸监测站',           'environment',  15, 5.5,  0.8,  72.8, 1015.20, NOW(), NOW(), NOW()),
(14, '兴凯湖水下感知-B1',     6, 'online',  '湖心区水下探测',       'underwater',   0,  5.2,  0.8,  72.8, 1015.20, NOW(), NOW(), NOW()),
-- 扎龙设备 (网关7 - 维护中)
(15, '扎龙湿地监测-A1',       7, 'error',   '设备固件升级中',       'environment',  0,  0,    -3.5, 80.1, 1016.00, NOW() - INTERVAL 2 DAY, NOW(), NOW()),
-- 五大连池设备 (网关8) — spot11 cap=45
(16, '五大连池环境监测-A1',   8, 'online',  '药泉湖监测站',         'environment',  3,  4.2,  -1.2, 68.9, 1014.80, NOW(), NOW(), NOW()),
-- 北京密云设备 (网关9)
(17, '密云水库环境监测-A1',   9, 'online',  '密云水库大坝监测站',   'environment',  18, 12.5, 15.3, 55.2, 1016.30, NOW(), NOW(), NOW()),
(18, '密云水库水下感知-B1',   9, 'online',  '密云水库深水区探测',   'underwater',   0,  10.8, 15.3, 55.2, 1016.30, NOW(), NOW(), NOW()),
-- 北京怀柔设备 (网关10)
(19, '怀柔水库环境监测-A1',  10, 'online',  '怀柔水库西岸监测站',   'environment',  10, 13.1, 16.0, 52.8, 1016.10, NOW(), NOW(), NOW()),
-- 石家庄岗南设备 (网关11)
(20, '岗南水库环境监测-A1',  11, 'online',  '岗南水库坝下监测站',   'environment',  14, 14.2, 18.5, 48.6, 1015.50, NOW(), NOW(), NOW()),
(21, '岗南水库水下感知-B1',  11, 'online',  '岗南水库深水区探测',   'underwater',   0,  12.8, 18.5, 48.6, 1015.50, NOW(), NOW(), NOW()),
-- 石家庄黄壁庄设备 (网关12)
(22, '黄壁庄水库环境监测-A1',12, 'online',  '黄壁庄水库北岸监测站', 'environment',  8,  14.8, 19.0, 46.5, 1015.20, NOW(), NOW(), NOW()),
-- 伊犁赛里木湖设备 (网关13)
(23, '赛里木湖环境监测-A1',  13, 'online',  '赛里木湖南岸监测站',   'environment',  6,  8.0,  10.2, 35.8, 1018.50, NOW(), NOW(), NOW()),
(24, '赛里木湖水下感知-B1',  13, 'online',  '赛里木湖冷水鱼探测',   'underwater',   0,  6.5,  10.2, 35.8, 1018.50, NOW(), NOW(), NOW()),
-- 伊犁河设备 (网关14)
(25, '伊犁河环境监测-A1',    14, 'online',  '伊犁河伊宁段监测站',   'environment',  11, 11.5, 14.8, 40.2, 1017.00, NOW(), NOW(), NOW());

-- ============================================================================
-- 5. 垂钓水域 (fishing_spots)
-- ============================================================================
INSERT INTO fishing_spots (id, name, region_id, description, latitude, longitude, water_type, capacity, status, bound_device_id, created_at, updated_at) VALUES
(1,  '松花江北岸钓场',         1, '松花江北岸，适合路亚和台钓，冬季冰钓热门地点',                     45.7772, 126.6177, 'river',     80,  'open',        1,    NOW(), NOW()),
(2,  '松花江公路大桥钓位',     1, '公路大桥下游，深水区鱼种丰富，常有大鲤鱼出没',                     45.7685, 126.6089, 'river',     40,  'open',        4,    NOW(), NOW()),
(3,  '松花江铁路桥野钓区',     1, '铁路桥附近自然钓场，水流较急，适合有经验的钓友',                   45.7710, 126.6250, 'river',     30,  'open',        5,    NOW(), NOW()),
(4,  '太阳岛西侧湖钓场',       2, '太阳岛西侧静水湖，水质清澈，适合休闲垂钓，环境宜人',             45.7891, 126.5812, 'lake',      60,  'open',        6,    NOW(), NOW()),
(5,  '太阳岛码头钓位',         2, '老码头附近，水深2-4米，鲫鱼鲤鱼密度高',                           45.7856, 126.5923, 'river',     30,  'open',        8,    NOW(), NOW()),
(6,  '龙凤湿地外围钓场',       3, '湿地外围开放区域，冬季关闭',                                       46.5812, 125.0934, 'lake',      50,  'closed',      9,    NOW(), NOW()),
(7,  '镜泊湖大坝钓场',         4, '镜泊湖大坝下游，大鱼频出，专业钓手天堂',                           44.0112, 128.9945, 'reservoir', 100, 'open',        10,   NOW(), NOW()),
(8,  '镜泊湖吊水楼瀑布钓区',   4, '吊水楼瀑布下方深潭，水流交汇处鱼群聚集',                         43.9856, 129.0123, 'reservoir', 60,  'open',        11,   NOW(), NOW()),
(9,  '兴凯湖北岸垂钓基地',     5, '兴凯湖北岸专业垂钓基地，盛产大白鱼、鲤鱼、鲢鳙',               45.3245, 132.4567, 'lake',      120, 'open',        13,   NOW(), NOW()),
(10, '扎龙湿地观鸟钓场',       6, '扎龙保护区外围，可边观鸟边垂钓，别有趣味',                         47.1789, 124.2856, 'pond',      40,  'maintenance', 15,   NOW(), NOW()),
(11, '五大连池药泉湖钓场',     7, '药泉湖垂钓区，矿泉水养殖鱼口感鲜美',                             48.7523, 126.1189, 'lake',      45,  'open',        16,   NOW(), NOW()),
(12, '佳木斯松花江外滩钓场',   8, '佳木斯外滩公园段，休闲垂钓好去处，交通便利',                     46.8012, 130.3678, 'river',     70,  'open',        NULL, NOW(), NOW()),
-- 北京 (region 9=密云, 10=怀柔)
(13, '密云水库大坝钓场',       9, '密云水库大坝下游，水深6-15米，大鱼频出，北京钓友圣地',           40.5280, 116.9750, 'reservoir', 100, 'open',        17,   NOW(), NOW()),
(14, '密云水库白河湾钓场',     9, '白河入库口，水流交汇处鱼群密集，适合路亚和台钓',               40.4950, 116.9200, 'reservoir', 60,  'open',        18,   NOW(), NOW()),
(15, '怀柔水库西坝钓场',      10, '怀柔水库西侧安静水域，周末钓友云集，环境清幽',                 40.3560, 116.6280, 'reservoir', 80,  'open',        19,   NOW(), NOW()),
-- 石家庄 (region 11=岗南, 12=黄壁庄)
(16, '岗南水库坝下钓场',      11, '岗南水库主坝下方，水深适中，鲤鱼草鱼活跃',                     38.3520, 114.1680, 'reservoir', 90,  'open',        20,   NOW(), NOW()),
(17, '岗南水库苇塘钓区',      11, '水库东侧苇塘区域，鲫鱼密度高，休闲台钓首选',                 38.3800, 114.2100, 'reservoir', 50,  'open',        21,   NOW(), NOW()),
(18, '黄壁庄水库北岸钓场',    12, '黄壁庄水库北岸开阔地带，水面宽广，适合抛竿远投',             38.2850, 114.0650, 'reservoir', 70,  'open',        22,   NOW(), NOW()),
-- 伊犁 (region 13=赛里木湖, 14=伊犁河)
(19, '赛里木湖南岸钓场',      13, '赛里木湖南岸，高山冷水湖，盛产高白鲑和虹鳟，景色绝美',     44.5800, 81.1600,  'lake',      40,  'open',        23,   NOW(), NOW()),
(20, '赛里木湖东岸野钓区',    13, '赛里木湖东岸原生态区域，人少鱼多，徒步可达',                 44.6200, 81.2500,  'lake',      30,  'open',        24,   NOW(), NOW()),
(21, '伊犁河伊宁段钓场',      14, '伊犁河流经伊宁市区段，河鲈和鲤鱼资源丰富，交通便利',       43.9150, 81.3200,  'river',     60,  'open',        25,   NOW(), NOW());

-- ============================================================================
-- 6. 历史数据 (historical_data) — 最近48小时，每小时一条
-- ============================================================================
-- 用存储过程批量生成
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS seed_historical_data()
BEGIN
    DECLARE i INT DEFAULT 47;
    DECLARE sid INT;
    DECLARE base_count INT;
    DECLARE hour_count INT;
    DECLARE spot_cap INT;
    DECLARE ts DATETIME;

    WHILE i >= 0 DO
        SET ts = NOW() - INTERVAL i HOUR;

        -- 给所有 open 状态的水域生成数据
        SET sid = 1;
        WHILE sid <= 21 DO
            -- 跳过 closed(6) 和 maintenance(10) 的水域
            IF sid NOT IN (6, 10) THEN
                -- 获取该水域容量
                SET spot_cap = CASE sid
                    WHEN 1 THEN 80  WHEN 2 THEN 40  WHEN 3 THEN 30
                    WHEN 4 THEN 60  WHEN 5 THEN 30  WHEN 7 THEN 100
                    WHEN 8 THEN 60  WHEN 9 THEN 120 WHEN 11 THEN 45
                    WHEN 12 THEN 70
                    WHEN 13 THEN 100 WHEN 14 THEN 60  WHEN 15 THEN 80
                    WHEN 16 THEN 90  WHEN 17 THEN 50  WHEN 18 THEN 70
                    WHEN 19 THEN 40  WHEN 20 THEN 30  WHEN 21 THEN 60
                    ELSE 50
                END;

                -- 模拟白天人多、夜晚人少的规律（比例控制偏低，更贴近实际）
                SET base_count = CASE
                    WHEN HOUR(ts) BETWEEN 5 AND 7   THEN FLOOR(spot_cap * 0.05) + (sid * 2) % GREATEST(1, FLOOR(spot_cap * 0.03))
                    WHEN HOUR(ts) BETWEEN 8 AND 10  THEN FLOOR(spot_cap * 0.12) + (sid * 3) % GREATEST(1, FLOOR(spot_cap * 0.05))
                    WHEN HOUR(ts) BETWEEN 11 AND 13 THEN FLOOR(spot_cap * 0.08) + (sid * 2) % GREATEST(1, FLOOR(spot_cap * 0.04))
                    WHEN HOUR(ts) BETWEEN 14 AND 17 THEN FLOOR(spot_cap * 0.15) + (sid * 3) % GREATEST(1, FLOOR(spot_cap * 0.06))
                    WHEN HOUR(ts) BETWEEN 18 AND 20 THEN FLOOR(spot_cap * 0.06) + (sid * 2) % GREATEST(1, FLOOR(spot_cap * 0.03))
                    ELSE FLOOR(spot_cap * 0.01) + sid % 2
                END;

                -- 小幅扰动（±3 以内），确保不超过容量且不低于 0
                SET hour_count = GREATEST(0, LEAST(spot_cap, base_count + (i * 7 + sid * 13) % 7 - 3));

                INSERT INTO historical_data (spot_id, fishing_count, timestamp, created_at, updated_at)
                VALUES (sid, hour_count, ts, ts, ts);
            END IF;
            SET sid = sid + 1;
        END WHILE;

        SET i = i - 1;
    END WHILE;
END//
DELIMITER ;

CALL seed_historical_data();
DROP PROCEDURE IF EXISTS seed_historical_data;

-- ============================================================================
-- 7. 环境数据 (environment_data) — 最近48小时
-- ============================================================================
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS seed_environment_data()
BEGIN
    DECLARE i INT DEFAULT 47;
    DECLARE spot_id INT;
    DECLARE ts DATETIME;
    DECLARE base_water_temp DOUBLE;
    DECLARE base_air_temp DOUBLE;
    DECLARE base_humidity DOUBLE;
    DECLARE base_pressure DOUBLE;

    WHILE i >= 0 DO
        SET ts = NOW() - INTERVAL i HOUR;

        SET spot_id = 1;
        WHILE spot_id <= 21 DO
            IF spot_id NOT IN (6, 10) THEN
                -- 根据水域位置和时间生成不同的基础温度
                SET base_water_temp = 5.0 + spot_id * 0.4 + SIN(i * 0.26) * 1.5;
                SET base_air_temp = base_water_temp - 4.0 + COS(i * 0.26) * 2.0;
                SET base_humidity = 58.0 + (spot_id * 7 + i * 3) % 25;
                SET base_pressure = 1010.0 + SIN(i * 0.13) * 3.0 + spot_id * 0.3;

                INSERT INTO environment_data (spot_id, water_temp, air_temp, humidity, pressure, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
                VALUES (
                    spot_id,
                    ROUND(base_water_temp, 1),
                    ROUND(base_air_temp, 1),
                    ROUND(base_humidity, 1),
                    ROUND(base_pressure, 2),
                    ROUND(6.8 + (spot_id + i) % 8 * 0.1, 1),
                    ROUND(7.5 + SIN(i * 0.3 + spot_id) * 1.5, 1),
                    ROUND(10.0 + (i * 3 + spot_id * 5) % 20, 1),
                    ts, ts, ts
                );
            END IF;
            SET spot_id = spot_id + 1;
        END WHILE;

        SET i = i - 1;
    END WHILE;
END//
DELIMITER ;

CALL seed_environment_data();
DROP PROCEDURE IF EXISTS seed_environment_data;

-- ============================================================================
-- 8. 水质数据 (water_quality_data) — 最近24小时，给有水下设备的写
-- ============================================================================
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS seed_water_quality_data()
BEGIN
    DECLARE i INT DEFAULT 23;
    DECLARE ts DATETIME;

    WHILE i >= 0 DO
        SET ts = NOW() - INTERVAL i HOUR;

        -- 设备2: 松花江水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (2, ROUND(7.0 + (i % 5) * 0.1, 1), ROUND(8.2 + SIN(i * 0.5) * 0.8, 1), ROUND(12.0 + (i * 3) % 8, 1), ts, ts, ts);

        -- 设备5: 松花江水下感知-B2
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (5, ROUND(7.1 + (i % 4) * 0.1, 1), ROUND(8.0 + SIN(i * 0.4) * 1.0, 1), ROUND(14.0 + (i * 5) % 10, 1), ts, ts, ts);

        -- 设备7: 太阳岛水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (7, ROUND(7.2 + (i % 6) * 0.08, 1), ROUND(8.5 + COS(i * 0.3) * 0.6, 1), ROUND(8.0 + (i * 2) % 7, 1), ts, ts, ts);

        -- 设备11: 镜泊湖水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (11, ROUND(7.3 + (i % 3) * 0.12, 1), ROUND(9.0 + SIN(i * 0.6) * 0.5, 1), ROUND(6.0 + (i * 4) % 9, 1), ts, ts, ts);

        -- 设备14: 兴凯湖水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (14, ROUND(7.0 + (i % 5) * 0.15, 1), ROUND(8.8 + COS(i * 0.5) * 0.7, 1), ROUND(10.0 + (i * 3) % 12, 1), ts, ts, ts);

        -- 设备18: 密云水库水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (18, ROUND(7.1 + (i % 4) * 0.12, 1), ROUND(8.6 + SIN(i * 0.4) * 0.9, 1), ROUND(9.0 + (i * 2) % 10, 1), ts, ts, ts);

        -- 设备21: 岗南水库水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (21, ROUND(7.2 + (i % 5) * 0.10, 1), ROUND(8.3 + COS(i * 0.3) * 0.8, 1), ROUND(11.0 + (i * 4) % 8, 1), ts, ts, ts);

        -- 设备24: 赛里木湖水下感知-B1
        INSERT INTO water_quality_data (device_id, ph, dissolved_oxygen, turbidity, timestamp, created_at, updated_at)
        VALUES (24, ROUND(7.4 + (i % 3) * 0.08, 1), ROUND(9.2 + SIN(i * 0.6) * 0.4, 1), ROUND(5.0 + (i * 2) % 6, 1), ts, ts, ts);

        SET i = i - 1;
    END WHILE;
END//
DELIMITER ;

CALL seed_water_quality_data();
DROP PROCEDURE IF EXISTS seed_water_quality_data;

-- ============================================================================
-- 9. 提醒 (reminders)
-- ============================================================================
INSERT INTO reminders (spot_id, level, reminder_type, message, timestamp, resolved, publicity, created_at, updated_at) VALUES
-- 天气类
(1,  1, 'weather',     '明日有大风预警，松花江北岸风力可达5-6级，建议做好防风准备',                   NOW() - INTERVAL 2 HOUR,  0, 1, NOW(), NOW()),
(7,  1, 'weather',     '镜泊湖区域未来3天有降雪预报，气温将降至-15℃以下，请注意防寒保暖',           NOW() - INTERVAL 5 HOUR,  0, 1, NOW(), NOW()),
(9,  0, 'weather',     '兴凯湖今日天气晴好，风力2-3级，非常适合出钓',                                 NOW() - INTERVAL 1 HOUR,  0, 1, NOW(), NOW()),
-- 垂钓类
(4,  0, 'fishing',     '太阳岛西侧湖今日适宜垂钓，水温适中，鱼口较好，建议使用红虫',               NOW() - INTERVAL 1 HOUR,  0, 1, NOW(), NOW()),
(1,  0, 'fishing',     '松花江北岸近期鲫鱼活跃，早6点-8点上鱼率最高',                                 NOW() - INTERVAL 3 HOUR,  0, 1, NOW(), NOW()),
(8,  0, 'fishing',     '镜泊湖吊水楼钓区今晨有钓友钓获12斤鲤鱼，鱼情大好',                         NOW() - INTERVAL 30 MINUTE, 0, 1, NOW(), NOW()),
-- 安全类
(6,  3, 'safety',      '龙凤湿地外围钓场因冰面不稳定暂时关闭，严禁进入！',                           NOW() - INTERVAL 30 MINUTE, 0, 1, NOW(), NOW()),
(3,  2, 'safety',      '松花江铁路桥野钓区水流较急，近日有险情上报，请佩戴救生衣',                   NOW() - INTERVAL 4 HOUR,  0, 1, NOW(), NOW()),
(10, 2, 'safety',      '扎龙湿地观鸟钓场设备维护中，暂不开放',                                       NOW() - INTERVAL 6 HOUR,  0, 1, NOW(), NOW()),
-- 环境类
(2,  2, 'environment', '松花江公路大桥段水质浑浊度升高至25NTU，可能影响钓获',                         NOW(), 0, 1, NOW(), NOW()),
(11, 1, 'environment', '五大连池药泉湖 pH 值偏高(8.2)，建议关注水质变化',                             NOW() - INTERVAL 8 HOUR,  0, 1, NOW(), NOW()),
-- 已解决的历史提醒
(1,  1, 'weather',     '昨日大风预警已解除，松花江北岸恢复正常垂钓',                                 NOW() - INTERVAL 26 HOUR, 1, 1, NOW(), NOW()),
(4,  2, 'environment', '太阳岛西侧湖溶解氧恢复正常水平',                                             NOW() - INTERVAL 20 HOUR, 1, 1, NOW(), NOW()),
-- 北京
(13, 0, 'fishing',     '密云水库大坝钓场近期鲤鱼活跃，水温12.5℃适宜垂钓，建议使用玉米粒打窝',       NOW() - INTERVAL 2 HOUR,  0, 1, NOW(), NOW()),
(15, 1, 'weather',     '怀柔水库周末有阵雨预报，请携带雨具，注意防滑',                               NOW() - INTERVAL 4 HOUR,  0, 1, NOW(), NOW()),
-- 石家庄
(16, 0, 'fishing',     '岗南水库坝下草鱼开口良好，推荐使用嫩玉米挂钩，钓3-4米深',                   NOW() - INTERVAL 1 HOUR,  0, 1, NOW(), NOW()),
(18, 1, 'environment', '黄壁庄水库北岸近日水位下降约0.5米，钓位需前移',                               NOW() - INTERVAL 6 HOUR,  0, 1, NOW(), NOW()),
-- 伊犁
(19, 0, 'fishing',     '赛里木湖南岸高白鲑活跃，水温8℃，路亚银色勺形亮片5g效果最佳',               NOW() - INTERVAL 3 HOUR,  0, 1, NOW(), NOW()),
(21, 1, 'weather',     '伊犁河伊宁段上游有融雪水汇入，水温偏低且水流加大，注意安全',               NOW() - INTERVAL 5 HOUR,  0, 1, NOW(), NOW());

-- ============================================================================
-- 10. 通知 (notices)
-- ============================================================================
INSERT INTO notices (id, title, content, timestamp, outdated, created_at, updated_at) VALUES
(1, '2026年冬季冰钓节活动通知',
    '哈尔滨市第十届冬季冰钓节将于2026年1月15日在松花江北岸举办，欢迎广大钓友参加！活动设有个人赛、团队赛等多个项目，总奖金池10万元。报名截止日期为2026年1月10日。',
    NOW() - INTERVAL 72 HOUR, 0, NOW(), NOW()),
(2, '垂钓安全须知更新',
    '各位钓友注意：冬季垂钓务必注意冰面安全，请在标识安全区域内活动。携带救生装备，结伴出行。设备有异常及时上报管理员。遇紧急情况请拨打平台热线 400-888-FISH。',
    NOW() - INTERVAL 48 HOUR, 0, NOW(), NOW()),
(3, '系统维护通知',
    '智钓蓝海平台将于本周日凌晨2:00-4:00进行系统升级维护，届时部分功能可能暂时不可用，敬请谅解。本次更新将优化水质监测数据的实时展示。',
    NOW() - INTERVAL 12 HOUR, 0, NOW(), NOW()),
(4, '新增兴凯湖垂钓基地',
    '好消息！智钓蓝海平台新增兴凯湖北岸垂钓基地，配备环境监测和水下探测设备，实时掌握鱼情动态。兴凯湖盛产大白鱼，欢迎钓友们前去体验！',
    NOW() - INTERVAL 24 HOUR, 0, NOW(), NOW()),
(5, '春季禁渔期提醒',
    '根据黑龙江省渔业管理规定，松花江流域每年4月1日至6月15日为禁渔期，届时相关水域将暂停垂钓服务。请钓友们提前安排行程，也可选择非禁渔水域。',
    NOW() - INTERVAL 6 HOUR, 0, NOW(), NOW()),
(6, '【已过期】元旦活动已结束',
    '2026年元旦垂钓嘉年华活动已圆满结束，感谢各位钓友的热情参与！获奖名单已在公告栏公布。',
    NOW() - INTERVAL 30 DAY, 1, NOW(), NOW()),
(7, '新增北京密云水库、怀柔水库钓场',
    '好消息！智钓蓝海平台新增北京密云水库和怀柔水库两大钓场，已部署环境监测和水下探测设备。密云水库盛产大鲤鱼，怀柔水库环境优美，欢迎北京钓友前往体验！',
    NOW() - INTERVAL 8 HOUR, 0, NOW(), NOW()),
(8, '新增石家庄岗南水库、黄壁庄水库钓场',
    '河北石家庄地区的岗南水库和黄壁庄水库钓场正式上线！两大水库水面开阔，鲤鱼草鱼资源丰富，配备全套智能监测设备。',
    NOW() - INTERVAL 10 HOUR, 0, NOW(), NOW()),
(9, '新增新疆伊犁赛里木湖、伊犁河钓场',
    '伊犁赛里木湖高山冷水钓场和伊犁河钓场已接入平台。赛里木湖海拔2073米，盛产高白鲑和虹鳟，是路亚爱好者的天堂。伊犁河段河鲈资源丰富，欢迎体验！',
    NOW() - INTERVAL 9 HOUR, 0, NOW(), NOW());

-- ============================================================================
-- 11. 通知关联水域 (spot_notices)
-- ============================================================================
INSERT INTO spot_notices (notice_id, fishing_spot_id) VALUES
-- 冰钓节关联松花江水域
(1, 1), (1, 2), (1, 3),
-- 安全须知关联所有开放水域
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 7), (2, 8), (2, 9), (2, 11), (2, 12),
(2, 13), (2, 14), (2, 15), (2, 16), (2, 17), (2, 18), (2, 19), (2, 20), (2, 21),
-- 系统维护通知不关联特定水域（全局）
-- 兴凯湖新增通知
(4, 9),
-- 禁渔期关联松花江水域
(5, 1), (5, 2), (5, 3), (5, 12),
-- 北京钓场通知
(7, 13), (7, 14), (7, 15),
-- 石家庄钓场通知
(8, 16), (8, 17), (8, 18),
-- 伊犁钓场通知
(9, 19), (9, 20), (9, 21);

-- ============================================================================
-- 12. 用户收藏 (user_favorites)
-- ============================================================================
INSERT INTO user_favorites (user_id, fishing_spot_id) VALUES
(4,  1), (4,  4), (4,  7),   -- fisher01 收藏松花江北岸、太阳岛西侧湖、镜泊湖大坝
(5,  1), (5,  9),             -- fisher02 收藏松花江北岸、兴凯湖
(6,  4), (6,  5), (6,  11),  -- fisher03 收藏太阳岛两个、五大连池
(7,  7), (7,  8), (7,  9),   -- fisher04 收藏镜泊湖两个、兴凯湖
(8,  1), (8,  2), (8,  3),   -- fisher05 收藏松花江三个
(9,  4), (9,  12),            -- fisher06 收藏太阳岛、佳木斯
(10, 7), (10, 9),             -- fisher07 收藏镜泊湖、兴凯湖
(4, 13), (5, 16), (7, 19);   -- fisher01收藏密云, fisher02收藏岗南, fisher04收藏赛里木湖

-- ============================================================================
-- 13. 垂钓建议 (fishing_suggestions)
-- ============================================================================
INSERT INTO fishing_suggestions (spot_id, user_id, suggestion_text, score, timestamp, created_at, updated_at) VALUES
(1,  4,    '当前松花江北岸水温8.5℃，建议使用红虫或蚯蚓做饵，主攻鲫鱼。早晨6-8点和傍晚4-6点是最佳垂钓时段。1.5号子线配4号袖钩效果最佳。',
    85.5, NOW(), NOW(), NOW()),
(4,  NULL, '太阳岛西侧湖水深适中，水质清澈，当前溶氧充足。推荐使用3.6米手竿，钓底层，搓饵料球。近期鲤鱼活跃度较高，建议用玉米粒打窝。',
    78.2, NOW() - INTERVAL 2 HOUR, NOW(), NOW()),
(7,  NULL, '镜泊湖大坝钓场水深8-12米，建议使用矶钓竿或6.3米以上长竿。深水层水温稳定，鲢鳙和大鲤鱼活跃。推荐使用发酵饵料配合雾化饵。',
    92.1, NOW() - INTERVAL 1 HOUR, NOW(), NOW()),
(9,  7,    '兴凯湖北岸当前白鱼群已开始活跃，水面温度5.5℃。推荐使用路亚拟饵在浅水区抛投，勺型亮片3-5克效果好。日出后1小时为最佳窗口。',
    88.7, NOW() - INTERVAL 3 HOUR, NOW(), NOW()),
(2,  5,    '松花江大桥段近期水质浑浊度偏高，建议使用味道浓烈的饵料（腥饵为主），钓位选择缓流区。2号主线配1号子线，伊势尼5号钩。',
    72.3, NOW() - INTERVAL 5 HOUR, NOW(), NOW()),
(5,  NULL, '太阳岛码头水深2-4米，底部淤泥较厚，建议调钓灵敏（调4钓2）。近期鲫鱼为主，偶有黄辣丁。红虫拉饵效果优于商品饵。',
    80.5, NOW() - INTERVAL 4 HOUR, NOW(), NOW()),
(11, 6,    '五大连池药泉湖矿物质含量高，鱼口偏轻。建议使用0.6号子线配3号金袖，浮漂选择吃铅小的芦苇漂，调目不宜过高。',
    76.8, NOW() - INTERVAL 6 HOUR, NOW(), NOW()),
(8,  NULL, '镜泊湖吊水楼瀑布钓区水流交汇处，大物频出。今日已有多名钓友中获5斤以上鲤鱼。建议加粗线组（3号主线+1.5号子线），准备抄网。',
    90.3, NOW() - INTERVAL 30 MINUTE, NOW(), NOW()),
(12, NULL, '佳木斯松花江段目前尚未绑定监测设备，暂无实时环境数据。根据历史数据，该段3月下旬开始回暖，鲫鱼和鲤鱼陆续开口。',
    65.0, NOW() - INTERVAL 12 HOUR, NOW(), NOW()),
(1,  6,    '综合近48小时数据分析：松花江北岸水温呈缓慢上升趋势，溶解氧充足，气压稳定。预计明日垂钓条件评分将达90+，强烈推荐出钓。',
    95.0, NOW() - INTERVAL 10 MINUTE, NOW(), NOW()),
-- 北京
(13, NULL, '密云水库大坝钓场当前水温12.5℃，春季回暖鲤鱼开口积极。建议使用4.5米手竿，主线2号配子线1号，伊势尼6号钩。搓饵或玉米粒打窝效果俱佳。',
    88.0, NOW() - INTERVAL 2 HOUR, NOW(), NOW()),
(14, 5,    '密云水库白河湾入水口处，路亚翘嘴效果极佳。推荐使用7-10g银色米诺，匀速收线带小抽。清晨6-8点为最佳窗口期。',
    82.5, NOW() - INTERVAL 4 HOUR, NOW(), NOW()),
(15, NULL, '怀柔水库西坝水深3-6米，底部碎石底质，适合钓鲫鱼。推荐3.6米短竿，0.8号子线配3号袖钩，红虫或蚯蚓拉饵。调4钓2较为灵敏。',
    79.3, NOW() - INTERVAL 6 HOUR, NOW(), NOW()),
-- 石家庄
(16, 4,    '岗南水库坝下当前水温14.2℃，草鱼进入活跃期。推荐5.4米竿，3号主线配1.5号子线，使用嫩玉米或青草做饵。钓3-4米深效果最好。',
    91.2, NOW() - INTERVAL 1 HOUR, NOW(), NOW()),
(18, NULL, '黄壁庄水库北岸近日水位略降，鱼群向深水区移动。建议使用6.3米以上长竿远投，或矶钓竿配合串钩。饵料以腥香为主，雾化饵打窝聚鱼。',
    75.8, NOW() - INTERVAL 8 HOUR, NOW(), NOW()),
-- 伊犁
(19, NULL, '赛里木湖南岸高白鲑活跃，水温8℃适合冷水鱼觅食。路亚推荐5g银色勺形亮片，慢速匀收。注意高原紫外线强烈，做好防晒。海拔2073米，体力消耗比平原大。',
    86.5, NOW() - INTERVAL 3 HOUR, NOW(), NOW()),
(21, 7,    '伊犁河伊宁段河鲈活跃，当前水温11.5℃。路亚推荐使用T尾软虫配铅头钩（7g），沿河岸缓流区搜索。台钓推荐蚯蚓挂钩钓底，主攻鲤鱼和鲫鱼。',
    83.7, NOW() - INTERVAL 5 HOUR, NOW(), NOW());

-- ============================================================================
-- 完成
-- ============================================================================
SELECT '✅ 种子数据插入完成！' AS status;
SELECT '用户' AS `表`, COUNT(*) AS `数量` FROM users
UNION ALL SELECT '区域', COUNT(*) FROM regions
UNION ALL SELECT '网关', COUNT(*) FROM gateways
UNION ALL SELECT '设备', COUNT(*) FROM devices
UNION ALL SELECT '垂钓水域', COUNT(*) FROM fishing_spots
UNION ALL SELECT '历史数据', COUNT(*) FROM historical_data
UNION ALL SELECT '环境数据', COUNT(*) FROM environment_data
UNION ALL SELECT '水质数据', COUNT(*) FROM water_quality_data
UNION ALL SELECT '提醒', COUNT(*) FROM reminders
UNION ALL SELECT '通知', COUNT(*) FROM notices
UNION ALL SELECT '垂钓建议', COUNT(*) FROM fishing_suggestions
UNION ALL SELECT '用户收藏', COUNT(*) FROM user_favorites
UNION ALL SELECT '通知-水域关联', COUNT(*) FROM spot_notices;
