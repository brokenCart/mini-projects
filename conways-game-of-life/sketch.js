const sketchWidth = 500;
const sketchHeight = 500;
const cellSize = 25;

let grid;
let playState = false;
let startStopButton;

function setup() {
  createCanvas(sketchWidth, sketchHeight);
  grid = new Grid(sketchWidth, sketchHeight, cellSize);
  grid.initiateCells();

  startStopButton = createButton("Start");
  startStopButton.mousePressed(() => {
    playState = !playState;
  });

  let restartButton = createButton("Restart State");
  restartButton.mousePressed(() => {
    playState = false;
    grid.restartState();
  });
}

function draw() {
  frameRate(5);
  background(255);
  if (playState) {
    startStopButton.html("Stop");
    grid.updateCells();
  } else {
    startStopButton.html("Start");
  }
  grid.drawCells();
}

function mouseClicked() {
  const cellYIdx = Math.floor(mouseY / grid.cellSize);
  const cellXIdx = Math.floor(mouseX / grid.cellSize);
  if (grid.validateCellIdx(cellYIdx, cellXIdx)) {
    grid.toggleCell(cellYIdx, cellXIdx);
  }
}
