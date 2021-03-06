package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	// . "github.com/smartystreets/goconvey/convey"

	"github.com/golang-devops/go-psexec/client"
	"github.com/golang-devops/go-psexec/shared"
)

func testsGetFilePath(fileName string) (string, error) {
	return filepath.Abs("../shared/testdata/" + fileName)
}

type tmpTestsLogger struct {
	sync.RWMutex
	ErrorList []string
}

func (t *tmpTestsLogger) Info(v ...interface{}) error {
	// fmt.Println(v...)
	return nil
}
func (t *tmpTestsLogger) Infof(frmt string, a ...interface{}) error {
	// fmt.Println(fmt.Sprintf(frmt, a...))
	return nil
}
func (t *tmpTestsLogger) Warning(v ...interface{}) error {
	t.Lock()
	defer t.Unlock()
	// fmt.Println(v...)
	t.ErrorList = append(t.ErrorList, fmt.Sprintln(v...))
	return nil
}
func (t *tmpTestsLogger) Warningf(frmt string, a ...interface{}) error {
	t.Warning(fmt.Sprintln(fmt.Sprintf(frmt, a...)))
	return nil
}
func (t *tmpTestsLogger) Error(v ...interface{}) error {
	t.Lock()
	defer t.Unlock()

	// fmt.Println(v...)
	t.ErrorList = append(t.ErrorList, fmt.Sprintln(v...))
	return nil
}
func (t *tmpTestsLogger) Errorf(frmt string, a ...interface{}) error {
	t.Error(fmt.Sprintln(fmt.Sprintf(frmt, a...)))
	return nil
}

func setupClient(clientPemFile string) (*client.Client, error) {
	pvtKey, err := shared.ReadPemKey(clientPemFile)
	if err != nil {
		return nil, err
	}
	return client.New(pvtKey), nil
}

func cleanFeedbackLine(line string) string {
	return strings.Trim(line, " \"'")
}

func doRequest(wg *sync.WaitGroup, logger *tmpTestsLogger, index int, cl *client.Client, serverBaseUrl string) {
	defer wg.Done()

	session, err := cl.RequestNewSession(serverBaseUrl)
	if err != nil {
		logger.Errorf("Index %d (RequestNewSession) err: %s", index, err.Error())
		return
	}

	echoStr := fmt.Sprintf("Hallo (%d)", index)
	resp, err := session.ExecRequestBuilder().Winshell().Exe("echo").Args(echoStr).BuildAndDoRequest()
	if err != nil {
		logger.Errorf("Index %d (ExecRequestBuilder echo) err: %s", index, err.Error())
		return
	}

	responseChannel, errChannel := resp.TextResponseChannel()

	lines := []string{}
	errors := []error{}
outerFor:
	for {
		select {
		case feedbackLine, ok := <-responseChannel:
			if !ok {
				break outerFor
			}
			lines = append(lines, feedbackLine)
		case errLine, ok := <-errChannel:
			if !ok {
				break outerFor
			}
			errors = append(errors, errLine)
		}
	}

	if len(errors) > 0 {
		errStrs := []string{}
		for _, e := range errors {
			errStrs = append(errStrs, e.Error())
		}
		logger.Errorf("ERRORS OCCURRED: %s", strings.Join(errStrs, "\n"))
	}

	expectedFeedback := []string{echoStr, shared.RESPONSE_EOF}
	if len(lines) != len(expectedFeedback) {
		logger.Errorf("Index %d expected was %#v, but actual was %#v", index, expectedFeedback, lines)
		return
	}

	for i, expLine := range expectedFeedback {
		if cleanFeedbackLine(lines[i]) != cleanFeedbackLine(expLine) {
			logger.Errorf("Index %d expected was %#v, but actual was %#v", index, expectedFeedback, lines)
			return
		}
	}
}

func TestHighLoad(t *testing.T) {
	/*Convey("Test HighLoad", t, func() {
		logger := &tmpTestsLogger{ErrorList: []string{}}
		a := &app{}

		port := "64040"
		serverAddress := "localhost:" + port
		serverBaseUrl := "http://localhost:" + port

		serverPemPath, err := testsGetFilePath("recipient.pem")
		So(err, ShouldBeNil)
		allowedKeysPath, err := testsGetFilePath("allowed_keys")
		So(err, ShouldBeNil)
		clientPemPath, err := testsGetFilePath("sender.pem")
		So(err, ShouldBeNil)

		addressFlag = &serverAddress
		serverPemFlag = &serverPemPath
		allowedPublicKeysFileFlag = &allowedKeysPath
		go a.Run(logger)

		time.Sleep(500 * time.Millisecond) //Give server time to start
		cl, err := setupClient(clientPemPath)
		So(err, ShouldBeNil)

		num := 300
		var wg sync.WaitGroup
		wg.Add(num)
		for i := 0; i < num; i++ {
			go doRequest(&wg, logger, i, cl, serverBaseUrl)
		}
		wg.Wait()

		for i, e := range logger.ErrorList {
			t.Errorf("ErrorList[%d]: %s", i, e)
		}
		So(logger.ErrorList, ShouldResemble, []string{})
	})*/
}
