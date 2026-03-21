package database

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"time"

	"smart-fish/back_end/models"

	"golang.org/x/crypto/bcrypt"
)

// sha256Hex 模拟前端 CryptoJS.SHA256(password).toString() 的行为
func sha256Hex(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

// Seed 填充示例数据（开发环境使用）
func Seed() {
	// 检查是否已有数据
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already has base data, checking SFR data...")
		seedSFRIfNeeded()
		return
	}

	log.Println("Seeding database with sample data...")

	now := time.Now()

	// ===== 1. 创建用户 =====
	// 前端登录时会先 SHA256 再发给后端，所以 seed 也要对 SHA256 后的值做 bcrypt
	adminHash, _ := bcrypt.GenerateFromPassword([]byte(sha256Hex("admin123")), bcrypt.DefaultCost)
	staffHash, _ := bcrypt.GenerateFromPassword([]byte(sha256Hex("staff123")), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte(sha256Hex("user123")), bcrypt.DefaultCost)
	fisherHash, _ := bcrypt.GenerateFromPassword([]byte(sha256Hex("fisher666")), bcrypt.DefaultCost)

	users := []models.User{
		{Username: "admin", PasswordHash: string(adminHash), Role: "admin", Phone: "13800000001", Email: "admin@smartfish.com"},
		{Username: "staff01", PasswordHash: string(staffHash), Role: "staff", Phone: "13800000002", Email: "staff01@smartfish.com"},
		{Username: "staff02", PasswordHash: string(staffHash), Role: "staff", Phone: "13800000010", Email: "staff02@smartfish.com"},
		{Username: "fisher01", PasswordHash: string(userHash), Role: "user", Phone: "13800000003", Email: "fisher01@smartfish.com"},
		{Username: "fisher02", PasswordHash: string(userHash), Role: "user", Phone: "13800000004", Email: "fisher02@qq.com"},
		{Username: "fisher03", PasswordHash: string(fisherHash), Role: "user", Phone: "13800000005", Email: "fisher03@163.com"},
		{Username: "fisher04", PasswordHash: string(fisherHash), Role: "user", Phone: "13800000006", Email: "fisher04@gmail.com"},
		{Username: "fisher05", PasswordHash: string(userHash), Role: "user", Phone: "13800000007", Email: "fisher05@smartfish.com"},
		{Username: "fisher06", PasswordHash: string(userHash), Role: "user", Phone: "13800000008", Email: "fisher06@qq.com"},
		{Username: "fisher07", PasswordHash: string(fisherHash), Role: "user", Phone: "13800000009", Email: "fisher07@163.com"},
	}
	DB.Create(&users)

	// ===== 2. 创建区域 =====
	regions := []models.Region{
		{Name: "哈尔滨市松花江水域", Province: "黑龙江", City: "哈尔滨", Description: "松花江哈尔滨段，冬季可冰钓，夏季鱼种丰富"},
		{Name: "哈尔滨市太阳岛水域", Province: "黑龙江", City: "哈尔滨", Description: "太阳岛周边水域，环境优美，适合休闲垂钓"},
		{Name: "大庆市龙凤湿地", Province: "黑龙江", City: "大庆", Description: "龙凤湿地自然保护区周边，生态资源丰富"},
		{Name: "镜泊湖垂钓区", Province: "黑龙江", City: "牡丹江", Description: "镜泊湖国家级风景名胜区，大型水库鱼种多样"},
		{Name: "兴凯湖垂钓区", Province: "黑龙江", City: "鸡西", Description: "中俄界湖兴凯湖，盛产大白鱼、鲤鱼"},
		{Name: "齐齐哈尔扎龙湿地水域", Province: "黑龙江", City: "齐齐哈尔", Description: "扎龙国家级自然保护区外围，丹顶鹤之乡"},
		{Name: "五大连池水域", Province: "黑龙江", City: "黑河", Description: "五大连池风景区，矿泉鱼闻名遐迩"},
		{Name: "佳木斯松花江段", Province: "黑龙江", City: "佳木斯", Description: "松花江下游佳木斯段，水面宽阔，鱼群密集"},
		// 北京
		{Name: "北京密云水库水域", Province: "北京", City: "北京", Description: "密云水库是北京最大水库，水质优良，盛产鲤鱼、鲢鳙、翘嘴"},
		{Name: "北京怀柔水库水域", Province: "北京", City: "北京", Description: "怀柔水库周边环境优美，适合休闲垂钓"},
		// 河北石家庄
		{Name: "石家庄岗南水库水域", Province: "河北", City: "石家庄", Description: "岗南水库位于滹沱河上游，是石家庄主要水源地之一，鱼种丰富"},
		{Name: "石家庄黄壁庄水库水域", Province: "河北", City: "石家庄", Description: "黄壁庄水库紧邻岗南水库，水面开阔，鲤鱼草鱼资源丰富"},
		// 新疆伊犁
		{Name: "伊犁赛里木湖水域", Province: "新疆", City: "伊犁", Description: "赛里木湖海拔2073米，高山冷水湖，盛产高白鲑和虹鳟"},
		{Name: "伊犁伊犁河水域", Province: "新疆", City: "伊犁", Description: "伊犁河是新疆水量最大河流，河鲈、鲤鱼资源丰富"},
	}
	DB.Create(&regions)

	// ===== 3. 创建网关 =====
	past2d := now.Add(-48 * time.Hour)
	gateways := []models.Gateway{
		{Name: "松花江边缘网关-01", Status: "online", Mode: "online", CPUUsage: 35.2, MemoryUsage: 48.5, DiskUsage: 22.1, BatteryLevel: 87, LastActiveAt: &now},
		{Name: "松花江边缘网关-02", Status: "online", Mode: "online", CPUUsage: 41.8, MemoryUsage: 52.1, DiskUsage: 25.3, BatteryLevel: 79, LastActiveAt: &now},
		{Name: "太阳岛边缘网关-01", Status: "online", Mode: "online", CPUUsage: 28.7, MemoryUsage: 42.3, DiskUsage: 18.6, BatteryLevel: 92, LastActiveAt: &now},
		{Name: "龙凤湿地网关-01", Status: "offline", Mode: "online", CPUUsage: 0, MemoryUsage: 0, DiskUsage: 15.3, BatteryLevel: 15},
		{Name: "镜泊湖边缘网关-01", Status: "online", Mode: "online", CPUUsage: 22.5, MemoryUsage: 38.9, DiskUsage: 19.8, BatteryLevel: 95, LastActiveAt: &now},
		{Name: "兴凯湖边缘网关-01", Status: "online", Mode: "online", CPUUsage: 30.1, MemoryUsage: 45.2, DiskUsage: 20.5, BatteryLevel: 88, LastActiveAt: &now},
		{Name: "扎龙湿地网关-01", Status: "maintenance", Mode: "maintenance", CPUUsage: 0, MemoryUsage: 0, DiskUsage: 30.2, BatteryLevel: 60, LastActiveAt: &past2d},
		{Name: "五大连池网关-01", Status: "online", Mode: "online", CPUUsage: 18.3, MemoryUsage: 35.6, DiskUsage: 16.7, BatteryLevel: 96, LastActiveAt: &now},
		// 北京
		{Name: "密云水库网关-01", Status: "online", Mode: "online", CPUUsage: 25.6, MemoryUsage: 40.8, DiskUsage: 18.2, BatteryLevel: 91, LastActiveAt: &now},
		{Name: "怀柔水库网关-01", Status: "online", Mode: "online", CPUUsage: 22.1, MemoryUsage: 38.5, DiskUsage: 16.5, BatteryLevel: 94, LastActiveAt: &now},
		// 石家庄
		{Name: "岗南水库网关-01", Status: "online", Mode: "online", CPUUsage: 31.5, MemoryUsage: 46.2, DiskUsage: 21.0, BatteryLevel: 85, LastActiveAt: &now},
		{Name: "黄壁庄水库网关-01", Status: "online", Mode: "online", CPUUsage: 27.3, MemoryUsage: 43.1, DiskUsage: 19.5, BatteryLevel: 88, LastActiveAt: &now},
		// 伊犁
		{Name: "赛里木湖网关-01", Status: "online", Mode: "online", CPUUsage: 19.8, MemoryUsage: 36.7, DiskUsage: 15.8, BatteryLevel: 97, LastActiveAt: &now},
		{Name: "伊犁河网关-01", Status: "online", Mode: "online", CPUUsage: 24.5, MemoryUsage: 41.3, DiskUsage: 17.9, BatteryLevel: 90, LastActiveAt: &now},
	}
	DB.Create(&gateways)

	// ===== 4. 创建设备 =====
	// 注意：FishingCount 是该设备监测到的实时垂钓人数，不能超过水域容量
	// 每个水域只绑定一个设备，设备的 fishing_count 就代表该水域的当前垂钓人数
	devices := []models.Device{
		// 松花江设备群 (网关1,2) — spot1 cap=80, spot2 cap=40, spot3 cap=30
		{Name: "松花江环境监测-A1", GatewayID: &gateways[0].ID, Status: "online", Description: "北岸主环境监测站", DeviceType: "environment", FishingCount: 7, WaterTemp: 8.5, AirTemp: 3.2, Humidity: 65.3, Pressure: 1013.25, LastActiveAt: &now},
		{Name: "松花江水下感知-B1", GatewayID: &gateways[0].ID, Status: "online", Description: "北岸水下探测节点", DeviceType: "underwater", FishingCount: 0, WaterTemp: 7.8, AirTemp: 3.2, Humidity: 65.3, Pressure: 1013.25, LastActiveAt: &now},
		{Name: "松花江探鱼辅助-C1", GatewayID: &gateways[0].ID, Status: "online", Description: "北岸声呐探鱼器", DeviceType: "fishfinder", FishingCount: 0, WaterTemp: 8.2, AirTemp: 3.1, Humidity: 64.8, Pressure: 1013.20, LastActiveAt: &now},
		{Name: "松花江环境监测-A2", GatewayID: &gateways[1].ID, Status: "online", Description: "大桥段环境监测站", DeviceType: "environment", FishingCount: 5, WaterTemp: 8.0, AirTemp: 2.8, Humidity: 66.1, Pressure: 1013.10, LastActiveAt: &now},
		{Name: "松花江水下感知-B2", GatewayID: &gateways[1].ID, Status: "online", Description: "大桥段水下探测节点", DeviceType: "underwater", FishingCount: 0, WaterTemp: 7.5, AirTemp: 2.8, Humidity: 66.1, Pressure: 1013.10, LastActiveAt: &now},
		// 太阳岛设备群 (网关3) — spot4 cap=60, spot5 cap=30
		{Name: "太阳岛环境监测-A1", GatewayID: &gateways[2].ID, Status: "online", Description: "西侧湖主监测站", DeviceType: "environment", FishingCount: 9, WaterTemp: 9.1, AirTemp: 4.5, Humidity: 62.1, Pressure: 1012.80, LastActiveAt: &now},
		{Name: "太阳岛水下感知-B1", GatewayID: &gateways[2].ID, Status: "online", Description: "西侧湖水下节点", DeviceType: "underwater", FishingCount: 0, WaterTemp: 8.9, AirTemp: 4.5, Humidity: 62.1, Pressure: 1012.80, LastActiveAt: &now},
		{Name: "太阳岛环境监测-A2", GatewayID: &gateways[2].ID, Status: "online", Description: "码头区域监测站", DeviceType: "environment", FishingCount: 4, WaterTemp: 9.3, AirTemp: 4.8, Humidity: 61.5, Pressure: 1012.75, LastActiveAt: &now},
		// 龙凤湿地设备群 (网关4 - 离线) — spot6 cap=50 closed
		{Name: "龙凤湿地监测-A1", GatewayID: &gateways[3].ID, Status: "offline", Description: "湿地外围监测站", DeviceType: "environment", FishingCount: 0, WaterTemp: 0, AirTemp: -2.1, Humidity: 78.5, Pressure: 1015.10},
		// 镜泊湖设备群 (网关5) — spot7 cap=100, spot8 cap=60
		{Name: "镜泊湖环境监测-A1", GatewayID: &gateways[4].ID, Status: "online", Description: "大坝下游监测站", DeviceType: "environment", FishingCount: 12, WaterTemp: 6.8, AirTemp: 1.5, Humidity: 70.2, Pressure: 1014.50, LastActiveAt: &now},
		{Name: "镜泊湖水下感知-B1", GatewayID: &gateways[4].ID, Status: "online", Description: "大坝深水区探测", DeviceType: "underwater", FishingCount: 0, WaterTemp: 6.2, AirTemp: 1.5, Humidity: 70.2, Pressure: 1014.50, LastActiveAt: &now},
		{Name: "镜泊湖探鱼辅助-C1", GatewayID: &gateways[4].ID, Status: "online", Description: "大坝声呐系统", DeviceType: "fishfinder", FishingCount: 0, WaterTemp: 6.5, AirTemp: 1.3, Humidity: 71.0, Pressure: 1014.55, LastActiveAt: &now},
		// 兴凯湖设备群 (网关6) — spot9 cap=120
		{Name: "兴凯湖环境监测-A1", GatewayID: &gateways[5].ID, Status: "online", Description: "北岸监测站", DeviceType: "environment", FishingCount: 15, WaterTemp: 5.5, AirTemp: 0.8, Humidity: 72.8, Pressure: 1015.20, LastActiveAt: &now},
		{Name: "兴凯湖水下感知-B1", GatewayID: &gateways[5].ID, Status: "online", Description: "湖心区水下探测", DeviceType: "underwater", FishingCount: 0, WaterTemp: 5.2, AirTemp: 0.8, Humidity: 72.8, Pressure: 1015.20, LastActiveAt: &now},
		// 扎龙设备 (网关7 - 维护中) — spot10 cap=40 maintenance
		{Name: "扎龙湿地监测-A1", GatewayID: &gateways[6].ID, Status: "error", Description: "设备固件升级中", DeviceType: "environment", FishingCount: 0, WaterTemp: 0, AirTemp: -3.5, Humidity: 80.1, Pressure: 1016.00, LastActiveAt: &past2d},
		// 五大连池设备 (网关8) — spot11 cap=45
		{Name: "五大连池环境监测-A1", GatewayID: &gateways[7].ID, Status: "online", Description: "药泉湖监测站", DeviceType: "environment", FishingCount: 3, WaterTemp: 4.2, AirTemp: -1.2, Humidity: 68.9, Pressure: 1014.80, LastActiveAt: &now},
		// 北京密云设备 (网关9)
		{Name: "密云水库环境监测-A1", GatewayID: &gateways[8].ID, Status: "online", Description: "密云水库大坝监测站", DeviceType: "environment", FishingCount: 18, WaterTemp: 12.5, AirTemp: 15.3, Humidity: 55.2, Pressure: 1016.30, LastActiveAt: &now},
		{Name: "密云水库水下感知-B1", GatewayID: &gateways[8].ID, Status: "online", Description: "密云水库深水区探测", DeviceType: "underwater", FishingCount: 0, WaterTemp: 10.8, AirTemp: 15.3, Humidity: 55.2, Pressure: 1016.30, LastActiveAt: &now},
		// 北京怀柔设备 (网关10)
		{Name: "怀柔水库环境监测-A1", GatewayID: &gateways[9].ID, Status: "online", Description: "怀柔水库西岸监测站", DeviceType: "environment", FishingCount: 10, WaterTemp: 13.1, AirTemp: 16.0, Humidity: 52.8, Pressure: 1016.10, LastActiveAt: &now},
		// 石家庄岗南设备 (网关11)
		{Name: "岗南水库环境监测-A1", GatewayID: &gateways[10].ID, Status: "online", Description: "岗南水库坝下监测站", DeviceType: "environment", FishingCount: 14, WaterTemp: 14.2, AirTemp: 18.5, Humidity: 48.6, Pressure: 1015.50, LastActiveAt: &now},
		{Name: "岗南水库水下感知-B1", GatewayID: &gateways[10].ID, Status: "online", Description: "岗南水库深水区探测", DeviceType: "underwater", FishingCount: 0, WaterTemp: 12.8, AirTemp: 18.5, Humidity: 48.6, Pressure: 1015.50, LastActiveAt: &now},
		// 石家庄黄壁庄设备 (网关12)
		{Name: "黄壁庄水库环境监测-A1", GatewayID: &gateways[11].ID, Status: "online", Description: "黄壁庄水库北岸监测站", DeviceType: "environment", FishingCount: 8, WaterTemp: 14.8, AirTemp: 19.0, Humidity: 46.5, Pressure: 1015.20, LastActiveAt: &now},
		// 伊犁赛里木湖设备 (网关13)
		{Name: "赛里木湖环境监测-A1", GatewayID: &gateways[12].ID, Status: "online", Description: "赛里木湖南岸监测站", DeviceType: "environment", FishingCount: 6, WaterTemp: 8.0, AirTemp: 10.2, Humidity: 35.8, Pressure: 1018.50, LastActiveAt: &now},
		{Name: "赛里木湖水下感知-B1", GatewayID: &gateways[12].ID, Status: "online", Description: "赛里木湖冷水鱼探测", DeviceType: "underwater", FishingCount: 0, WaterTemp: 6.5, AirTemp: 10.2, Humidity: 35.8, Pressure: 1018.50, LastActiveAt: &now},
		// 伊犁河设备 (网关14)
		{Name: "伊犁河环境监测-A1", GatewayID: &gateways[13].ID, Status: "online", Description: "伊犁河伊宁段监测站", DeviceType: "environment", FishingCount: 11, WaterTemp: 11.5, AirTemp: 14.8, Humidity: 40.2, Pressure: 1017.00, LastActiveAt: &now},
	}
	DB.Create(&devices)

	// ===== 5. 创建垂钓水域 =====
	spots := []models.FishingSpot{
		{Name: "松花江北岸钓场", RegionID: regions[0].ID, Description: "松花江北岸，适合路亚和台钓，冬季冰钓热门地点", Latitude: 45.8050, Longitude: 126.5450, WaterType: "river", Capacity: 80, Status: "open", BoundDeviceID: &devices[0].ID},
		{Name: "松花江公路大桥钓位", RegionID: regions[0].ID, Description: "公路大桥下游，深水区鱼种丰富，常有大鲤鱼出没", Latitude: 45.7550, Longitude: 126.6700, WaterType: "river", Capacity: 40, Status: "open", BoundDeviceID: &devices[3].ID},
		{Name: "松花江铁路桥野钓区", RegionID: regions[0].ID, Description: "铁路桥附近自然钓场，水流较急，适合有经验的钓友", Latitude: 45.7200, Longitude: 126.7800, WaterType: "river", Capacity: 30, Status: "open", BoundDeviceID: &devices[4].ID},
		{Name: "太阳岛西侧湖钓场", RegionID: regions[1].ID, Description: "太阳岛西侧静水湖，水质清澈，适合休闲垂钓", Latitude: 45.8400, Longitude: 126.4800, WaterType: "lake", Capacity: 60, Status: "open", BoundDeviceID: &devices[5].ID},
		{Name: "太阳岛码头钓位", RegionID: regions[1].ID, Description: "老码头附近，水深2-4米，鲫鱼鲤鱼密度高", Latitude: 45.7750, Longitude: 126.5100, WaterType: "river", Capacity: 30, Status: "open", BoundDeviceID: &devices[7].ID},
		{Name: "龙凤湿地外围钓场", RegionID: regions[2].ID, Description: "湿地外围开放区域，冬季关闭", Latitude: 46.5812, Longitude: 125.0934, WaterType: "lake", Capacity: 50, Status: "closed", BoundDeviceID: &devices[8].ID},
		{Name: "镜泊湖大坝钓场", RegionID: regions[3].ID, Description: "镜泊湖大坝下游，大鱼频出，专业钓手天堂", Latitude: 44.0112, Longitude: 128.9945, WaterType: "reservoir", Capacity: 100, Status: "open", BoundDeviceID: &devices[9].ID},
		{Name: "镜泊湖吊水楼瀑布钓区", RegionID: regions[3].ID, Description: "吊水楼瀑布下方深潭，水流交汇处鱼群聚集", Latitude: 43.9856, Longitude: 129.0123, WaterType: "reservoir", Capacity: 60, Status: "open", BoundDeviceID: &devices[10].ID},
		{Name: "兴凯湖北岸垂钓基地", RegionID: regions[4].ID, Description: "兴凯湖北岸专业垂钓基地，盛产大白鱼、鲤鱼、鲢鳙", Latitude: 45.3245, Longitude: 132.4567, WaterType: "lake", Capacity: 120, Status: "open", BoundDeviceID: &devices[12].ID},
		{Name: "扎龙湿地观鸟钓场", RegionID: regions[5].ID, Description: "扎龙保护区外围，可边观鸟边垂钓，别有趣味", Latitude: 47.1789, Longitude: 124.2856, WaterType: "pond", Capacity: 40, Status: "maintenance", BoundDeviceID: &devices[14].ID},
		{Name: "五大连池药泉湖钓场", RegionID: regions[6].ID, Description: "药泉湖垂钓区，矿泉水养殖鱼口感鲜美", Latitude: 48.7523, Longitude: 126.1189, WaterType: "lake", Capacity: 45, Status: "open", BoundDeviceID: &devices[15].ID},
		{Name: "佳木斯松花江外滩钓场", RegionID: regions[7].ID, Description: "佳木斯外滩公园段，休闲垂钓好去处，交通便利", Latitude: 46.8012, Longitude: 130.3678, WaterType: "river", Capacity: 70, Status: "open"},
		// 北京 (regions[8]=密云, regions[9]=怀柔)
		{Name: "密云水库大坝钓场", RegionID: regions[8].ID, Description: "密云水库大坝下游，水深6-15米，大鱼频出，北京钓友圣地", Latitude: 40.5280, Longitude: 116.9750, WaterType: "reservoir", Capacity: 100, Status: "open", BoundDeviceID: &devices[16].ID},
		{Name: "密云水库白河湾钓场", RegionID: regions[8].ID, Description: "白河入库口，水流交汇处鱼群密集，适合路亚和台钓", Latitude: 40.4950, Longitude: 116.9200, WaterType: "reservoir", Capacity: 60, Status: "open", BoundDeviceID: &devices[17].ID},
		{Name: "怀柔水库西坝钓场", RegionID: regions[9].ID, Description: "怀柔水库西侧安静水域，周末钓友云集，环境清幽", Latitude: 40.3560, Longitude: 116.6280, WaterType: "reservoir", Capacity: 80, Status: "open", BoundDeviceID: &devices[18].ID},
		// 石家庄 (regions[10]=岗南, regions[11]=黄壁庄)
		{Name: "岗南水库坝下钓场", RegionID: regions[10].ID, Description: "岗南水库主坝下方，水深适中，鲤鱼草鱼活跃", Latitude: 38.3520, Longitude: 114.1680, WaterType: "reservoir", Capacity: 90, Status: "open", BoundDeviceID: &devices[19].ID},
		{Name: "岗南水库苇塘钓区", RegionID: regions[10].ID, Description: "水库东侧苇塘区域，鲫鱼密度高，休闲台钓首选", Latitude: 38.3800, Longitude: 114.2100, WaterType: "reservoir", Capacity: 50, Status: "open", BoundDeviceID: &devices[20].ID},
		{Name: "黄壁庄水库北岸钓场", RegionID: regions[11].ID, Description: "黄壁庄水库北岸开阔地带，水面宽广，适合抛竿远投", Latitude: 38.2850, Longitude: 114.0650, WaterType: "reservoir", Capacity: 70, Status: "open", BoundDeviceID: &devices[21].ID},
		// 伊犁 (regions[12]=赛里木湖, regions[13]=伊犁河)
		{Name: "赛里木湖南岸钓场", RegionID: regions[12].ID, Description: "赛里木湖南岸，高山冷水湖，盛产高白鲑和虹鳟，景色绝美", Latitude: 44.5800, Longitude: 81.1600, WaterType: "lake", Capacity: 40, Status: "open", BoundDeviceID: &devices[22].ID},
		{Name: "赛里木湖东岸野钓区", RegionID: regions[12].ID, Description: "赛里木湖东岸原生态区域，人少鱼多，徒步可达", Latitude: 44.6200, Longitude: 81.2500, WaterType: "lake", Capacity: 30, Status: "open", BoundDeviceID: &devices[23].ID},
		{Name: "伊犁河伊宁段钓场", RegionID: regions[13].ID, Description: "伊犁河流经伊宁市区段，河鲈和鲤鱼资源丰富，交通便利", Latitude: 43.9150, Longitude: 81.3200, WaterType: "river", Capacity: 60, Status: "open", BoundDeviceID: &devices[24].ID},
	}
	DB.Create(&spots)

	// open 状态的水域索引（用于生成环境/历史数据）
	// spots: 0-4(open),5(closed),6-8(open),9(maintenance),10-11(open), 12-20(新增全部open)
	openSpotIndices := []int{0, 1, 2, 3, 4, 6, 7, 8, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	// ===== 6. 创建历史数据（最近 48 小时，每小时一条）=====
	var historicalData []models.HistoricalData
	for i := 47; i >= 0; i-- {
		ts := now.Add(-time.Duration(i) * time.Hour)
		hour := ts.Hour()

		for _, si := range openSpotIndices {
			spot := spots[si]
			cap := spot.Capacity

			// 模拟白天人多、夜晚人少的规律（比例控制偏低，更贴近实际）
			var ratio float64
			switch {
			case hour >= 5 && hour <= 7:
				ratio = 0.05 // 清晨早钓，少数人
			case hour >= 8 && hour <= 10:
				ratio = 0.12 // 上午黄金时段
			case hour >= 11 && hour <= 13:
				ratio = 0.08 // 中午较少
			case hour >= 14 && hour <= 17:
				ratio = 0.15 // 下午高峰
			case hour >= 18 && hour <= 20:
				ratio = 0.06 // 傍晚收竿
			default:
				ratio = 0.01 // 夜间几乎无人
			}

			base := int(float64(cap) * ratio)
			// 小幅扰动（±3 以内）
			noise := (i*7 + si*13) % 7 - 3
			count := base + noise
			if count < 0 {
				count = 0
			}
			if count > cap {
				count = cap
			}

			historicalData = append(historicalData, models.HistoricalData{
				SpotID:       spot.ID,
				FishingCount: count,
				Timestamp:    ts,
			})
		}
	}
	DB.CreateInBatches(&historicalData, 200)

	// ===== 7. 创建环境数据（最近 48 小时，每小时一条）=====
	var envData []models.EnvironmentData
	for i := 47; i >= 0; i-- {
		ts := now.Add(-time.Duration(i) * time.Hour)

		for _, si := range openSpotIndices {
			spot := spots[si]
			sid := int(spot.ID)

			// 根据水域位置和时间生成不同的基础温度
			baseWaterTemp := 5.0 + float64(sid)*0.4 + math.Sin(float64(i)*0.26)*1.5
			baseAirTemp := baseWaterTemp - 4.0 + math.Cos(float64(i)*0.26)*2.0
			baseHumidity := 58.0 + float64((sid*7+i*3)%25)
			basePressure := 1010.0 + math.Sin(float64(i)*0.13)*3.0 + float64(sid)*0.3

			ph := 6.8 + float64((sid+i)%8)*0.1
			do := 7.5 + math.Sin(float64(i)*0.3+float64(sid))*1.5
			turbidity := 10.0 + float64((i*3+sid*5)%20)

			envData = append(envData, models.EnvironmentData{
				SpotID:          spot.ID,
				WaterTemp:       math.Round(baseWaterTemp*10) / 10,
				AirTemp:         math.Round(baseAirTemp*10) / 10,
				Humidity:        math.Round(baseHumidity*10) / 10,
				Pressure:        math.Round(basePressure*100) / 100,
				PH:              math.Round(ph*10) / 10,
				DissolvedOxygen: math.Round(do*10) / 10,
				Turbidity:       math.Round(turbidity*10) / 10,
				Timestamp:       ts,
			})
		}
	}
	DB.CreateInBatches(&envData, 200)

	// ===== 8. 创建水质数据（最近 24 小时）=====
	// 给有水下设备的水域写水质数据
	// devices[1]=松花江水下B1, devices[4]=松花江水下B2, devices[6]=太阳岛水下B1, devices[10]=镜泊湖水下B1, devices[13]=兴凯湖水下B1
	// devices[17]=密云水下B1, devices[20]=岗南水下B1, devices[23]=赛里木湖水下B1
	underwaterDeviceIDs := []uint{devices[1].ID, devices[4].ID, devices[6].ID, devices[10].ID, devices[13].ID, devices[17].ID, devices[20].ID, devices[23].ID}
	var wqData []models.WaterQualityData
	for i := 23; i >= 0; i-- {
		ts := now.Add(-time.Duration(i) * time.Hour)
		for idx, devID := range underwaterDeviceIDs {
			ph := 7.0 + float64(i%(4+idx))*0.1
			do := 8.2 + math.Sin(float64(i)*0.4+float64(idx))*0.8
			turb := 10.0 + float64((i*3+idx*5)%12)
			wqData = append(wqData, models.WaterQualityData{
				DeviceID:        devID,
				PH:              math.Round(ph*10) / 10,
				DissolvedOxygen: math.Round(do*10) / 10,
				Turbidity:       math.Round(turb*10) / 10,
				Timestamp:       ts,
			})
		}
	}
	DB.CreateInBatches(&wqData, 100)

	// ===== 9. 创建提醒 =====
	reminders := []models.Reminder{
		// 天气类
		{SpotID: spots[0].ID, Level: 1, ReminderType: "weather", Message: "明日有大风预警，松花江北岸风力可达5-6级，建议做好防风准备", Timestamp: now.Add(-2 * time.Hour), Publicity: true},
		{SpotID: spots[6].ID, Level: 1, ReminderType: "weather", Message: "镜泊湖区域未来3天有降雪预报，气温将降至-15℃以下，请注意防寒保暖", Timestamp: now.Add(-5 * time.Hour), Publicity: true},
		{SpotID: spots[8].ID, Level: 0, ReminderType: "weather", Message: "兴凯湖今日天气晴好，风力2-3级，非常适合出钓", Timestamp: now.Add(-1 * time.Hour), Publicity: true},
		// 垂钓类
		{SpotID: spots[3].ID, Level: 0, ReminderType: "fishing", Message: "太阳岛西侧湖今日适宜垂钓，水温适中，鱼口较好，建议使用红虫", Timestamp: now.Add(-1 * time.Hour), Publicity: true},
		{SpotID: spots[0].ID, Level: 0, ReminderType: "fishing", Message: "松花江北岸近期鲫鱼活跃，早6点-8点上鱼率最高", Timestamp: now.Add(-3 * time.Hour), Publicity: true},
		{SpotID: spots[7].ID, Level: 0, ReminderType: "fishing", Message: "镜泊湖吊水楼钓区今晨有钓友钓获12斤鲤鱼，鱼情大好", Timestamp: now.Add(-30 * time.Minute), Publicity: true},
		// 安全类
		{SpotID: spots[5].ID, Level: 3, ReminderType: "safety", Message: "龙凤湿地外围钓场因冰面不稳定暂时关闭，严禁进入！", Timestamp: now.Add(-30 * time.Minute), Publicity: true},
		{SpotID: spots[2].ID, Level: 2, ReminderType: "safety", Message: "松花江铁路桥野钓区水流较急，近日有险情上报，请佩戴救生衣", Timestamp: now.Add(-4 * time.Hour), Publicity: true},
		{SpotID: spots[9].ID, Level: 2, ReminderType: "safety", Message: "扎龙湿地观鸟钓场设备维护中，暂不开放", Timestamp: now.Add(-6 * time.Hour), Publicity: true},
		// 环境类
		{SpotID: spots[1].ID, Level: 2, ReminderType: "environment", Message: "松花江公路大桥段水质浑浊度升高至25NTU，可能影响钓获", Timestamp: now, Publicity: true},
		{SpotID: spots[10].ID, Level: 1, ReminderType: "environment", Message: "五大连池药泉湖 pH 值偏高(8.2)，建议关注水质变化", Timestamp: now.Add(-8 * time.Hour), Publicity: true},
		// 已解决的历史提醒
		{SpotID: spots[0].ID, Level: 1, ReminderType: "weather", Message: "昨日大风预警已解除，松花江北岸恢复正常垂钓", Timestamp: now.Add(-26 * time.Hour), Resolved: true, Publicity: true},
		{SpotID: spots[3].ID, Level: 2, ReminderType: "environment", Message: "太阳岛西侧湖溶解氧恢复正常水平", Timestamp: now.Add(-20 * time.Hour), Resolved: true, Publicity: true},
		// 北京
		{SpotID: spots[12].ID, Level: 0, ReminderType: "fishing", Message: "密云水库大坝钓场近期鲤鱼活跃，水温12.5℃适宜垂钓，建议使用玉米粒打窝", Timestamp: now.Add(-2 * time.Hour), Publicity: true},
		{SpotID: spots[14].ID, Level: 1, ReminderType: "weather", Message: "怀柔水库周末有阵雨预报，请携带雨具，注意防滑", Timestamp: now.Add(-4 * time.Hour), Publicity: true},
		// 石家庄
		{SpotID: spots[15].ID, Level: 0, ReminderType: "fishing", Message: "岗南水库坝下草鱼开口良好，推荐使用嫩玉米挂钩，钓3-4米深", Timestamp: now.Add(-1 * time.Hour), Publicity: true},
		{SpotID: spots[17].ID, Level: 1, ReminderType: "environment", Message: "黄壁庄水库北岸近日水位下降约0.5米，钓位需前移", Timestamp: now.Add(-6 * time.Hour), Publicity: true},
		// 伊犁
		{SpotID: spots[18].ID, Level: 0, ReminderType: "fishing", Message: "赛里木湖南岸高白鲑活跃，水温8℃，路亚银色勺形亮片5g效果最佳", Timestamp: now.Add(-3 * time.Hour), Publicity: true},
		{SpotID: spots[20].ID, Level: 1, ReminderType: "weather", Message: "伊犁河伊宁段上游有融雪水汇入，水温偏低且水流加大，注意安全", Timestamp: now.Add(-5 * time.Hour), Publicity: true},
	}
	DB.Create(&reminders)

	// ===== 10. 创建通知 =====
	notices := []models.Notice{
		{Title: "2026年冬季冰钓节活动通知", Content: "哈尔滨市第十届冬季冰钓节将于2026年1月15日在松花江北岸举办，欢迎广大钓友参加！活动设有个人赛、团队赛等多个项目，总奖金池10万元。", Timestamp: now.Add(-72 * time.Hour)},
		{Title: "垂钓安全须知更新", Content: "各位钓友注意：冬季垂钓务必注意冰面安全，请在标识安全区域内活动。携带救生装备，结伴出行。设备有异常及时上报管理员。", Timestamp: now.Add(-48 * time.Hour)},
		{Title: "系统维护通知", Content: "智钓蓝海平台将于本周日凌晨2:00-4:00进行系统升级维护，届时部分功能可能暂时不可用，敬请谅解。本次更新将优化水质监测数据的实时展示。", Timestamp: now.Add(-12 * time.Hour)},
		{Title: "新增兴凯湖垂钓基地", Content: "好消息！智钓蓝海平台新增兴凯湖北岸垂钓基地，配备环境监测和水下探测设备，实时掌握鱼情动态。兴凯湖盛产大白鱼，欢迎钓友们前去体验！", Timestamp: now.Add(-24 * time.Hour)},
		{Title: "春季禁渔期提醒", Content: "根据黑龙江省渔业管理规定，松花江流域每年4月1日至6月15日为禁渔期，届时相关水域将暂停垂钓服务。请钓友们提前安排行程。", Timestamp: now.Add(-6 * time.Hour)},
		{Title: "【已过期】元旦活动已结束", Content: "2026年元旦垂钓嘉年华活动已圆满结束，感谢各位钓友的热情参与！获奖名单已在公告栏公布。", Timestamp: now.Add(-30 * 24 * time.Hour), Outdated: true},
		{Title: "新增北京密云水库、怀柔水库钓场", Content: "好消息！智钓蓝海平台新增北京密云水库和怀柔水库两大钓场，已部署环境监测和水下探测设备。密云水库盛产大鲤鱼，怀柔水库环境优美，欢迎北京钓友前往体验！", Timestamp: now.Add(-8 * time.Hour)},
		{Title: "新增石家庄岗南水库、黄壁庄水库钓场", Content: "河北石家庄地区的岗南水库和黄壁庄水库钓场正式上线！两大水库水面开阔，鲤鱼草鱼资源丰富，配备全套智能监测设备。", Timestamp: now.Add(-10 * time.Hour)},
		{Title: "新增新疆伊犁赛里木湖、伊犁河钓场", Content: "伊犁赛里木湖高山冷水钓场和伊犁河钓场已接入平台。赛里木湖海拔2073米，盛产高白鲑和虹鳟，是路亚爱好者的天堂。伊犁河段河鲈资源丰富，欢迎体验！", Timestamp: now.Add(-9 * time.Hour)},
	}
	DB.Create(&notices)

	// 关联通知到水域
	DB.Model(&notices[0]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2])
	DB.Model(&notices[1]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2], &spots[3], &spots[4], &spots[6], &spots[7], &spots[8], &spots[10], &spots[11])
	DB.Model(&notices[3]).Association("RelatedSpots").Append(&spots[8])
	DB.Model(&notices[4]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2], &spots[11])
	DB.Model(&notices[6]).Association("RelatedSpots").Append(&spots[12], &spots[13], &spots[14])
	DB.Model(&notices[7]).Association("RelatedSpots").Append(&spots[15], &spots[16], &spots[17])
	DB.Model(&notices[8]).Association("RelatedSpots").Append(&spots[18], &spots[19], &spots[20])

	// ===== 11. 用户收藏 =====
	DB.Model(&users[3]).Association("Favorites").Append(&spots[0], &spots[3], &spots[6])   // fisher01
	DB.Model(&users[4]).Association("Favorites").Append(&spots[0], &spots[8])               // fisher02
	DB.Model(&users[5]).Association("Favorites").Append(&spots[3], &spots[4], &spots[10])   // fisher03
	DB.Model(&users[6]).Association("Favorites").Append(&spots[6], &spots[7], &spots[8])    // fisher04
	DB.Model(&users[7]).Association("Favorites").Append(&spots[0], &spots[1], &spots[2])    // fisher05
	DB.Model(&users[8]).Association("Favorites").Append(&spots[3], &spots[11])              // fisher06
	DB.Model(&users[9]).Association("Favorites").Append(&spots[6], &spots[8])               // fisher07

	// ===== 12. 垂钓建议 =====
	suggestions := []models.FishingSuggestion{
		{SpotID: spots[0].ID, UserID: &users[3].ID, SuggestionText: "当前松花江北岸水温8.5℃，建议使用红虫或蚯蚓做饵，主攻鲫鱼。早晨6-8点和傍晚4-6点是最佳垂钓时段。1.5号子线配4号袖钩效果最佳。", Score: 85.5, Timestamp: now},
		{SpotID: spots[3].ID, SuggestionText: "太阳岛西侧湖水深适中，水质清澈，当前溶氧充足。推荐使用3.6米手竿，钓底层，搓饵料球。近期鲤鱼活跃度较高，建议用玉米粒打窝。", Score: 78.2, Timestamp: now.Add(-2 * time.Hour)},
		{SpotID: spots[6].ID, SuggestionText: "镜泊湖大坝钓场水深8-12米，建议使用矶钓竿或6.3米以上长竿。深水层水温稳定，鲢鳙和大鲤鱼活跃。推荐使用发酵饵料配合雾化饵。", Score: 92.1, Timestamp: now.Add(-1 * time.Hour)},
		{SpotID: spots[8].ID, UserID: &users[6].ID, SuggestionText: "兴凯湖北岸当前白鱼群已开始活跃，水面温度5.5℃。推荐使用路亚拟饵在浅水区抛投，勺型亮片3-5克效果好。日出后1小时为最佳窗口。", Score: 88.7, Timestamp: now.Add(-3 * time.Hour)},
		{SpotID: spots[1].ID, UserID: &users[4].ID, SuggestionText: "松花江大桥段近期水质浑浊度偏高，建议使用味道浓烈的饵料（腥饵为主），钓位选择缓流区。2号主线配1号子线，伊势尼5号钩。", Score: 72.3, Timestamp: now.Add(-5 * time.Hour)},
		{SpotID: spots[4].ID, SuggestionText: "太阳岛码头水深2-4米，底部淤泥较厚，建议调钓灵敏（调4钓2）。近期鲫鱼为主，偶有黄辣丁。红虫拉饵效果优于商品饵。", Score: 80.5, Timestamp: now.Add(-4 * time.Hour)},
		{SpotID: spots[10].ID, UserID: &users[5].ID, SuggestionText: "五大连池药泉湖矿物质含量高，鱼口偏轻。建议使用0.6号子线配3号金袖，浮漂选择吃铅小的芦苇漂，调目不宜过高。", Score: 76.8, Timestamp: now.Add(-6 * time.Hour)},
		{SpotID: spots[7].ID, SuggestionText: "镜泊湖吊水楼瀑布钓区水流交汇处，大物频出。今日已有多名钓友中获5斤以上鲤鱼。建议加粗线组，准备抄网。", Score: 90.3, Timestamp: now.Add(-30 * time.Minute)},
		{SpotID: spots[11].ID, SuggestionText: "佳木斯松花江段目前尚未绑定监测设备，暂无实时环境数据。根据历史数据，该段3月下旬开始回暖，鲫鱼和鲤鱼陆续开口。", Score: 65.0, Timestamp: now.Add(-12 * time.Hour)},
		{SpotID: spots[0].ID, UserID: &users[5].ID, SuggestionText: "综合近48小时数据分析：松花江北岸水温呈缓慢上升趋势，溶解氧充足，气压稳定。预计明日垂钓条件评分将达90+，强烈推荐出钓。", Score: 95.0, Timestamp: now.Add(-10 * time.Minute)},
		// 北京
		{SpotID: spots[12].ID, SuggestionText: "密云水库大坝钓场当前水温12.5℃，春季回暖鲤鱼开口积极。建议使用4.5米手竿，主线2号配子线1号，伊势尼6号钩。搓饵或玉米粒打窝效果俱佳。", Score: 88.0, Timestamp: now.Add(-2 * time.Hour)},
		{SpotID: spots[13].ID, UserID: &users[4].ID, SuggestionText: "密云水库白河湾入水口处，路亚翘嘴效果极佳。推荐使用7-10g银色米诺，匀速收线带小抽。清晨6-8点为最佳窗口期。", Score: 82.5, Timestamp: now.Add(-4 * time.Hour)},
		{SpotID: spots[14].ID, SuggestionText: "怀柔水库西坝水深3-6米，底部碎石底质，适合钓鲫鱼。推荐3.6米短竿，0.8号子线配3号袖钩，红虫或蚯蚓拉饵。调4钓2较为灵敏。", Score: 79.3, Timestamp: now.Add(-6 * time.Hour)},
		// 石家庄
		{SpotID: spots[15].ID, UserID: &users[3].ID, SuggestionText: "岗南水库坝下当前水温14.2℃，草鱼进入活跃期。推荐5.4米竿，3号主线配1.5号子线，使用嫩玉米或青草做饵。钓3-4米深效果最好。", Score: 91.2, Timestamp: now.Add(-1 * time.Hour)},
		{SpotID: spots[17].ID, SuggestionText: "黄壁庄水库北岸近日水位略降，鱼群向深水区移动。建议使用6.3米以上长竿远投，或矶钓竿配合串钩。饵料以腥香为主，雾化饵打窝聚鱼。", Score: 75.8, Timestamp: now.Add(-8 * time.Hour)},
		// 伊犁
		{SpotID: spots[18].ID, SuggestionText: "赛里木湖南岸高白鲑活跃，水温8℃适合冷水鱼觅食。路亚推荐5g银色勺形亮片，慢速匀收。注意高原紫外线强烈，做好防晒。海拔2073米，体力消耗比平原大。", Score: 86.5, Timestamp: now.Add(-3 * time.Hour)},
		{SpotID: spots[20].ID, UserID: &users[6].ID, SuggestionText: "伊犁河伊宁段河鲈活跃，当前水温11.5℃。路亚推荐使用T尾软虫配铅头钩（7g），沿河岸缓流区搜索。台钓推荐蚯蚓挂钩钓底，主攻鲤鱼和鲫鱼。", Score: 83.7, Timestamp: now.Add(-5 * time.Hour)},
	}
	DB.Create(&suggestions)

	// ===== 13. SFR 社区帖子 =====
	posts := []models.Post{
		{UserID: users[3].ID, Title: "今日松花江北岸冰钓，收获满满", Body: "今天天气不错，零下5度，风力2级。早上6点到的松花江北岸，冰层厚度约30cm，非常安全。\n\n用的红虫做饵，1.5号子线配4号袖钩，调4钓2。\n\n早上7点开始连续上鱼，到10点一共钓了12条鲫鱼，最大的有半斤多。中间还跑了一条大鲤鱼，估计有3斤左右，可惜子线太细被切了。\n\n总结：松花江北岸鱼情确实不错，推荐大家来试试！", Tag: "钓鱼日记"},
		{UserID: users[4].ID, Title: "路亚新手入坑指南 —— 装备选购篇", Body: "最近很多钓友问我路亚入门买什么装备好，这里整理一份新手指南：\n\n1. 竿子：推荐ML调性的直柄竿，6.6尺-7尺，适合抛投3-15g的饵。品牌推荐达亿瓦、禧玛诺入门款。\n\n2. 轮子：纺车轮2000-2500型号，选带微调刹车的，新手容易上手。\n\n3. 线组：PE线0.6-0.8号配碳素前导线6-8磅。\n\n4. 假饵：先买一些亮片（3-7g银色/金色各几个）、软虫（T尾和卷尾）、米诺（5-7cm）。\n\n预算1000元左右就能搞定一套入门装备，不要一上来就买太贵的。", Tag: "经验分享"},
		{UserID: users[5].ID, Title: "镜泊湖大坝钓场实测报告", Body: "上周末去了镜泊湖大坝钓场，这里水深8-12米，鱼种丰富。\n\n使用6.3米矶钓竿，3号主线配1.5号子线，伊势尼7号钩。饵料用的发酵玉米粒配合商品饵搓饵。\n\n第一天上了3条鲤鱼（最大的8斤）和5条鲢鳙。第二天鱼口稍差，但也有2条鲤鱼入账。\n\n这个钓场真的适合想搞大物的钓友，设备监测数据也很准确，水温和气压信息很有参考价值。", Tag: "钓点推荐"},
		{UserID: users[6].ID, Title: "兴凯湖白鱼季开始了！路亚爆护", Body: "兴凯湖北岸的白鱼季终于来了！今天用路亚拟饵在浅水区抛投，3-5克银色亮片，疯狂上鱼！\n\n从日出开始到上午10点，一共中了15条白鱼，最大的有1.2斤。白鱼拉力很足，用ML竿手感极佳。\n\n提醒钓友们：白鱼嘴薄，中鱼后不要太急着收线，保持竿稍弹性。", Tag: "钓鱼日记"},
		{UserID: users[7].ID, Title: "求助：台钓新手调漂总是调不准", Body: "各位大佬好，我是刚入坑台钓的新手，有个问题想请教。\n\n我买的是一支3.6米手竿，用的是吃铅1.2g的芦苇浮漂。按照教程调4钓2，但是实际钓的时候总感觉不对劲，鱼口反应不明显，经常空竿。\n\n请问是浮漂选的不对，还是调法有问题？我的钓点水深大约2米。用的红虫拉饵。\n\n希望有经验的钓友指点一二，谢谢！", Tag: "问答求助"},
		{UserID: users[8].ID, Title: "达亿瓦 23款天弓鲤测评", Body: "入手了一根达亿瓦23款天弓鲤4.5米，使用一个月来分享下感受。\n\n优点：\n- 腰力非常好，博大鱼不吃力\n- 涂装漂亮，做工精细\n- 自重轻，长时间持竿不累\n\n缺点：\n- 价格偏高（1200+）\n- 竿稍偏硬，小鱼手感一般\n\n总体来说是一款非常优秀的综合竿，适合野钓和黑坑。如果预算够的话强烈推荐。打分：8.5/10", Tag: "装备测评"},
		{UserID: users[9].ID, Title: "太阳岛西侧湖夜钓记", Body: "昨晚在太阳岛西侧湖夜钓，从下午5点一直钓到凌晨1点。\n\n环境很好，水面平静，蚊虫也不多（带了驱蚊灯）。用的夜光浮漂+蓝光灯，看漂非常清楚。\n\n饵料：前期用腥饵打窝聚鱼，后期换拉饵钓。\n\n战果：鲫鱼23条、鲤鱼2条、黄辣丁4条，还有一条不到半斤的小鲶鱼。\n\n太阳岛夜钓体验极佳，推荐！记得带厚衣服，夜里降温明显。", Tag: "钓鱼日记"},
	}
	DB.Create(&posts)

	// ===== 14. 评论 =====
	comments := []models.Comment{
		{PostID: posts[0].PostID, UserID: users[4].ID, Body: "太厉害了！松花江北岸确实出鱼，我上次也去了，不过没有你钓得多"},
		{PostID: posts[0].PostID, UserID: users[5].ID, Body: "请问冰钓用什么型号的冰钻？厚度30cm的话普通手钻能打透吗？"},
		{PostID: posts[0].PostID, UserID: users[6].ID, Body: "大鲤鱼跑了太可惜了，下次建议用1号子线+伊势尼钩"},
		{PostID: posts[1].PostID, UserID: users[3].ID, Body: "写得很详细，新手友好！我就是按照这个思路入门的路亚"},
		{PostID: posts[1].PostID, UserID: users[7].ID, Body: "补充一下，新手PE线可以先买编织线，比较耐磨，不容易炸线"},
		{PostID: posts[2].PostID, UserID: users[4].ID, Body: "镜泊湖大坝确实出大物，我之前去钓过一条10斤的鲤鱼"},
		{PostID: posts[2].PostID, UserID: users[8].ID, Body: "请问那边住宿方便吗？打算周末过去"},
		{PostID: posts[3].PostID, UserID: users[3].ID, Body: "羡慕！白鱼路亚手感确实好，我也想去兴凯湖试试"},
		{PostID: posts[4].PostID, UserID: users[5].ID, Body: "新手调漂建议先在家里用水盆练习，找到半水调目。你的浮漂吃铅1.2g偏大了，建议换0.8g左右的"},
		{PostID: posts[4].PostID, UserID: users[6].ID, Body: "调4钓2是经典调法，但实际钓的时候要根据鱼口调整。鱼口轻就调高钓低，鱼口重就调低钓高"},
		{PostID: posts[5].PostID, UserID: users[9].ID, Body: "天弓鲤确实不错，我用的3.9米款，博5斤鲤鱼游刃有余"},
		{PostID: posts[6].PostID, UserID: users[4].ID, Body: "太阳岛夜钓氛围确实好，下次组队一起去"},
	}
	DB.Create(&comments)

	// ===== 15. 楼中楼子评论 =====
	cocs := []models.CommentOnComments{
		{CommentID: comments[0].CommentID, UserID: users[3].ID, Body: "哈哈主要是今天鱼口好，运气成分也有"},
		{CommentID: comments[1].CommentID, UserID: users[3].ID, Body: "我用的手摇冰钻，一般15分钟就能打穿30cm。电钻也行但是重"},
		{CommentID: comments[2].CommentID, UserID: users[3].ID, Body: "有道理，下次升级线组再战！"},
		{CommentID: comments[4].CommentID, UserID: users[4].ID, Body: "确实，编织PE线不容易打结，新手友好很多"},
		{CommentID: comments[8].CommentID, UserID: users[7].ID, Body: "谢谢！我去买个0.8g的浮漂试试"},
		{CommentID: comments[9].CommentID, UserID: users[7].ID, Body: "明白了，我试试调5钓3看看效果"},
	}
	DB.Create(&cocs)

	// ===== 16. 帖子点赞 =====
	postLikes := []models.LikeOnPosts{
		{PostID: posts[0].PostID, UserID: users[4].ID},
		{PostID: posts[0].PostID, UserID: users[5].ID},
		{PostID: posts[0].PostID, UserID: users[6].ID},
		{PostID: posts[0].PostID, UserID: users[7].ID},
		{PostID: posts[1].PostID, UserID: users[3].ID},
		{PostID: posts[1].PostID, UserID: users[5].ID},
		{PostID: posts[1].PostID, UserID: users[7].ID},
		{PostID: posts[2].PostID, UserID: users[4].ID},
		{PostID: posts[2].PostID, UserID: users[8].ID},
		{PostID: posts[3].PostID, UserID: users[3].ID},
		{PostID: posts[3].PostID, UserID: users[5].ID},
		{PostID: posts[3].PostID, UserID: users[4].ID},
		{PostID: posts[4].PostID, UserID: users[5].ID},
		{PostID: posts[4].PostID, UserID: users[6].ID},
		{PostID: posts[5].PostID, UserID: users[9].ID},
		{PostID: posts[6].PostID, UserID: users[3].ID},
		{PostID: posts[6].PostID, UserID: users[4].ID},
		{PostID: posts[6].PostID, UserID: users[8].ID},
	}
	DB.Create(&postLikes)

	// ===== 17. IoT 设备（用户绑定的智能浮漂/钓箱） =====
	iotDevices := []models.IoTDevice{
		{DeviceID: "SF-FLOAT-001", Temperature: 8.2, Humidity: 64.5, Pulling: 0.35, Pressure: 1013.8, GpsInfo: "45.8050,126.5450", ImuData: "stable", LastUpdate: now},
		{DeviceID: "SF-FLOAT-002", Temperature: 6.5, Humidity: 70.1, Pulling: 0.0, Pressure: 1014.2, GpsInfo: "44.0112,128.9945", ImuData: "stable", LastUpdate: now.Add(-2 * time.Hour)},
		{DeviceID: "SF-BOX-001", Temperature: 5.8, Humidity: 71.3, Pulling: 1.25, Pressure: 1015.0, GpsInfo: "45.3245,132.4567", ImuData: "active", LastUpdate: now.Add(-30 * time.Minute)},
		{DeviceID: "SF-FLOAT-003", Temperature: 9.1, Humidity: 61.8, Pulling: 0.0, Pressure: 1012.6, GpsInfo: "45.8400,126.4800", ImuData: "idle", LastUpdate: now.Add(-5 * time.Hour)},
		// 北京/石家庄/伊犁
		{DeviceID: "SF-FLOAT-004", Temperature: 12.5, Humidity: 55.2, Pulling: 0.42, Pressure: 1016.3, GpsInfo: "40.5280,116.9750", ImuData: "active", LastUpdate: now.Add(-1 * time.Hour)},
		{DeviceID: "SF-BOX-002", Temperature: 14.2, Humidity: 48.6, Pulling: 0.85, Pressure: 1015.5, GpsInfo: "38.3520,114.1680", ImuData: "active", LastUpdate: now.Add(-3 * time.Hour)},
		{DeviceID: "SF-FLOAT-005", Temperature: 8.0, Humidity: 35.8, Pulling: 0.0, Pressure: 1018.5, GpsInfo: "44.5800,81.1600", ImuData: "stable", LastUpdate: now.Add(-4 * time.Hour)},
	}
	DB.Create(&iotDevices)

	// ===== 18. 垂钓记录 =====
	fishingRecords := []models.FishingRecord{
		{UserID: users[3].ID, DeviceID: "SF-FLOAT-001", StartTime: now.Add(-10 * time.Hour), EndTime: now.Add(-6 * time.Hour), Latitude: 45.8050, Longitude: 126.5450},
		{UserID: users[3].ID, DeviceID: "", StartTime: now.Add(-34 * time.Hour), EndTime: now.Add(-30 * time.Hour), Latitude: 45.7550, Longitude: 126.6700},
		{UserID: users[4].ID, DeviceID: "SF-FLOAT-002", StartTime: now.Add(-26 * time.Hour), EndTime: now.Add(-22 * time.Hour), Latitude: 44.0112, Longitude: 128.9945},
		{UserID: users[5].ID, DeviceID: "", StartTime: now.Add(-48 * time.Hour), EndTime: now.Add(-44 * time.Hour), Latitude: 44.0112, Longitude: 128.9945},
		{UserID: users[6].ID, DeviceID: "SF-BOX-001", StartTime: now.Add(-8 * time.Hour), EndTime: now.Add(-4 * time.Hour), Latitude: 45.3245, Longitude: 132.4567},
		{UserID: users[7].ID, DeviceID: "", StartTime: now.Add(-72 * time.Hour), EndTime: now.Add(-68 * time.Hour), Latitude: 45.7891, Longitude: 126.5812},
		{UserID: users[8].ID, DeviceID: "SF-FLOAT-003", StartTime: now.Add(-14 * time.Hour), EndTime: now.Add(-8 * time.Hour), Latitude: 45.7891, Longitude: 126.5812},
		// 北京/石家庄/伊犁
		{UserID: users[4].ID, DeviceID: "SF-FLOAT-004", StartTime: now.Add(-12 * time.Hour), EndTime: now.Add(-8 * time.Hour), Latitude: 40.5280, Longitude: 116.9750},
		{UserID: users[5].ID, DeviceID: "SF-BOX-002", StartTime: now.Add(-20 * time.Hour), EndTime: now.Add(-16 * time.Hour), Latitude: 38.3520, Longitude: 114.1680},
		{UserID: users[6].ID, DeviceID: "SF-FLOAT-005", StartTime: now.Add(-18 * time.Hour), EndTime: now.Add(-14 * time.Hour), Latitude: 44.5800, Longitude: 81.1600},
	}
	DB.Create(&fishingRecords)

	// ===== 19. 渔获 =====
	fishCaught := []models.FishCaught{
		// fisher01 第一次记录
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-9 * time.Hour), FishType: "鲫鱼", Weight: 0.35, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.5},
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-8*time.Hour - 30*time.Minute), FishType: "鲫鱼", Weight: 0.42, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.5},
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-8 * time.Hour), FishType: "鲤鱼", Weight: 2.10, BaitType: "蚯蚓", BaitWeight: 0.02, FishingDepth: 2.0},
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-7*time.Hour - 15*time.Minute), FishType: "鲫鱼", Weight: 0.28, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.5},
		// fisher01 第二次记录
		{RecordID: fishingRecords[1].RecordID, CaughtTime: now.Add(-33 * time.Hour), FishType: "鲫鱼", Weight: 0.30, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.8},
		{RecordID: fishingRecords[1].RecordID, CaughtTime: now.Add(-32 * time.Hour), FishType: "鲤鱼", Weight: 1.80, BaitType: "玉米粒", BaitWeight: 0.05, FishingDepth: 2.5},
		// fisher02 镜泊湖
		{RecordID: fishingRecords[2].RecordID, CaughtTime: now.Add(-25 * time.Hour), FishType: "鲤鱼", Weight: 3.50, BaitType: "发酵饵", BaitWeight: 0.10, FishingDepth: 8.0},
		{RecordID: fishingRecords[2].RecordID, CaughtTime: now.Add(-24 * time.Hour), FishType: "鲢鳙", Weight: 5.20, BaitType: "雾化饵", BaitWeight: 0.15, FishingDepth: 6.0},
		{RecordID: fishingRecords[2].RecordID, CaughtTime: now.Add(-23 * time.Hour), FishType: "鲤鱼", Weight: 4.10, BaitType: "发酵饵", BaitWeight: 0.10, FishingDepth: 9.0},
		// fisher03 镜泊湖
		{RecordID: fishingRecords[3].RecordID, CaughtTime: now.Add(-47 * time.Hour), FishType: "鲤鱼", Weight: 8.00, BaitType: "发酵玉米", BaitWeight: 0.15, FishingDepth: 10.0},
		{RecordID: fishingRecords[3].RecordID, CaughtTime: now.Add(-46 * time.Hour), FishType: "鲢鳙", Weight: 4.50, BaitType: "雾化饵", BaitWeight: 0.12, FishingDepth: 7.0},
		// fisher04 兴凯湖
		{RecordID: fishingRecords[4].RecordID, CaughtTime: now.Add(-7 * time.Hour), FishType: "大白鱼", Weight: 0.60, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 1.0},
		{RecordID: fishingRecords[4].RecordID, CaughtTime: now.Add(-6*time.Hour - 40*time.Minute), FishType: "大白鱼", Weight: 0.75, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 0.8},
		{RecordID: fishingRecords[4].RecordID, CaughtTime: now.Add(-6 * time.Hour), FishType: "大白鱼", Weight: 1.20, BaitType: "金色亮片", BaitWeight: 0.005, FishingDepth: 1.2},
		// fisher05 太阳岛
		{RecordID: fishingRecords[5].RecordID, CaughtTime: now.Add(-71 * time.Hour), FishType: "鲫鱼", Weight: 0.25, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 2.0},
		{RecordID: fishingRecords[5].RecordID, CaughtTime: now.Add(-70 * time.Hour), FishType: "黄辣丁", Weight: 0.15, BaitType: "蚯蚓", BaitWeight: 0.02, FishingDepth: 2.5},
		// fisher06 太阳岛夜钓
		{RecordID: fishingRecords[6].RecordID, CaughtTime: now.Add(-12 * time.Hour), FishType: "鲫鱼", Weight: 0.30, BaitType: "拉饵", BaitWeight: 0.02, FishingDepth: 2.0},
		{RecordID: fishingRecords[6].RecordID, CaughtTime: now.Add(-11 * time.Hour), FishType: "鲤鱼", Weight: 1.50, BaitType: "搓饵", BaitWeight: 0.05, FishingDepth: 3.0},
		{RecordID: fishingRecords[6].RecordID, CaughtTime: now.Add(-10 * time.Hour), FishType: "黄辣丁", Weight: 0.12, BaitType: "蚯蚓", BaitWeight: 0.02, FishingDepth: 2.5},
		// fisher02 北京密云水库
		{RecordID: fishingRecords[7].RecordID, CaughtTime: now.Add(-11 * time.Hour), FishType: "鲤鱼", Weight: 4.20, BaitType: "玉米粒", BaitWeight: 0.08, FishingDepth: 6.0},
		{RecordID: fishingRecords[7].RecordID, CaughtTime: now.Add(-10 * time.Hour), FishType: "鲢鳙", Weight: 3.80, BaitType: "雾化饵", BaitWeight: 0.12, FishingDepth: 5.0},
		{RecordID: fishingRecords[7].RecordID, CaughtTime: now.Add(-9 * time.Hour), FishType: "翘嘴", Weight: 1.50, BaitType: "银色米诺", BaitWeight: 0.01, FishingDepth: 2.0},
		// fisher03 石家庄岗南水库
		{RecordID: fishingRecords[8].RecordID, CaughtTime: now.Add(-19 * time.Hour), FishType: "草鱼", Weight: 6.50, BaitType: "嫩玉米", BaitWeight: 0.10, FishingDepth: 3.5},
		{RecordID: fishingRecords[8].RecordID, CaughtTime: now.Add(-18 * time.Hour), FishType: "鲤鱼", Weight: 3.20, BaitType: "发酵饵", BaitWeight: 0.08, FishingDepth: 4.0},
		// fisher04 伊犁赛里木湖
		{RecordID: fishingRecords[9].RecordID, CaughtTime: now.Add(-17 * time.Hour), FishType: "高白鲑", Weight: 1.80, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 3.0},
		{RecordID: fishingRecords[9].RecordID, CaughtTime: now.Add(-16 * time.Hour), FishType: "虹鳟", Weight: 2.50, BaitType: "勺形亮片", BaitWeight: 0.005, FishingDepth: 4.0},
		{RecordID: fishingRecords[9].RecordID, CaughtTime: now.Add(-15 * time.Hour), FishType: "高白鲑", Weight: 1.20, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 2.5},
	}
	DB.Create(&fishCaught)

	log.Println("Database seeded successfully!")
}

// seedSFRIfNeeded 检查 SFR 表是否已有数据，若没有则从已有用户中提取并插入
func seedSFRIfNeeded() {
	var postCount int64
	DB.Model(&models.Post{}).Count(&postCount)
	if postCount > 0 {
		log.Println("SFR data already exists, skipping SFR seed")
		return
	}

	log.Println("Seeding SFR (community/fishing) data...")

	now := time.Now()

	// 从数据库读取已有用户（至少需要7个普通用户 index 3~9）
	var users []models.User
	DB.Order("id ASC").Find(&users)
	if len(users) < 10 {
		log.Println("Not enough users for SFR seed, need at least 10. Skipping.")
		return
	}

	// ===== 帖子 =====
	posts := []models.Post{
		{UserID: users[3].ID, Title: "今日松花江北岸冰钓，收获满满", Body: "今天天气不错，零下5度，风力2级。早上6点到的松花江北岸，冰层厚度约30cm，非常安全。\n\n用的红虫做饵，1.5号子线配4号袖钩，调4钓2。\n\n早上7点开始连续上鱼，到10点一共钓了12条鲫鱼，最大的有半斤多。中间还跑了一条大鲤鱼，估计有3斤左右，可惜子线太细被切了。\n\n总结：松花江北岸鱼情确实不错，推荐大家来试试！", Tag: "钓鱼日记"},
		{UserID: users[4].ID, Title: "路亚新手入坑指南 —— 装备选购篇", Body: "最近很多钓友问我路亚入门买什么装备好，这里整理一份新手指南：\n\n1. 竿子：推荐ML调性的直柄竿，6.6尺-7尺，适合抛投3-15g的饵。\n2. 轮子：纺车轮2000-2500型号。\n3. 线组：PE线0.6-0.8号配碳素前导线6-8磅。\n4. 假饵：先买一些亮片、软虫、米诺。\n\n预算1000元左右就能搞定一套入门装备，不要一上来就买太贵的。", Tag: "经验分享"},
		{UserID: users[5].ID, Title: "镜泊湖大坝钓场实测报告", Body: "上周末去了镜泊湖大坝钓场，这里水深8-12米，鱼种丰富。使用6.3米矶钓竿，3号主线配1.5号子线。第一天上了3条鲤鱼（最大的8斤）和5条鲢鳙。设备监测数据也很准确。", Tag: "钓点推荐"},
		{UserID: users[6].ID, Title: "兴凯湖白鱼季开始了！路亚爆护", Body: "兴凯湖北岸的白鱼季终于来了！今天用路亚拟饵在浅水区抛投，3-5克银色亮片，疯狂上鱼！从日出到上午10点，一共中了15条白鱼。", Tag: "钓鱼日记"},
		{UserID: users[7].ID, Title: "求助：台钓新手调漂总是调不准", Body: "我是刚入坑台钓的新手。按照教程调4钓2，但实际钓的时候总感觉不对劲，鱼口反应不明显。请问是浮漂选的不对，还是调法有问题？", Tag: "问答求助"},
		{UserID: users[8].ID, Title: "达亿瓦 23款天弓鲤测评", Body: "入手了一根达亿瓦23款天弓鲤4.5米，使用一个月分享感受。优点：腰力好、涂装漂亮、自重轻。缺点：价格偏高、竿稍偏硬。总体打分8.5/10。", Tag: "装备测评"},
		{UserID: users[9].ID, Title: "太阳岛西侧湖夜钓记", Body: "昨晚在太阳岛西侧湖夜钓，从下午5点一直钓到凌晨1点。战果：鲫鱼23条、鲤鱼2条、黄辣丁4条。太阳岛夜钓体验极佳，推荐！", Tag: "钓鱼日记"},
	}
	DB.Create(&posts)

	// ===== 评论 =====
	comments := []models.Comment{
		{PostID: posts[0].PostID, UserID: users[4].ID, Body: "太厉害了！松花江北岸确实出鱼"},
		{PostID: posts[0].PostID, UserID: users[5].ID, Body: "请问冰钓用什么型号的冰钻？"},
		{PostID: posts[0].PostID, UserID: users[6].ID, Body: "大鲤鱼跑了太可惜了，下次建议用1号子线"},
		{PostID: posts[1].PostID, UserID: users[3].ID, Body: "写得很详细，新手友好！"},
		{PostID: posts[1].PostID, UserID: users[7].ID, Body: "补充一下，新手PE线可以先买编织线"},
		{PostID: posts[2].PostID, UserID: users[4].ID, Body: "镜泊湖大坝确实出大物"},
		{PostID: posts[3].PostID, UserID: users[3].ID, Body: "羡慕！白鱼路亚手感确实好"},
		{PostID: posts[4].PostID, UserID: users[5].ID, Body: "新手调漂建议先在家里用水盆练习"},
		{PostID: posts[4].PostID, UserID: users[6].ID, Body: "调4钓2是经典调法，但要根据鱼口调整"},
		{PostID: posts[5].PostID, UserID: users[9].ID, Body: "天弓鲤确实不错，博5斤鲤鱼游刃有余"},
		{PostID: posts[6].PostID, UserID: users[4].ID, Body: "太阳岛夜钓氛围确实好，下次组队去"},
	}
	DB.Create(&comments)

	// ===== 子评论 =====
	cocs := []models.CommentOnComments{
		{CommentID: comments[0].CommentID, UserID: users[3].ID, Body: "哈哈主要是今天鱼口好"},
		{CommentID: comments[1].CommentID, UserID: users[3].ID, Body: "我用的手摇冰钻，15分钟就能打穿30cm"},
		{CommentID: comments[7].CommentID, UserID: users[7].ID, Body: "谢谢！我去试试换0.8g的浮漂"},
	}
	DB.Create(&cocs)

	// ===== 点赞 =====
	postLikes := []models.LikeOnPosts{
		{PostID: posts[0].PostID, UserID: users[4].ID},
		{PostID: posts[0].PostID, UserID: users[5].ID},
		{PostID: posts[0].PostID, UserID: users[6].ID},
		{PostID: posts[1].PostID, UserID: users[3].ID},
		{PostID: posts[1].PostID, UserID: users[5].ID},
		{PostID: posts[2].PostID, UserID: users[4].ID},
		{PostID: posts[3].PostID, UserID: users[3].ID},
		{PostID: posts[3].PostID, UserID: users[5].ID},
		{PostID: posts[4].PostID, UserID: users[5].ID},
		{PostID: posts[5].PostID, UserID: users[9].ID},
		{PostID: posts[6].PostID, UserID: users[3].ID},
		{PostID: posts[6].PostID, UserID: users[4].ID},
	}
	DB.Create(&postLikes)

	// ===== IoT 设备 =====
	iotDevices := []models.IoTDevice{
		{DeviceID: "SF-FLOAT-001", Temperature: 8.2, Humidity: 64.5, Pulling: 0.35, Pressure: 1013.8, GpsInfo: "45.8050,126.5450", ImuData: "stable", LastUpdate: now},
		{DeviceID: "SF-FLOAT-002", Temperature: 6.5, Humidity: 70.1, Pulling: 0.0, Pressure: 1014.2, GpsInfo: "44.0112,128.9945", ImuData: "stable", LastUpdate: now.Add(-2 * time.Hour)},
		{DeviceID: "SF-BOX-001", Temperature: 5.8, Humidity: 71.3, Pulling: 1.25, Pressure: 1015.0, GpsInfo: "45.3245,132.4567", ImuData: "active", LastUpdate: now.Add(-30 * time.Minute)},
		{DeviceID: "SF-FLOAT-003", Temperature: 9.1, Humidity: 61.8, Pulling: 0.0, Pressure: 1012.6, GpsInfo: "45.8400,126.4800", ImuData: "idle", LastUpdate: now.Add(-5 * time.Hour)},
		// 北京/石家庄/伊犁
		{DeviceID: "SF-FLOAT-004", Temperature: 12.5, Humidity: 55.2, Pulling: 0.42, Pressure: 1016.3, GpsInfo: "40.5280,116.9750", ImuData: "active", LastUpdate: now.Add(-1 * time.Hour)},
		{DeviceID: "SF-BOX-002", Temperature: 14.2, Humidity: 48.6, Pulling: 0.85, Pressure: 1015.5, GpsInfo: "38.3520,114.1680", ImuData: "active", LastUpdate: now.Add(-3 * time.Hour)},
		{DeviceID: "SF-FLOAT-005", Temperature: 8.0, Humidity: 35.8, Pulling: 0.0, Pressure: 1018.5, GpsInfo: "44.5800,81.1600", ImuData: "stable", LastUpdate: now.Add(-4 * time.Hour)},
	}
	DB.Create(&iotDevices)

	// ===== 垂钓记录 =====
	fishingRecords := []models.FishingRecord{
		{UserID: users[3].ID, DeviceID: "SF-FLOAT-001", StartTime: now.Add(-10 * time.Hour), EndTime: now.Add(-6 * time.Hour), Latitude: 45.8050, Longitude: 126.5450},
		{UserID: users[3].ID, DeviceID: "", StartTime: now.Add(-34 * time.Hour), EndTime: now.Add(-30 * time.Hour), Latitude: 45.7550, Longitude: 126.6700},
		{UserID: users[4].ID, DeviceID: "SF-FLOAT-002", StartTime: now.Add(-26 * time.Hour), EndTime: now.Add(-22 * time.Hour), Latitude: 44.0112, Longitude: 128.9945},
		{UserID: users[5].ID, DeviceID: "", StartTime: now.Add(-48 * time.Hour), EndTime: now.Add(-44 * time.Hour), Latitude: 44.0112, Longitude: 128.9945},
		{UserID: users[6].ID, DeviceID: "SF-BOX-001", StartTime: now.Add(-8 * time.Hour), EndTime: now.Add(-4 * time.Hour), Latitude: 45.3245, Longitude: 132.4567},
		// 北京/石家庄/伊犁
		{UserID: users[4].ID, DeviceID: "SF-FLOAT-004", StartTime: now.Add(-12 * time.Hour), EndTime: now.Add(-8 * time.Hour), Latitude: 40.5280, Longitude: 116.9750},
		{UserID: users[5].ID, DeviceID: "SF-BOX-002", StartTime: now.Add(-20 * time.Hour), EndTime: now.Add(-16 * time.Hour), Latitude: 38.3520, Longitude: 114.1680},
		{UserID: users[6].ID, DeviceID: "SF-FLOAT-005", StartTime: now.Add(-18 * time.Hour), EndTime: now.Add(-14 * time.Hour), Latitude: 44.5800, Longitude: 81.1600},
	}
	DB.Create(&fishingRecords)

	// ===== 渔获 =====
	fishCaught := []models.FishCaught{
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-9 * time.Hour), FishType: "鲫鱼", Weight: 0.35, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.5},
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-8 * time.Hour), FishType: "鲤鱼", Weight: 2.10, BaitType: "蚯蚓", BaitWeight: 0.02, FishingDepth: 2.0},
		{RecordID: fishingRecords[0].RecordID, CaughtTime: now.Add(-7 * time.Hour), FishType: "鲫鱼", Weight: 0.28, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.5},
		{RecordID: fishingRecords[1].RecordID, CaughtTime: now.Add(-33 * time.Hour), FishType: "鲫鱼", Weight: 0.30, BaitType: "红虫", BaitWeight: 0.01, FishingDepth: 1.8},
		{RecordID: fishingRecords[1].RecordID, CaughtTime: now.Add(-32 * time.Hour), FishType: "鲤鱼", Weight: 1.80, BaitType: "玉米粒", BaitWeight: 0.05, FishingDepth: 2.5},
		{RecordID: fishingRecords[2].RecordID, CaughtTime: now.Add(-25 * time.Hour), FishType: "鲤鱼", Weight: 3.50, BaitType: "发酵饵", BaitWeight: 0.10, FishingDepth: 8.0},
		{RecordID: fishingRecords[2].RecordID, CaughtTime: now.Add(-24 * time.Hour), FishType: "鲢鳙", Weight: 5.20, BaitType: "雾化饵", BaitWeight: 0.15, FishingDepth: 6.0},
		{RecordID: fishingRecords[3].RecordID, CaughtTime: now.Add(-47 * time.Hour), FishType: "鲤鱼", Weight: 8.00, BaitType: "发酵玉米", BaitWeight: 0.15, FishingDepth: 10.0},
		{RecordID: fishingRecords[4].RecordID, CaughtTime: now.Add(-7 * time.Hour), FishType: "大白鱼", Weight: 0.60, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 1.0},
		{RecordID: fishingRecords[4].RecordID, CaughtTime: now.Add(-6 * time.Hour), FishType: "大白鱼", Weight: 1.20, BaitType: "金色亮片", BaitWeight: 0.005, FishingDepth: 1.2},
		// 北京密云
		{RecordID: fishingRecords[5].RecordID, CaughtTime: now.Add(-11 * time.Hour), FishType: "鲤鱼", Weight: 4.20, BaitType: "玉米粒", BaitWeight: 0.08, FishingDepth: 6.0},
		{RecordID: fishingRecords[5].RecordID, CaughtTime: now.Add(-10 * time.Hour), FishType: "翘嘴", Weight: 1.50, BaitType: "银色米诺", BaitWeight: 0.01, FishingDepth: 2.0},
		// 石家庄岗南
		{RecordID: fishingRecords[6].RecordID, CaughtTime: now.Add(-19 * time.Hour), FishType: "草鱼", Weight: 6.50, BaitType: "嫩玉米", BaitWeight: 0.10, FishingDepth: 3.5},
		{RecordID: fishingRecords[6].RecordID, CaughtTime: now.Add(-18 * time.Hour), FishType: "鲤鱼", Weight: 3.20, BaitType: "发酵饵", BaitWeight: 0.08, FishingDepth: 4.0},
		// 伊犁赛里木湖
		{RecordID: fishingRecords[7].RecordID, CaughtTime: now.Add(-17 * time.Hour), FishType: "高白鲑", Weight: 1.80, BaitType: "银色亮片", BaitWeight: 0.005, FishingDepth: 3.0},
		{RecordID: fishingRecords[7].RecordID, CaughtTime: now.Add(-16 * time.Hour), FishType: "虹鳟", Weight: 2.50, BaitType: "勺形亮片", BaitWeight: 0.005, FishingDepth: 4.0},
	}
	DB.Create(&fishCaught)

	log.Println("SFR community data seeded successfully!")
}
