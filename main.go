package main

import (
    "fmt"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "time"
    "image"
    "image/draw"
    "image/png"
    "os"
)

type Player struct {
    x       int
    y       int
    width   int
    height  int
    frameX  int
    frameY  int
    cyclesX int
    upY     int
    downY   int
    leftY   int
    rightY  int
    speed   int
    xMov    int
    yMov    int
}

type Obstacle struct {
    x      int
    y      int
    width  int
    height int
}


type Game struct {
    canvasWidth  float32
    canvasHeight float32
    fps          int
    then         int64
    margin       int
}

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

func main() {
    myApp := app.New()
    w := myApp.NewWindow("Game")

    obstacleImage := load("img/ramos.png")

    background := load("img/background.png")
    playerSprites := load("img/messi.png")

    now := time.Now().UnixMilli()
    game := &Game{
        800,
        500,
        10,
        now,
        4,
    }

    fpsInterval := int64(1000 / game.fps)

    obstacles := []Obstacle{
        {300, 100, 30, 30},
        {500, 250, 40, 40},
    }    

    player := &Player{100, 200, 40, 72, 0, 0, 4, 3, 0, 1, 2, 9, 0, 0}

    img := canvas.NewImageFromImage(background)
    img.FillMode = canvas.ImageFillOriginal

    sprite := image.NewRGBA(background.Bounds())

    playerImg := canvas.NewRasterFromImage(sprite)
    spriteSize := image.Pt(player.width, player.height)

    c := container.New(layout.NewMaxLayout(), img, playerImg)
    w.SetContent(c)

    w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
        switch k.Name {
        case fyne.KeyDown:
            if player.y+player.speed+player.height <= int(game.canvasHeight)-player.height-game.margin {
                player.yMov = player.speed
            }
            player.frameY = player.downY
        case fyne.KeyUp:
            if player.y-player.speed >= 0 {
                player.yMov = -player.speed
            }
            player.frameY = player.upY
        case fyne.KeyLeft:
            if player.x-player.speed >= game.margin {
                player.xMov = -player.speed
            }
            player.frameY = player.leftY
        case fyne.KeyRight:
            if player.x+player.speed+player.width <= int(game.canvasWidth)-game.margin {
                player.xMov = player.speed
            }
            player.frameY = player.rightY
        }

        for _, obstacle := range obstacles {
            obstacleRect := image.Rect(obstacle.x, obstacle.y, obstacle.x+obstacle.width, obstacle.y+obstacle.height)
            draw.Draw(sprite, obstacleRect, obstacleImage, image.Point{}, draw.Over)
        }
        
        playerRect := image.Rect(player.x, player.y, player.x+player.width, player.y+player.height)

        for _, obstacle := range obstacles {
            obstacleRect := image.Rect(obstacle.x, obstacle.y, obstacle.x+obstacle.width, obstacle.y+obstacle.height)
        
            if playerRect.Overlaps(obstacleRect) {
                // Aquí maneja la colisión, por ejemplo, resta puntos al jugador o reinicia el juego
            }
        }
    })

    go func() {
        for {
            time.Sleep(time.Millisecond)

            now := time.Now().UnixMilli()
            elapsed := now - game.then

            if elapsed > fpsInterval {
                game.then = now

                spriteDP := image.Pt(player.width*player.frameX, player.height*player.frameY)
                sr := image.Rectangle{spriteDP, spriteDP.Add(spriteSize)}

                dp := image.Pt(player.x, player.y)
                r := image.Rectangle{dp, dp.Add(spriteSize)}

                draw.Draw(sprite, sprite.Bounds(), image.Transparent, image.ZP, draw.Src)
                draw.Draw(sprite, r, playerSprites, sr.Min, draw.Src)
                playerImg = canvas.NewRasterFromImage(sprite)

                if player.xMov != 0 || player.yMov != 0 {
                    player.x += player.xMov
                    player.y += player.yMov
                    player.frameX = (player.frameX + 1) % player.cyclesX
                    player.xMov = 0
                    player.yMov = 0
                } else {
                    player.frameX = 0
                }

                c.Refresh()
            }
        }
    }()

    w.CenterOnScreen()
    w.ShowAndRun()
}
