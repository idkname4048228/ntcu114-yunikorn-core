# NTCU 114級 畢業專題

## 現在的目標： AGA (Ant-Grasshopper Algorithm)
已經完成

同時完成的還有：
- GOA 的重構

- GOA 的問題 (go 與 python 結果不相同)
  結果是因為 colab 前面的程式碼沒有執行到，所以結果會是之前的測試資料

- CAO 的開發
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


## 之前的目標： CAO (螞蟻演算法)
被要求加速進度，而且 CAO 有需設計公式的問題，短時間內解決不了

### 試著利用外接的程式檔案
實作請參考[筆記](https://hackmd.io/@ntcu-k8s/S1tny9dG0)

ACS110131 寫完了，這裡是[commit 連結](https://github.com/YouKaiWu/ntcu114-yunikorn-core/commit/16225ab0b28464b65699f67b0daae6f77e84f60f)

git commit message 請參考[這篇文章](https://wadehuanglearning.blogspot.com/2019/05/commit-commit-commit-why-what-commit.html)

[GOA 的 colab](https://colab.research.google.com/drive/1lP3TBnIsf6nGDB-Ft_8SyCyWjQOEv4wq#scrollTo=B-e9lphkix5z)

### 試著將 GOA 應用在 YuniKorn 排程
成功搞定

將 addUser, addNode 加在 partition 裡面
去 entrypoint 初始化 GOA
在 schedule 那邊執行

#### 遇到的小問題
將 colab 上 python 的 GOA 移植到 go 上面後，兩者的最終結果有些差距。
python 的可以穩定算出最佳解(目前以公平性來看)，而 go 是沒辦法算出最佳解的，甚至沒有一個相對穩定的解答。

以 user 要求 [1, 4] 及 [3, 1] 來看，當資源限制為 [9, 18] 時，都能計算出 [3, 2] 的答案。
但將資源限制提高到 [18, 28] 時，卻只有 python 可以穩定計算出 [6, 4] 的答案。

下面是 python 的視覺化
![pyGOA](https://hackmd.io/_uploads/SkGfy5AOC.gif)

下面是 go 的視覺化
![go_GOA](https://hackmd.io/_uploads/Hyu-k9CdA.gif)

之後找個時間解決這個問題。

#### GOA 應用使用的 yaml 檔
```
apiVersion: batch/v1
kind: Job
metadata:
    name: pi
    namespace: testjob
spec:
    template:
        metadata:
            labels:
                applicationId: "app1"
            annotations:
                yunikorn.apache.org/user.info: "
                {
                    \"user\": \"user1\",
                    \"groups\": [
                        \"developers\",
                        \"devops\"
                    ]
                }"
        spec:
            schedulerName: yunikorn
            containers:
              - name: pi
                image: perl:5.34.0
                resources:
                    requests:
                        memory: "256Mi"
                        cpu: "0.5"
                    limits:
                        memory: "512Mi"
                        cpu: "1"
                command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(1000)"]
            restartPolicy: Never
```
以及
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      schedulerName: yunikorn
      containers:
      - name: nginx
        image: nginx:1.14.2
        resources:
          requests:
              memory: "256Mi"
              cpu: "0.5"
          limits:
              memory: "512Mi"
              cpu: "1"
        ports:
        - containerPort: 80
```

#### 遇到的小問題
應該要重構一下 GOA ，像是 userCount 就一堆地方都有

有嘗試在 entrypoint 那邊修改排程的間隔，但好像沒有用，要再研究看看要在哪裡改

現在只排程 AllocationAsk ，應該要排程 Application ，要記得修改