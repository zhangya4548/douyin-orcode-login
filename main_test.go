// @Time : 2024/7/16 下午3:32
// @Author : zhangguangqiang
// @File : main_test.go
// @Software: GoLand

package main

import "testing"

func Test_getSession(t *testing.T) {

	session, err := getSession()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(session)
}

func Test_getSession2(t *testing.T) {

	session, err := getSession2()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(session)
}

func Test_Wg(t *testing.T) {
	Wg()
}
