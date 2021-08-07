

change-position:
	curl -XPOST http://localhost:7070/api/cmds/position -d '{"drone_id":"drone-1", "latitude": 81.231, "longitude": 43.1231, "altitude": 2301.1, "current_speed": 41.3, "heading_cardinal": 1}'

query-position:
	curl http://localhost:7070/api/cmds/query-position/drone-1


add-telemetry:
	curl -XPOST http://localhost:7070/api/cmds/telemetry -d '{"drone_id":"drone-1", "battery": 72, "uptime": 6941, "core_temp": 21 }'

query-telemetry:
	curl http://localhost:7070/api/cmds/query-telemetry/drone-1


add-alert:
	curl -XPOST http://localhost:7070/api/cmds/alertsignal -d '{"drone_id":"drone-1", "fault_code": 12, "description": "all the things are failing"}'

query-alert:
	curl http://localhost:7070/api/cmds/query-alertsignal/drone-1
