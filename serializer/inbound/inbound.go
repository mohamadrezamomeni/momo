package inbound

import "time"

type InboundSerializer struct {
	UserID       string    `json:"user_id"`
	Protocol     string    `json:"protocol"`
	Tag          string    `json:"tag"`
	Domain       string    `json:"domain"`
	Port         string    `json:"port"`
	IsNotified   bool      `json:"is_notified"`
	IsBlock      bool      `json:"is_block"`
	IsAssigned   bool      `json:"is_assigned"`
	Start        time.Time `json:"start_time"`
	End          time.Time `json:"end_time"`
	ChargeCount  uint32    `json:"charhe_count"`
	TrafficUsage uint32    `json:"traffic_usage"`
	TrafficLimit uint32    `json:"traffic_limit"`
}
