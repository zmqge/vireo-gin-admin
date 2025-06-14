package utils

/*
 * @Author: 27
 * @Description: 通用工具类
 * @Date: 2024-08-20 16:51:23
 * @LastEditors: 27
 * @LastEditTime: 2024-08-20 16:51:23
 */

import (
	"errors"
	"strconv"
)

// parseUintID 将字符串ID转换为uint类型，并返回友好的错误信息
func ParseUintID(idStr string) (uint, error) {
	// 检查输入是否为空
	if idStr == "" {
		return 0, errors.New("ID不能为空")
	}

	// 转换为无符号整数（支持十进制，64位）
	num, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 解析错误时，区分不同错误类型
		if numErr, ok := err.(*strconv.NumError); ok {
			switch numErr.Err {
			case strconv.ErrSyntax: // 语法错误（非数字字符）
				return 0, errors.New("ID格式错误，必须为正整数")
			case strconv.ErrRange: // 数值超出范围（如超过uint最大值）
				return 0, errors.New("ID数值过大，超出允许范围")
			}
		}
		// 其他未知错误
		return 0, errors.New("ID无效，请检查格式")
	}

	// 转换成功，返回uint类型
	return uint(num), nil
}
