package bazi

import (
	"fmt"
	"log"
)

// 八字
type TBazi struct {
	SolarDate   TDate       // 新历日期
	LunarDate   TDate       // 农历日期
	BaziDate    TDate       // 八字日期
	PreviousJie TDate       // 上一个节(气)
	NextJie     TDate       // 下一个节(气)
	SiZhu       TSiZhu      // 四柱
	XiYong      TXiYong     // 喜用神
	DaYun       TDaYun      // 大运
	HeHuaChong  THeHuaChong // 合化冲
}

// 计算
func calc(bazi *TBazi, nSex int) {
	// 通过立春获取当年的年份
	bazi.BaziDate.Year = GetLiChun2(bazi.SolarDate)
	// 通过节气获取当前后的两个节
	bazi.PreviousJie, bazi.NextJie = GetJieQi(bazi.SolarDate)
	// 八字所在的节气是上一个的节气
	bazi.BaziDate.JieQi = bazi.PreviousJie.JieQi
	// 节气0 是立春 是1月
	bazi.BaziDate.Month = bazi.BaziDate.JieQi/2 + 1

	// 通过八字年来获取年柱
	bazi.SiZhu.YearZhu = GetZhuFromYear(bazi.BaziDate.Year)
	// 通过年干支和八字月
	bazi.SiZhu.MonthZhu = GetZhuFromMonth(bazi.BaziDate.Month, bazi.SiZhu.YearZhu.Gan.Value)

	// 通过公历 年月日计算日柱
	bazi.SiZhu.DayZhu = GetZhuFromDay(bazi.SolarDate.Year, bazi.SolarDate.Month, bazi.SolarDate.Day, bazi.SolarDate.Hour)

	//获取时柱
	bazi.SiZhu.HourZhu = GetZhuFromHour(bazi.SolarDate.Hour, bazi.SiZhu.DayZhu.Gan.Value)

	// 计算十神
	CalcShiShen(&bazi.SiZhu)
	// 计算纳音
	CalcNaYin(&bazi.SiZhu)

	// 检查合化冲
	CheckHeHuaChong(&bazi.SiZhu, &bazi.HeHuaChong)

	// 计算大运
	bazi.DaYun = CalcDaYun(&bazi.SiZhu, nSex)

	// 计算起运时间
	CalcQiYun(&bazi.DaYun, bazi.PreviousJie, bazi.NextJie, bazi.SolarDate)

	// 计算喜用神
	bazi.XiYong = CalcXiYong(&bazi.SiZhu)
}

// 从新历获取八字(年, 月, 日, 时, 分, 秒, 性别男1,女0)
// Get the eight characters from the new calendar (year, month, day, hour, minute, second, gender male 1, female 0)
func GetBazi(nYear, nMonth, nDay, nHour, nMinute, nSecond, nSex int) TBazi {
	var bazi TBazi

	if !GetDateIsValid(nYear, nMonth, nDay) {
		log.Println("无效的日期", nYear, nMonth, nDay)
		return bazi
	}

	// 新历年
	bazi.SolarDate.Year = nYear
	bazi.SolarDate.Month = nMonth
	bazi.SolarDate.Day = nDay
	bazi.SolarDate.Hour = nHour
	bazi.SolarDate.Minute = nMinute
	bazi.SolarDate.Second = nSecond

	// 转农历
	var nTimeStamp = Get64TimeStamp(nYear, nMonth, nDay, nHour, nMinute, nSecond)
	bazi.LunarDate = GetLunarDateFrom64TimeStamp(nTimeStamp)

	// 进行计算
	calc(&bazi, nSex)

	return bazi
}

// 从农历获取八字
// Get the eight characters from the lunar calendar
func GetBaziFromLunar(nYear, nMonth, nDay, nHour, nMinute, nSecond, nSex int, isLeap bool) TBazi {
	nYear, nMonth = ChangeLunarLeap(nYear, nMonth, isLeap)

	var bazi TBazi

	if !GetLunarDateIsValid(nYear, nMonth, nDay) {
		log.Println("无效的日期", nYear, nMonth, nDay)
		return bazi
	}

	// 农历年
	bazi.LunarDate.Year = nYear
	bazi.LunarDate.Month = nMonth
	bazi.LunarDate.Day = nDay
	bazi.LunarDate.Hour = nHour
	bazi.LunarDate.Minute = nMinute
	bazi.LunarDate.Second = nSecond

	// 转新历
	var nTimeStamp = GetLunar64TimeStamp(nYear, nMonth, nDay, nHour, nMinute, nSecond)
	bazi.LunarDate = GetDateFrom64TimeStamp(nTimeStamp)

	// 进行计算
	calc(&bazi, nSex)

	return bazi

}

func (bazi *TBazi) String() string {
	return fmt.Sprintf("%s %s %s %s",
		bazi.SiZhu.YearZhu.GanZhi.ToString(),
		bazi.SiZhu.MonthZhu.GanZhi.ToString(),
		bazi.SiZhu.DayZhu.GanZhi.ToString(),
		bazi.SiZhu.HourZhu.GanZhi.ToString())
}

func (bazi *TBazi) Print() {
	PrintBazi(bazi)
}

// 按照特殊格式化输出(未完成)
// Output according to special format (unfinished)
func PrintBazi(bazi *TBazi) {
	log.Println("======================================================================")
	// New Calendar of Birth Date
	log.Println("出生日期新历 (New Calendar of Birth Date)： ", bazi.SolarDate.Year, "年(year)",
		bazi.SolarDate.Month, "月(month)",
		bazi.SolarDate.Day, "日(day)  ",
		bazi.SolarDate.Hour, ":",
		bazi.SolarDate.Minute, ":",
		bazi.SolarDate.Second,
	)
	// Basic horoscope:
	log.Println("基本八字 (Basic horoscope)：",
		bazi.SiZhu.YearZhu.GanZhi.ToString(),
		bazi.SiZhu.MonthZhu.GanZhi.ToString(),
		bazi.SiZhu.DayZhu.GanZhi.ToString(),
		bazi.SiZhu.HourZhu.GanZhi.ToString())

	//	ife chart analysis:
	log.Println("命盘解析 (Life chart analysis:)：")
	log.Println(
		bazi.SiZhu.YearZhu.Gan.ToString(), "(",
		bazi.SiZhu.YearZhu.Gan.WuXing.ToString(), ")[",
		bazi.SiZhu.YearZhu.Gan.ShiShen.ToString(), "]\t",
		bazi.SiZhu.MonthZhu.Gan.ToString(), "(",
		bazi.SiZhu.MonthZhu.Gan.WuXing.ToString(), ")[",
		bazi.SiZhu.MonthZhu.Gan.ShiShen.ToString(), "]\t",
		bazi.SiZhu.DayZhu.Gan.ToString(), "(",
		bazi.SiZhu.DayZhu.Gan.WuXing.ToString(), ")[日主(Sun Lord)]\t",
		bazi.SiZhu.HourZhu.Gan.ToString(), "(",
		bazi.SiZhu.HourZhu.Gan.WuXing.ToString(), ")[",
		bazi.SiZhu.HourZhu.Gan.ShiShen.ToString(), "]")
	log.Println(
		bazi.SiZhu.YearZhu.Zhi.ToString(), "(",
		bazi.SiZhu.YearZhu.Zhi.WuXing.ToString(), ")",
		bazi.SiZhu.YearZhu.Zhi.CangGan[0].ShiShen.ToString(),
		bazi.SiZhu.YearZhu.Zhi.CangGan[1].ShiShen.ToString(),
		bazi.SiZhu.YearZhu.Zhi.CangGan[2].ShiShen.ToString(),
		"\t",
		bazi.SiZhu.MonthZhu.Zhi.ToString(), "(",
		bazi.SiZhu.MonthZhu.Zhi.WuXing.ToString(), ")",
		bazi.SiZhu.MonthZhu.Zhi.CangGan[0].ShiShen.ToString(),
		bazi.SiZhu.MonthZhu.Zhi.CangGan[1].ShiShen.ToString(),
		bazi.SiZhu.MonthZhu.Zhi.CangGan[2].ShiShen.ToString(),
		"\t",
		bazi.SiZhu.DayZhu.Zhi.ToString(), "(",
		bazi.SiZhu.DayZhu.Zhi.WuXing.ToString(), ")",
		bazi.SiZhu.DayZhu.Zhi.CangGan[0].ShiShen.ToString(),
		bazi.SiZhu.DayZhu.Zhi.CangGan[1].ShiShen.ToString(),
		bazi.SiZhu.DayZhu.Zhi.CangGan[2].ShiShen.ToString(),
		"\t",
		bazi.SiZhu.HourZhu.Zhi.ToString(), "(",
		bazi.SiZhu.HourZhu.Zhi.WuXing.ToString(), ")",
		bazi.SiZhu.HourZhu.Zhi.CangGan[0].ShiShen.ToString(),
		bazi.SiZhu.HourZhu.Zhi.CangGan[1].ShiShen.ToString(),
		bazi.SiZhu.HourZhu.Zhi.CangGan[2].ShiShen.ToString())

	log.Println(
		bazi.SiZhu.YearZhu.GanZhi.NaYin.ToString(), "\t\t",
		bazi.SiZhu.MonthZhu.GanZhi.NaYin.ToString(), "\t\t",
		bazi.SiZhu.DayZhu.GanZhi.NaYin.ToString(), "\t\t",
		bazi.SiZhu.HourZhu.GanZhi.NaYin.ToString(), "\t\t",
	)
	// 天干五合
	// Tianganwuhe
	log.Println(
		"天干五合(Tianganwuhe)：",
		bazi.HeHuaChong.TgWuHe[0].Str,
		bazi.HeHuaChong.TgWuHe[1].Str,
		bazi.HeHuaChong.TgWuHe[2].Str,
		bazi.HeHuaChong.TgWuHe[3].Str,
	)
	// 地支三会
	// Earthly Branch Three
	log.Println(
		"地支三会(Earthly Branch Three)：",
		bazi.HeHuaChong.DzSanHui[0].Str,
		bazi.HeHuaChong.DzSanHui[1].Str,
	)
	// 地支三合
	// Earthly Branch Sanhe
	log.Println(
		"地支三合(Earthly Branch Sanhe)：",
		bazi.HeHuaChong.DzSanHe[0].Str,
		bazi.HeHuaChong.DzSanHe[1].Str,
	)

	// 地支六冲
	// Earthly Branches and Six Chongs
	log.Println(
		"地支六冲(Earthly Branches and Six Chongs)：",
		bazi.HeHuaChong.DzLiuChong[0].Str,
		bazi.HeHuaChong.DzLiuChong[1].Str,
		bazi.HeHuaChong.DzLiuChong[2].Str,
		bazi.HeHuaChong.DzLiuChong[3].Str,
	)
	// 地支六害
	// Six evils
	log.Println(
		"地支六害(Six evils)：",
		bazi.HeHuaChong.DzLiuHai[0].Str,
		bazi.HeHuaChong.DzLiuHai[1].Str,
		bazi.HeHuaChong.DzLiuHai[2].Str,
		bazi.HeHuaChong.DzLiuHai[3].Str,
	)

	log.Println("----------------------------------------------------------------------")
	log.Println("所属节令(Season)：")
	log.Println(GetJieQiFromNumber(bazi.PreviousJie.JieQi), bazi.PreviousJie.Year, "年(year)",
		bazi.PreviousJie.Month, "月(month)",
		bazi.PreviousJie.Day, "日(year)  ",
		bazi.PreviousJie.Hour, ":",
		bazi.PreviousJie.Minute, ":",
		bazi.PreviousJie.Second)
	log.Println(GetJieQiFromNumber(bazi.NextJie.JieQi), bazi.NextJie.Year, "年(year)",
		bazi.NextJie.Month, "月(month)",
		bazi.NextJie.Day, "日(day)  ",
		bazi.NextJie.Hour, ":",
		bazi.NextJie.Minute, ":",
		bazi.NextJie.Second)

	var szDaYun = "大运(Universiade)："
	for i := 0; i < 10; i++ {
		szDaYun = szDaYun + " " + bazi.DaYun.Zhu[i].GanZhi.ToString()
	}
	log.Println(szDaYun)
	log.Println("起运时间(Departure Time)", bazi.DaYun.QiYun.Year, "年(year)",
		bazi.DaYun.QiYun.Month, "月(month)",
		bazi.DaYun.QiYun.Day, "日(day)  ",
		bazi.DaYun.QiYun.Hour, ":",
		bazi.DaYun.QiYun.Minute, ":",
		bazi.DaYun.QiYun.Second)
	log.Println("======================================================================")
}
