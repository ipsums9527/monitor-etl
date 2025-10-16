package server

import (
	"github.com/bitly/go-simplejson"
	"github.com/labstack/echo/v4"
)

/*
netdata apis

"/api/v1/data?chart=system.cpu"
"/api/v1/data?chart=system.ram"
"/api/v1/data?chart=system.uptime"
"/api/v1/data?context=system.hw.sensor.temperature.input"
"/api/v1/data?context=system.net"

*/

func (s *Server) Data(c echo.Context) error {
	chart := c.QueryParam("chart")
	ctx := c.QueryParam("context")

	js := simplejson.New()

	if chart != "" {
		switch chart {
		case "system.ram":
			js.Set("dimension_ids", []string{"free", "used", "cached", "buffers"})
			js.Set("latest_values", []float64{s.data.memFree.Load().(float64), s.data.memUsed.Load().(float64), 0, 0})
			return c.JSON(200, js)

		case "system.uptime":
			js.Set("view_latest_values", []float64{s.data.uptime.Load().(float64)})
			return c.JSON(200, js)

		case "system.cpu":
			js.Set("dimension_ids", []string{"user", "system"})
			js.Set("latest_values", []float64{s.data.cpuUser.Load().(float64), s.data.cpuSys.Load().(float64)})
			return c.JSON(200, js)

		default:
			return c.HTML(404, "No metrics where matched to query.")
		}
	}

	if ctx != "" {
		switch ctx {
		case "system.net":
			// download upload
			js.Set("view_latest_values", []float64{s.data.netDown.Load().(float64), s.data.netUp.Load().(float64)})
			return c.JSON(200, js)

		case "system.hw.sensor.temperature.input":
			js.Set("chart_ids", []string{"sensors.fake"})
			js.Set("view_latest_values", []float64{s.data.temp.Load().(float64)})
			return c.JSON(200, js)

		default:
			return c.HTML(404, "No metrics where matched to query.")
		}
	}
	return nil
}
