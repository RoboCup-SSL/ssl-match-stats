package matchstats

import (
	"github.com/RoboCup-SSL/ssl-match-stats/internal/vision"
)

func (g *Generator) OnNewVisionMessage(visionMsg *vision.SSL_WrapperPacket) {
	if visionMsg.Detection == nil {
		return
	}

	timestamp := *visionMsg.Detection.TCapture
	g.updateRobotCount(TeamColor_TEAM_YELLOW, visionMsg.Detection.RobotsYellow, timestamp)
	g.updateRobotCount(TeamColor_TEAM_BLUE, visionMsg.Detection.RobotsBlue, timestamp)
}

func (g *Generator) updateRobotCount(teamColor TeamColor, detections []*vision.SSL_DetectionRobot, timestamp float64) {
	firstDetections := g.robotFirstDetection[teamColor]
	lastDetections := g.robotLastDetection[teamColor]
	for _, d := range detections {
		if d.RobotId == nil {
			// optional in protobuf, so this can happen
			continue
		}
		if _, ok := firstDetections[*d.RobotId]; !ok {
			firstDetections[*d.RobotId] = timestamp
		}
		lastDetections[*d.RobotId] = timestamp
	}
	for id, lastSeen := range lastDetections {
		// remove if older than one second
		if lastSeen+1 < timestamp {
			delete(firstDetections, id)
			delete(lastDetections, id)
		}
	}
}
