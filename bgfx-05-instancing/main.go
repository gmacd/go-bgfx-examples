package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/james4k/go-bgfx"
	"github.com/james4k/go-bgfx-examples/assets"
	"github.com/james4k/go-bgfx-examples/example"
)

type PosColorVertex struct {
	X, Y, Z float32
	ABGR    uint32
}

var vertices = []PosColorVertex{
	{-1.0, 1.0, 1.0, 0xff000000},
	{1.0, 1.0, 1.0, 0xff0000ff},
	{-1.0, -1.0, 1.0, 0xff00ff00},
	{1.0, -1.0, 1.0, 0xff00ffff},
	{-1.0, 1.0, -1.0, 0xffff0000},
	{1.0, 1.0, -1.0, 0xffff00ff},
	{-1.0, -1.0, -1.0, 0xffffff00},
	{1.0, -1.0, -1.0, 0xffffffff},
}

var indices = []uint16{
	0, 1, 2, // 0
	1, 3, 2,
	4, 6, 5, // 2
	5, 6, 7,
	0, 2, 4, // 4
	4, 2, 6,
	1, 5, 3, // 6
	5, 7, 3,
	0, 4, 1, // 8
	4, 5, 1,
	2, 3, 6, // 10
	6, 3, 7,
}

func main() {
	app := example.Open()
	defer app.Close()
	bgfx.Init()
	defer bgfx.Shutdown()

	bgfx.Reset(app.Width, app.Height, bgfx.ResetVSync)
	bgfx.SetDebug(bgfx.DebugText)
	bgfx.SetViewClear(
		0,
		bgfx.ClearColor|bgfx.ClearDepth,
		0x303030ff,
		1.0,
		0,
	)

	caps := bgfx.Caps()

	var vd bgfx.VertexDecl
	vd.Begin()
	vd.Add(bgfx.AttribPosition, 3, bgfx.AttribTypeFloat, false, false)
	vd.Add(bgfx.AttribColor0, 4, bgfx.AttribTypeUint8, true, false)
	vd.End()
	vb := bgfx.CreateVertexBuffer(vertices, vd)
	defer bgfx.DestroyVertexBuffer(vb)
	ib := bgfx.CreateIndexBuffer(indices)
	defer bgfx.DestroyIndexBuffer(ib)
	prog, err := loadProgram("vs_instancing", "fs_instancing")
	if err != nil {
		log.Fatalln(err)
	}
	defer bgfx.DestroyProgram(prog)

	for app.Continue() {
		var (
			eye = mgl32.Vec3{0, 0, -35.0}
			at  = mgl32.Vec3{0, 0, 0}
			up  = mgl32.Vec3{1, 0, 0}
		)
		view := [16]float32(mgl32.LookAtV(eye, at, up))
		proj := [16]float32(mgl32.Perspective(
			mgl32.DegToRad(60),
			float32(app.Width)/float32(app.Height),
			0.1, 100.0,
		))
		bgfx.SetViewTransform(0, view, proj)
		bgfx.SetViewRect(0, 0, 0, app.Width, app.Height)
		bgfx.DebugTextClear()
		bgfx.DebugTextPrintf(0, 1, 0x4f, app.Title)
		bgfx.DebugTextPrintf(0, 2, 0x6f, "Description: Geometry instancing.")
		bgfx.DebugTextPrintf(0, 3, 0x0f, "Frame: % 7.3f[ms]", app.DeltaTime*1000.0)
		bgfx.Submit(0)

		if caps.Supported&bgfx.CapsInstancing == 0 {
			color := uint8(0x01)
			if uint32(app.Time*2)&1 != 0 {
				color = 0x1f
			}
			bgfx.DebugTextPrintf(0, 5, color, " Instancing is not supported by GPU. ")
			bgfx.Frame()
			continue
		}

		const stride = 80
		idb := bgfx.AllocInstanceDataBuffer(11*11, stride)
		// Submit 11x11 cubes
		time64 := float64(app.Time)
		for y := 0; y < 11; y++ {
			for x := 0; x < 11; x++ {
				mtx := mgl32.HomogRotate3DX(app.Time + float32(x)*0.21)
				mtx = mtx.Mul4(mgl32.HomogRotate3DY(app.Time + float32(y)*0.37))
				mtx[12] = -15 + float32(x)*3
				mtx[13] = -15 + float32(y)*3
				mtx[14] = 0
				color := [4]float32{
					float32(math.Sin(time64+float64(x)/11.0)*0.5 + 0.5),
					float32(math.Cos(time64+float64(y)/11.0)*0.5 + 0.5),
					float32(math.Sin(time64*3.0)*0.5 + 0.5),
					1.0,
				}
				binary.Write(&idb, binary.LittleEndian, mtx)
				binary.Write(&idb, binary.LittleEndian, color)
			}
		}

		bgfx.SetProgram(prog)
		bgfx.SetVertexBuffer(vb)
		bgfx.SetIndexBuffer(ib)
		bgfx.SetInstanceDataBuffer(idb)
		bgfx.SetState(bgfx.StateDefault)
		bgfx.Submit(0)

		bgfx.Frame()
	}
}

func loadProgram(vsh, fsh string) (bgfx.Program, error) {
	v, err := loadShader(vsh)
	if err != nil {
		return bgfx.Program{}, err
	}
	f, err := loadShader(fsh)
	if err != nil {
		return bgfx.Program{}, err
	}
	return bgfx.CreateProgram(v, f, true), nil
}

func loadShader(name string) (bgfx.Shader, error) {
	f, err := assets.Open(filepath.Join("shaders/glsl", name+".bin"))
	if err != nil {
		return bgfx.Shader{}, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return bgfx.Shader{}, err
	}
	return bgfx.CreateShader(data), nil
}
