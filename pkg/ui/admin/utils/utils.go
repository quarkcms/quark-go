package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/quarkcms/quark-go/pkg/framework/token"
)

// 获取管理员Token
func GetAdminToken(c *fiber.Ctx) string {
	header := c.GetReqHeaders()
	getToken := strings.Split(header["Authorization"], " ")

	if len(getToken) != 2 {
		return ""
	}

	return getToken[1]
}

// 获取当前登录管理员信息
func Admin(c *fiber.Ctx, field string) interface{} {
	getToken := GetAdminToken(c)
	userInfo, err := token.Parse(getToken)
	if err != nil {
		return nil
	}

	return userInfo[field]
}

// 数据集转换成Tree
func ListToTree(list []interface{}, pk string, pid string, child string, root float64) []interface{} {
	var treeList []interface{}
	for _, v := range list {
		if v.(map[string]interface{})["pid"] == root {
			childNode := ListToTree(list, pk, pid, child, v.(map[string]interface{})[pk].(float64))
			if childNode != nil {
				v.(map[string]interface{})[child] = childNode
			}
			treeList = append(treeList, v)
		}
	}

	return treeList
}

// Tree转换为有序列表
func TreeToOrderedList(tree []interface{}, level int, field string, child string) []interface{} {
	var list []interface{}
	for _, v := range tree {
		v.(map[string]string)[field] = strings.Repeat("—", level) + v.(map[string]string)[field]
		if v.(map[string][]interface{})[child] != nil {
			children := TreeToOrderedList(v.(map[string][]interface{})[child], level, field, child)
			v = append(v.([]interface{}), children)
		}
		list = append(list, v)
	}

	return list
}

// struct转map
func StructToMap(v any) any {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	var mapResult any
	json.Unmarshal(jsonBytes, &mapResult)

	return mapResult
}

// 存储权限
var Permissions []string

// 设置权限
func SetPermissions(permissions []string) {
	Permissions = permissions
}

// 获取权限
func GetPermissions() []string {
	return Permissions
}
