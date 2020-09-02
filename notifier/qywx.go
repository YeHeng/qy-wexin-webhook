package notifier

import (
	"bytes"
	"encoding/json"
	"github.com/YeHeng/qy-wexin-webhook/model"
	"github.com/YeHeng/qy-wexin-webhook/transformer"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Send send markdown message to dingtalk
func Send(notification model.Notification, defaultRobot string, logger *logrus.Logger) (model.ResultVo, error) {

	markdown, robotURL, err := transformer.TransformToMarkdown(notification)

	if err != nil {
		return model.ResultVo{
				Code:    400,
				Message: "marshal json fail " + err.Error(),
			},
			nil
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return model.ResultVo{
				Code:    400,
				Message: "marshal json fail " + err.Error(),
			},
			nil
	}

	var qywxRobotURL string

	if robotURL != "" {
		qywxRobotURL = robotURL
	} else {
		qywxRobotURL = defaultRobot
	}

	if len(qywxRobotURL) == 0 {
		return model.ResultVo{
				Code:    404,
				Message: "robot url is nil",
			},
			nil
	}

	req, err := http.NewRequest(
		"POST",
		qywxRobotURL,
		bytes.NewBuffer(data))

	if err != nil {
		return model.ResultVo{
				Code:    400,
				Message: "request robot url fail " + err.Error(),
			},
			nil
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return model.ResultVo{
				Code:    404,
				Message: "request wx api url fail " + err.Error(),
			},
			nil
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	bodyString := string(bodyBytes)
	logger.Debugf("response: %s, header: %s", bodyString, resp.Header)

	return model.ResultVo{
		Code:    resp.StatusCode,
		Message: bodyString,
	}, nil
}
