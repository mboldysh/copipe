package main

type LaunchConfig struct {
	Steps []Step
}

type Step struct {
	Image string
	Name string
	Script []string
}