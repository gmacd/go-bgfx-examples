package main

import (
	"github.com/gmacd/go-bgfx"
	"github.com/gmacd/go-bgfx-examples/assets"
	"github.com/gmacd/go-bgfx-examples/example"
	"j4k.co/cgm"
	"j4k.co/cgm/mat4"
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

	var vd bgfx.VertexDecl
	vd.Begin()
	vd.Add(bgfx.AttribPosition, 3, bgfx.AttribTypeFloat, false, false)
	vd.Add(bgfx.AttribColor0, 4, bgfx.AttribTypeUint8, true, false)
	vd.End()
	vb := bgfx.CreateVertexBuffer(vertices, vd)
	defer bgfx.DestroyVertexBuffer(vb)
	ib := bgfx.CreateIndexBuffer(indices)
	defer bgfx.DestroyIndexBuffer(ib)
	prog := assets.LoadProgram("vs_cubes", "fs_cubes")
	defer bgfx.DestroyProgram(prog)

	for app.Continue() {
		t := app.Time
		dt := app.DeltaTime
		var (
			eye = [3]float32{0, 0, -35.0}
			at  = [3]float32{0, 0, 0}
			up  = [3]float32{1, 0, 0}
		)
		view := mat4.LookAtLH(eye, at, up)
		proj := mat4.PerspectiveLH(
			cgm.ToRadians(60),
			float32(app.Width)/float32(app.Height),
			0.1, 100,
		)
		bgfx.SetViewTransform(0, view, proj)
		bgfx.SetViewRect(0, 0, 0, app.Width, app.Height)
		bgfx.DebugTextClear()
		bgfx.DebugTextPrintf(0, 1, 0x4f, app.Title)
		bgfx.DebugTextPrintf(0, 2, 0x6f, "Description: Rendering simple static mesh.")
		bgfx.DebugTextPrintf(0, 3, 0x0f, "Frame: % 7.3f[ms]", dt*1000.0)
		bgfx.Submit(0)

		// Submit 11x11 cubes
		for y := 0; y < 11; y++ {
			for x := 0; x < 11; x++ {
				mtx := mat4.RotateXYZ(
					cgm.Radians(t)+cgm.Radians(x)*0.21,
					cgm.Radians(t)+cgm.Radians(y)*0.37,
					0,
				)
				mtx[12] = -15 + float32(x)*3
				mtx[13] = -15 + float32(y)*3
				mtx[14] = 0

				bgfx.SetTransform(mtx)
				bgfx.SetProgram(prog)
				bgfx.SetVertexBuffer(vb)
				bgfx.SetIndexBuffer(ib)
				bgfx.SetState(bgfx.StateDefault)
				bgfx.Submit(0)
			}
		}

		bgfx.Frame()
	}
}
