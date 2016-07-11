package util

import (
	"errors"
	"fmt"
	"logger"
	"regexp"
	"strconv"
	"strings"
)

func GetCidrInfo(network string, logger logger.Logger) (map[string]string, error) {
	network = strings.TrimSpace(network)
	result := make(map[string]string)
	list := strings.Split(network, "/")
	if len(list) != 2 {
		return result, errors.New("网段格式不正确")
	}

	isValidate, err := regexp.MatchString("^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$", list[0])
	if err != nil {
		return result, err
	}

	if !isValidate {
		return result, errors.New("IP格式不正确")
	}

	logger.Debugf("get network info:%s", network)

	minIp, maxIp := GetCidrIpRange(network)
	logger.Debugf("get ip range:%s~%s", minIp, maxIp)
	result["MinIP"] = minIp
	result["MaxIP"] = maxIp

	mask, err := strconv.Atoi(list[1])
	if err != nil {
		return result, err
	}

	logger.Debugf("get mask:%d", mask)
	if mask <= 0 || mask >= 32 {
		logger.Error("掩码位不正确!")
		return result, errors.New("掩码位不正确!")
	}

	result["Mask"] = GetCidrIpMask(mask)
	result["IPNum"] = fmt.Sprintf("%d", GetCidrHostNum(mask))
	logger.Debugf("ip nums:%s", result["IPNum"])
	return result, nil
}

func GetCidrIpRange(cidr string) (string, string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIp, seg3MaxIp := GetIpSeg3Range(ipSegs, maskLen)
	seg4MinIp, seg4MaxIp := GetIpSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return ipPrefix + strconv.Itoa(seg3MinIp) + "." + strconv.Itoa(seg4MinIp),
		ipPrefix + strconv.Itoa(seg3MaxIp) + "." + strconv.Itoa(seg4MaxIp)
}

//计算得到CIDR地址范围内可拥有的主机数量
func GetCidrHostNum(maskLen int) uint {
	cidrIpNum := uint(0)
	var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
}

//获取Cidr的掩码
func GetCidrIpMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	//fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIpSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSegs[2])
		return segIp, segIp
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return GetIpSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func GetIpSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIp, segMaxIp := GetIpSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp - 1
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func GetIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}

func GetIPListByMinAndMaxIP(min string, max string) ([]string, error) {
	var result []string
	list1 := strings.Split(min, ".")
	list2 := strings.Split(max, ".")
	min3, err := strconv.Atoi(list1[2])
	if err != nil {
		return result, err
	}

	min4, err := strconv.Atoi(list1[3])
	if err != nil {
		return result, err
	}

	max3, err := strconv.Atoi(list2[2])
	if err != nil {
		return result, err
	}

	max4, err := strconv.Atoi(list2[3])
	if err != nil {
		return result, err
	}

	for i := min3; i <= max3; i++ {
		for j := min4; j <= max4; j++ {
			ip := list1[0] + "." + list1[1] + "." + fmt.Sprintf("%d", i) + "." + fmt.Sprintf("%d", j)
			result = append(result, ip)
		}
	}
	return result, nil
}
