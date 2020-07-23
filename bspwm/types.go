package bspwm

type Client struct {
	InstanceName string
	ClassName string
}

type Node struct {
	Id          int
	FirstChild  *Node
	SecondChild *Node
	Client Client
}

type DesktopInfo struct {
	Name string
	Id   int
	Root *Node
	FocusedNodeId int
}

type MonitorInfo struct {
	Name             string
	Id               int
	FocusedDesktopId int
	Desktops         []DesktopInfo
}