# Brian's Brain

Brian's brain是康威生命游戏的一个变种，也是一个细胞自动机的经典游戏。它在康威游戏的基础上增加了一个状态，即每个cell有3个状态：

* ready
* firing
* refactory

## 规则

在一个N*M的矩阵中，每个单元格代表一个细胞，它有3种状态，同生命游戏一样，每个细胞下一阶段的状态由这一阶段它自身的状态和它周围的八个邻居的状态决定。具体规则如下：

1. 当前状态是ready，则如果周围刚好有2个firing状态的邻居，则状态变为firing，否则保持不变
1. 当前状态是firing，则变为refactory
1. 当前状态是refactory，则变为ready

详见[维基百科](https://en.wikipedia.org/wiki/Brian%27s_Brain)

## D3JS版本

* [Online Demo](//blog.zhangyu.so/brianbrain/index.html)

## Golang终端版本

* go build -o bb bb.go
* ./bb -i 200 # 默认间隔100ms

建议深色终端背景。

执行效果：

[![asciicast](https://asciinema.org/a/f1SPnBZdVPVPZAWbgnNYHOstX.svg)](https://asciinema.org/a/f1SPnBZdVPVPZAWbgnNYHOstX)
