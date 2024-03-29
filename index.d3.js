const ready = 0
const firing = 1
const refactory = 2


class BrianBrain {
  constructor(rows, cols) {
    this.rows = rows
    this.cols = cols
    this._lives = []
    this._next_lives = []
  }

  initBoard = () => {
    this._lives = _.range(this.rows)
      .map(() => (
        _.range(this.cols).map(() => Math.floor(Math.random() * 3))
      ))
    this._next_lives = _.range(this.rows).map(() => _.range(this.cols))
  }

  getLives = () => this._lives

  nextRound = () => {
    for(let x = 0; x < this.rows; x++) {
      for(let y = 0; y < this.cols; y++) {
        switch(this._lives[x][y]) {
          case ready:
            const aliveNbCount = _([
              this._lives[x-1] && this._lives[x-1][y-1],
              this._lives[x-1] && this._lives[x-1][y],
              this._lives[x-1] && this._lives[x-1][y+1],
              this._lives[x][y-1],
              this._lives[x][y+1],
              this._lives[x+1] && this._lives[x+1][y-1],
              this._lives[x+1] && this._lives[x+1][y],
              this._lives[x+1] && this._lives[x+1][y+1],
            ]).filter(e => e === firing).size()
            if (aliveNbCount === 2) {
              this._next_lives[x][y] = firing
            } else {
              this._next_lives[x][y] = ready
            }
            break
          case firing:
            this._next_lives[x][y] = refactory
            break
          case refactory:
            this._next_lives[x][y] = ready
            break
          default:
            console.log('invalid cell: ' + this._lives[x][y])
        }
      }
    }

    this._lives = this._next_lives
    this._next_lives = _.range(this.rows).map(() => _.range(this.cols))
  }
}


class Board {
  constructor(container) {
    this.chart = {left: 20, top: 20, r: 4, border: 1}
    this.circleSize = this.chart.r * 2 + this.chart.border
    this.svg = d3.select(container).append('svg')
  }

  arrayToObject = (arr) => (
    _.flatten(_.range(arr.length).map(x => (
      _.range(arr[0].length).map(y => {
        return {x: x, y: y, v: arr[x][y]}
      })
    )))
  )

  render = (lives) => {
    const data = this.arrayToObject(lives)

    this.svg.attr('height', this.circleSize * lives.length + this.chart.top)
    // Enter
    this.svg.selectAll('circle')
      .data(data)
      .enter()
        .append('circle')
        .attr('cx', d => (
            d.y * this.circleSize + this.chart.left
          ))
        .attr('cy', d => (
            d.x * this.circleSize + this.chart.top
          ))
        .attr('r', this.chart.r)
        .attr('data', d => d.v)

    // Update
    this.svg.selectAll('circle')
      .data(data)
      .attr('data', d => d.v)

    // Exit
    this.svg.selectAll('circle')
      .data(data)
      .exit()
        .remove()
  }
}
