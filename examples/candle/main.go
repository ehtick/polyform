package main

import (
	"image"
	"image/color"
	"math"

	"github.com/EliCDavis/polyform/formats/obj"
	"github.com/EliCDavis/polyform/modeling"
	"github.com/EliCDavis/polyform/modeling/extrude"
	"github.com/EliCDavis/vector/vector2"
	"github.com/EliCDavis/vector/vector3"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

func bevel(startWidth, startHeight, endWidth, endHeight, uvStart, uvEnd, uvThickness float64, resolution int) []extrude.ExtrusionPoint {
	halfPi := math.Pi / 2.

	points := make([]extrude.ExtrusionPoint, 0)

	heightDelta := endHeight - startHeight
	widthDelta := startWidth - endWidth
	uvDelta := uvEnd - uvStart

	for i := 0; i < resolution; i++ {
		percent := float64(i+1) / float64(resolution)

		sinResult := math.Sin(percent * halfPi)
		cosResult := math.Cos(percent * halfPi)

		height := sinResult * float64(heightDelta)

		points = append(points, extrude.ExtrusionPoint{
			Point:     vector3.Up[float64]().Scale(height + startHeight),
			Thickness: (cosResult * widthDelta) + endWidth,
			UV: &extrude.ExtrusionPointUV{
				Point:     vector2.New(0.5, uvStart+(uvDelta*sinResult)),
				Thickness: uvThickness,
			},
		})
	}

	return points
}

func candleBody(height, width, rimWidth, percentUsed, wickWidth, wickHeight float64) modeling.Mesh {
	points := bevel(0, -.1, width, 0, 0, 0.1, 1, 10)

	startDome := bevel(width, height, width-(rimWidth/2), height+.1, 0.7, 0.75, 1, 10)
	points = append(points, startDome...)

	endDome := bevel(width-(rimWidth/2), height+.1, width-rimWidth, height, 0.75, 0.8, 1, 10)
	points = append(points, endDome...)

	heightToWax := percentUsed * height

	points = append(
		points,
		extrude.ExtrusionPoint{
			Point:     vector3.Up[float64]().Scale(heightToWax),
			Thickness: width - rimWidth,
			UV: &extrude.ExtrusionPointUV{
				Thickness: 1,
				Point:     vector2.New(0.5, 0.9),
			},
		},
		extrude.ExtrusionPoint{
			Point:     vector3.Up[float64]().Scale(heightToWax),
			Thickness: wickWidth,
			UV: &extrude.ExtrusionPointUV{
				Thickness: 1,
				Point:     vector2.New(0.5, 0.95),
			},
		},
		extrude.ExtrusionPoint{
			Point:     vector3.Up[float64]().Scale(heightToWax + wickHeight),
			Thickness: wickWidth,
			UV: &extrude.ExtrusionPointUV{
				Thickness: 1,
				Point:     vector2.New(0.5, 0.975),
			},
		},
		extrude.ExtrusionPoint{
			Point:     vector3.Up[float64]().Scale(heightToWax + wickHeight),
			Thickness: 0,
			UV: &extrude.ExtrusionPointUV{
				Thickness: 1,
				Point:     vector2.New(0.5, 1.0),
			},
		},
	)

	return extrude.Polygon(
		30,
		points,
	)
}

func candleTexture(containerColor, waxColor color.Color, logoPath, outPath string) {
	S := 1024
	LogoSize := S / 4

	im, err := gg.LoadJPG(logoPath)
	if err != nil {
		panic(err)
	}
	logo := image.NewRGBA(image.Rect(0, 0, LogoSize, int(float64(LogoSize)*1.5)))
	draw.ApproxBiLinear.Scale(logo, logo.Rect, im, im.Bounds(), draw.Over, nil)
	// draw.ApproxBiLinear
	dc := gg.NewContext(S, S)

	waxSize := int(math.Round(float64(S) * 0.05))
	wickSize := int(math.Round(float64(S) * 0.05))

	candleColorHeight := S - waxSize - wickSize
	dc.SetColor(containerColor)
	dc.DrawRectangle(0, float64(waxSize+wickSize), float64(S), float64(candleColorHeight))
	dc.Fill()

	dc.SetColor(waxColor)
	dc.DrawRectangle(0, float64(wickSize), float64(S), float64(waxSize))
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, float64(S), float64(wickSize))
	dc.Fill()

	dc.SetColor(color.Black)
	dc.DrawRectangle(0, 0, float64(S), float64(wickSize)/1.5)
	dc.Fill()

	dc.DrawImageAnchored(logo, S/2, int(float64(S)/3.5), 0.5, -0.3)
	err = dc.SavePNG(outPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	logoPath := "candlelogo.jpg"
	texturePath := "candle-diffuse.png"
	candleTexture(
		color.RGBA{R: 250, G: 244, B: 230, A: 255},
		color.RGBA{R: 255, G: 242, B: 161, A: 255},
		logoPath,
		texturePath,
	)

	err := obj.Save("candle.obj", obj.Scene{
		Objects: []obj.Object{
			{
				Entries: []obj.Entry{
					{
						Mesh: candleBody(1, 0.5, 0.1, 0.9, 0.0125, 0.1),
						Material: &obj.Material{
							Name:            "Candle",
							ColorTextureURI: &texturePath,
							DiffuseColor:    color.White,
							Transparency:    0,
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
