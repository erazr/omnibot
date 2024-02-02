package weather_widget

import (
	"image"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/erazr/omnibot/internal/weather"
	"github.com/fogleman/gg"
)

// Draws a widget based on weather data
func DrawWeatherWidget(data *weather.WeatherResponse) (*os.File, error) {
	var bgImage, _ = gg.LoadPNG("assets/images/bg.png")
	dc := gg.NewContext(bgImage.Bounds().Dx(), bgImage.Bounds().Dy())

	err := PrepareImage(dc, bgImage, data.Current.Condition.Icon)
	if err != nil {
		return nil, err
	}

	dc.LoadFontFace("assets/fonts/Inconsolata-Regular.ttf", 36)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(strconv.FormatFloat(data.Current.Temp_c, 'f', 1, 64)+" °C", 467, 164, 0.5, 0.5)

	wind := "Wind: " + strconv.FormatFloat(data.Current.Wind_kph, 'f', 1, 64)
	precipIn := "Precip in: " + strconv.FormatFloat(data.Current.Precip_in, 'f', 1, 64)
	pressureIn := "Pressure: " + strconv.FormatFloat(data.Current.Pressure_in, 'f', 1, 64)

	dc.LoadFontFace("assets/fonts/Inconsolata-Regular.ttf", 20)

	dc.DrawStringWrapped(data.Current.Condition.Text, 170, 120, 0, 0.5, 170, 2, gg.AlignLeft)
	dc.DrawStringAnchored(wind, 467, 64, 0.5, 0.5)
	dc.DrawStringAnchored(precipIn, 467, 94, 0.5, 0.5)
	dc.DrawStringAnchored(pressureIn, 467, 124, 0.5, 0.5)

	for i, v := range data.Forecast.Days[1:] {
		tmp := strconv.FormatFloat(v.Day.Avgtemp_c, 'f', 1, 64) + " °C"

		offset := bgImage.Bounds().Dx() / len(data.Forecast.Days)
		x := offset * (i + 1)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			icon, _ := GetIcon(v.Day.Condition.Icon)
			dd, _ := time.Parse(time.DateOnly, v.Date)
			dc.DrawStringAnchored(dd.Weekday().String(), float64(x), 239, 0.5, 0.5)
			dc.DrawStringAnchored(tmp, float64(x), 326, 0.5, 0.5)
			dc.DrawImageAnchored(icon, x, 282, 0.5, 0.5)
		}()

		wg.Wait()

		if err != nil {
			return nil, err
		}
	}

	widgetPath := "assets/images/widget.png"

	dc.SavePNG(widgetPath)

	reader, err := os.Open(widgetPath)

	return reader, nil
}

func PrepareImage(dc *gg.Context, image image.Image, iconUrl string) error {
	icon, err := GetIcon(iconUrl)

	if err != nil {
		return err
	}

	scaled := imaging.Resize(icon, 86, 0, imaging.CatmullRom)

	dc.DrawImage(image, 0, 0)
	dc.DrawImageAnchored(scaled, 117, 119, 0.5, 0.5)

	return nil
}

func GetIcon(iconUrl string) (image.Image, error) {
	res, err := http.Get("https:" + iconUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	file, _ := os.Create("assets/icons/icon.png")
	defer file.Close()

	io.Copy(file, res.Body)

	icon, _ := gg.LoadImage(file.Name())

	return icon, nil
}
