## 現在的目標： AGA (Ant-Grasshopper Algorithm)
已經完成

同時完成的還有：
- 現在是排程 Application 了

- GOA 的重構

- GOA 的問題 (go 與 python 結果不相同)
  結果是因為 colab 前面的程式碼沒有執行到，所以結果會是之前的測試資料

- ACO 的開發
  python colab 在 [這裡](https://colab.research.google.com/drive/1XJ8S49GLkDJ2flS6QoQmsuqi6kbpRcPA?usp=sharing)

- scheduler 的排程間隔
  在 `pkg/scheduler.go` 下面的 `internalSchedule` 會看到我將 `case <-s.activityPending` 給 comment 起來，是因為如果讓他執行的話，他會將所有有關 RMEvent 的動作記錄下來，然後有改變的話就會觸發 scheduler
  所以我將他 comment 掉後，現在排程只會通過固定頻率排程了。缺點是排一個 pod 可能要很久 

### 遇到的問題
- 公平性的基準
  由於之前計算的是「每個使用者在每台機器上的公平分配」，所以會有一個剛好長度是 users $\times$ nodes 的座標可以讓我當「公平線」
  但問題是，在每台機器上都是公平不代表在整體系統上公平，所以所謂的「公平線」還要再研究一下

- 測試資料
  由於直接上 YuniKorn 可能不知道到底有沒有符合標準，所以還是要想一個測試資料看有沒有達到我預想的情況
  - 目前預想的情況應該會是「如果 fairness 的評分大於 0.9 ，則選擇 effect 最小的；否則，將最小 $\frac{effect}{fairness}$ 當成最佳解」
    0.9 的計算結果是暫時的，我需要的公平性情境是「解裡面最大值必須小於最小值的兩倍」，而用了幾個測資發現 threshold 是 0.9 ，所以先設定這個數字，之後再看有沒有數學證明

### 2024/09/19 修正
原先因為初始值是從直接取亂數，最大是相應使用者的 ask count，但這會導致初始解所在的位置的評分幾乎都是 inf ，也就是初始解幾乎都在沒用的地方。
因此，這次修正將亂數的值乘上 (目前第幾隻螞蟻 / 螞蟻總數)，這樣的修改會讓原始值的範圍由小到大產生，進而避免到幾乎全部都沒用的情況。
在修正後，執行結果確實從幾乎每次都全 0 ，變為正常排程數量。

這個問題的發現是因為我將 ACO 執行完的最佳解印出來，發現幾乎每次評分都是 inf ，尤其需要排成的數量越多，就越不可能出現初始解不為 inf 的情況。