package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/dialog"
    "time"
    "image"
    "image/draw"
    "math/rand"
    "image/png"
    "os"
    
    "juego2/models"
)

func load(filePath string) image.Image {
    imgFile, err := os.Open(filePath)
    defer imgFile.Close()
    if err != nil {
        fmt.Println("Cannot read file:", err)
    }

    imgData, err := png.Decode(imgFile)
    if err != nil {
        fmt.Println("Cannot decode file:", err)
    }
    return imgData.(image.Image)
}

func resetPlayerPosition(player *models.Player) {
    player.X = 100
    player.Y = 200
}

func main() {
    myApp := app.New()
    w := myApp.NewWindow("Game")

    obstacleImage := load("img/ramos.png")
    player := models.Player{}

    background := load("img/background3.png")
    playerSprites := load("img/messi.png")
    pointsImage := load("img/pelota.png")

    points := &models.Points{X: 400, Y: 300, Width: 40, Height: 72, Collected: false}

    now := time.Now().UnixMilli()
    game := &models.Game{
        CanvasWidth:  800,
        CanvasHeight: 500,
        FPS:          10,
        Then:         now,
        Margin:       4,
    }

    fpsInterval := int64(1000 / game.FPS)

    obstacles := []models.Obstacle{
        models.Obstacle{X: 300, Y: 100, Width: 40, Height: 72, FrameX: 0, FrameY: 0, CyclesX: 4, UpY: 3, DownY: 0, LeftY: 1, RightY: 2, Speed: 9, XMov: 0, YMov: 0},
        models.Obstacle{X: 500, Y: 250, Width: 40, Height: 72, FrameX: 0, FrameY: 0, CyclesX: 4, UpY: 3, DownY: 0, LeftY: 1, RightY: 2, Speed: 9, XMov: 0, YMov: 0},
    }

    player = models.Player{X: 100, Y: 200, Width: 40, Height: 72, FrameX: 0, FrameY: 0, CyclesX: 4, UpY: 3, DownY: 0, LeftY: 1, RightY: 2, Speed: 9, XMov: 0, YMov: 0}

    img := canvas.NewImageFromImage(background)
    img.FillMode = canvas.ImageFillOriginal

    sprite := image.NewRGBA(background.Bounds())

    playerImg := canvas.NewRasterFromImage(sprite)
    spriteSize := image.Pt(player.Width, player.Height)

    puntos := 0
    puntosText := widget.NewLabel(fmt.Sprintf("Puntos: %d", puntos))
    puntosText.Move(fyne.NewPos(10, 10)) // Posición en la esquina superior izquierda
    puntosText.TextStyle = fyne.TextStyle{Bold: true}
    c := container.New(layout.NewMaxLayout(), img, playerImg, puntosText)
    w.SetContent(c)

    gameActions := make(chan fyne.KeyEvent)
    updateScreen := make(chan struct{})

    go func() {
        for {
            select {
            case k := <-gameActions:
                // Manejar las acciones del jugador
                // Actualizar la lógica del juego en respuesta a la entrada del usuario
                switch k.Name {
                case fyne.KeyDown:
                    if player.Y+player.Speed+player.Height <= int(game.CanvasHeight)-player.Height-game.Margin {
                        player.YMov = player.Speed
                    }
                    player.FrameY = player.DownY
                case fyne.KeyUp:
                    if player.Y-player.Speed >= 0 {
                        player.YMov = -player.Speed
                    }
                    player.FrameY = player.UpY
                case fyne.KeyLeft:
                    if player.X-player.Speed >= game.Margin {
                        player.XMov = -player.Speed
                    }
                    player.FrameY = player.LeftY
                case fyne.KeyRight:
                    if player.X+player.Speed+player.Width <= int(game.CanvasWidth)-game.Margin {
                        player.XMov = player.Speed
                    }
                    player.FrameY = player.RightY
                }

                playerRect := image.Rect(player.X, player.Y, player.X+player.Width, player.Y+player.Height)
                pointsRect := image.Rect(points.X, points.Y, points.X+points.Width, points.Y+points.Height)

                for _, obstacle := range obstacles {
                    obstacleRect := image.Rect(obstacle.X, obstacle.Y, obstacle.X+obstacle.Width, obstacle.Y+obstacle.Height)
                    draw.Draw(sprite, obstacleRect, obstacleImage, image.Point{}, draw.Over)
                    if playerRect.Overlaps(obstacleRect) {
                        resetPlayerPosition(&player)
                        puntos--
                        if puntos == -1 {
                            dialog.ShowInformation("Juego Terminado", "Perdiste el juego. Tu puntuación es -1.", w)
                            break
                        }
                        puntosText.SetText(fmt.Sprintf("Puntos: %d", puntos))
                    }
                }

                if !points.Collected && playerRect.Overlaps(pointsRect) {
                    rand.Seed(time.Now().UnixNano())
                    newX := rand.Intn(int(game.CanvasWidth - float32(points.Width)))
                    newY := rand.Intn(int(game.CanvasHeight - float32(points.Height)))

                    points.X = newX
                    points.Y = newY

                    puntos++
                    puntosText.SetText(fmt.Sprintf("Puntos: %d", puntos))
                }
            }
        }
    }()

    go func() {
        for {
            time.Sleep(time.Millisecond)

            now := time.Now().UnixMilli()
            elapsed := now - game.Then

            if elapsed > fpsInterval {
                game.Then = now

                spriteDP := image.Pt(player.Width*player.FrameX, player.Height*player.FrameY)
                sr := image.Rectangle{spriteDP, spriteDP.Add(spriteSize)}

                if puntos >= 5 {
                    dialog.ShowInformation("¡Felicidades!", "¡Ganaste el juego con 5 puntos!", w)
                    break
                }

                dp := image.Pt(player.X, player.Y)
                r := image.Rectangle{dp, dp.Add(spriteSize)}

                draw.Draw(sprite, sprite.Bounds(), image.Transparent, image.ZP, draw.Src)

                for _, obstacle := range obstacles {
                    obstacleRect := image.Rect(obstacle.X, obstacle.Y, obstacle.X+obstacle.Width, obstacle.Y+obstacle.Height)
                    draw.Draw(sprite, obstacleRect, obstacleImage, image.Point{}, draw.Over)
                }

                draw.Draw(sprite, r, playerSprites, sr.Min, draw.Src)
                playerImg = canvas.NewRasterFromImage(sprite)

                if player.XMov != 0 || player.YMov != 0 {
                    player.X += player.XMov
                    player.Y += player.YMov
                    player.FrameX = (player.FrameX + 1) % player.CyclesX
                    player.XMov = 0
                    player.YMov = 0
                } else {
                    player.FrameX = 0
                }
                if puntos == -1 {
                    dialog.ShowInformation("Juego Terminado", "Perdiste el juego. Tu puntuación es -1.", w)
                    break 
                }
                updateScreen <- struct{}{}
            }
            if !points.Collected {
                pointsRect := image.Rect(points.X, points.Y, points.X+points.Width, points.Y+points.Height)
                draw.Draw(sprite, pointsRect, pointsImage, image.Point{}, draw.Over)
            }
        }
    }()
    
    w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
        gameActions <- *k
    })

    go func() {
        for {
            <-updateScreen
            c.Refresh()
        }
    }()

    w.CenterOnScreen()
    w.ShowAndRun()
}
