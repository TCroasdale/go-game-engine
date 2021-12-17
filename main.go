package main

import (
	"fmt"
	"go-game-engine/gl"
	"go-game-engine/gl/gltf"
	"go-game-engine/log"
)

const (
	width  = 500
	height = 500
)

func main() {
	log.Start(500)
	defer log.Close()

	models, err := gltf.Load("./data/models/cube1.glb")
	if err != nil {
		log.Msgf(log.ERROR, "Error: %v", err)
		return
	}

	_, window, program, err := gl.CreateWindow(width, height)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gl.TerminateGLFW()

	gl.GenMeshVAO(&models[0].Meshes[0])

	for !window.ShouldClose() {
		gl.Draw(models[0], window, program)
	}
}
