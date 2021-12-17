package gl

import (
	"fmt"
	"runtime"
	"strings"

	"go-game-engine/log"
	"go-game-engine/model"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	triangle = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}

	square = []float32{
		0.5, 0.5, 0,
		0.5, -0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
	}

	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(0, 0.5, 1, 1);
    }
` + "\x00"
)

func CreateWindow(width, height int) (uint32, *glfw.Window, uint32, error) {
	runtime.LockOSThread()

	window, err := initGlfw(width, height)
	if err != nil {
		return 0, nil, 0, err
	}

	program, err := initOpenGL()
	if err != nil {
		return 0, nil, 0, err
	}
	vao := makeVao(square)

	return vao, window, program, nil
}

func TerminateGLFW() {
	glfw.Terminate()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func GenMeshVAO(mesh *model.Mesh) {
	gl.GenBuffers(1, &mesh.VAO)
	gl.BindVertexArray(mesh.VAO)

	gl.GenBuffers(1, &mesh.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Vertices)*4, gl.Ptr(mesh.Vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.GenBuffers(1, &mesh.IndexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.IndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*2, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)
}

func Draw(mdl model.Model, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// Draw the triangles !
	gl.BindVertexArray(mdl.Meshes[0].VAO)
	gl.DrawElements(gl.TRIANGLE_FAN, int32(len(mdl.Meshes[0].Indices)), gl.UNSIGNED_SHORT, gl.PtrOffset(0))

	// gl.BindVertexArray(mdl.Meshes[0].VAO)
	// gl.DrawArrays(gl.TRIANGLE_FAN, 0, 24)

	glfw.PollEvents()
	window.SwapBuffers()
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() (uint32, error) {
	err := gl.Init()
	if err != nil {
		return 0, fmt.Errorf("could not initialise gl: %v", err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Msgf(0, "OpenGL version: %v", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile shader: %v", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("could not compile shader: %v", err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog, nil
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(width, height int) (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("Could not init glfw: %v", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not create window: %v", err)
	}
	window.MakeContextCurrent()

	return window, nil
}
