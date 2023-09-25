package models

type Obstacle struct {
    X       int
    Y       int
    Width   int
    Height  int
    FrameX  int
    FrameY  int
    CyclesX int
    UpY     int
    DownY   int
    LeftY   int
    RightY  int
    Speed   int
    XMov    int
    YMov    int
}
