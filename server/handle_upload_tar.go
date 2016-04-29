package main

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/tar_io"
)

func (h *handler) handleUploadTarFunc(c *echo.Context) error {
	req := c.Request()
	basePath := req.Header.Get(shared.BASE_PATH_HTTP_HEADER_NAME)
	if strings.TrimSpace(basePath) == "" {
		return fmt.Errorf("Request does not contain header '%s'", shared.BASE_PATH_HTTP_HEADER_NAME)
	}

	h.logger.Infof("Now starting to receive dir '%s'", basePath)
	tarReceiver := tar_io.Factories.TarReceiver.Dir(basePath)
	return tar_io.SaveTarToReceiver(req.Body, tarReceiver)
}