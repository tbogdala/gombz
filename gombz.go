// Copyright 2015, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package gombz

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"

	mgl "github.com/go-gl/mathgl/mgl32"
	"gopkg.in/mgo.v2/bson"
)

const (
	// MaxUVChannelCount is the maximum number of UV channels supported.
	MaxUVChannelCount = 8
)

// MeshFace is an alias for a vector of 3 unsigned ints.
// Note: currently mathgl doesn't have support for integer vectors
type MeshFace [3]uint32

// AnimationQuatKey is a animation channel key that contains
// a quaternion based animation
type AnimationQuatKey struct {
	// Time specifies the time in the animation for this current Quat value (e.g. rotation).
	Time float32

	// Key is the Quat value (e.g. rotation) at the given time in animation.
	Key mgl.Quat
}

// AnimationVec3Key is a animation channel key that contains
// a vector based animation
type AnimationVec3Key struct {
	// Time specifies the time in the animation for this current vector value (e.g. position).
	Time float32

	// Key is the vector value (e.g. position) at the given time in animation.
	Key mgl.Vec3
}

// AnimationChannel is an object that contains data required to transform
// a bone in an animation.
type AnimationChannel struct {
	// Name is the name of the animation channel
	Name string

	// PositionKeys is a slice of vector keys describing bone position at a given time in an animation.
	PositionKeys []AnimationVec3Key

	// ScaleKeys is a slice of vector keys describing bone scale at a given time in an animation.
	ScaleKeys []AnimationVec3Key

	// RotationKeys is a slice of Quat keys describing bone rotation at a given time in an animation.
	RotationKeys []AnimationQuatKey
}

// Animation represents data required to transform bones in an animation.
type Animation struct {
	// Name is the name of the animation
	Name string

	// Duration is the length of the animation in ticks
	Duration float32

	// TicksPerSecond is the number of ticks per second for the animation as designed.
	TicksPerSecond float32

	// Transform is the transformation matrix.
	Transform mgl.Mat4

	// Channels is a slice of AnimationChannel objects that deform the bones
	Channels []AnimationChannel
}

// Bone is a struct that contains transform and skeletal information
// for skeletal animation.
type Bone struct {
	// Name is the name of the bone
	Name string

	// Id is a number assigned to the bone that should be >= 0
	Id int32

	// Parent is the Id of the parent bone in the skeleton. It will be -1 for the root bone.
	Parent int32

	// Offset is the matrix that transforms from mesh space to bone space in bind pose
	Offset mgl.Mat4

	// Transform is the matrix that is the transformation relative to the node's parent.
	Transform mgl.Mat4
}

// Mesh contains the data for a given mesh such as vertices and faces.
type Mesh struct {
	// FaceCount is the number of faces in the mesh
	FaceCount uint32

	// BoneCount is the number of bones in the mesh
	BoneCount uint32

	// VertexCount is the number of vertices in the mesh
	VertexCount uint32

	// Vertices is a slice of vertices that are represented with Vec3 float vectors.
	Vertices []mgl.Vec3

	// Normals is a slice of Vec3 floats representing a normal for each vertex.
	Normals []mgl.Vec3

	// Tangents is a slice of Vec3 floats representing a tangent for each vertex.
	Tangents []mgl.Vec3

	// Faces are triplets of unsigned ints that specify the vertices used in a face
	Faces []MeshFace

	// UVChannels is an array (of size MaxUVChannelCount) of slices that are Vec2
	// float arrays representing the UV coordinate for a vertex.
	UVChannels [MaxUVChannelCount][]mgl.Vec2

	// Bones is slice of Bone objects that form the skeleton of the mesh
	Bones []Bone

	// VertexWeightIds is a slice that has a Vec4 for each vertex in Vertices
	// which can contain up to four bone ids that will modify the vertex.
	// Note: stored as a float since that's how it will be passed to shaders
	VertexWeightIds []mgl.Vec4

	// VertexWeights is a slice that has a Vec4 for each vertex in Vertices
	// which can contain up to four bone weights used to determine how much
	// the bone specified by VertexWeightIds affects the vertex position.
	VertexWeights []mgl.Vec4

	// Animations is a slice of Animation objects that represent all animations that
	// can deform the mesh's Bones.
	Animations []Animation
}

// Encode takes a given mesh and encodes it to binary with bson
// and then compresses it with zlib and returns the result -- or
// returns a non-nil err on fail.
func (mesh *Mesh) Encode() (out []byte, err error) {
	// encode
	bs, err := bson.Marshal(mesh)
	if err != nil {
		return nil, err
	}

	// compress
	gzBuffer := new(bytes.Buffer)
	gz, err := zlib.NewWriterLevel(gzBuffer, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err = gz.Write(bs); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}

	return gzBuffer.Bytes(), nil
}

// DecodeMesh takes a byte stream and decompresses it with zlib and
// then decodes it with bson and returns the result -- or returns
// a non-nil err on fail.
func DecodeMesh(bs []byte) (outMesh *Mesh, err error) {
	// load up the buffer
	gzBuffer := bytes.NewBuffer(bs)

	// decompress
	gzReader, err := zlib.NewReader(gzBuffer)
	if err != nil {
		return nil, err
	}
	decompBytes, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	// decode
	outMesh = new(Mesh)
	err = bson.Unmarshal(decompBytes, outMesh)
	if err != nil {
		return nil, err
	}

	return outMesh, nil
}
