package model

type Model struct {
	Name   string
	Meshes []Mesh
}

type Mesh struct {
	Vertices       []float32
	Indices        []uint16
	VAO            uint32
	VBO            uint32
	IndexBuffer    uint32
	VertexShader   uint32
	FragmentShader uint32
}
