package httptool

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-irain/logger"
)

//Post http post下发
func Post(logid, url, param string, flags ...bool) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*10) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(30 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	alert := false
	if len(flags) > 0 {
		alert = flags[0]
	}
	retry := false
	if len(flags) > 1 {
		retry = flags[1]
	}
	logger.Debug(logid, "Post: url: ", url, " param : ", param)
	reqest, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		if alert {
			logger.Error(logid, "remote", err.Error())
		} else {
			logger.Error(logid, err.Error())
		}
		if retry {
			Put(logid, url, param, 1)
		}
		return "", errors.New("server_inner_error")
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	response, err := client.Do(reqest)
	if err != nil {
		if alert {
			logger.Error(logid, "remote", err.Error())
		} else {
			logger.Error(logid, err.Error())
		}
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("request_time_out")
		} else {
			err = errors.New("server_inner_error")
		}
		if retry {
			Put(logid, url, param, 1)
		}
		return "", err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if alert {
			logger.Error(logid, "remote", err.Error())
		} else {
			logger.Error(logid, err.Error())
		}
		if retry {
			Put(logid, url, param, 1)
		}
		return "", errors.New("server_inner_error")
	}
	return string(result), nil
}

//PostData http rowdata请求
func PostData(url, param string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*5) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				c.SetDeadline(time.Now().Add(10 * time.Second)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
	logger.Debug("Post: url: ", url, " param : ", param)
	request, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		logger.Error("remote", err.Error())
		return "", errors.New("SERVER_INNER_ERROR")
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("User-Agent", "iRainService")
	request.Header.Set("Accept-Version", "v2.0")
	response, err := client.Do(request)
	if err != nil {
		logger.Error("remote", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("REQUEST_TIME_OUT")
		} else {
			err = errors.New("SERVER_INNER_ERROR")
		}
		return "", err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("remote", err.Error())
		return "", errors.New("SERVER_INNER_ERROR")
	}
	return string(result), nil
}

//Get http get请求
func Get(logid, url, param string, retrys ...bool) (string, error) {
	retry := false
	if len(retrys) > 0 {
		retry = retrys[0]
	}
	uri := url + "?" + param
	resp, err := http.Get(uri)
	if err != nil {
		logger.Error("remote", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("request_time_out")
		} else {
			err = errors.New("server_inner_error")
		}
		if retry {
			Put(logid, url, param, 2)
		}
		return "", err
	}
	defer resp.Body.Close()
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("remote", err.Error())
		if retry {
			Put(logid, url, param, 2)
		}
		return "", errors.New("server_inner_error")
	}

	return string(reply), nil
}

//SimplePost 简单http post 请求
func SimplePost(url, param string) (string, error) {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(param))
	if err != nil {
		logger.Error("remote", err.Error())
		if strings.Contains(err.Error(), "timeout") {
			err = errors.New("request_time_out")
		} else {
			err = errors.New("server_inner_error")
		}
		return "", err
	}

	defer resp.Body.Close()
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("remote", err.Error())
		return "", errors.New("server_inner_error")
	}

	return string(reply), nil
}
