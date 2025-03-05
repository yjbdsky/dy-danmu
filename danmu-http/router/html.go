package router

import "embed"

//go:embed dist/index.html
var Html []byte

//go:embed dist/* dist/assets/*
var Static embed.FS
