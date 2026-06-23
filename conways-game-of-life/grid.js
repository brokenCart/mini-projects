class Grid {
  constructor(screenWidth, screenHeight, cellSize) {
    this.screenWidth = screenWidth;
    this.screenHeight = screenHeight;
    this.cellSize = cellSize;
    this.cells = [];
  }

  initiateCells() {
    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      let row = [];
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        row.push(false);
      }
      this.cells.push(row);
    }
  }

  restartState() {
    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        this.cells[i][j] = false;
      }
    }
  }

  drawCells() {
    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        if (this.cells[i][j]) {
          fill("#9ACD32");
        } else {
          fill("#FFF");
        }
        square(j * this.cellSize, i * this.cellSize, this.cellSize);
      }
    }
  }

  toggleCell(cellYIdx, cellXIdx) {
    this.cells[cellYIdx][cellXIdx] = !this.cells[cellYIdx][cellXIdx];
  }

  validateCellIdx(cellYIdx, cellXIdx) {
    return (
      cellYIdx >= 0 &&
      cellYIdx < this.screenHeight / this.cellSize &&
      cellXIdx >= 0 &&
      cellXIdx < this.screenWidth / this.cellSize
    );
  }

  getActiveNeighboursCount(cellYIdx, cellXIdx) {
    let count = 0;
    let possibleCellsRelativeIdx = [
      [-1, -1],
      [-1, 0],
      [-1, 1],
      [0, -1],
      [0, 1],
      [1, -1],
      [1, 0],
      [1, 1],
    ]; // formatted as [y, x]
    for (let relativeIdx of possibleCellsRelativeIdx) {
      if (
        this.validateCellIdx(
          cellYIdx + relativeIdx[0],
          cellXIdx + relativeIdx[1]
        ) &&
        this.cells[cellYIdx + relativeIdx[0]][cellXIdx + relativeIdx[1]]
      ) {
        count++;
      }
    }
    return count;
  }

  updateCells() {
    let nextFrame = [];
    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      let row = [];
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        row.push(this.cells[i][j]);
      }
      nextFrame.push(row);
    }

    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        let activeNeighbours = this.getActiveNeighboursCount(i, j);
        if (this.cells[i][j]) {
          if (activeNeighbours > 3 || activeNeighbours < 2) {
            nextFrame[i][j] = false;
          }
        } else {
          if (activeNeighbours === 3) {
            nextFrame[i][j] = true;
          }
        }
      }
    }

    for (let i = 0; i < this.screenHeight / this.cellSize; i++) {
      for (let j = 0; j < this.screenWidth / this.cellSize; j++) {
        this.cells[i][j] = nextFrame[i][j];
      }
    }
  }
}
