<?php

class point {
    public int $x;
    public int $y;

    public function __construct($x, $y) {
        $this->x = $x;
        $this->y = $y;
    }

    public function moveUp() {
        return new point($this->x, $this->y+1);
    }
    public function moveDown() {
        return new point($this->x, $this->y-1);
    }
    public function moveLeft() {
        return new point($this->x-1, $this->y);
    }
    public function moveRight() {
        return new point($this->x+1, $this->y);
    }
    public function hitsBoundary($size) {
        return $this->x == 0 || $this->x == $size-1
            || $this->y == 0 || $this->y == $size-1;
    }
}

$size = 100;
$current = new point(50, 50);
$points = [$current];

while (true) {
    if ($current->hitsBoundary($size)) {
        break;
    }

    $step = rand(0, 100) % 4;

    if ($step === 0) {
        $current = $current->moveLeft();
        $points[] = $current;
    } elseif ($step === 1) {
        $current = $current->moveRight();
        $points[] = $current;
    } elseif ($step === 2) {
        $current = $current->moveUp();
        $points[] = $current;
    } elseif ($step === 3) {
        $current = $current->moveDown();
        $points[] = $current;
    }
}

echo "taken ".count($points)." steps\n";
$grid = [];
for ($i = 0; $i < $size; $i++) {
    for ($j = 0; $j < $size; $j++) {
        $grid[$i][$j] = '-';
    }
}
foreach ($points as $point) {
    $grid[$point->x][$point->y] = "*";
}
$grid[50][50] = "S";
$grid[$current->x][$current->y] = "E";

for ($i = 0; $i < $size; $i++) {
    for ($j = 0; $j < $size; $j++) {
        echo $grid[$i][$j];
    }
    echo "\n";
}
