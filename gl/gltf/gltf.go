package gltf

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"go-game-engine/log"
	"go-game-engine/math"
	"go-game-engine/model"

	"github.com/qmuntal/gltf"
)

const (
	GLFloat = 5126
	GLUInt  = 5123
)

const (
	GLScalar = iota
	GLVec2
	GLVec3
	GLVec4
	GlMat2
	GLMat3
	GLMat4
)

func Load(path string) ([]model.Model, error) {
	log.Msgf(log.INFO, "reading gltf file from path: %v", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file from %q: %v", path, err)
	}

	var doc gltf.Document
	reader := bytes.NewReader(b)
	gltf.NewDecoder(reader).Decode(&doc)
	models, _ := process(&doc)
	return models, nil
}

func process(doc *gltf.Document) ([]model.Model, error) {
	log.Msgf(log.INFO, "processing gltf document: %v", doc)

	models := make([]model.Model, len(doc.Meshes))
	for i, mesh := range doc.Meshes {
		posView := mesh.Primitives[0].Attributes["POSITION"]
		// normView := mesh.Primitives[0].Attributes["NORMAL"]
		// uvView := mesh.Primitives[0].Attributes["TEXCOORD_0"]
		indicesView := *mesh.Primitives[0].Indices

		// log.Msgf(0, "posView %v, normView: %v, uvView: %v", posView, normView, uvView)
		waitChan := make(chan struct{}, 3)

		positions := make([]float32, doc.Accessors[posView].Count)
		go getArrayAsync(doc, posView, &positions, waitChan)

		// normals := make([]float32, doc.Accessors[normView].Count)
		// go getArrayAsync(doc, normView, &normals, waitChan)

		// uvs := make([]float32, doc.Accessors[uvView].Count)
		// go getArrayAsync(doc, uvView, &uvs, waitChan)

		indices := make([]uint16, doc.Accessors[indicesView].Count)
		go getArrayAsync(doc, indicesView, &indices, waitChan)

		<-waitChan
		// <-waitChan
		// <-waitChan
		<-waitChan

		vCount := len(positions)
		count := len(positions) // + len(normals)
		vertices := make([]float32, count)

		log.Msgf(0, "found %v positions: %v", len(positions), positions)
		// log.Msgf(0, "found %v normals: %v", len(normals), normals)
		log.Msgf(0, "found %v indices: %v", len(indices), indices)
		// log.Msgf(0, "found %v uvs: %v", len(uvs), uvs)

		for x := 0; x < 3; x++ {
			vertices[x] = positions[x]
			// vertices[3+x] = normals[x]
		}
		for i := 3; i < vCount-1; i += 3 {
			for x := 0; x < 3; x++ {
				// vertices[(i*2)+x] = positions[i+x]
				vertices[(i)+x] = positions[i+x]
				// vertices[3+(i*2)+x] = normals[i+x]
			}
		}

		mesh := model.Mesh{positions, indices, 0, 0, 0, 0, 0}
		models[i] = model.Model{"", []model.Mesh{mesh}}
		log.Msgf(0, "found %v positions: %v", len(vertices), vertices)
	}

	return models, nil
}

func getArrayAsync(doc *gltf.Document, view uint32, result interface{}, done chan<- struct{}) {
	byteLength, offset, glType, cType, bufferData := getBufferValues(doc, view)
	var err error
	if cType == GLFloat {
		res, ok := result.(*[]float32)
		if !ok {
			log.Msgf(0, "failed to convert array into float32s: %v", result)
			done <- struct{}{}
			return
		}

		switch glType {
		case GLVec2:
			err = bufferToVec2Array(byteLength, offset, *bufferData, res)
		case GLVec3:
			err = bufferToVec3Array(byteLength, offset, *bufferData, res)
		}
	} else if cType == GLUInt {
		res, ok := result.(*[]uint16)
		if !ok {
			log.Msgf(0, "failed to convert array into uint16s: %v", result)
			done <- struct{}{}
			return
		}

		switch glType {
		case GLScalar:
			err = bufferToScalarArray(byteLength, offset, *bufferData, res)
		}
	}
	if err != nil {
		log.Msgf(log.INFO, "could not get array async: %v", err)
	}

	done <- struct{}{}
}

// Returns Byte length, Offset, buffer type, conversionType and the Buffer to use
func getBufferValues(doc *gltf.Document, view uint32) (uint32, uint32, int, int, *[]byte) {
	bufferViewIndex := *doc.Accessors[view].BufferView
	bufferView := doc.BufferViews[bufferViewIndex]
	byteLength := bufferView.ByteLength
	offset := bufferView.ByteOffset
	buffer := doc.Buffers[bufferView.Buffer]
	glType := int(doc.Accessors[view].Type)
	convType := gltfComponentTypeToGLDataType(int(doc.Accessors[view].ComponentType))

	return byteLength, offset, glType, convType, &buffer.Data
}

func bufferToVec3Array(byteLength, offset uint32, fullBuffer []byte, result *[]float32) error {
	count := uint32(len(*result))
	byteSize := byteLength / count

	for i := uint32(0); i < count-1; i += 3 {
		index := offset + (i * byteSize)
		v3, err := Vector3FromBytes(fullBuffer[index : index+byteSize])
		if err != nil {
			return fmt.Errorf("cannot convert buffer to vec3: %v", err)
		}

		(*result)[i] = v3[0]
		(*result)[i+1] = v3[1]
		(*result)[i+2] = v3[2]
	}
	return nil
}

func bufferToVec2Array(byteLength, offset uint32, fullBuffer []byte, result *[]float32) error {
	count := uint32(len(*result))
	byteSize := byteLength / count

	for i := uint32(0); i < count-1; i += 2 {
		index := offset + (i * byteSize)
		v2, err := Vector2FromBytes(fullBuffer[index : index+byteSize])
		if err != nil {
			return fmt.Errorf("cannot convert buffer to vec2: %v", err)
		}

		(*result)[i] = v2[0]
		(*result)[i+1] = v2[1]
	}
	return nil
}

func bufferToScalarArray(byteLength, offset uint32, fullBuffer []byte, result *[]uint16) error {
	count := uint32(len(*result))
	byteSize := byteLength / count

	for i := uint32(0); i < count; i++ {
		index := offset + (i * byteSize)
		v, err := math.BytesToUInt16(fullBuffer[index : index+byteSize])
		if err != nil {
			return fmt.Errorf("cannot convert buffer to scalar: %v", err)
		}

		(*result)[i] = v
	}
	return nil
}

func Vector3FromBytes(bytes []byte) ([]float32, error) {
	s := len(bytes) / 3

	f1, err := math.BytesToFloat32(bytes[:s])
	if err != nil {
		return nil, fmt.Errorf("cannot turn []byte into 3 float32s, failed on f1: %v", err)
	}
	f2, err := math.BytesToFloat32(bytes[s : s+s])
	if err != nil {
		return nil, fmt.Errorf("cannot turn []byte into 3 float32s, failed on f2: %v", err)
	}
	f3, err := math.BytesToFloat32(bytes[s+s:])
	if err != nil {
		return nil, fmt.Errorf("cannot turn []byte into 3 float32s, failed on f3: %v", err)
	}

	return []float32{f1, f2, f3}, nil
}

func Vector2FromBytes(bytes []byte) ([]float32, error) {
	s := len(bytes) / 2

	f1, err := math.BytesToFloat32(bytes[:s])
	if err != nil {
		return nil, fmt.Errorf("cannot turn []byte into 2 float32s, failed on f1: %v", err)
	}
	f2, err := math.BytesToFloat32(bytes[s:])
	if err != nil {
		return nil, fmt.Errorf("cannot turn []byte into 2 float32s, failed on f2: %v", err)
	}

	return []float32{f1, f2}, nil
}

func gltfComponentTypeToGLDataType(gltf int) int {
	switch gltf {
	case 0:
		return GLFloat
	case 4:
		return GLUInt
	}
	return -1
}
