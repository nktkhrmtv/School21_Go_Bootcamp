package main

import (
	"github.com/fogleman/gg"
	"image/color"
)

func main() {
	const width = 300
	const height = 300

	dc := gg.NewContext(width, height)

	grad := gg.NewLinearGradient(0, 0, width, height)
	grad.AddColorStop(0, color.RGBA{255, 105, 180, 255}) 
	grad.AddColorStop(1, color.RGBA{0, 191, 255, 255})  
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()

	dc.DrawCircle(width/2, height/2, 100)
	dc.SetColor(color.RGBA{200, 255, 200, 255}) 
	dc.Fill()

	dc.SetColor(color.RGBA{0, 0, 0, 255}) 
	dc.DrawStringAnchored("Meteoriw Blog", width/2, height/2, 0.5, 0.5)

	dc.SavePNG("amazing_logo.png")
}