package database

import (
	"log"
	"math"
	"time"

	"smart-fish/back_end/models"

	"golang.org/x/crypto/bcrypt"
)

// Seed 填充示例数据（开发环境使用）
func Seed() {
	// 检查是否已有数据
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already has data, skipping seed")
		return
	}

	log.Println("Seeding database with sample data...")

	now := time.Now()

	// ===== 1. 创建用户 =====
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	staffHash, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	fisherHash, _ := bcrypt.GenerateFromPassword([]byte("fisher666"), bcrypt.DefaultCost)

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
	}
	DB.Create(&devices)

	// ===== 5. 创建垂钓水域 =====
	spots := []models.FishingSpot{
		{Name: "松花江北岸钓场", RegionID: regions[0].ID, Description: "松花江北岸，适合路亚和台钓，冬季冰钓热门地点", Latitude: 45.7772, Longitude: 126.6177, WaterType: "river", Capacity: 80, Status: "open", BoundDeviceID: &devices[0].ID},
		{Name: "松花江公路大桥钓位", RegionID: regions[0].ID, Description: "公路大桥下游，深水区鱼种丰富，常有大鲤鱼出没", Latitude: 45.7685, Longitude: 126.6089, WaterType: "river", Capacity: 40, Status: "open", BoundDeviceID: &devices[3].ID},
		{Name: "松花江铁路桥野钓区", RegionID: regions[0].ID, Description: "铁路桥附近自然钓场，水流较急，适合有经验的钓友", Latitude: 45.7710, Longitude: 126.6250, WaterType: "river", Capacity: 30, Status: "open", BoundDeviceID: &devices[4].ID},
		{Name: "太阳岛西侧湖钓场", RegionID: regions[1].ID, Description: "太阳岛西侧静水湖，水质清澈，适合休闲垂钓", Latitude: 45.7891, Longitude: 126.5812, WaterType: "lake", Capacity: 60, Status: "open", BoundDeviceID: &devices[5].ID},
		{Name: "太阳岛码头钓位", RegionID: regions[1].ID, Description: "老码头附近，水深2-4米，鲫鱼鲤鱼密度高", Latitude: 45.7856, Longitude: 126.5923, WaterType: "river", Capacity: 30, Status: "open", BoundDeviceID: &devices[7].ID},
		{Name: "龙凤湿地外围钓场", RegionID: regions[2].ID, Description: "湿地外围开放区域，冬季关闭", Latitude: 46.5812, Longitude: 125.0934, WaterType: "lake", Capacity: 50, Status: "closed", BoundDeviceID: &devices[8].ID},
		{Name: "镜泊湖大坝钓场", RegionID: regions[3].ID, Description: "镜泊湖大坝下游，大鱼频出，专业钓手天堂", Latitude: 44.0112, Longitude: 128.9945, WaterType: "reservoir", Capacity: 100, Status: "open", BoundDeviceID: &devices[9].ID},
		{Name: "镜泊湖吊水楼瀑布钓区", RegionID: regions[3].ID, Description: "吊水楼瀑布下方深潭，水流交汇处鱼群聚集", Latitude: 43.9856, Longitude: 129.0123, WaterType: "reservoir", Capacity: 60, Status: "open", BoundDeviceID: &devices[10].ID},
		{Name: "兴凯湖北岸垂钓基地", RegionID: regions[4].ID, Description: "兴凯湖北岸专业垂钓基地，盛产大白鱼、鲤鱼、鲢鳙", Latitude: 45.3245, Longitude: 132.4567, WaterType: "lake", Capacity: 120, Status: "open", BoundDeviceID: &devices[12].ID},
		{Name: "扎龙湿地观鸟钓场", RegionID: regions[5].ID, Description: "扎龙保护区外围，可边观鸟边垂钓，别有趣味", Latitude: 47.1789, Longitude: 124.2856, WaterType: "pond", Capacity: 40, Status: "maintenance", BoundDeviceID: &devices[14].ID},
		{Name: "五大连池药泉湖钓场", RegionID: regions[6].ID, Description: "药泉湖垂钓区，矿泉水养殖鱼口感鲜美", Latitude: 48.7523, Longitude: 126.1189, WaterType: "lake", Capacity: 45, Status: "open", BoundDeviceID: &devices[15].ID},
		{Name: "佳木斯松花江外滩钓场", RegionID: regions[7].ID, Description: "佳木斯外滩公园段，休闲垂钓好去处，交通便利", Latitude: 46.8012, Longitude: 130.3678, WaterType: "river", Capacity: 70, Status: "open"},
	}
	DB.Create(&spots)

	// open 状态的水域索引（用于生成环境/历史数据）
	// spots: 0(open),1(open),2(open),3(open),4(open),5(closed),6(open),7(open),8(open),9(maintenance),10(open),11(open)
	openSpotIndices := []int{0, 1, 2, 3, 4, 6, 7, 8, 10, 11}

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
	underwaterDeviceIDs := []uint{devices[1].ID, devices[4].ID, devices[6].ID, devices[10].ID, devices[13].ID}
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
	}
	DB.Create(&notices)

	// 关联通知到水域
	DB.Model(&notices[0]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2])
	DB.Model(&notices[1]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2], &spots[3], &spots[4], &spots[6], &spots[7], &spots[8], &spots[10], &spots[11])
	DB.Model(&notices[3]).Association("RelatedSpots").Append(&spots[8])
	DB.Model(&notices[4]).Association("RelatedSpots").Append(&spots[0], &spots[1], &spots[2], &spots[11])

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
	}
	DB.Create(&suggestions)

	log.Println("Database seeded successfully!")
}
